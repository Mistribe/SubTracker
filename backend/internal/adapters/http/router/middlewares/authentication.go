package middlewares

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/pkg/ginx"
)

type AuthenticationMiddleware struct {
	idp ports.IdentityProvider
}

func NewAuthenticationMiddleware(idp ports.IdentityProvider) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{
		idp: idp,
	}
}

func getSessionToken(c *gin.Context) string {
	// First try to get the token from the Authorization header (cross-origin)
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		return strings.TrimPrefix(authHeader, "Bearer ")
	}

	// If not found in header, try to get from __session cookie (same-origin)
	cookie, err := c.Cookie("__session")
	if err != nil {
		return ""
	}
	return cookie
}

func (m AuthenticationMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		sessionToken := getSessionToken(c)
		if sessionToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		identity, err := m.idp.ReadSessionToken(c, sessionToken)
		if err != nil {
			ginx.FromError(c, err)
			return
		}

		// Store claims in context for use in handlers
		c.Set(auth.ContextIdentityKey, identity)
		c.Set(auth.ContextUserIdKey, identity.Id)

		c.Next()
	}
}
