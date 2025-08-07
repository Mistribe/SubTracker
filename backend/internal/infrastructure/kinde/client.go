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
)

type Client interface {
	GetToken(ctx context.Context) (string, error)
	UpdateUser(ctx context.Context, givenName, familyName, email string) error
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

func (k kindeService) GetToken(ctx context.Context) (string, error) {
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

func NewClient(cfg config.Configuration) Client {
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
