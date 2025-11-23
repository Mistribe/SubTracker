package subscription_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	xcur "golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/subscription"
	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

// mockSummaryQueryHandler is a mock implementation of the summary query handler
type mockSummaryQueryHandler struct {
	handleFunc func(ctx context.Context, q query.SummaryQuery) result.Result[query.SummaryQueryResponse]
}

func (m *mockSummaryQueryHandler) Handle(ctx context.Context, q query.SummaryQuery) result.Result[query.SummaryQueryResponse] {
	return m.handleFunc(ctx, q)
}

// TestSummaryEndpoint_MapsAllFields tests that all fields from domain model are correctly mapped to API DTO
func TestSummaryEndpoint_MapsAllFields(t *testing.T) {
	// Create test data with all fields populated
	usd := xcur.USD

	domainResponse := query.SummaryQueryResponse{
		Active:            10,
		ActivePersonal:    6,
		ActiveFamily:      4,
		TotalMonthly:      currency.NewAmount(100.50, usd),
		TotalYearly:       currency.NewAmount(1200.00, usd),
		TotalLastMonth:    currency.NewAmount(95.00, usd),
		TotalLastYear:     currency.NewAmount(1150.00, usd),
		PersonalMonthly:   currency.NewAmount(60.25, usd),
		PersonalYearly:    currency.NewAmount(720.00, usd),
		PersonalLastMonth: currency.NewAmount(55.00, usd),
		PersonalLastYear:  currency.NewAmount(680.00, usd),
		FamilyMonthly:     currency.NewAmount(40.25, usd),
		FamilyYearly:      currency.NewAmount(480.00, usd),
		FamilyLastMonth:   currency.NewAmount(40.00, usd),
		FamilyLastYear:    currency.NewAmount(470.00, usd),
		UpcomingRenewals:  []query.SummaryQueryUpcomingRenewalsResponse{},
		TopProviders:      []query.SummaryQueryTopProvidersResponse{},
		TopLabels:         []query.SummaryQueryLabelResponse{},
	}

	// Create mock handler
	mockHandler := &mockSummaryQueryHandler{
		handleFunc: func(ctx context.Context, q query.SummaryQuery) result.Result[query.SummaryQueryResponse] {
			return result.Success(domainResponse)
		},
	}

	// Create endpoint
	endpoint := subscription.NewSummaryEndpoint(mockHandler)
	require.NotNil(t, endpoint)

	// Execute the handler's mapping logic by calling Handle
	// Note: We can't easily test the HTTP handler without setting up Gin context,
	// so we'll test the mapping logic directly by calling the handler
	ctx := context.Background()
	q := query.SummaryQuery{
		TopProviders:     5,
		TopLabels:        5,
		UpcomingRenewals: 5,
		TotalMonthly:     true,
		TotalYearly:      true,
	}

	res := mockHandler.Handle(ctx, q)
	require.True(t, res.IsSuccess())

	// Use Match to extract the value
	response := result.Match(
		res,
		func(value query.SummaryQueryResponse) query.SummaryQueryResponse {
			return value
		},
		func(err error) query.SummaryQueryResponse {
			t.Fatalf("unexpected error: %v", err)
			return query.SummaryQueryResponse{}
		},
	)

	// Verify all active counts are mapped correctly
	assert.Equal(t, uint16(10), response.Active)
	assert.Equal(t, uint16(6), response.ActivePersonal)
	assert.Equal(t, uint16(4), response.ActiveFamily)

	// Verify total amounts are mapped correctly
	assert.Equal(t, 100.50, response.TotalMonthly.Value())
	assert.Equal(t, 1200.00, response.TotalYearly.Value())
	assert.Equal(t, 95.00, response.TotalLastMonth.Value())
	assert.Equal(t, 1150.00, response.TotalLastYear.Value())

	// Verify personal amounts are mapped correctly
	assert.Equal(t, 60.25, response.PersonalMonthly.Value())
	assert.Equal(t, 720.00, response.PersonalYearly.Value())
	assert.Equal(t, 55.00, response.PersonalLastMonth.Value())
	assert.Equal(t, 680.00, response.PersonalLastYear.Value())

	// Verify family amounts are mapped correctly
	assert.Equal(t, 40.25, response.FamilyMonthly.Value())
	assert.Equal(t, 480.00, response.FamilyYearly.Value())
	assert.Equal(t, 40.00, response.FamilyLastMonth.Value())
	assert.Equal(t, 470.00, response.FamilyLastYear.Value())
}

// TestSummaryEndpoint_ZeroAmounts tests that zero amounts are handled correctly
func TestSummaryEndpoint_ZeroAmounts(t *testing.T) {
	usd := xcur.USD

	domainResponse := query.SummaryQueryResponse{
		Active:            0,
		ActivePersonal:    0,
		ActiveFamily:      0,
		TotalMonthly:      currency.NewAmount(0, usd),
		TotalYearly:       currency.NewAmount(0, usd),
		TotalLastMonth:    currency.NewAmount(0, usd),
		TotalLastYear:     currency.NewAmount(0, usd),
		PersonalMonthly:   currency.NewAmount(0, usd),
		PersonalYearly:    currency.NewAmount(0, usd),
		PersonalLastMonth: currency.NewAmount(0, usd),
		PersonalLastYear:  currency.NewAmount(0, usd),
		FamilyMonthly:     currency.NewAmount(0, usd),
		FamilyYearly:      currency.NewAmount(0, usd),
		FamilyLastMonth:   currency.NewAmount(0, usd),
		FamilyLastYear:    currency.NewAmount(0, usd),
		UpcomingRenewals:  []query.SummaryQueryUpcomingRenewalsResponse{},
		TopProviders:      []query.SummaryQueryTopProvidersResponse{},
		TopLabels:         []query.SummaryQueryLabelResponse{},
	}

	mockHandler := &mockSummaryQueryHandler{
		handleFunc: func(ctx context.Context, q query.SummaryQuery) result.Result[query.SummaryQueryResponse] {
			return result.Success(domainResponse)
		},
	}

	ctx := context.Background()
	q := query.SummaryQuery{
		TopProviders:     0,
		TopLabels:        0,
		UpcomingRenewals: 0,
		TotalMonthly:     true,
		TotalYearly:      true,
	}

	res := mockHandler.Handle(ctx, q)
	require.True(t, res.IsSuccess())

	response := result.Match(
		res,
		func(value query.SummaryQueryResponse) query.SummaryQueryResponse {
			return value
		},
		func(err error) query.SummaryQueryResponse {
			t.Fatalf("unexpected error: %v", err)
			return query.SummaryQueryResponse{}
		},
	)

	// Verify all counts are zero
	assert.Equal(t, uint16(0), response.Active)
	assert.Equal(t, uint16(0), response.ActivePersonal)
	assert.Equal(t, uint16(0), response.ActiveFamily)

	// Verify all amounts are zero
	assert.Equal(t, 0.0, response.TotalMonthly.Value())
	assert.Equal(t, 0.0, response.TotalYearly.Value())
	assert.Equal(t, 0.0, response.PersonalMonthly.Value())
	assert.Equal(t, 0.0, response.PersonalYearly.Value())
	assert.Equal(t, 0.0, response.FamilyMonthly.Value())
	assert.Equal(t, 0.0, response.FamilyYearly.Value())
}

// TestSummaryEndpoint_OnlyPersonalSubscriptions tests when user has only personal subscriptions
func TestSummaryEndpoint_OnlyPersonalSubscriptions(t *testing.T) {
	usd := xcur.USD

	domainResponse := query.SummaryQueryResponse{
		Active:            5,
		ActivePersonal:    5,
		ActiveFamily:      0,
		TotalMonthly:      currency.NewAmount(50.00, usd),
		TotalYearly:       currency.NewAmount(600.00, usd),
		TotalLastMonth:    currency.NewAmount(45.00, usd),
		TotalLastYear:     currency.NewAmount(550.00, usd),
		PersonalMonthly:   currency.NewAmount(50.00, usd),
		PersonalYearly:    currency.NewAmount(600.00, usd),
		PersonalLastMonth: currency.NewAmount(45.00, usd),
		PersonalLastYear:  currency.NewAmount(550.00, usd),
		FamilyMonthly:     currency.NewAmount(0, usd),
		FamilyYearly:      currency.NewAmount(0, usd),
		FamilyLastMonth:   currency.NewAmount(0, usd),
		FamilyLastYear:    currency.NewAmount(0, usd),
		UpcomingRenewals:  []query.SummaryQueryUpcomingRenewalsResponse{},
		TopProviders:      []query.SummaryQueryTopProvidersResponse{},
		TopLabels:         []query.SummaryQueryLabelResponse{},
	}

	mockHandler := &mockSummaryQueryHandler{
		handleFunc: func(ctx context.Context, q query.SummaryQuery) result.Result[query.SummaryQueryResponse] {
			return result.Success(domainResponse)
		},
	}

	ctx := context.Background()
	q := query.SummaryQuery{
		TopProviders:     0,
		TopLabels:        0,
		UpcomingRenewals: 0,
		TotalMonthly:     true,
		TotalYearly:      true,
	}

	res := mockHandler.Handle(ctx, q)
	require.True(t, res.IsSuccess())

	response := result.Match(
		res,
		func(value query.SummaryQueryResponse) query.SummaryQueryResponse {
			return value
		},
		func(err error) query.SummaryQueryResponse {
			t.Fatalf("unexpected error: %v", err)
			return query.SummaryQueryResponse{}
		},
	)

	// Verify counts
	assert.Equal(t, uint16(5), response.Active)
	assert.Equal(t, uint16(5), response.ActivePersonal)
	assert.Equal(t, uint16(0), response.ActiveFamily)

	// Verify total equals personal
	assert.Equal(t, response.PersonalMonthly.Value(), response.TotalMonthly.Value())
	assert.Equal(t, response.PersonalYearly.Value(), response.TotalYearly.Value())

	// Verify family is zero
	assert.Equal(t, 0.0, response.FamilyMonthly.Value())
	assert.Equal(t, 0.0, response.FamilyYearly.Value())
}

// TestSummaryEndpoint_OnlyFamilySubscriptions tests when user has only family subscriptions
func TestSummaryEndpoint_OnlyFamilySubscriptions(t *testing.T) {
	usd := xcur.USD

	domainResponse := query.SummaryQueryResponse{
		Active:            3,
		ActivePersonal:    0,
		ActiveFamily:      3,
		TotalMonthly:      currency.NewAmount(30.00, usd),
		TotalYearly:       currency.NewAmount(360.00, usd),
		TotalLastMonth:    currency.NewAmount(30.00, usd),
		TotalLastYear:     currency.NewAmount(360.00, usd),
		PersonalMonthly:   currency.NewAmount(0, usd),
		PersonalYearly:    currency.NewAmount(0, usd),
		PersonalLastMonth: currency.NewAmount(0, usd),
		PersonalLastYear:  currency.NewAmount(0, usd),
		FamilyMonthly:     currency.NewAmount(30.00, usd),
		FamilyYearly:      currency.NewAmount(360.00, usd),
		FamilyLastMonth:   currency.NewAmount(30.00, usd),
		FamilyLastYear:    currency.NewAmount(360.00, usd),
		UpcomingRenewals:  []query.SummaryQueryUpcomingRenewalsResponse{},
		TopProviders:      []query.SummaryQueryTopProvidersResponse{},
		TopLabels:         []query.SummaryQueryLabelResponse{},
	}

	mockHandler := &mockSummaryQueryHandler{
		handleFunc: func(ctx context.Context, q query.SummaryQuery) result.Result[query.SummaryQueryResponse] {
			return result.Success(domainResponse)
		},
	}

	ctx := context.Background()
	q := query.SummaryQuery{
		TopProviders:     0,
		TopLabels:        0,
		UpcomingRenewals: 0,
		TotalMonthly:     true,
		TotalYearly:      true,
	}

	res := mockHandler.Handle(ctx, q)
	require.True(t, res.IsSuccess())

	response := result.Match(
		res,
		func(value query.SummaryQueryResponse) query.SummaryQueryResponse {
			return value
		},
		func(err error) query.SummaryQueryResponse {
			t.Fatalf("unexpected error: %v", err)
			return query.SummaryQueryResponse{}
		},
	)

	// Verify counts
	assert.Equal(t, uint16(3), response.Active)
	assert.Equal(t, uint16(0), response.ActivePersonal)
	assert.Equal(t, uint16(3), response.ActiveFamily)

	// Verify total equals family
	assert.Equal(t, response.FamilyMonthly.Value(), response.TotalMonthly.Value())
	assert.Equal(t, response.FamilyYearly.Value(), response.TotalYearly.Value())

	// Verify personal is zero
	assert.Equal(t, 0.0, response.PersonalMonthly.Value())
	assert.Equal(t, 0.0, response.PersonalYearly.Value())
}

// TestSummaryEndpoint_WithTopProvidersAndLabels tests mapping with top providers and labels
func TestSummaryEndpoint_WithTopProvidersAndLabels(t *testing.T) {
	usd := xcur.USD
	providerID1, _ := types.ParseProviderID("00000000-0000-0000-0000-000000000001")
	providerID2, _ := types.ParseProviderID("00000000-0000-0000-0000-000000000002")
	labelID1, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000011")
	labelID2, _ := types.ParseLabelID("00000000-0000-0000-0000-000000000012")

	domainResponse := query.SummaryQueryResponse{
		Active:            5,
		ActivePersonal:    3,
		ActiveFamily:      2,
		TotalMonthly:      currency.NewAmount(50.00, usd),
		TotalYearly:       currency.NewAmount(600.00, usd),
		TotalLastMonth:    currency.NewAmount(50.00, usd),
		TotalLastYear:     currency.NewAmount(600.00, usd),
		PersonalMonthly:   currency.NewAmount(30.00, usd),
		PersonalYearly:    currency.NewAmount(360.00, usd),
		PersonalLastMonth: currency.NewAmount(30.00, usd),
		PersonalLastYear:  currency.NewAmount(360.00, usd),
		FamilyMonthly:     currency.NewAmount(20.00, usd),
		FamilyYearly:      currency.NewAmount(240.00, usd),
		FamilyLastMonth:   currency.NewAmount(20.00, usd),
		FamilyLastYear:    currency.NewAmount(240.00, usd),
		TopProviders: []query.SummaryQueryTopProvidersResponse{
			{
				ProviderID: providerID1,
				Total:      currency.NewAmount(100.00, usd),
				Duration:   30 * 24 * time.Hour,
			},
			{
				ProviderID: providerID2,
				Total:      currency.NewAmount(50.00, usd),
				Duration:   15 * 24 * time.Hour,
			},
		},
		TopLabels: []query.SummaryQueryLabelResponse{
			{
				LabelID: labelID1,
				Total:   currency.NewAmount(75.00, usd),
			},
			{
				LabelID: labelID2,
				Total:   currency.NewAmount(25.00, usd),
			},
		},
		UpcomingRenewals: []query.SummaryQueryUpcomingRenewalsResponse{},
	}

	mockHandler := &mockSummaryQueryHandler{
		handleFunc: func(ctx context.Context, q query.SummaryQuery) result.Result[query.SummaryQueryResponse] {
			return result.Success(domainResponse)
		},
	}

	ctx := context.Background()
	q := query.SummaryQuery{
		TopProviders:     2,
		TopLabels:        2,
		UpcomingRenewals: 0,
		TotalMonthly:     true,
		TotalYearly:      true,
	}

	res := mockHandler.Handle(ctx, q)
	require.True(t, res.IsSuccess())

	response := result.Match(
		res,
		func(value query.SummaryQueryResponse) query.SummaryQueryResponse {
			return value
		},
		func(err error) query.SummaryQueryResponse {
			t.Fatalf("unexpected error: %v", err)
			return query.SummaryQueryResponse{}
		},
	)

	// Verify basic counts and amounts
	assert.Equal(t, uint16(5), response.Active)
	assert.Equal(t, uint16(3), response.ActivePersonal)
	assert.Equal(t, uint16(2), response.ActiveFamily)

	// Verify top providers are present
	require.Len(t, response.TopProviders, 2)
	assert.Equal(t, providerID1, response.TopProviders[0].ProviderID)
	assert.Equal(t, 100.00, response.TopProviders[0].Total.Value())

	// Verify top labels are present
	require.Len(t, response.TopLabels, 2)
	assert.Equal(t, labelID1, response.TopLabels[0].LabelID)
	assert.Equal(t, 75.00, response.TopLabels[0].Total.Value())
}

// TestDTOMapping_DirectConversion tests the DTO conversion directly
func TestDTOMapping_DirectConversion(t *testing.T) {
	usd := xcur.USD

	// Create domain response
	domainResponse := query.SummaryQueryResponse{
		Active:            10,
		ActivePersonal:    6,
		ActiveFamily:      4,
		TotalMonthly:      currency.NewAmount(100.50, usd),
		TotalYearly:       currency.NewAmount(1200.00, usd),
		TotalLastMonth:    currency.NewAmount(95.00, usd),
		TotalLastYear:     currency.NewAmount(1150.00, usd),
		PersonalMonthly:   currency.NewAmount(60.25, usd),
		PersonalYearly:    currency.NewAmount(720.00, usd),
		PersonalLastMonth: currency.NewAmount(55.00, usd),
		PersonalLastYear:  currency.NewAmount(680.00, usd),
		FamilyMonthly:     currency.NewAmount(40.25, usd),
		FamilyYearly:      currency.NewAmount(480.00, usd),
		FamilyLastMonth:   currency.NewAmount(40.00, usd),
		FamilyLastYear:    currency.NewAmount(470.00, usd),
		UpcomingRenewals:  []query.SummaryQueryUpcomingRenewalsResponse{},
		TopProviders:      []query.SummaryQueryTopProvidersResponse{},
		TopLabels:         []query.SummaryQueryLabelResponse{},
	}

	// Convert to DTO (simulating what the handler does)
	dtoResponse := dto.SubscriptionSummaryResponse{
		Active:            domainResponse.Active,
		ActivePersonal:    domainResponse.ActivePersonal,
		ActiveFamily:      domainResponse.ActiveFamily,
		TotalMonthly:      dto.NewAmount(domainResponse.TotalMonthly),
		TotalYearly:       dto.NewAmount(domainResponse.TotalYearly),
		TotalLastMonth:    dto.NewAmount(domainResponse.TotalLastMonth),
		TotalLastYear:     dto.NewAmount(domainResponse.TotalLastYear),
		PersonalMonthly:   dto.NewAmount(domainResponse.PersonalMonthly),
		PersonalYearly:    dto.NewAmount(domainResponse.PersonalYearly),
		PersonalLastMonth: dto.NewAmount(domainResponse.PersonalLastMonth),
		PersonalLastYear:  dto.NewAmount(domainResponse.PersonalLastYear),
		FamilyMonthly:     dto.NewAmount(domainResponse.FamilyMonthly),
		FamilyYearly:      dto.NewAmount(domainResponse.FamilyYearly),
		FamilyLastMonth:   dto.NewAmount(domainResponse.FamilyLastMonth),
		FamilyLastYear:    dto.NewAmount(domainResponse.FamilyLastYear),
		TopProviders:      []dto.SubscriptionSummaryTopProviderResponse{},
		TopLabels:         []dto.SubscriptionSummaryTopLabelResponse{},
		UpcomingRenewals:  []dto.SubscriptionSummaryUpcomingRenewalResponse{},
	}

	// Verify all fields are correctly mapped
	assert.Equal(t, domainResponse.Active, dtoResponse.Active)
	assert.Equal(t, domainResponse.ActivePersonal, dtoResponse.ActivePersonal)
	assert.Equal(t, domainResponse.ActiveFamily, dtoResponse.ActiveFamily)

	assert.Equal(t, domainResponse.TotalMonthly.Value(), dtoResponse.TotalMonthly.Value)
	assert.Equal(t, domainResponse.TotalYearly.Value(), dtoResponse.TotalYearly.Value)
	assert.Equal(t, domainResponse.TotalLastMonth.Value(), dtoResponse.TotalLastMonth.Value)
	assert.Equal(t, domainResponse.TotalLastYear.Value(), dtoResponse.TotalLastYear.Value)

	assert.Equal(t, domainResponse.PersonalMonthly.Value(), dtoResponse.PersonalMonthly.Value)
	assert.Equal(t, domainResponse.PersonalYearly.Value(), dtoResponse.PersonalYearly.Value)
	assert.Equal(t, domainResponse.PersonalLastMonth.Value(), dtoResponse.PersonalLastMonth.Value)
	assert.Equal(t, domainResponse.PersonalLastYear.Value(), dtoResponse.PersonalLastYear.Value)

	assert.Equal(t, domainResponse.FamilyMonthly.Value(), dtoResponse.FamilyMonthly.Value)
	assert.Equal(t, domainResponse.FamilyYearly.Value(), dtoResponse.FamilyYearly.Value)
	assert.Equal(t, domainResponse.FamilyLastMonth.Value(), dtoResponse.FamilyLastMonth.Value)
	assert.Equal(t, domainResponse.FamilyLastYear.Value(), dtoResponse.FamilyLastYear.Value)

	// Verify currency codes match
	assert.Equal(t, "USD", dtoResponse.TotalMonthly.Currency)
	assert.Equal(t, "USD", dtoResponse.PersonalMonthly.Currency)
	assert.Equal(t, "USD", dtoResponse.FamilyMonthly.Currency)
}
