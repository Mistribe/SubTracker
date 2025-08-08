package kinde

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Oleexo/config-go"

	openapi "github.com/oleexo/subtracker/internal/infrastructure/kinde/gen"
)

// TokenGenerator is an interface for generating authentication tokens and providing API clients
// The implementation includes token caching based on the ExpiresIn value:
// - Tokens with ExpiresIn > 0 are cached until they expire (with a small buffer to ensure they don't expire during use)
// - Tokens with ExpiresIn = 0 are not cached, and each call to GetToken() will make a new request
type TokenGenerator interface {
	GetToken() (string, error)
	GetClient() *openapi.APIClient
}

func MakeRequest[TOut any](
	generator TokenGenerator,
	call func(ctx context.Context, client *openapi.APIClient) (TOut, *http.Response, error)) (TOut, error) {
	client := generator.GetClient()
	ctx := context.Background()
	token, err := generator.GetToken()
	if err != nil {
		var rOut TOut
		return rOut, err
	}

	ctx = context.WithValue(ctx, openapi.ContextAccessToken, token)

	out, httpResp, err := call(ctx, client)
	if err != nil {
		var rOut TOut
		return rOut, err
	}
	defer httpResp.Body.Close()

	return out, nil
}

type kindeService struct {
	domain       string
	clientId     string
	clientSecret string
	audiences    []string

	// Token caching
	cachedToken    string
	tokenExpiresAt time.Time
	tokenMutex     sync.Mutex
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

func (k *kindeService) GetClient() *openapi.APIClient {
	// Create configuration
	cfg := openapi.NewConfiguration()

	// Set server URL (customize the subdomain)
	cfg.Servers = openapi.ServerConfigurations{
		{
			URL:         k.domain,
			Description: "Kinde API",
		},
	}

	// Create API client
	client := openapi.NewAPIClient(cfg)

	return client
}

func (k *kindeService) GetToken() (string, error) {
	// Lock for thread safety
	k.tokenMutex.Lock()
	defer k.tokenMutex.Unlock()

	// Check if we have a cached token that's still valid
	if k.cachedToken != "" && time.Now().Before(k.tokenExpiresAt) {
		return k.cachedToken, nil
	}

	// No valid cached token, request a new one
	tokenURL := fmt.Sprintf("%s/oauth2/token", k.domain)

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", k.clientId)
	data.Set("client_secret", k.clientSecret)
	data.Set("audience", strings.Join(k.audiences, " "))

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	// Cache the token with its expiration time only if ExpiresIn is greater than 0
	if tokenResp.ExpiresIn > 0 {
		// Apply a small buffer (10 seconds) to ensure we don't use a token that's about to expire
		k.cachedToken = tokenResp.AccessToken
		k.tokenExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn-10) * time.Second)
	} else {
		// If ExpiresIn is 0 or negative, don't cache the token
		k.cachedToken = ""
		k.tokenExpiresAt = time.Time{}
	}

	return tokenResp.AccessToken, nil
}

func NewTokenGenerator(cfg config.Configuration) TokenGenerator {
	domain := cfg.GetString("KINDE_AUTH_DOMAIN")
	clientId := cfg.GetString("KINDE_API_CLIENT_ID")
	clientSecret := cfg.GetString("KINDE_API_CLIENT_SECRET")
	audienceValue := cfg.GetStringOrDefault("KINDE_API_AUDIENCES", "https://api.kinde.com/api")
	audiences := strings.Split(audienceValue, ",")
	return &kindeService{
		domain:       domain,
		clientId:     clientId,
		clientSecret: clientSecret,
		audiences:    audiences,
		// Initialize cache-related fields with zero values
		cachedToken:    "",
		tokenExpiresAt: time.Time{},
	}
}
