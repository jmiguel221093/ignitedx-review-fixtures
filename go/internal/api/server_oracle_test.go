//go:build fixture_oracle

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/cache"
	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/profile"
)

type oracleContextKey struct{}

func TestOracle_GO_CTX_002_RequestContextReachesProfileFetcher(t *testing.T) {
	const marker = "request-scope"
	observed := ""
	server := NewServer(stubProfiles{fetch: func(ctx context.Context, id string) (profile.Profile, error) {
		observed, _ = ctx.Value(oracleContextKey{}).(string)
		return profile.Profile{ID: id}, nil
	}}, cache.NewStore())

	request := httptest.NewRequest(http.MethodGet, "/profiles/user-1", nil)
	request = request.WithContext(context.WithValue(request.Context(), oracleContextKey{}, marker))
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("GET /profiles status = %d, want 200", response.Code)
	}
	if observed != marker {
		t.Fatalf("profile fetch context marker = %q, want %q", observed, marker)
	}
}
