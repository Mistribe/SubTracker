package kinde

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

// mockConfig implements a simple configuration for testing
type mockConfig struct {
	values map[string]string
}

func (m *mockConfig) GetString(key string) string {
	return m.values[key]
}

func (m *mockConfig) GetStringOrDefault(key, defaultValue string) string {
	if val, ok := m.values[key]; ok {
		return val
	}
	return defaultValue
}

// Additional methods to satisfy the config.Configuration interface
func (m *mockConfig) GetInt(key string) int64 {
	if val, ok := m.values[key]; ok {
		i, _ := strconv.ParseInt(val, 10, 64)
		return i
	}
	return 0
}

func (m *mockConfig) GetIntOrDefault(key string, defaultValue int64) int64 {
	if val, ok := m.values[key]; ok {
		i, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			return i
		}
	}
	return defaultValue
}

func (m *mockConfig) GetBool(key string) bool {
	if val, ok := m.values[key]; ok {
		b, _ := strconv.ParseBool(val)
		return b
	}
	return false
}

func (m *mockConfig) GetBoolOrDefault(key string, defaultValue bool) bool {
	if val, ok := m.values[key]; ok {
		b, err := strconv.ParseBool(val)
		if err == nil {
			return b
		}
	}
	return defaultValue
}

func (m *mockConfig) GetFloat(key string) float64 {
	if val, ok := m.values[key]; ok {
		f, _ := strconv.ParseFloat(val, 64)
		return f
	}
	return 0
}

func (m *mockConfig) GetFloatOrDefault(key string, defaultValue float64) float64 {
	if val, ok := m.values[key]; ok {
		f, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return f
		}
	}
	return defaultValue
}

func TestTokenCaching(t *testing.T) {
	// Create a mock HTTP server that simulates the token endpoint
	var requestCount int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		// Return a token with 60 seconds expiration
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"access_token":"test-token-%d","token_type":"Bearer","expires_in":60}`, requestCount)
	}))
	defer server.Close()

	// Create a mock configuration
	cfg := &mockConfig{
		values: map[string]string{
			"KINDE_AUTH_DOMAIN":       server.URL,
			"KINDE_API_CLIENT_ID":     "test-client-id",
			"KINDE_API_CLIENT_SECRET": "test-client-secret",
		},
	}

	// Create a token generator
	generator := NewTokenGenerator(cfg)

	// First request should make an HTTP call
	token1, err := generator.GetToken()
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if token1 != "test-token-1" {
		t.Errorf("Expected token 'test-token-1', got '%s'", token1)
	}
	if requestCount != 1 {
		t.Errorf("Expected 1 HTTP request, got %d", requestCount)
	}

	// Second immediate request should use the cached token
	token2, err := generator.GetToken()
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if token2 != "test-token-1" {
		t.Errorf("Expected cached token 'test-token-1', got '%s'", token2)
	}
	if requestCount != 1 {
		t.Errorf("Expected still 1 HTTP request, got %d", requestCount)
	}

	// Access the kindeService directly to manipulate the expiration time
	ks, ok := generator.(*kindeService)
	if !ok {
		t.Fatal("Failed to cast generator to kindeService")
	}

	// Set the token as expired
	ks.tokenExpiresAt = time.Now().Add(-1 * time.Second)

	// Third request after expiration should make a new HTTP call
	token3, err := generator.GetToken()
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if token3 != "test-token-2" {
		t.Errorf("Expected new token 'test-token-2', got '%s'", token3)
	}
	if requestCount != 2 {
		t.Errorf("Expected 2 HTTP requests, got %d", requestCount)
	}

	// Test with a zero expiration time using a new server
	server.Close()

	// Create a new server for the zero expiration test
	zeroExpRequestCount := 0
	zeroExpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		zeroExpRequestCount++
		// Return a token with 0 seconds expiration
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"access_token":"zero-token-%d","token_type":"Bearer","expires_in":0}`, zeroExpRequestCount)
	}))
	defer zeroExpServer.Close()

	// Update the configuration to use the new server
	cfg.values["KINDE_AUTH_DOMAIN"] = zeroExpServer.URL

	// Create a new token generator with the updated config
	zeroExpGenerator := NewTokenGenerator(cfg)

	// First request should make an HTTP call
	token4, err := zeroExpGenerator.GetToken()
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if token4 != "zero-token-1" {
		t.Errorf("Expected token 'zero-token-1', got '%s'", token4)
	}
	if zeroExpRequestCount != 1 {
		t.Errorf("Expected 1 HTTP request, got %d", zeroExpRequestCount)
	}

	// Second request should also make a new HTTP call since the token wasn't cached (expires_in=0)
	token5, err := zeroExpGenerator.GetToken()
	if err != nil {
		t.Fatalf("Failed to get token: %v", err)
	}
	if token5 != "zero-token-2" {
		t.Errorf("Expected new token 'zero-token-2', got '%s'", token5)
	}
	if zeroExpRequestCount != 2 {
		t.Errorf("Expected 2 HTTP requests, got %d", zeroExpRequestCount)
	}
}
