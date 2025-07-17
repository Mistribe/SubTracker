package middlewares

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"sync"

	"github.com/Oleexo/config-go"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticationMiddleware struct {
	config KindeConfig
}

func NewAuthenticationMiddleware(cfg config.Configuration) *AuthenticationMiddleware {
	return &AuthenticationMiddleware{config: KindeConfig{
		Domain:   "",
		Audience: "",
		Issuer:   "",
	}}
}

type KindeConfig struct {
	Domain   string
	Audience string
	Issuer   string
}

func (m AuthenticationMiddleware) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := extractToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
			c.Abort()
			return
		}

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			cert, err := getCertificate(m.config.Domain, token.Header["kid"].(string))
			if err != nil {
				return nil, err
			}

			return cert, nil
		})

		if err != nil || !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		if !validateClaims(claims, m.config) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Next()
	}
}

func extractToken(c *gin.Context) (string, error) {
	bearerToken := c.GetHeader("Authorization")
	if bearerToken == "" {
		return "", errors.New("no token found")
	}

	parts := strings.Split(bearerToken, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid token format")
	}

	return parts[1], nil
}

func validateClaims(claims jwt.MapClaims, config KindeConfig) bool {
	aud, ok := claims["aud"].(string)
	if !ok || aud != config.Audience {
		return false
	}

	iss, ok := claims["iss"].(string)
	if !ok || iss != config.Issuer {
		return false
	}

	return true
}

func getCertificate(domain, kid string) (interface{}, error) {
	// Implements certificate retrieval and caching from the Kinde JWKS endpoint
	jwksURL := fmt.Sprintf("https://%s/.well-known/jwks.json", domain)
	key, err := getPublicKeyFromJWKS(jwksURL, kid)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// JWK and JWKS related types
type jwkKey struct {
	Kid string   `json:"kid"`
	Kty string   `json:"kty"`
	Alg string   `json:"alg"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

type jwks struct {
	Keys []jwkKey `json:"keys"`
}

// Caching of JWKS by URL
var (
	jwksCache   = make(map[string]jwks)
	jwksCacheMu sync.RWMutex
)

// Get public key for kid from JWKS endpoint
func getPublicKeyFromJWKS(jwksURL, kid string) (*rsa.PublicKey, error) {
	jwksCacheMu.RLock()
	jwksData, cached := jwksCache[jwksURL]
	jwksCacheMu.RUnlock()

	// Only fetch if missing or expired (TODO: Expiry/refresh logic, for now only in-memory cache)
	if !cached {
		resp, err := http.Get(jwksURL)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch JWKS: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected JWKS response status: %d", resp.StatusCode)
		}

		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&jwksData); err != nil {
			return nil, fmt.Errorf("failed to parse JWKS: %w", err)
		}

		jwksCacheMu.Lock()
		jwksCache[jwksURL] = jwksData
		jwksCacheMu.Unlock()
	}

	// Find the correct key by kid
	for _, key := range jwksData.Keys {
		if key.Kid == kid {
			// Use X5C if present (recommended by Kinde docs), fall back to N/E if missing
			if len(key.X5c) > 0 {
				// The certificate is base64 DER encoded, decode and parse
				certDER, err := base64.StdEncoding.DecodeString(key.X5c[0])
				if err != nil {
					return nil, fmt.Errorf("failed to decode x5c cert: %w", err)
				}
				cert, err := x509.ParseCertificate(certDER)
				if err != nil {
					return nil, fmt.Errorf("failed to parse x5c cert: %w", err)
				}
				pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
				if !ok {
					return nil, errors.New("public key from x5c is not RSA")
				}
				return pubKey, nil
			}
			// Else: use N/E (should not happen with Kinde, but implemented for completeness)
			nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
			if err != nil {
				return nil, fmt.Errorf("failed to decode n: %w", err)
			}
			eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
			if err != nil {
				return nil, fmt.Errorf("failed to decode e: %w", err)
			}
			eInt := 0
			for i := 0; i < len(eBytes); i++ {
				eInt = eInt<<8 + int(eBytes[i])
			}
			pubKey := &rsa.PublicKey{
				N: new(big.Int).SetBytes(nBytes),
				E: eInt,
			}
			return pubKey, nil
		}
	}
	return nil, fmt.Errorf("key with kid='%s' not found in JWKS", kid)
}
