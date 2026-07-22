package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/api"
	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/cache"
	"github.com/jmiguel221093/ignitedx-review-fixtures/go/internal/profile"
)

func main() {
	if err := run(); err != nil {
		slog.Error("fixture API stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	listenAddress := envOrDefault("LISTEN_ADDR", ":8080")
	upstreamURL := envOrDefault("PROFILE_UPSTREAM_URL", "http://127.0.0.1:9090")
	profiles := profile.NewClient(upstreamURL, &http.Client{Timeout: 5 * time.Second})
	handler := api.NewServer(profiles, cache.NewStore()).Handler()
	server := &http.Server{
		Addr:              listenAddress,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrors := make(chan error, 1)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return err
	}

	err := <-serverErrors
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
