package warehouse

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_FetchStock_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"sku":"X1","quantity":3}]`))
	}))
	defer srv.Close()

	c := NewClient(srv.URL, 2*time.Second, 2)
	items, err := c.FetchStock(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].SKU != "X1" {
		t.Fatalf("unexpected items: %+v", items)
	}
}

func TestClient_FetchStock_RetriesThenFails(t *testing.T) {
	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, 2*time.Second, 2) // 1 initial + 2 retries = 3 calls
	_, err := c.FetchStock(context.Background())
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	if calls != 3 {
		t.Fatalf("expected 3 attempts, got %d", calls)
	}
}
