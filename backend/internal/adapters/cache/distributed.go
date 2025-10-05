package cache

import (
	"github.com/Oleexo/config-go"

	"github.com/mistribe/subtracker/internal/ports"
)

type DistributedCache interface {
	ports.CacheLeveled
}

func NewDistributed(cfg config.Configuration) DistributedCache {
	return &distributed{}
}

type distributed struct {
}

func (d distributed) Set(key string, value interface{}, options ...func(*ports.CacheOptions)) {
	// do nothing for now
}

func (d distributed) Get(key string) interface{} {
	// do nothing for now
	return nil
}
