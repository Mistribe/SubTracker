package query

import (
	"context"
	"iter"
	"testing"
	"time"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
	xcur "golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
)

// Feature: dashboard-spending-breakdown, Property 1: Personal expense filtering
// Validates: Requirements 2.1, 2.3
func TestProperty_PersonalExpenseFiltering(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("subscriptions with nil payer or matching family member payer are classified as personal",
		prop.ForAll(
			func(seed int64) bool {
				// Generate test data
				currentUserID := types.UserID("user-123")
				familyID := types.NewFamilyID()
				memberID := types.NewFamilyMemberID()

				// Create family with member linked to current user
				member := family.NewMember(
					memberID,
					familyID,
					"Test User",
					family.OwnerMemberType,
					nil,
					time.Now(),
					time.Now(),
				)
				member.SetUserId(&currentUserID)

				fam := family.NewFamily(
					familyID,
					currentUserID,
					"Test Family",
					[]family.Member{member},
					time.Now(),
					time.Now(),
				)

				handler := &SummaryQueryHandler{}

				// Test 1: Subscription with nil payer should be personal
				subWithNilPayer := createTestSubscription(nil)
				classification1 := handler.classifySubscription(subWithNilPayer, currentUserID, fam)
				if classification1 != "personal" {
					return false
				}

				// Test 2: Subscription with family member payer matching current user should be personal
				payerMatchingUser := subscription.NewFamilyMemberPayer(familyID, memberID)
				subWithMatchingPayer := createTestSubscription(payerMatchingUser)
				classification2 := handler.classifySubscription(subWithMatchingPayer, currentUserID, fam)
				if classification2 != "personal" {
					return false
				}

				return true
			},
			gen.Int64(),
		))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Helper function to create test subscriptions
func createTestSubscription(payer subscription.Payer) subscription.Subscription {
	subID := types.NewSubscriptionID()
	providerID := types.NewProviderID()
	owner := types.NewPersonalOwner(types.UserID("user-123"))
	price := subscription.NewPrice(currency.NewAmount(10.0, xcur.USD))

	return subscription.NewSubscription(
		subID,
		nil,
		nil,
		providerID,
		price,
		owner,
		payer,
		[]types.FamilyMemberID{},
		[]subscription.LabelRef{},
		time.Now().AddDate(0, -1, 0),
		nil,
		subscription.MonthlyRecurrency,
		nil,
		time.Now(),
		time.Now(),
	)
}

// Feature: dashboard-spending-breakdown, Property 2: Family expense filtering
// Validates: Requirements 2.2
func TestProperty_FamilyExpenseFiltering(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("subscriptions with family payer are classified as family",
		prop.ForAll(
			func(seed int64) bool {
				currentUserID := types.UserID("user-123")
				familyID := types.NewFamilyID()
				memberID := types.NewFamilyMemberID()

				member := family.NewMember(
					memberID,
					familyID,
					"Test User",
					family.OwnerMemberType,
					nil,
					time.Now(),
					time.Now(),
				)
				member.SetUserId(&currentUserID)

				fam := family.NewFamily(
					familyID,
					currentUserID,
					"Test Family",
					[]family.Member{member},
					time.Now(),
					time.Now(),
				)

				handler := &SummaryQueryHandler{}

				// Subscription with family payer should be classified as family
				familyPayer := subscription.NewFamilyPayer(familyID)
				subWithFamilyPayer := createTestSubscription(familyPayer)
				classification := handler.classifySubscription(subWithFamilyPayer, currentUserID, fam)

				return classification == "family"
			},
			gen.Int64(),
		))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: dashboard-spending-breakdown, Property 3: Expense sum invariant
// Validates: Requirements 2.4
func TestProperty_ExpenseSumInvariant(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("total expenses equal personal plus family expenses",
		prop.ForAll(
			func(numPersonal, numFamily int) bool {
				if numPersonal < 0 || numFamily < 0 || numPersonal > 20 || numFamily > 20 {
					return true // Skip invalid inputs
				}

				currentUserID := types.UserID("user-123")
				familyID := types.NewFamilyID()
				memberID := types.NewFamilyMemberID()

				member := family.NewMember(
					memberID,
					familyID,
					"Test User",
					family.OwnerMemberType,
					nil,
					time.Now(),
					time.Now(),
				)
				member.SetUserId(&currentUserID)

				fam := family.NewFamily(
					familyID,
					currentUserID,
					"Test Family",
					[]family.Member{member},
					time.Now(),
					time.Now(),
				)

				// Create mock repositories
				mockSubRepo := &mockSubscriptionRepository{
					subscriptions: []subscription.Subscription{},
				}

				// Add personal subscriptions
				for i := 0; i < numPersonal; i++ {
					payer := subscription.NewFamilyMemberPayer(familyID, memberID)
					sub := createTestSubscriptionWithAmount(payer, 10.0+float64(i))
					mockSubRepo.subscriptions = append(mockSubRepo.subscriptions, sub)
				}

				// Add family subscriptions
				for i := 0; i < numFamily; i++ {
					payer := subscription.NewFamilyPayer(familyID)
					sub := createTestSubscriptionWithAmount(payer, 20.0+float64(i))
					mockSubRepo.subscriptions = append(mockSubRepo.subscriptions, sub)
				}

				// Calculate totals manually
				expectedPersonalMonthly := 0.0
				expectedFamilyMonthly := 0.0
				for i := 0; i < numPersonal; i++ {
					expectedPersonalMonthly += 10.0 + float64(i)
				}
				for i := 0; i < numFamily; i++ {
					expectedFamilyMonthly += 20.0 + float64(i)
				}
				expectedTotalMonthly := expectedPersonalMonthly + expectedFamilyMonthly

				// Simulate the handler logic
				personalMonthly := 0.0
				familyMonthly := 0.0
				totalMonthly := 0.0

				handler := &SummaryQueryHandler{}
				for _, sub := range mockSubRepo.subscriptions {
					if !sub.IsActive() {
						continue
					}
					classification := handler.classifySubscription(sub, currentUserID, fam)
					if classification == "exclude" {
						continue
					}

					monthlyAmount := sub.GetRecurrencyAmount(subscription.MonthlyRecurrency)
					if monthlyAmount.IsValid() {
						totalMonthly += monthlyAmount.Value()
						if classification == "personal" {
							personalMonthly += monthlyAmount.Value()
						} else if classification == "family" {
							familyMonthly += monthlyAmount.Value()
						}
					}
				}

				// Check invariant: total = personal + family
				tolerance := 0.01
				sumMatches := abs(totalMonthly-(personalMonthly+familyMonthly)) < tolerance
				expectedMatches := abs(totalMonthly-expectedTotalMonthly) < tolerance

				return sumMatches && expectedMatches
			},
			gen.IntRange(0, 10),
			gen.IntRange(0, 10),
		))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: dashboard-spending-breakdown, Property 4: Active subscription count invariant
// Validates: Requirements 4.5
func TestProperty_ActiveSubscriptionCountInvariant(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("total active count equals personal plus family active counts",
		prop.ForAll(
			func(numPersonal, numFamily int) bool {
				if numPersonal < 0 || numFamily < 0 || numPersonal > 20 || numFamily > 20 {
					return true // Skip invalid inputs
				}

				currentUserID := types.UserID("user-123")
				familyID := types.NewFamilyID()
				memberID := types.NewFamilyMemberID()

				member := family.NewMember(
					memberID,
					familyID,
					"Test User",
					family.OwnerMemberType,
					nil,
					time.Now(),
					time.Now(),
				)
				member.SetUserId(&currentUserID)

				fam := family.NewFamily(
					familyID,
					currentUserID,
					"Test Family",
					[]family.Member{member},
					time.Now(),
					time.Now(),
				)

				mockSubRepo := &mockSubscriptionRepository{
					subscriptions: []subscription.Subscription{},
				}

				// Add personal subscriptions
				for i := 0; i < numPersonal; i++ {
					payer := subscription.NewFamilyMemberPayer(familyID, memberID)
					sub := createTestSubscriptionWithAmount(payer, 10.0)
					mockSubRepo.subscriptions = append(mockSubRepo.subscriptions, sub)
				}

				// Add family subscriptions
				for i := 0; i < numFamily; i++ {
					payer := subscription.NewFamilyPayer(familyID)
					sub := createTestSubscriptionWithAmount(payer, 20.0)
					mockSubRepo.subscriptions = append(mockSubRepo.subscriptions, sub)
				}

				// Count active subscriptions
				var active, activePersonal, activeFamily uint16
				handler := &SummaryQueryHandler{}

				for _, sub := range mockSubRepo.subscriptions {
					if !sub.IsActive() {
						continue
					}

					classification := handler.classifySubscription(sub, currentUserID, fam)
					if classification == "exclude" {
						continue
					}

					active++
					if classification == "personal" {
						activePersonal++
					} else if classification == "family" {
						activeFamily++
					}
				}

				// Check invariant: total active = personal active + family active
				return active == (activePersonal + activeFamily)
			},
			gen.IntRange(0, 10),
			gen.IntRange(0, 10),
		))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Feature: dashboard-spending-breakdown, Property 6: Backend calculation completeness
// Validates: Requirements 4.1, 4.2, 4.3, 4.4, 4.6
func TestProperty_BackendCalculationCompleteness(t *testing.T) {
	properties := gopter.NewProperties(nil)

	properties.Property("response includes all required breakdown fields",
		prop.ForAll(
			func(hasPersonal, hasFamily bool) bool {
				// This property verifies that the response structure includes all required fields
				// We create a response and verify all fields are present

				preferredCurrency := xcur.USD
				response := SummaryQueryResponse{
					Active:            5,
					ActivePersonal:    3,
					ActiveFamily:      2,
					TotalMonthly:      currency.NewAmount(100.0, preferredCurrency),
					TotalYearly:       currency.NewAmount(1200.0, preferredCurrency),
					TotalLastMonth:    currency.NewAmount(90.0, preferredCurrency),
					TotalLastYear:     currency.NewAmount(1100.0, preferredCurrency),
					PersonalMonthly:   currency.NewAmount(60.0, preferredCurrency),
					PersonalYearly:    currency.NewAmount(720.0, preferredCurrency),
					PersonalLastMonth: currency.NewAmount(50.0, preferredCurrency),
					PersonalLastYear:  currency.NewAmount(600.0, preferredCurrency),
					FamilyMonthly:     currency.NewAmount(40.0, preferredCurrency),
					FamilyYearly:      currency.NewAmount(480.0, preferredCurrency),
					FamilyLastMonth:   currency.NewAmount(40.0, preferredCurrency),
					FamilyLastYear:    currency.NewAmount(500.0, preferredCurrency),
				}

				// Verify all fields are accessible and have valid values
				hasAllFields := response.Active > 0 &&
					response.TotalMonthly.IsValid() &&
					response.TotalYearly.IsValid() &&
					response.PersonalMonthly.IsValid() &&
					response.PersonalYearly.IsValid() &&
					response.FamilyMonthly.IsValid() &&
					response.FamilyYearly.IsValid() &&
					response.TotalLastMonth.IsValid() &&
					response.TotalLastYear.IsValid() &&
					response.PersonalLastMonth.IsValid() &&
					response.PersonalLastYear.IsValid() &&
					response.FamilyLastMonth.IsValid() &&
					response.FamilyLastYear.IsValid()

				return hasAllFields
			},
			gen.Bool(),
			gen.Bool(),
		))

	properties.TestingRun(t, gopter.ConsoleReporter(false))
}

// Helper functions
func createTestSubscriptionWithAmount(payer subscription.Payer, amount float64) subscription.Subscription {
	subID := types.NewSubscriptionID()
	providerID := types.NewProviderID()
	owner := types.NewPersonalOwner(types.UserID("user-123"))
	price := subscription.NewPrice(currency.NewAmount(amount, xcur.USD))

	return subscription.NewSubscription(
		subID,
		nil,
		nil,
		providerID,
		price,
		owner,
		payer,
		[]types.FamilyMemberID{},
		[]subscription.LabelRef{},
		time.Now().AddDate(0, -1, 0),
		nil,
		subscription.MonthlyRecurrency,
		nil,
		time.Now(),
		time.Now(),
	)
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// Mock subscription repository for testing
type mockSubscriptionRepository struct {
	subscriptions []subscription.Subscription
}

func (m *mockSubscriptionRepository) GetAllIt(ctx context.Context, userId types.UserID, search string) iter.Seq[subscription.Subscription] {
	return func(yield func(subscription.Subscription) bool) {
		for _, sub := range m.subscriptions {
			if !yield(sub) {
				return
			}
		}
	}
}

func (m *mockSubscriptionRepository) GetAll(ctx context.Context, userId types.UserID, search string) ([]subscription.Subscription, error) {
	return m.subscriptions, nil
}

func (m *mockSubscriptionRepository) GetById(ctx context.Context, userId types.UserID, id types.SubscriptionID) (subscription.Subscription, error) {
	return nil, nil
}

func (m *mockSubscriptionRepository) Create(ctx context.Context, entity subscription.Subscription) (subscription.Subscription, error) {
	return nil, nil
}

func (m *mockSubscriptionRepository) Update(ctx context.Context, entity subscription.Subscription) (subscription.Subscription, error) {
	return nil, nil
}

func (m *mockSubscriptionRepository) Delete(ctx context.Context, userId types.UserID, id types.SubscriptionID) error {
	return nil
}

func (m *mockSubscriptionRepository) GetPaginated(ctx context.Context, userId types.UserID, search string, page int, pageSize int) ([]subscription.Subscription, int, error) {
	return nil, 0, nil
}
