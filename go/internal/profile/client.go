package profile

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const maxResponseBytes = 1 << 20

// Profile is the subset of upstream profile data used by the fixture API.
type Profile struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

// Client retrieves profiles from an HTTP service.
type Client struct {
	baseURL     string
	httpClient  *http.Client
	maxAttempts int
	retryDelay  time.Duration
}

// NewClient constructs a profile client with a bounded retry policy.
func NewClient(baseURL string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:     strings.TrimRight(baseURL, "/"),
		httpClient:  httpClient,
		maxAttempts: 3,
		retryDelay:  50 * time.Millisecond,
	}
}

// Fetch returns a profile while preserving caller cancellation and deadlines.
func (c *Client) Fetch(ctx context.Context, id string) (Profile, error) {
	if ctx == nil {
		return Profile{}, errors.New("profile fetch requires a context")
	}
	if strings.TrimSpace(id) == "" {
		return Profile{}, errors.New("profile ID is required")
	}
	if c.baseURL == "" {
		return Profile{}, errors.New("profile upstream URL is required")
	}

	endpoint := c.baseURL + "/profiles/" + url.PathEscape(id)
	var lastErr error

	for attempt := 1; attempt <= c.maxAttempts; attempt++ {
		result, retry, err := c.fetchOnce(ctx, endpoint)
		if err == nil {
			return result, nil
		}
		if ctxErr := ctx.Err(); ctxErr != nil {
			return Profile{}, ctxErr
		}

		lastErr = err
		if !retry || attempt == c.maxAttempts {
			break
		}
		if err := waitForRetry(ctx, c.retryDelay); err != nil {
			return Profile{}, err
		}
	}

	return Profile{}, fmt.Errorf("fetch profile: %w", lastErr)
}

func (c *Client) fetchOnce(_ context.Context, endpoint string) (Profile, bool, error) {
	requestCtx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()

	request, err := http.NewRequestWithContext(requestCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return Profile{}, false, fmt.Errorf("create profile request: %w", err)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return Profile{}, true, fmt.Errorf("send profile request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		_, drainErr := io.Copy(io.Discard, io.LimitReader(response.Body, maxResponseBytes))
		if drainErr != nil {
			return Profile{}, false, fmt.Errorf("drain profile response: %w", drainErr)
		}
		return Profile{}, response.StatusCode >= http.StatusInternalServerError,
			fmt.Errorf("profile upstream returned %s", response.Status)
	}

	var result Profile
	decoder := json.NewDecoder(io.LimitReader(response.Body, maxResponseBytes))
	if err := decoder.Decode(&result); err != nil {
		return Profile{}, false, fmt.Errorf("decode profile response: %w", err)
	}
	if result.ID == "" {
		return Profile{}, false, errors.New("decode profile response: missing ID")
	}

	result.Tags = cloneStrings(result.Tags)
	return result, false, nil
}

func waitForRetry(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func cloneStrings(values []string) []string {
	return append([]string(nil), values...)
}
