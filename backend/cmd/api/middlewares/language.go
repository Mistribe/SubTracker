package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/domain/lang"
)

type LanguageMiddleware struct{}

func NewLanguageMiddleware() *LanguageMiddleware {
	return &LanguageMiddleware{}
}

func (m *LanguageMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLanguageHeader := c.Request.Header.Get("Accept-Language")
		var err error
		var info lang.Info
		info, err = lang.ParseAcceptLanguage(acceptLanguageHeader)
		if err != nil {
			info = lang.GetDefault()
		}

		c.Set(lang.ContextKey, info)

		c.Next()
	}
}
