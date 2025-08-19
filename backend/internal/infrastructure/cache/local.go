package cache

import (
	"sync"
	"time"

	"github.com/Oleexo/config-go"
)

type LocalOptions struct {
	Duration time.Duration
}

func WithDuration(duration time.Duration) func(*LocalOptions) {
	return func(opts *LocalOptions) {
		opts.Duration = duration
	}
}

type Local interface {
	Set(key string, value interface{}, options ...func(*LocalOptions))
	Get(key string) interface{}
}

type item struct {
	value     interface{}
	expiresAt time.Time // zero time means no expiration
}

type local struct {
	mu         sync.RWMutex
	items      map[string]item
	defaultTTL time.Duration
	stopCh     chan struct{}
}

// NewLocal creates a new in-memory cache with a default TTL for all entries
// and a cleanup interval for purging expired items.
// If defaultTTL <= 0, items never expire.
// If cleanupInterval > 0, a background janitor will periodically purge expired items.
func NewLocal(cfg config.Configuration) Local {
	defaultTTL := time.Duration(cfg.GetIntOrDefault("CACHE_DEFAULT_TTL", 0))
	cleanupInterval := time.Duration(cfg.GetIntOrDefault("CACHE_CLEANUP_INTERVAL", int64(1*time.Hour)))
	l := &local{
		items:      make(map[string]item),
		defaultTTL: defaultTTL,
		stopCh:     make(chan struct{}),
	}
	if cleanupInterval > 0 {
		go l.janitor(cleanupInterval)
	}
	return l
}

func (l *local) janitor(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			l.purgeExpired()
		case <-l.stopCh:
			return
		}
	}
}

func (l *local) purgeExpired() {
	now := time.Now()
	l.mu.Lock()
	for k, it := range l.items {
		if !it.expiresAt.IsZero() && now.After(it.expiresAt) {
			delete(l.items, k)
		}
	}
	l.mu.Unlock()
}

func (l *local) Set(key string, value interface{}, optionsFunc ...func(*LocalOptions)) {
	var exp time.Time
	var options LocalOptions
	for _, opt := range optionsFunc {
		opt(&options)
	}
	if options.Duration > 0 {
		exp = time.Now().Add(options.Duration)
	} else if l.defaultTTL > 0 {
		exp = time.Now().Add(l.defaultTTL)
	}

	l.mu.Lock()
	l.items[key] = item{value: value, expiresAt: exp}
	l.mu.Unlock()
}

func (l *local) Get(key string) interface{} {
	l.mu.RLock()
	it, ok := l.items[key]
	l.mu.RUnlock()
	if !ok {
		return nil
	}
	if !it.expiresAt.IsZero() && time.Now().After(it.expiresAt) {
		// delete expired item and miss
		l.mu.Lock()
		// verify same item still present
		cur, ok := l.items[key]
		if ok && cur.expiresAt.Equal(it.expiresAt) {
			delete(l.items, key)
		}
		l.mu.Unlock()
		return nil
	}
	return it.value
}
