package horoscopes

import (
	"context"
	"log/slog"
	"sync"
	"time"
)

type PublicHoroscope struct {
	ID        string    `json:"id"`
	Horoscope string    `json:"horoscope"`
	CreatedAt time.Time `json:"created_at"`
}

type cacheEntry struct {
	Horoscope PublicHoroscope
	expiresAt time.Time
}

type HoroscopeCache struct {
	entries map[string]cacheEntry
	mu      sync.RWMutex
	ttl     time.Duration
	maxSize int
}

func NewHoroscopeCache(maxSize int, ttl time.Duration) *HoroscopeCache {
	return &HoroscopeCache{
		maxSize: maxSize,
		ttl:     ttl,
		entries: make(map[string]cacheEntry),
	}
}
func (hc *HoroscopeCache) Get(id string) (PublicHoroscope, bool) {
	hc.mu.RLock()
	defer hc.mu.RUnlock()

	entry, ok := hc.entries[id]
	if !ok || time.Now().After(entry.expiresAt) {
		return PublicHoroscope{}, false // let the background job delete it
	}

	slog.Info("Getting new horoscope from cache", "id", id)
	return PublicHoroscope{
		ID:        id,
		Horoscope: entry.Horoscope.Horoscope,
		CreatedAt: entry.Horoscope.CreatedAt,
	}, true

}

func (hc *HoroscopeCache) Set(h PublicHoroscope) {
	hc.mu.Lock()
	defer hc.mu.Unlock()

	// Check if we hit the limit size, evict something else at random
	if len(hc.entries) >= hc.maxSize {
		for id := range hc.entries {
			slog.Info("Cache is full, deleting horoscope in cache ", "id", id, "horoscope", hc.entries[id])
			delete(hc.entries, id)
			break
		}
	}

	slog.Info("Inserting new horoscope in cache ", "id", h.ID, "horoscope", h)

	// Insert the new one
	hc.entries[h.ID] = cacheEntry{expiresAt: time.Now().Add(hc.ttl), Horoscope: h}
}

func (hc *HoroscopeCache) BackgroundCleanup(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 30)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			hc.mu.Lock()

			for id, entry := range hc.entries {
				if entry.expiresAt.Before(time.Now()) {
					slog.Info("Deleting horoscope from cache due to expiry", "id", id)
					delete(hc.entries, id)
				}
			}
			hc.mu.Unlock()

		case <-ctx.Done():
			return
		}
	}
}
