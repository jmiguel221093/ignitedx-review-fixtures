package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/cache"
	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/profile"
)

type stubProfiles struct {
	fetch func(context.Context, string) (profile.Profile, error)
}

func (s stubProfiles) Fetch(ctx context.Context, id string) (profile.Profile, error) {
	return s.fetch(ctx, id)
}

func TestHandlerFetchesAndCachesProfile(t *testing.T) {
	t.Parallel()

	fetchCount := 0
	server := NewServer(stubProfiles{fetch: func(_ context.Context, id string) (profile.Profile, error) {
		fetchCount++
		return profile.Profile{ID: id, Name: "Ada"}, nil
	}}, cache.NewStore())

	for range 2 {
		request := httptest.NewRequest(http.MethodGet, "/profiles/user-1", nil)
		response := httptest.NewRecorder()
		server.Handler().ServeHTTP(response, request)
		if response.Code != http.StatusOK {
			t.Fatalf("GET /profiles status = %d, want 200", response.Code)
		}
	}
	if fetchCount != 1 {
		t.Fatalf("Fetch() called %d times, want 1", fetchCount)
	}
}

func TestHandlerMapsCancellationToRequestTimeout(t *testing.T) {
	t.Parallel()

	server := NewServer(stubProfiles{fetch: func(context.Context, string) (profile.Profile, error) {
		return profile.Profile{}, context.Canceled
	}}, cache.NewStore())
	request := httptest.NewRequest(http.MethodGet, "/profiles/user-1", nil)
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusRequestTimeout {
		t.Fatalf("GET /profiles status = %d, want 408", response.Code)
	}
}
