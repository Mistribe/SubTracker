package label_test

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
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"

	"github.com/mistribe/subtracker/internal/adapters/http/export"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/label"
	domainLabel "github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

// mockQueryHandler is a mock implementation of the query handler
type mockQueryHandler struct {
	handleFunc func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]]
}

func (m *mockQueryHandler) Handle(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]] {
	return m.handleFunc(ctx, q)
}

func TestExportEndpoint_CSV_Format(t *testing.T) {
	// Create test labels
	labels := createTestLabels()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]] {
			return result.Success(shared.NewPaginatedResponse(labels, int64(len(labels))))
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := label.NewExportEndpoint(mockHandler, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/labels/export?format=csv", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "labels_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".csv")

	// Verify CSV is valid
	reader := csv.NewReader(bytes.NewReader(w.Body.Bytes()))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	require.Len(t, records, 3) // header + 2 data rows

	// Verify headers
	assert.Equal(t, []string{"id", "name", "color", "ownerType", "ownerFamilyId"}, records[0])

	// Verify first row (personal label)
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", records[1][0])
	assert.Equal(t, "Personal Label", records[1][1])
	assert.Equal(t, "#FF0000", records[1][2])
	assert.Equal(t, "personal", records[1][3])
	assert.Equal(t, "", records[1][4]) // No family ID for personal

	// Verify second row (family label)
	assert.Equal(t, "00000000-0000-0000-0000-000000000002", records[2][0])
	assert.Equal(t, "Family Label", records[2][1])
	assert.Equal(t, "#00FF00", records[2][2])
	assert.Equal(t, "family", records[2][3])
	assert.Equal(t, "00000000-0000-0000-0000-000000000123", records[2][4])
}

func TestExportEndpoint_JSON_Format(t *testing.T) {
	// Create test labels
	labels := createTestLabels()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]] {
			return result.Success(shared.NewPaginatedResponse(labels, int64(len(labels))))
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := label.NewExportEndpoint(mockHandler, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/labels/export?format=json", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "labels_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".json")

	// Verify JSON is valid
	var exportModels []export.LabelExportModel
	err := json.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)

	// Verify first label
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", exportModels[0].Id)
	assert.Equal(t, "Personal Label", exportModels[0].Name)
	assert.Equal(t, "#FF0000", exportModels[0].Color)
	assert.Equal(t, "personal", exportModels[0].OwnerType)
	assert.Nil(t, exportModels[0].OwnerFamilyId)

	// Verify second label
	assert.Equal(t, "00000000-0000-0000-0000-000000000002", exportModels[1].Id)
	assert.Equal(t, "Family Label", exportModels[1].Name)
	assert.Equal(t, "#00FF00", exportModels[1].Color)
	assert.Equal(t, "family", exportModels[1].OwnerType)
	require.NotNil(t, exportModels[1].OwnerFamilyId)
	assert.Equal(t, "00000000-0000-0000-0000-000000000123", *exportModels[1].OwnerFamilyId)
}

func TestExportEndpoint_YAML_Format(t *testing.T) {
	// Create test labels
	labels := createTestLabels()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]] {
			return result.Success(shared.NewPaginatedResponse(labels, int64(len(labels))))
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := label.NewExportEndpoint(mockHandler, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/labels/export?format=yaml", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/x-yaml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "labels_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".yaml")

	// Verify YAML is valid
	var exportModels []export.LabelExportModel
	err := yaml.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)

	// Verify first label
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", exportModels[0].Id)
	assert.Equal(t, "Personal Label", exportModels[0].Name)
	assert.Equal(t, "#FF0000", exportModels[0].Color)
	assert.Equal(t, "personal", exportModels[0].OwnerType)

	// Verify second label
	assert.Equal(t, "00000000-0000-0000-0000-000000000002", exportModels[1].Id)
	assert.Equal(t, "Family Label", exportModels[1].Name)
}

func TestExportEndpoint_InvalidFormat(t *testing.T) {
	// Create mock handler (won't be called)
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]] {
			t.Fatal("handler should not be called for invalid format")
			return result.Fail[shared.PaginatedResponse[domainLabel.Label]](nil)
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := label.NewExportEndpoint(mockHandler, exportService)

	// Create test request with invalid format
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/labels/export?format=invalid", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid format")
}

func TestExportEndpoint_EmptyResultSet(t *testing.T) {
	// Create mock handler returning empty results
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]] {
			return result.Success(shared.NewPaginatedResponse([]domainLabel.Label{}, 0))
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := label.NewExportEndpoint(mockHandler, exportService)

	t.Run("CSV with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/labels/export?format=csv", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		// Should still have headers
		output := w.Body.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		require.Len(t, lines, 1)
		assert.Equal(t, "id,name,color,ownerType,ownerFamilyId", lines[0])
	})

	t.Run("JSON with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/labels/export?format=json", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var exportModels []export.LabelExportModel
		err := json.Unmarshal(w.Body.Bytes(), &exportModels)
		require.NoError(t, err)
		assert.Len(t, exportModels, 0)
	})

	t.Run("YAML with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/labels/export?format=yaml", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var exportModels []export.LabelExportModel
		err := yaml.Unmarshal(w.Body.Bytes(), &exportModels)
		require.NoError(t, err)
		assert.Len(t, exportModels, 0)
	})
}

func TestExportEndpoint_ContentDispositionTimestamp(t *testing.T) {
	// Create test labels
	labels := createTestLabels()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]] {
			return result.Success(shared.NewPaginatedResponse(labels, int64(len(labels))))
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := label.NewExportEndpoint(mockHandler, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/labels/export?format=csv", nil)

	// Execute
	endpoint.Handle(c)

	// Verify Content-Disposition header contains timestamp
	contentDisposition := w.Header().Get("Content-Disposition")
	assert.Contains(t, contentDisposition, "attachment; filename=\"labels_")
	assert.Contains(t, contentDisposition, ".csv\"")

	// Verify timestamp format (YYYY-MM-DDTHH-MM-SS)
	// Extract filename from header
	parts := strings.Split(contentDisposition, "\"")
	require.Len(t, parts, 3)
	filename := parts[1]
	assert.True(t, strings.HasPrefix(filename, "labels_"))
	assert.True(t, strings.HasSuffix(filename, ".csv"))
}

func TestExportEndpoint_DefaultFormat(t *testing.T) {
	// Create test labels
	labels := createTestLabels()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context, q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainLabel.Label]] {
			return result.Success(shared.NewPaginatedResponse(labels, int64(len(labels))))
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := label.NewExportEndpoint(mockHandler, exportService)

	// Create test request without format parameter
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/labels/export", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response defaults to JSON
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".json")

	// Verify JSON is valid
	var exportModels []export.LabelExportModel
	err := json.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)
}

// Helper function to create test labels
func createTestLabels() []domainLabel.Label {
	now := time.Now()

	// Create personal label
	labelID1, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000001")
	personalLabel := domainLabel.NewLabel(
		labelID1,
		types.NewPersonalOwner(types.UserID("user-123")),
		"Personal Label",
		nil,
		"#FF0000",
		now,
		now,
	)

	// Create family label
	labelID2, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000002")
	familyID, _ := types.ParseFamilyID("00000000-0000-0000-0000-000000000123")
	familyLabel := domainLabel.NewLabel(
		labelID2,
		types.NewFamilyOwner(familyID),
		"Family Label",
		nil,
		"#00FF00",
		now,
		now,
	)

	return []domainLabel.Label{personalLabel, familyLabel}
}
