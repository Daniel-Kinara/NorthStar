package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"northstar-sync/internal/cache"
)

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

func (h *Handlers) AllStock(w http.ResponseWriter, r *http.Request) {
	items, updatedAt := h.store.All()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items":      items,
		"updated_at": updatedAt,
		"count":      len(items),
	})
}

// Health is liveness: "is the process running at all". Always 200 if the server can respond.
func (h *Handlers) Health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

// Ready is readiness: "has this instance ever successfully populated its cache".
// A load balancer should stop sending traffic here if this returns non-200.
func (h *Handlers) Ready(w http.ResponseWriter, r *http.Request) {
	if !h.ready() {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("not ready"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

// LoggingMiddleware logs method, path, status, and latency for every request.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start),
		)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
