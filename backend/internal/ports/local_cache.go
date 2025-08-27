package ports

import (
	"time"
)

type CacheOptions struct {
	Duration time.Duration
}

func WithDuration(duration time.Duration) func(*CacheOptions) {
	return func(opts *CacheOptions) {
		opts.Duration = duration
	}
}

type LocalCache interface {
	Set(key string, value interface{}, options ...func(*CacheOptions))
	Get(key string) interface{}
	PurgeExpired()
}
