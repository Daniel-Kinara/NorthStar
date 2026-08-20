package cache

import (
	"sync"
	"time"

	"northstar-sync/internal/warehouse"
)

// Store holds the last-known stock snapshot in memory.
type Store struct {
	mu        sync.RWMutex
	items     map[string]warehouse.StockItem
	updatedAt time.Time
}

func NewStore() *Store {
	return &Store{
		items: make(map[string]warehouse.StockItem),
	}
}

// Replace swaps the whole cache with a fresh snapshot from the warehouse.
func (s *Store) Replace(items []warehouse.StockItem) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fresh := make(map[string]warehouse.StockItem, len(items))
	for _, item := range items {
		fresh[item.SKU] = item
	}
	s.items = fresh
	s.updatedAt = time.Now()
}

func (s *Store) Get(sku string) (warehouse.StockItem, bool, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, ok := s.items[sku]
	return item, ok, s.updatedAt
}

func (s *Store) All() ([]warehouse.StockItem, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]warehouse.StockItem, 0, len(s.items))
	for _, item := range s.items {
		out = append(out, item)
	}
	return out, s.updatedAt
}
