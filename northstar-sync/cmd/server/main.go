package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"northstar-sync/internal/api"
	"northstar-sync/internal/cache"
	"northstar-sync/internal/config"
	"northstar-sync/internal/poller"
	"northstar-sync/internal/warehouse"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		slog.Error("invalid configuration", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	store := cache.NewStore()
	client := warehouse.NewClient(cfg.WarehouseAPIURL, cfg.RequestTimeout, cfg.MaxRetries)
	p := poller.New(client, store, cfg.PollInterval)

	go p.Start(ctx)

	h := api.NewHandlers(store, store.Ready)
	mux := http.NewServeMux()
	mux.HandleFunc("/stock", h.StockBySKU)
	mux.HandleFunc("/stock/all", h.AllStock)
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/ready", h.Ready)

	server := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      api.LoggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		slog.Info("server starting",
			"port", cfg.ServerPort,
			"warehouseURL", cfg.WarehouseAPIURL,
			"pollInterval", cfg.PollInterval)

		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	slog.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	slog.Info("shutdown complete")
}
