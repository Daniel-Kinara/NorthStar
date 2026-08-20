package poller

import (
	"log"
	"time"

	"northstar-sync/internal/cache"
	"northstar-sync/internal/warehouse"
)

// Poller repeatedly pulls stock from the warehouse and refreshes the cache.
type Poller struct {
	client   *warehouse.Client
	store    *cache.Store
	interval time.Duration
	stopCh   chan struct{}
}

func New(client *warehouse.Client, store *cache.Store, interval time.Duration) *Poller {
	return &Poller{
		client:   client,
		store:    store,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start runs one poll immediately, then keeps polling on a ticker in the background.
func (p *Poller) Start() {
	p.runOnce()

	ticker := time.NewTicker(p.interval)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				p.runOnce()
			case <-p.stopCh:
				return
			}
		}
	}()
}

func (p *Poller) Stop() {
	close(p.stopCh)
}

func (p *Poller) runOnce() {
	items, err := p.client.FetchStock()
	if err != nil {
		log.Printf("poll failed: %v", err)
		return
	}
	p.store.Replace(items)
	log.Printf("poll succeeded: cached %d items", len(items))
}
