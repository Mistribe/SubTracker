package middlewares

import (
	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/shared/i18n"
)

type LanguageMiddleware struct{}

func NewLanguageMiddleware() *LanguageMiddleware {
	return &LanguageMiddleware{}
}

func (m *LanguageMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		acceptLanguageHeader := c.Request.Header.Get("Accept-Language")
		var err error
		var info i18n.Info
		info, err = i18n.ParseAcceptLanguage(acceptLanguageHeader)
		if err != nil {
			info = i18n.GetDefault()
		}

		c.Set(i18n.ContextKey, info)

		c.Next()
	}
}
