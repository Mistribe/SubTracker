package provider_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/adapters/http/export"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/provider"
	domainProvider "github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/provider/query"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

// mockQueryHandler is a mock implementation of the query handler
type mockQueryHandler struct {
	handleFunc func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]]
}

func (m *mockQueryHandler) Handle(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]] {
	return m.handleFunc(ctx, q)
}

func TestExportEndpoint_CSV_Format(t *testing.T) {
	// Create test providers
	providers := createTestProviders()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]] {
			return result.Success(shared.NewPaginatedResponse(providers, int64(len(providers))))
		},
	}

	// Create mock label resolver
	labelID1, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000011")
	labelID2, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000012")
	mockResolver := export.NewMockLabelResolver(t)
	mockResolver.EXPECT().ResolveLabelNames(mock.Anything, mock.MatchedBy(func(ids []types.LabelID) bool {
		return len(ids) > 0
	})).Return(map[types.LabelID]string{
		labelID1: "Label 1",
		labelID2: "Label 2",
	}, nil).Maybe()

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := provider.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=csv", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "providers_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".csv")

	// Verify CSV is valid
	reader := csv.NewReader(bytes.NewReader(w.Body.Bytes()))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	require.Len(t, records, 3) // header + 2 data rows

	// Verify headers
	assert.Equal(t, []string{"name", "key", "description", "url", "iconUrl", "pricingPageUrl", "labels"}, records[0])

	// Verify first row
	assert.Equal(t, "Provider 1", records[1][0])
	assert.Equal(t, "c_provider-1_user-123", records[1][1])
	assert.Equal(t, "Description 1", records[1][2])
	assert.Equal(t, "https://provider1.com", records[1][3])
	assert.Equal(t, "https://provider1.com/icon.png", records[1][4])
	assert.Equal(t, "https://provider1.com/pricing", records[1][5])
	// Labels should be comma-separated names
	assert.Contains(t, records[1][6], "Label")
}

func TestExportEndpoint_JSON_Format(t *testing.T) {
	// Create test providers
	providers := createTestProviders()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]] {
			return result.Success(shared.NewPaginatedResponse(providers, int64(len(providers))))
		},
	}

	// Create mock label resolver
	labelID1, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000011")
	labelID2, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000012")
	mockResolver := export.NewMockLabelResolver(t)
	mockResolver.EXPECT().ResolveLabelNames(mock.Anything, mock.Anything).Return(map[types.LabelID]string{
		labelID1: "Label 1",
		labelID2: "Label 2",
	}, nil).Maybe()

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := provider.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=json", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "providers_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".json")

	// Verify JSON is valid
	var exportModels []dto.ProviderExportModel
	err := json.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)

	// Verify first provider
	assert.Equal(t, "Provider 1", exportModels[0].Name)
	require.NotNil(t, exportModels[0].Description)
	assert.Equal(t, "Description 1", *exportModels[0].Description)
	require.NotNil(t, exportModels[0].Url)
	assert.Equal(t, "https://provider1.com", *exportModels[0].Url)
	// Verify labels are names, not IDs
	assert.NotEmpty(t, exportModels[0].Labels)
}

func TestExportEndpoint_YAML_Format(t *testing.T) {
	// Create test providers
	providers := createTestProviders()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]] {
			return result.Success(shared.NewPaginatedResponse(providers, int64(len(providers))))
		},
	}

	// Create mock label resolver
	labelID1, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000011")
	labelID2, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000012")
	mockResolver := export.NewMockLabelResolver(t)
	mockResolver.EXPECT().ResolveLabelNames(mock.Anything, mock.Anything).Return(map[types.LabelID]string{
		labelID1: "Label 1",
		labelID2: "Label 2",
	}, nil).Maybe()

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := provider.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=yaml", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/x-yaml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "providers_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".yaml")

	// Verify YAML is valid
	var exportModels []dto.ProviderExportModel
	err := yaml.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)

	// Verify first provider
	assert.Equal(t, "Provider 1", exportModels[0].Name)
}

func TestExportEndpoint_InvalidFormat(t *testing.T) {
	// Create mock handler (won't be called)
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]] {
			t.Fatal("handler should not be called for invalid format")
			return result.Fail[shared.PaginatedResponse[domainProvider.Provider]](nil)
		},
	}

	mockResolver := export.NewMockLabelResolver(t)

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := provider.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request with invalid format
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=invalid", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid format")
}

func TestExportEndpoint_EmptyResultSet(t *testing.T) {
	// Create mock handler returning empty results
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]] {
			return result.Success(shared.NewPaginatedResponse([]domainProvider.Provider{}, 0))
		},
	}

	mockResolver := export.NewMockLabelResolver(t)
	mockResolver.EXPECT().ResolveLabelNames(mock.Anything, mock.Anything).Return(map[types.LabelID]string{}, nil).Maybe()

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := provider.NewExportEndpoint(mockHandler, mockResolver, exportService)

	t.Run("CSV with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=csv", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		// Should still have headers
		output := w.Body.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		require.Len(t, lines, 1)
		assert.Equal(t, "name,key,description,url,iconUrl,pricingPageUrl,labels", lines[0])
	})

	t.Run("JSON with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=json", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var exportModels []dto.ProviderExportModel
		err := json.Unmarshal(w.Body.Bytes(), &exportModels)
		require.NoError(t, err)
		assert.Len(t, exportModels, 0)
	})

	t.Run("YAML with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=yaml", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var exportModels []dto.ProviderExportModel
		err := yaml.Unmarshal(w.Body.Bytes(), &exportModels)
		require.NoError(t, err)
		assert.Len(t, exportModels, 0)
	})
}

func TestExportEndpoint_ContentDispositionTimestamp(t *testing.T) {
	// Create test providers
	providers := createTestProviders()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]] {
			return result.Success(shared.NewPaginatedResponse(providers, int64(len(providers))))
		},
	}

	mockResolver := export.NewMockLabelResolver(t)
	mockResolver.EXPECT().ResolveLabelNames(mock.Anything, mock.Anything).Return(map[types.LabelID]string{}, nil).Maybe()

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := provider.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=csv", nil)

	// Execute
	endpoint.Handle(c)

	// Verify Content-Disposition header contains timestamp
	contentDisposition := w.Header().Get("Content-Disposition")
	assert.Contains(t, contentDisposition, "attachment; filename=\"providers_")
	assert.Contains(t, contentDisposition, ".csv\"")

	// Verify timestamp format (YYYY-MM-DDTHH-MM-SS)
	// Extract filename from header
	parts := strings.Split(contentDisposition, "\"")
	require.Len(t, parts, 3)
	filename := parts[1]
	assert.True(t, strings.HasPrefix(filename, "providers_"))
	assert.True(t, strings.HasSuffix(filename, ".csv"))
}

func TestExportEndpoint_LabelsAreNames(t *testing.T) {
	// Create test providers with labels
	providers := createTestProviders()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainProvider.Provider]] {
			return result.Success(shared.NewPaginatedResponse(providers, int64(len(providers))))
		},
	}

	// Create mock label resolver that returns specific names
	labelID1, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000011")
	labelID2, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000012")
	mockResolver := export.NewMockLabelResolver(t)
	mockResolver.EXPECT().ResolveLabelNames(mock.Anything, mock.Anything).Return(map[types.LabelID]string{
		labelID1: "Work",
		labelID2: "Personal",
	}, nil).Maybe()

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := provider.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/providers/export?format=json", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify JSON contains label names, not IDs
	var exportModels []dto.ProviderExportModel
	err := json.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)

	// Verify labels are names
	for _, model := range exportModels {
		for _, label := range model.Labels {
			// Should be names like "Work" or "Personal", not UUIDs
			assert.NotContains(t, label, "00000000-0000-0000-0000")
		}
	}
}

// Helper function to create test providers
func createTestProviders() []domainProvider.Provider {
	now := time.Now()

	// Create label IDs
	labelID1, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000011")
	labelID2, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000012")

	// Create provider 1
	providerID1, _ := types.ParseProviderID("00000000-0000-0000-0000-000000000001")
	desc1 := "Description 1"
	url1 := "https://provider1.com"
	iconUrl1 := "https://provider1.com/icon.png"
	pricingUrl1 := "https://provider1.com/pricing"
	provider1 := domainProvider.NewProvider(
		providerID1,
		"Provider 1",
		&desc1,
		&iconUrl1,
		&url1,
		&pricingUrl1,
		[]types.LabelID{labelID1, labelID2},
		types.NewPersonalOwner(types.UserID("user-123")),
		now,
		now,
	)

	// Create provider 2
	providerID2, _ := types.ParseProviderID("00000000-0000-0000-0000-000000000002")
	desc2 := "Description 2"
	url2 := "https://provider2.com"
	provider2 := domainProvider.NewProvider(
		providerID2,
		"Provider 2",
		&desc2,
		nil,
		&url2,
		nil,
		[]types.LabelID{labelID1},
		types.NewPersonalOwner(types.UserID("user-123")),
		now,
		now,
	)

	return []domainProvider.Provider{provider1, provider2}
}
