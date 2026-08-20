package api

import (
	"encoding/json"
	"net/http"

	"northstar-sync/internal/cache"
)

type Handlers struct {
	store *cache.Store
}

func NewHandlers(store *cache.Store) *Handlers {
	return &Handlers{store: store}
}

// GET /stock?sku=NS-1001
func (h *Handlers) StockBySKU(w http.ResponseWriter, r *http.Request) {
	sku := r.URL.Query().Get("sku")
	if sku == "" {
		http.Error(w, "missing sku query parameter", http.StatusBadRequest)
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

// GET /stock/all
func (h *Handlers) AllStock(w http.ResponseWriter, r *http.Request) {
	items, updatedAt := h.store.All()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":      items,
		"updated_at": updatedAt,
		"count":      len(items),
	})
}

// GET /health
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
