//go:build fixture_oracle

package jobs

import (
	"context"
	"testing"
	"time"
)

func TestOracle_GO_CONC_001_CancellationClosesResults(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	started := make(chan struct{})
	results := (Runner{}).Start(ctx, []Job{{
		ID: "wait-for-cancel",
		Run: func(ctx context.Context) error {
			close(started)
			<-ctx.Done()
			return ctx.Err()
		},
	}})

	<-started
	cancel()
	assertChannelCloses(t, results)
}

func TestOracle_GO_CONC_002_ProducerClosesResults(t *testing.T) {
	results := (Runner{}).Start(context.Background(), []Job{{
		ID:  "complete",
		Run: func(context.Context) error { return nil },
	}})
	assertChannelCloses(t, results)
}

func assertChannelCloses(t *testing.T, results <-chan Result) {
	t.Helper()
	timeout := time.After(time.Second)
	for {
		select {
		case _, ok := <-results:
			if !ok {
				return
			}
		case <-timeout:
			t.Fatal("results channel did not close")
		}
	}
}
