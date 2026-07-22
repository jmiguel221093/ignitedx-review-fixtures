package jobs

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

// Job is a unit of context-aware work.
type Job struct {
	ID  string
	Run func(context.Context) error
}

// Result captures the outcome of one job.
type Result struct {
	ID  string `json:"id"`
	Err error  `json:"-"`
}

// Runner starts jobs and owns the lifecycle of the returned results channel.
type Runner struct{}

// Start executes jobs concurrently and closes the result channel after every
// worker exits. The buffer prevents an abandoned consumer from trapping a
// worker while cancellation propagates.
func (Runner) Start(ctx context.Context, pending []Job) <-chan Result {
	results := make(chan Result, len(pending))
	var workers sync.WaitGroup
	workers.Add(len(pending))

	for _, pendingJob := range pending {
		job := pendingJob
		go func() {
			defer workers.Done()

			result := Result{ID: job.ID}
			if job.Run == nil {
				result.Err = errors.New("job has no runner")
			} else {
				result.Err = job.Run(ctx)
			}

			select {
			case results <- result:
			case <-ctx.Done():
			}
		}()
	}

	go func() {
		workers.Wait()
		close(results)
	}()

	return results
}

// Collect drains a producer-owned result channel until completion or caller
// cancellation.
func Collect(ctx context.Context, results <-chan Result) ([]Result, error) {
	collected := make([]Result, 0)
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result, ok := <-results:
			if !ok {
				return collected, nil
			}
			collected = append(collected, result)
		}
	}
}

// Validate returns a combined error for failed jobs.
func Validate(results []Result) error {
	for _, result := range results {
		if result.Err != nil {
			return fmt.Errorf("job %q failed: %w", result.ID, result.Err)
		}
	}
	return nil
}
