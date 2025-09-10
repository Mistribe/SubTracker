package cache

import (
	"context"

	"github.com/mistribe/subtracker/internal/ports"
)

func New(localCache LocalCache,
	distributedCache DistributedCache) ports.Cache {
	return &cache{
		localCache:       localCache,
		distributedCache: distributedCache,
	}
}

type cache struct {
	localCache       LocalCache
	distributedCache DistributedCache
}

func (c cache) From(ctx context.Context, level ports.CacheLevel) ports.CacheLeveled {
	switch level {
	case ports.CacheLevelRequest:
		return fromContext(ctx)
	case ports.CacheLevelServer:
		return c.localCache
	case ports.CacheLevelDistributed:
		return c.distributedCache
	default:
		panic("unknown cache level")
	}
}

func (c cache) Set(ctx context.Context, key string, value interface{}, options ...func(*ports.CacheOptions)) {
	c.From(ctx, ports.CacheLevelRequest).Set(key, value, options...)
	c.From(ctx, ports.CacheLevelServer).Set(key, value, options...)
	c.From(ctx, ports.CacheLevelDistributed).Set(key, value, options...)
}

func (c cache) Get(ctx context.Context, key string) interface{} {
	var value interface{}
	value = c.From(ctx, ports.CacheLevelRequest).Get(key)
	if value != nil {
		return value
	}
	value = c.From(ctx, ports.CacheLevelServer).Get(key)
	if value != nil {
		return value
	}
	return c.From(ctx, ports.CacheLevelDistributed).Get(key)
}
