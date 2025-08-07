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

	"github.com/Oleexo/config-go"

	openapi "github.com/oleexo/subtracker/internal/infrastructure/kinde/gen"
)

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
			URL:         k.domain + "/api",
			Description: "Kinde API",
		},
	}

	// Create API client
	client := openapi.NewAPIClient(cfg)

	return client
}

func (k *kindeService) GetToken() (string, error) {
	tokenURL := fmt.Sprintf("https://%s/oauth2/token", k.domain)

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

	return tokenResp.AccessToken, nil

}

func NewTokenGenerator(cfg config.Configuration) TokenGenerator {
	domain := cfg.GetString("KINDE_AUTH_DOMAIN")
	clientId := cfg.GetString("KINDE_AUTH_CLIENT_ID")
	clientSecret := cfg.GetString("KINDE_AUTH_CLIENT_SECRET")
	audiences := cfg.GetStringOrDefault("KINDE_AUTH_AUDIENCES", "https://api.kinde.com/api")

	return &kindeService{
		domain:       domain,
		clientId:     clientId,
		clientSecret: clientSecret,
		audiences:    []string{audiences},
	}
}
