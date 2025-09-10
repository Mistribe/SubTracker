package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/cache"
)

type CacheMiddleware struct{}

func NewCacheMiddleware() *CacheMiddleware {
	return &CacheMiddleware{}
}

func (m *CacheMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(cache.RequestCacheKey, cache.NewRequest())
		c.Next()
	}
}
