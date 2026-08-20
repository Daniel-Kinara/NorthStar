package warehouse

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"net/http"
	"time"
)

type StockItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

type Client struct {
	baseURL    string
	httpClient *http.Client
	maxRetries int
}

func NewClient(baseURL string, timeout time.Duration, maxRetries int) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{Timeout: timeout},
		maxRetries: maxRetries,
	}
}

// FetchStock retries transient failures with exponential backoff:
// attempt 1 immediate, then waits 500ms, 1s, 2s, ... capped at 8s.
func (c *Client) FetchStock(ctx context.Context) ([]StockItem, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(math.Min(float64(500*time.Millisecond)*math.Pow(2, float64(attempt-1)), float64(8*time.Second)))
			slog.Warn("retrying warehouse fetch", "attempt", attempt, "backoff", backoff, "lastError", lastErr)

			select {
			case <-time.After(backoff):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}

		items, err := c.doFetch(ctx)
		if err == nil {
			return items, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("warehouse fetch failed after %d attempts: %w", c.maxRetries+1, lastErr)
}

func (c *Client) doFetch(ctx context.Context) ([]StockItem, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL, nil)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("warehouse request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("warehouse returned status %d", resp.StatusCode)
	}

	var items []StockItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("decoding warehouse response: %w", err)
	}

	return items, nil
}
