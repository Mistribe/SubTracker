package cache

import (
	"context"
	"sync"
	"time"

	"github.com/mistribe/subtracker/internal/ports"
)

const (
	RequestCacheKey = "X-Cache-Request"
)

type RequestCache interface {
	ports.CacheLeveled
}

func NewRequest() RequestCache {
	return &request{
		items: make(map[string]item),
	}
}

type request struct {
	mu    sync.RWMutex
	items map[string]item
}

func fromContext(ctx context.Context) RequestCache {
	return ctx.Value(RequestCacheKey).(RequestCache)
}

func (r *request) Set(key string, value interface{}, optionsFunc ...func(*ports.CacheOptions)) {
	var options ports.CacheOptions
	for _, opt := range optionsFunc {
		opt(&options)
	}
	exp := options.ExpiresAt(0)

	r.mu.Lock()
	r.items[key] = item{
		value:     value,
		expiresAt: exp,
	}
	r.mu.Unlock()
}

func (r *request) Get(key string) interface{} {
	r.mu.RLock()
	it, ok := r.items[key]
	r.mu.RUnlock()
	if !ok {
		return nil
	}
	if !it.expiresAt.IsZero() && time.Now().After(it.expiresAt) {
		r.mu.Lock()
		cur, ok := r.items[key]
		if ok && cur.expiresAt.Equal(it.expiresAt) {
			delete(r.items, key)
		}
		r.mu.Unlock()
		return nil
	}
	return it.value
}
