package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/Oleexo/config-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/oleexo/subtracker/internal/domain/auth"
)

type AuthenticationMiddleware struct {
	config KindeConfig
}

func NewAuthenticationMiddleware(cfg config.Configuration) *AuthenticationMiddleware {
	domain := cfg.GetString("KINDE_AUTH_DOMAIN")
	audience := cfg.GetString("KINDE_AUDIENCE")
	return &AuthenticationMiddleware{
		config: NewKindeConfig(domain, audience),
	}
}

type KindeConfig struct {
	Domain   string
	Audience string
	Keyfunc  keyfunc.Keyfunc
}

func NewKindeConfig(domain, audience string) KindeConfig {
	jwksURL := domain + "/.well-known/jwks"

	keyFunc, err := keyfunc.NewDefaultCtx(context.Background(), []string{jwksURL})
	if err != nil {
		panic(fmt.Sprintf("Failed to get the JWKS: %s", err.Error()))
	}

	return KindeConfig{
		Domain:   domain,
		Audience: audience,
		Keyfunc:  keyFunc,
	}
}

func (m AuthenticationMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Check if the header starts with "Bearer "
		const bearerPrefix = "Bearer "
		if !strings.HasPrefix(authHeader, bearerPrefix) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header must start with Bearer"})
			c.Abort()
			return
		}

		// Extract the token
		tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
			c.Abort()
			return
		}

		// Parse and validate the token
		token, err := jwt.Parse(tokenString, m.config.Keyfunc.Keyfunc)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Failed to parse token: %s", err.Error())})
			c.Abort()
			return
		}

		// Check if token is valid
		if !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// Extract and validate claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Validate audience if specified
		if m.config.Audience != "" {
			if !m.validateAudience(claims, m.config.Audience) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid audience"})
				c.Abort()
				return
			}
		}

		// Validate issuer
		if !m.validateIssuer(claims, m.config.Domain) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid issuer"})
			c.Abort()
			return
		}

		// Store claims in context for use in handlers
		c.Set("claims", claims)
		c.Set(auth.ContextKey, claims["sub"])

		c.Next()
	}
}

// validateAudience checks if the token audience matches the expected audience
func (m AuthenticationMiddleware) validateAudience(claims jwt.MapClaims, expectedAudience string) bool {
	aud, ok := claims["aud"]
	if !ok {
		return false
	}

	switch v := aud.(type) {
	case string:
		return v == expectedAudience
	case []interface{}:
		for _, a := range v {
			if str, ok := a.(string); ok && str == expectedAudience {
				return true
			}
		}
	}
	return false
}

// validateIssuer checks if the token issuer matches the expected domain
func (m AuthenticationMiddleware) validateIssuer(claims jwt.MapClaims, expectedDomain string) bool {
	iss, ok := claims["iss"]
	if !ok {
		return false
	}

	issuer, ok := iss.(string)
	if !ok {
		return false
	}

	// The issuer should match the Kinde domain
	return issuer == expectedDomain
}
