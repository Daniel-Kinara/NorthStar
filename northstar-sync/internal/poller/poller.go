package poller

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"northstar-sync/internal/cache"
	"northstar-sync/internal/warehouse"
)

type Poller struct {
	client   *warehouse.Client
	store    *cache.Store
	interval time.Duration

	mu      sync.RWMutex
	lastErr error
}

func New(client *warehouse.Client, store *cache.Store, interval time.Duration) *Poller {
	return &Poller{client: client, store: store, interval: interval}
}

// Start runs one poll immediately, then keeps polling until ctx is cancelled.
// It blocks the caller until the background goroutine has exited — call it with `go p.Start(ctx)`.
func (p *Poller) Start(ctx context.Context) {
	p.runOnce(ctx)

	ticker := time.NewTicker(p.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.runOnce(ctx)
		case <-ctx.Done():
			slog.Info("poller stopping", "reason", ctx.Err())
			return
		}
	}
}

func (p *Poller) LastError() error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.lastErr
}

func (p *Poller) runOnce(ctx context.Context) {
	items, err := p.client.FetchStock(ctx)

	p.mu.Lock()
	p.lastErr = err
	p.mu.Unlock()

	if err != nil {
		slog.Error("poll failed", "error", err)
		return
	}
	p.store.Replace(items)
	slog.Info("poll succeeded", "itemCount", len(items))
}
