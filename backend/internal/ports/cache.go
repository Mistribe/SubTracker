package ports

import (
	"context"
	"time"
)

type CacheLevel uint8

const (
	CacheLevelRequest CacheLevel = iota
	CacheLevelServer
	CacheLevelDistributed
)

type CacheOptions struct {
	Duration time.Duration
}

func (opts *CacheOptions) ExpiresAt(defaultTTL time.Duration) time.Time {
	if opts.Duration == 0 {
		if defaultTTL == 0 {
			return time.Time{}
		}
		return time.Now().Add(defaultTTL)
	}
	return time.Now().Add(opts.Duration)
}

func WithDuration(duration time.Duration) func(*CacheOptions) {
	return func(opts *CacheOptions) {
		opts.Duration = duration
	}
}

type CacheLeveled interface {
	Set(key string, value interface{}, options ...func(*CacheOptions))
	Get(key string) interface{}
}

type Cache interface {
	From(ctx context.Context, level CacheLevel) CacheLeveled
	Set(ctx context.Context, key string, value interface{}, options ...func(*CacheOptions))
	Get(ctx context.Context, key string) interface{}
}
