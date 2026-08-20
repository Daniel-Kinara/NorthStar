package cache

import (
	"testing"

	"northstar-sync/internal/warehouse"
)

func TestStore_ReplaceAndGet(t *testing.T) {
	s := NewStore()

	if s.Ready() {
		t.Fatal("expected store to be not-ready before any Replace call")
	}

	s.Replace([]warehouse.StockItem{
		{SKU: "A1", Quantity: 10},
		{SKU: "A2", Quantity: 0},
	})

	if !s.Ready() {
		t.Fatal("expected store to be ready after Replace")
	}

	item, found, _ := s.Get("A1")
	if !found || item.Quantity != 10 {
		t.Fatalf("expected A1 quantity 10, got found=%v quantity=%d", found, item.Quantity)
	}

	_, found, _ = s.Get("missing-sku")
	if found {
		t.Fatal("expected missing-sku to not be found")
	}
}

func TestStore_ReplaceOverwritesPreviousSnapshot(t *testing.T) {
	s := NewStore()
	s.Replace([]warehouse.StockItem{{SKU: "OLD", Quantity: 5}})
	s.Replace([]warehouse.StockItem{{SKU: "NEW", Quantity: 1}})

	_, found, _ := s.Get("OLD")
	if found {
		t.Fatal("expected OLD sku to be gone after second Replace")
	}
}
