//go:build fixture_oracle

package profile

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestOracle_GO_ERR_001_RequestContextPropagates(t *testing.T) {
	started := make(chan struct{})
	observedCancellation := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		close(started)
		<-r.Context().Done()
		close(observedCancellation)
	}))
	defer server.Close()

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		_, err := NewClient(server.URL, server.Client()).Fetch(ctx, "user-1")
		result <- err
	}()

	select {
	case <-started:
	case <-time.After(time.Second):
		t.Fatal("upstream request did not start")
	}
	cancel()

	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Fetch() error = %v, want context.Canceled", err)
		}
	case <-time.After(time.Second):
		t.Fatal("Fetch() did not return after cancellation")
	}

	select {
	case <-observedCancellation:
	case <-time.After(time.Second):
		t.Fatal("upstream request context was not cancelled")
	}
}

func TestOracle_GO_ERR_002_DecodeFailureIsReturned(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte(`{"id":`))
	}))
	defer server.Close()

	profile, err := NewClient(server.URL, server.Client()).Fetch(context.Background(), "user-1")
	if err == nil {
		t.Fatalf("Fetch() = %#v, nil; want decode error", profile)
	}
}

func TestOracle_GO_CTX_001_RetryWaitHonorsCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "temporary failure", http.StatusServiceUnavailable)
	}))
	defer server.Close()

	client := NewClient(server.URL, server.Client())
	client.retryDelay = time.Second
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	started := time.Now()
	_, err := client.Fetch(ctx, "user-1")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Fetch() error = %v, want context.DeadlineExceeded", err)
	}
	if elapsed := time.Since(started); elapsed > 250*time.Millisecond {
		t.Fatalf("Fetch() respected retry delay instead of cancellation: %s", elapsed)
	}
}
