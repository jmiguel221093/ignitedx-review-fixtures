package profile

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClientFetchReturnsProfile(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/profiles/user-1" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"user-1","name":"Ada","tags":["admin"]}`))
	}))
	defer server.Close()

	result, err := NewClient(server.URL, server.Client()).Fetch(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("Fetch() error = %v", err)
	}
	if result.ID != "user-1" || result.Name != "Ada" {
		t.Fatalf("Fetch() result = %#v", result)
	}
}

func TestClientFetchReturnsCancelledContext(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := NewClient("http://127.0.0.1", nil).Fetch(ctx, "user-1")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Fetch() error = %v, want context.Canceled", err)
	}
}
