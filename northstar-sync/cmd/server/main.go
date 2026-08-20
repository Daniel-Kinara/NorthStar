package main

import (
	"log"
	"net/http"

	"northstar-sync/internal/api"
	"northstar-sync/internal/cache"
	"northstar-sync/internal/config"
	"northstar-sync/internal/poller"
	"northstar-sync/internal/warehouse"
)

func main() {
	cfg := config.Load()

	store := cache.NewStore()
	client := warehouse.NewClient(cfg.WarehouseAPIURL)
	p := poller.New(client, store, cfg.PollInterval)
	p.Start()
	defer p.Stop()

	h := api.NewHandlers(store)
	mux := http.NewServeMux()
	mux.HandleFunc("/stock", h.StockBySKU)
	mux.HandleFunc("/stock/all", h.AllStock)
	mux.HandleFunc("/health", h.Health)

	log.Printf("Meridian sync service listening on :%s (polling %s every %s)",
		cfg.ServerPort, cfg.WarehouseAPIURL, cfg.PollInterval)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, mux))
}
