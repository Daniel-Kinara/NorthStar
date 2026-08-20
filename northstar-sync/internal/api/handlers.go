package api

import (
	"encoding/json"
	"net/http"
	"regexp"

	"northstar-sync/internal/cache"
)

var skuPattern = regexp.MustCompile(`^[A-Za-z0-9\-]{1,40}$`)

type Handlers struct {
	store *cache.Store
	ready func() bool
}

func NewHandlers(store *cache.Store, ready func() bool) *Handlers {
	return &Handlers{store: store, ready: ready}
}

func (h *Handlers) StockBySKU(w http.ResponseWriter, r *http.Request) {
	sku := r.URL.Query().Get("sku")
	if sku == "" {
		http.Error(w, `{"error":"missing sku query parameter"}`, http.StatusBadRequest)
		return
	}
	if !skuPattern.MatchString(sku) {
		http.Error(w, `{"error":"invalid sku format"}`, http.StatusBadRequest)
		return
	}

	item, found, updatedAt := h.store.Get(sku)
	w.Header().Set("Content-Type", "application/json")

	if !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "sku not found in cache"})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"sku":        item.SKU,
		"quantity":   item.Quantity,
		"updated_at": updatedAt,
	})
}

func (h *Handlers) AllStock(w http.ResponseWriter, r *http.Request) {
	items, updatedAt := h.store.All()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":      items,
		"updated_at": updatedAt,
		"count":      len(items),
	})
}

func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handlers) Ready(w http.ResponseWriter, r *http.Request) {
	if !h.ready() {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}
