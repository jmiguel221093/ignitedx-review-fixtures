package jobs

import (
	"context"
	"errors"
	"testing"
)

func TestRunnerCollectsResults(t *testing.T) {
	t.Parallel()

	pending := []Job{
		{ID: "one", Run: func(context.Context) error { return nil }},
		{ID: "two", Run: func(context.Context) error { return errors.New("failed") }},
	}
	resultChannel := (Runner{}).Start(context.Background(), pending)
	results := []Result{<-resultChannel, <-resultChannel}
	if len(results) != len(pending) {
		t.Fatalf("Runner returned %d results, want %d", len(results), len(pending))
	}
	if err := Validate(results); err == nil {
		t.Fatal("Validate() error = nil, want job failure")
	}
}

func TestRunnerRejectsMissingFunction(t *testing.T) {
	t.Parallel()

	result := <-(Runner{}).Start(
		context.Background(),
		[]Job{{ID: "missing"}},
	)
	if err := Validate([]Result{result}); err == nil {
		t.Fatal("Validate() error = nil, want missing runner error")
	}
}
