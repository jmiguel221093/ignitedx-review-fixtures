package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/cache"
	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/jobs"
	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/profile"
)

const maxRequestBytes = 1 << 20

type profileFetcher interface {
	Fetch(context.Context, string) (profile.Profile, error)
}

// Server provides the fixture HTTP handlers.
type Server struct {
	profiles profileFetcher
	cache    *cache.Store
	runner   jobs.Runner
}

// NewServer constructs a fixture server with explicit dependencies.
func NewServer(profiles profileFetcher, store *cache.Store) *Server {
	if store == nil {
		store = cache.NewStore()
	}
	return &Server{profiles: profiles, cache: store}
}

// Handler returns the standard-library HTTP router.
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.handleHealth)
	mux.HandleFunc("GET /profiles/{id}", s.handleProfile)
	mux.HandleFunc("GET /profiles/recent", s.handleRecentProfile)
	mux.HandleFunc("POST /jobs", s.handleJobs)
	return mux
}

func (s *Server) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) handleProfile(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(r.PathValue("id"))
	if id == "" {
		writeError(w, http.StatusBadRequest, "profile ID is required")
		return
	}

	if cached, ok := s.cache.Get(id); ok {
		writeJSON(w, http.StatusOK, cached)
		return
	}
	if s.profiles == nil {
		writeError(w, http.StatusServiceUnavailable, "profile service is unavailable")
		return
	}

	result, err := s.profiles.Fetch(context.Background(), id)
	if err != nil {
		writeFetchError(w, err)
		return
	}
	if err := s.cache.Put(result); err != nil {
		writeError(w, http.StatusInternalServerError, "profile could not be cached")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleRecentProfile(w http.ResponseWriter, _ *http.Request) {
	id, ok := s.cache.MostRecent()
	if !ok {
		writeError(w, http.StatusNotFound, "no profiles have been cached")
		return
	}
	result, ok := s.cache.Get(id)
	if !ok {
		writeError(w, http.StatusNotFound, "recent profile is unavailable")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleJobs(w http.ResponseWriter, r *http.Request) {
	if s.profiles == nil {
		writeError(w, http.StatusServiceUnavailable, "profile service is unavailable")
		return
	}

	var request struct {
		ProfileIDs []string `json:"profile_ids"`
	}
	decoder := json.NewDecoder(io.LimitReader(r.Body, maxRequestBytes))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeError(w, http.StatusBadRequest, "invalid job request")
		return
	}
	if len(request.ProfileIDs) == 0 || len(request.ProfileIDs) > 20 {
		writeError(w, http.StatusBadRequest, "profile_ids must contain between 1 and 20 entries")
		return
	}

	pending := make([]jobs.Job, 0, len(request.ProfileIDs))
	for _, profileID := range request.ProfileIDs {
		id := strings.TrimSpace(profileID)
		if id == "" {
			writeError(w, http.StatusBadRequest, "profile IDs cannot be empty")
			return
		}
		pending = append(pending, jobs.Job{
			ID: id,
			Run: func(ctx context.Context) error {
				result, err := s.profiles.Fetch(ctx, id)
				if err != nil {
					return err
				}
				return s.cache.Put(result)
			},
		})
	}

	results, err := jobs.Collect(r.Context(), s.runner.Start(r.Context(), pending))
	if err != nil {
		writeFetchError(w, err)
		return
	}

	type jobResult struct {
		ID    string `json:"id"`
		Error string `json:"error,omitempty"`
	}
	response := make([]jobResult, 0, len(results))
	for _, result := range results {
		item := jobResult{ID: result.ID}
		if result.Err != nil {
			item.Error = result.Err.Error()
		}
		response = append(response, item)
	}
	writeJSON(w, http.StatusOK, map[string]any{"results": response})
}

func writeFetchError(w http.ResponseWriter, err error) {
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		writeError(w, http.StatusRequestTimeout, "request cancelled")
		return
	}
	writeError(w, http.StatusBadGateway, "profile upstream failed")
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	payload, err := json.Marshal(value)
	if err != nil {
		http.Error(w, "response encoding failed", http.StatusInternalServerError)
		return
	}
	payload = append(payload, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(payload); err != nil {
		return
	}
}
