package warehouse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// StockItem mirrors one row from the warehouse API response.
type StockItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

// Client fetches current stock levels from the warehouse API.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// FetchStock calls the warehouse API and returns the current stock list.
func (c *Client) FetchStock() ([]StockItem, error) {
	resp, err := c.httpClient.Get(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("warehouse request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("warehouse returned status %d", resp.StatusCode)
	}

	var items []StockItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, fmt.Errorf("failed to decode warehouse response: %w", err)
	}

	return items, nil
}