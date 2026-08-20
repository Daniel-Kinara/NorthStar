package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type StockItem struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
}

func main() {
	http.HandleFunc("/stock", func(w http.ResponseWriter, r *http.Request) {
		items := []StockItem{
			{SKU: "NS-1001", Quantity: 42},
			{SKU: "NS-1002", Quantity: 0},
			{SKU: "NS-1003", Quantity: 17},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
	})

	log.Println("Mock warehouse API listening on :9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}