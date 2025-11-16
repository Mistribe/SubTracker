package subscription_test

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
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/subscription"
	"github.com/mistribe/subtracker/internal/domain/currency"
	domainSubscription "github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

// mockQueryHandler is a mock implementation of the query handler
type mockQueryHandler struct {
	handleFunc func(ctx context.Context,
		q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]]
}

func (m *mockQueryHandler) Handle(ctx context.Context,
	q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]] {
	return m.handleFunc(ctx, q)
}

// mockLabelResolver is a mock implementation of the label resolver
type mockLabelResolver struct {
	resolveFunc func(ctx context.Context, labelIds []types.LabelID) ([]string, error)
}

func (m *mockLabelResolver) ResolveLabelNames(ctx context.Context, labelIds []types.LabelID) ([]string, error) {
	return m.resolveFunc(ctx, labelIds)
}

func TestExportEndpoint_CSV_Format(t *testing.T) {
	// Create test subscriptions
	subscriptions := createTestSubscriptions()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context,
			q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]] {
			return result.Success(shared.NewPaginatedResponse(subscriptions, int64(len(subscriptions))))
		},
	}

	// Create mock label resolver
	mockResolver := &mockLabelResolver{
		resolveFunc: func(ctx context.Context, labelIds []types.LabelID) ([]string, error) {
			names := make([]string, len(labelIds))
			for i, id := range labelIds {
				if id.String() == "00000000-0000-0000-0000-000000000011" {
					names[i] = "Label 1"
				} else if id.String() == "00000000-0000-0000-0000-000000000012" {
					names[i] = "Label 2"
				}
			}
			return names, nil
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := subscription.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=csv", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "text/csv; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "subscriptions_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".csv")

	// Verify CSV is valid
	reader := csv.NewReader(bytes.NewReader(w.Body.Bytes()))
	records, err := reader.ReadAll()
	require.NoError(t, err)
	require.Len(t, records, 3) // header + 2 data rows

	// Verify headers
	expectedHeaders := []string{"id", "providerId", "friendlyName", "startDate", "endDate", "recurrency",
		"customRecurrency", "customPriceAmount", "customPriceCurrency", "ownerType", "freeTrialStartDate",
		"freeTrialEndDate", "labels"}
	assert.Equal(t, expectedHeaders, records[0])

	// Verify first row
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", records[1][0])
	assert.Equal(t, "00000000-0000-0000-0000-000000000101", records[1][1])
	assert.Equal(t, "My Subscription", records[1][2])
	// Date should be in ISO 8601 format (YYYY-MM-DD)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, records[1][3])
	assert.Equal(t, "monthly", records[1][5])
	assert.Equal(t, "9.99", records[1][7])
	assert.Equal(t, "USD", records[1][8])
	assert.Equal(t, "personal", records[1][9])
}

func TestExportEndpoint_JSON_Format(t *testing.T) {
	// Create test subscriptions
	subscriptions := createTestSubscriptions()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context,
			q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]] {
			return result.Success(shared.NewPaginatedResponse(subscriptions, int64(len(subscriptions))))
		},
	}

	// Create mock label resolver
	mockResolver := &mockLabelResolver{
		resolveFunc: func(ctx context.Context, labelIds []types.LabelID) ([]string, error) {
			names := make([]string, len(labelIds))
			for i, id := range labelIds {
				if id.String() == "00000000-0000-0000-0000-000000000011" {
					names[i] = "Label 1"
				} else if id.String() == "00000000-0000-0000-0000-000000000012" {
					names[i] = "Label 2"
				}
			}
			return names, nil
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := subscription.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=json", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "subscriptions_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".json")

	// Verify JSON is valid
	var exportModels []export.SubscriptionExportModel
	err := json.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)

	// Verify first subscription
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", exportModels[0].Id)
	assert.Equal(t, "00000000-0000-0000-0000-000000000101", exportModels[0].ProviderId)
	require.NotNil(t, exportModels[0].FriendlyName)
	assert.Equal(t, "My Subscription", *exportModels[0].FriendlyName)
	// Date should be in ISO 8601 format (YYYY-MM-DD)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, exportModels[0].StartDate)
	assert.Equal(t, "monthly", exportModels[0].Recurrency)
	assert.Equal(t, 9.99, exportModels[0].Amount)
	assert.Equal(t, "USD", exportModels[0].Currency)
	assert.Equal(t, "personal", exportModels[0].OwnerType)
	// Verify labels are names, not IDs
	assert.NotEmpty(t, exportModels[0].Labels)
}

func TestExportEndpoint_YAML_Format(t *testing.T) {
	// Create test subscriptions
	subscriptions := createTestSubscriptions()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context,
			q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]] {
			return result.Success(shared.NewPaginatedResponse(subscriptions, int64(len(subscriptions))))
		},
	}

	// Create mock label resolver
	mockResolver := &mockLabelResolver{
		resolveFunc: func(ctx context.Context, labelIds []types.LabelID) ([]string, error) {
			names := make([]string, len(labelIds))
			for i, id := range labelIds {
				if id.String() == "00000000-0000-0000-0000-000000000011" {
					names[i] = "Label 1"
				} else if id.String() == "00000000-0000-0000-0000-000000000012" {
					names[i] = "Label 2"
				}
			}
			return names, nil
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := subscription.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=yaml", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/x-yaml; charset=utf-8", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "subscriptions_")
	assert.Contains(t, w.Header().Get("Content-Disposition"), ".yaml")

	// Verify YAML is valid
	var exportModels []export.SubscriptionExportModel
	err := yaml.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)

	// Verify first subscription
	assert.Equal(t, "00000000-0000-0000-0000-000000000001", exportModels[0].Id)
	assert.Equal(t, "00000000-0000-0000-0000-000000000101", exportModels[0].ProviderId)
}

func TestExportEndpoint_InvalidFormat(t *testing.T) {
	// Create mock handler (won't be called)
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context,
			q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]] {
			t.Fatal("handler should not be called for invalid format")
			return result.Fail[shared.PaginatedResponse[domainSubscription.Subscription]](nil)
		},
	}

	mockResolver := &mockLabelResolver{}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := subscription.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request with invalid format
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=invalid", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "invalid format")
}

func TestExportEndpoint_EmptyResultSet(t *testing.T) {
	// Create mock handler returning empty results
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context,
			q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]] {
			return result.Success(shared.NewPaginatedResponse([]domainSubscription.Subscription{}, 0))
		},
	}

	mockResolver := &mockLabelResolver{
		resolveFunc: func(ctx context.Context, labelIds []types.LabelID) ([]string, error) {
			return []string{}, nil
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := subscription.NewExportEndpoint(mockHandler, mockResolver, exportService)

	t.Run("CSV with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=csv", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		// Should still have headers
		output := w.Body.String()
		lines := strings.Split(strings.TrimSpace(output), "\n")
		require.Len(t, lines, 1)
		assert.Equal(t,
			"id,providerId,friendlyName,startDate,endDate,recurrency,customRecurrency,customPriceAmount,customPriceCurrency,ownerType,freeTrialStartDate,freeTrialEndDate,labels",
			lines[0])
	})

	t.Run("JSON with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=json", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var exportModels []export.SubscriptionExportModel
		err := json.Unmarshal(w.Body.Bytes(), &exportModels)
		require.NoError(t, err)
		assert.Len(t, exportModels, 0)
	})

	t.Run("YAML with empty data", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=yaml", nil)

		endpoint.Handle(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var exportModels []export.SubscriptionExportModel
		err := yaml.Unmarshal(w.Body.Bytes(), &exportModels)
		require.NoError(t, err)
		assert.Len(t, exportModels, 0)
	})
}

func TestExportEndpoint_FreeTrialDates(t *testing.T) {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)

	// Create subscription with free trial
	subID1, _ := types.ParseSubscriptionID("00000000-0000-0000-0000-000000000001")
	providerID1, _ := types.ParseProviderID("00000000-0000-0000-0000-000000000101")
	friendlyName := "Subscription with Trial"

	freeTrialStart := now.AddDate(0, -1, 0)
	freeTrialEnd := now.AddDate(0, 0, 7)
	freeTrial := domainSubscription.NewFreeTrial(freeTrialStart, freeTrialEnd)

	price := domainSubscription.NewPrice(currency.NewAmount(9.99, currency.USD))

	sub1 := domainSubscription.NewSubscription(
		subID1,
		&friendlyName,
		freeTrial,
		providerID1,
		price,
		types.NewPersonalOwner(types.UserID("user-123")),
		nil,
		nil,
		[]domainSubscription.LabelRef{},
		startDate,
		nil,
		domainSubscription.MonthlyRecurrency,
		nil,
		now,
		now,
	)

	// Create subscription without free trial
	subID2, _ := types.ParseSubscriptionID("00000000-0000-0000-0000-000000000002")
	providerID2, _ := types.ParseProviderID("00000000-0000-0000-0000-000000000102")

	sub2 := domainSubscription.NewSubscription(
		subID2,
		nil,
		nil,
		providerID2,
		price,
		types.NewPersonalOwner(types.UserID("user-123")),
		nil,
		nil,
		[]domainSubscription.LabelRef{},
		startDate,
		nil,
		domainSubscription.MonthlyRecurrency,
		nil,
		now,
		now,
	)

	subscriptions := []domainSubscription.Subscription{sub1, sub2}

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context,
			q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]] {
			return result.Success(shared.NewPaginatedResponse(subscriptions, int64(len(subscriptions))))
		},
	}

	mockResolver := &mockLabelResolver{
		resolveFunc: func(ctx context.Context, labelIds []types.LabelID) ([]string, error) {
			return []string{}, nil
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := subscription.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=json", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	var exportModels []export.SubscriptionExportModel
	err := json.Unmarshal(w.Body.Bytes(), &exportModels)
	require.NoError(t, err)
	require.Len(t, exportModels, 2)

	// First subscription should have free trial dates
	assert.NotNil(t, exportModels[0].FreeTrialStartDate)
	assert.NotNil(t, exportModels[0].FreeTrialEndDate)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, *exportModels[0].FreeTrialStartDate)
	assert.Regexp(t, `^\d{4}-\d{2}-\d{2}$`, *exportModels[0].FreeTrialEndDate)

	// Second subscription should not have free trial dates
	assert.Nil(t, exportModels[1].FreeTrialStartDate)
	assert.Nil(t, exportModels[1].FreeTrialEndDate)
}

func TestExportEndpoint_LabelsAreNames(t *testing.T) {
	// Create test subscriptions with labels
	subscriptions := createTestSubscriptions()

	// Create mock handler
	mockHandler := &mockQueryHandler{
		handleFunc: func(ctx context.Context,
			q query.FindAllQuery) result.Result[shared.PaginatedResponse[domainSubscription.Subscription]] {
			return result.Success(shared.NewPaginatedResponse(subscriptions, int64(len(subscriptions))))
		},
	}

	// Create mock label resolver that returns specific names
	mockResolver := &mockLabelResolver{
		resolveFunc: func(ctx context.Context, labelIds []types.LabelID) ([]string, error) {
			names := make([]string, len(labelIds))
			for i, id := range labelIds {
				if id.String() == "00000000-0000-0000-0000-000000000011" {
					names[i] = "Work"
				} else if id.String() == "00000000-0000-0000-0000-000000000012" {
					names[i] = "Personal"
				}
			}
			return names, nil
		},
	}

	// Create endpoint
	exportService := export.NewExportService()
	endpoint := subscription.NewExportEndpoint(mockHandler, mockResolver, exportService)

	// Create test request
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodGet, "/subscriptions/export?format=json", nil)

	// Execute
	endpoint.Handle(c)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify JSON contains label names, not IDs
	var exportModels []export.SubscriptionExportModel
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

// Helper function to create test subscriptions
func createTestSubscriptions() []domainSubscription.Subscription {
	now := time.Now()
	startDate := now.AddDate(0, -1, 0)

	// Create label IDs
	labelID1, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000011")
	labelID2, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000012")

	// Create subscription 1
	subID1, _ := types.ParseSubscriptionID("00000000-0000-0000-0000-000000000001")
	providerID1, _ := types.ParseProviderID("00000000-0000-0000-0000-000000000101")
	friendlyName1 := "My Subscription"

	price1 := domainSubscription.NewPrice(currency.NewAmount(9.99, currency.USD))

	labelRefs1 := []domainSubscription.LabelRef{
		{LabelId: labelID1, Source: domainSubscription.LabelSourceSubscription},
		{LabelId: labelID2, Source: domainSubscription.LabelSourceSubscription},
	}

	sub1 := domainSubscription.NewSubscription(
		subID1,
		&friendlyName1,
		nil,
		providerID1,
		price1,
		types.NewPersonalOwner(types.UserID("user-123")),
		nil,
		nil,
		labelRefs1,
		startDate,
		nil,
		domainSubscription.MonthlyRecurrency,
		nil,
		now,
		now,
	)

	// Create subscription 2
	subID2, _ := types.ParseSubscriptionID("00000000-0000-0000-0000-000000000002")
	providerID2, _ := types.ParseProviderID("00000000-0000-0000-0000-000000000102")

	price2 := domainSubscription.NewPrice(currency.NewAmount(19.99, currency.EUR))

	labelRefs2 := []domainSubscription.LabelRef{
		{LabelId: labelID1, Source: domainSubscription.LabelSourceSubscription},
	}

	endDate := now.AddDate(1, 0, 0)
	familyID, _ := types.ParseFamilyID("00000000-0000-0000-0000-000000000456")
	sub2 := domainSubscription.NewSubscription(
		subID2,
		nil,
		nil,
		providerID2,
		price2,
		types.NewFamilyOwner(familyID),
		nil,
		nil,
		labelRefs2,
		startDate,
		&endDate,
		domainSubscription.YearlyRecurrency,
		nil,
		now,
		now,
	)

	return []domainSubscription.Subscription{sub1, sub2}
}
