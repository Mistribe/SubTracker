package query

import (
	"context"
	"sort"
	"time"

	xcur "golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

type SummaryQuery struct {
	TopProviders     uint8
	TopLabels        uint8
	UpcomingRenewals uint8
	TotalMonthly     bool
	TotalYearly      bool
}

type SummaryQueryUpcomingRenewalsResponse struct {
	ProviderId     types.ProviderID
	SubscriptionId types.SubscriptionID
	At             time.Time
	Total          currency.Amount
	Source         *currency.Amount
}

type SummaryQueryTopProvidersResponse struct {
	ProviderID types.ProviderID
	Total      currency.Amount
	Duration   time.Duration
}

type SummaryQueryLabelResponse struct {
	LabelID types.LabelID
	Total   currency.Amount
}

type SummaryQueryResponse struct {
	Active            uint16
	ActivePersonal    uint16
	ActiveFamily      uint16
	TotalMonthly      currency.Amount
	TotalLastMonth    currency.Amount
	TotalYearly       currency.Amount
	TotalLastYear     currency.Amount
	PersonalMonthly   currency.Amount
	PersonalLastMonth currency.Amount
	PersonalYearly    currency.Amount
	PersonalLastYear  currency.Amount
	FamilyMonthly     currency.Amount
	FamilyLastMonth   currency.Amount
	FamilyYearly      currency.Amount
	FamilyLastYear    currency.Amount
	UpcomingRenewals  []SummaryQueryUpcomingRenewalsResponse
	TopProviders      []SummaryQueryTopProvidersResponse
	TopLabels         []SummaryQueryLabelResponse
}

type SummaryQueryHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	currencyRepository     ports.CurrencyRepository
	familyRepository       ports.FamilyRepository
	accountService         ports.AccountService
	authService            ports.Authentication
	exchange               ports.Exchange
}

func NewSummaryQueryHandler(
	subscriptionRepository ports.SubscriptionRepository,
	currencyRepository ports.CurrencyRepository,
	familyRepository ports.FamilyRepository,
	accountService ports.AccountService,
	authService ports.Authentication,
	exchange ports.Exchange) *SummaryQueryHandler {
	return &SummaryQueryHandler{
		subscriptionRepository: subscriptionRepository,
		currencyRepository:     currencyRepository,
		familyRepository:       familyRepository,
		accountService:         accountService,
		authService:            authService,
		exchange:               exchange,
	}
}

func (h SummaryQueryHandler) convertToCurrency(
	ctx context.Context,
	amount currency.Amount,
	preferredCurrency xcur.Unit,
	date time.Time) currency.Amount {
	v, err := h.exchange.ToCurrencyAt(ctx, amount, preferredCurrency, time.Now())
	if err != nil {
		return currency.NewInvalidAmount()
	}
	return v
}

// classifySubscription determines if a subscription is personal, family, or should be excluded
// Returns: "personal", "family", or "exclude"
func (h SummaryQueryHandler) classifySubscription(sub subscription.Subscription, currentUserID types.UserID,
	fam family.Family) string {
	payer := sub.Payer()

	// If no payer, default to personal
	if payer == nil {
		return "personal"
	}

	switch payer.Type() {
	case subscription.FamilyPayer:
		return "family"
	case subscription.FamilyMemberPayer:
		// Find the family member with this member ID and check if their user ID matches
		member := fam.GetMember(payer.MemberId())
		if member != nil && member.UserId() != nil && *member.UserId() == currentUserID {
			return "personal"
		}
		// Another family member is paying, exclude from both
		return "exclude"
	default:
		return "personal"
	}
}

func (h SummaryQueryHandler) Handle(ctx context.Context, query SummaryQuery) result.Result[SummaryQueryResponse] {
	connectedAccount := h.authService.MustGetConnectedAccount(ctx)
	currentUserID := connectedAccount.UserID()
	preferredCurrency := h.accountService.GetPreferredCurrency(ctx, currentUserID)
	currencyRates, err := h.currencyRepository.GetRatesByDate(ctx, time.Now())
	if err != nil {
		return result.Fail[SummaryQueryResponse](err)
	}

	currencyRates = currencyRates.WithReverse()

	// Get family for payer classification
	family, err := h.familyRepository.GetAccountFamily(ctx, currentUserID)
	if err != nil {
		return result.Fail[SummaryQueryResponse](err)
	}

	topProviders := make(map[types.ProviderID]SummaryQueryTopProvidersResponse)
	topLabels := make(map[types.LabelID]currency.Amount)
	var upcomingRenewals []SummaryQueryUpcomingRenewalsResponse

	// Separate accumulators for personal, family, and total
	totalMonthly := 0.0
	totalYearly := 0.0
	totalLastMonth := 0.0
	totalLastYear := 0.0
	personalMonthly := 0.0
	personalYearly := 0.0
	personalLastMonth := 0.0
	personalLastYear := 0.0
	familyMonthly := 0.0
	familyYearly := 0.0
	familyLastMonth := 0.0
	familyLastYear := 0.0

	lastMonth := time.Now().AddDate(0, -1, 0)
	lastYear := time.Now().AddDate(-1, 0, 0)

	var active uint16
	var activePersonal uint16
	var activeFamily uint16

	for sub := range h.subscriptionRepository.GetAllIt(ctx, currentUserID, "") {
		// Classify subscription
		classification := h.classifySubscription(sub, currentUserID, family)

		// Skip subscriptions paid by other family members
		if classification == "exclude" {
			continue
		}

		isActive := sub.IsActive()
		isActiveLastMonth := sub.IsActiveAt(lastMonth)
		isActiveLastYear := sub.IsActiveAt(lastYear)

		// Count active subscriptions
		if isActive {
			active++
			if classification == "personal" {
				activePersonal++
			} else if classification == "family" {
				activeFamily++
			}
		}

		// Calculate monthly amounts
		if query.TotalMonthly {
			if isActive {
				monthlyAmount := h.convertToCurrency(ctx,
					sub.GetRecurrencyAmount(subscription.MonthlyRecurrency),
					preferredCurrency,
					sub.StartDate())
				if monthlyAmount.IsValid() {
					totalMonthly += monthlyAmount.Value()
					if classification == "personal" {
						personalMonthly += monthlyAmount.Value()
					} else if classification == "family" {
						familyMonthly += monthlyAmount.Value()
					}
				}
			}
			if isActiveLastMonth {
				lastMonthAmount := h.convertToCurrency(ctx,
					sub.GetRecurrencyAmount(subscription.MonthlyRecurrency),
					preferredCurrency,
					sub.StartDate())
				if lastMonthAmount.IsValid() {
					totalLastMonth += lastMonthAmount.Value()
					if classification == "personal" {
						personalLastMonth += lastMonthAmount.Value()
					} else if classification == "family" {
						familyLastMonth += lastMonthAmount.Value()
					}
				}
			}
		}

		// Calculate yearly amounts
		if query.TotalYearly {
			if isActive {
				yearlyAmount := h.convertToCurrency(ctx,
					sub.GetRecurrencyAmount(subscription.YearlyRecurrency),
					preferredCurrency,
					sub.StartDate())
				if yearlyAmount.IsValid() {
					totalYearly += yearlyAmount.Value()
					if classification == "personal" {
						personalYearly += yearlyAmount.Value()
					} else if classification == "family" {
						familyYearly += yearlyAmount.Value()
					}
				}
			}
			if isActiveLastYear {
				lastYearAmount := h.convertToCurrency(ctx,
					sub.GetRecurrencyAmount(subscription.YearlyRecurrency),
					preferredCurrency,
					sub.StartDate())
				if lastYearAmount.IsValid() {
					totalLastYear += lastYearAmount.Value()
					if classification == "personal" {
						personalLastYear += lastYearAmount.Value()
					} else if classification == "family" {
						familyLastYear += lastYearAmount.Value()
					}
				}
			}
		}

		if query.UpcomingRenewals > 0 {
			if isActive {
				renewalDate := sub.GetNextRenewalDate()
				price := h.convertToCurrency(ctx,
					sub.GetPrice(),
					preferredCurrency,
					sub.StartDate())
				if renewalDate != nil && price.IsValid() {
					upcomingRenewals = append(upcomingRenewals, SummaryQueryUpcomingRenewalsResponse{
						ProviderId:     sub.ProviderId(),
						SubscriptionId: sub.Id(),
						At:             *renewalDate,
						Total:          price,
						Source:         &price,
					})
				}
			}
		}

		if query.TopProviders > 0 || query.TopLabels > 0 {
			if sub.IsStarted() {
				totalSpent := h.convertToCurrency(ctx,
					sub.GetTotalSpent(),
					preferredCurrency,
					sub.StartDate())
				if totalSpent.IsValid() && !totalSpent.IsZero() {
					if query.TopProviders > 0 {
						existingTopProvider, ok := topProviders[sub.ProviderId()]
						if ok {
							existingTopProvider.Total = existingTopProvider.Total.Add(totalSpent.Value())
							existingTopProvider.Duration += sub.GetTotalDuration()
							topProviders[sub.ProviderId()] = existingTopProvider
						} else {
							topProviders[sub.ProviderId()] = SummaryQueryTopProvidersResponse{
								ProviderID: sub.ProviderId(),
								Total:      totalSpent,
								Duration:   sub.GetTotalDuration(),
							}
						}
					}
					if query.TopLabels > 0 {
						for l := range sub.Labels().It() {
							existingTopLabel, ok := topLabels[l.LabelId]
							if ok {
								existingTopLabel = existingTopLabel.Add(totalSpent.Value())
								topLabels[l.LabelId] = existingTopLabel
							} else {
								topLabels[l.LabelId] = totalSpent
							}
						}
					}
				}
			}
		}
	}

	response := SummaryQueryResponse{
		Active:            active,
		ActivePersonal:    activePersonal,
		ActiveFamily:      activeFamily,
		TotalMonthly:      currency.NewAmount(totalMonthly, preferredCurrency),
		TotalYearly:       currency.NewAmount(totalYearly, preferredCurrency),
		TotalLastMonth:    currency.NewAmount(totalLastMonth, preferredCurrency),
		TotalLastYear:     currency.NewAmount(totalLastYear, preferredCurrency),
		PersonalMonthly:   currency.NewAmount(personalMonthly, preferredCurrency),
		PersonalYearly:    currency.NewAmount(personalYearly, preferredCurrency),
		PersonalLastMonth: currency.NewAmount(personalLastMonth, preferredCurrency),
		PersonalLastYear:  currency.NewAmount(personalLastYear, preferredCurrency),
		FamilyMonthly:     currency.NewAmount(familyMonthly, preferredCurrency),
		FamilyYearly:      currency.NewAmount(familyYearly, preferredCurrency),
		FamilyLastMonth:   currency.NewAmount(familyLastMonth, preferredCurrency),
		FamilyLastYear:    currency.NewAmount(familyLastYear, preferredCurrency),
		UpcomingRenewals:  upcomingRenewals,
		TopProviders: herd.MapToArr(topProviders,
			func(providerId types.ProviderID, res SummaryQueryTopProvidersResponse) SummaryQueryTopProvidersResponse {
				return res
			}),
		TopLabels: herd.MapToArr(topLabels,
			func(labelId types.LabelID, res currency.Amount) SummaryQueryLabelResponse {
				return SummaryQueryLabelResponse{
					LabelID: labelId,
					Total:   res,
				}
			}),
	}
	sort.Slice(response.UpcomingRenewals, func(i, j int) bool {
		return response.UpcomingRenewals[i].At.Before(response.UpcomingRenewals[j].At)
	})
	response.UpcomingRenewals = herd.Take(response.UpcomingRenewals, int(query.UpcomingRenewals))
	sort.Slice(response.TopProviders, func(i, j int) bool {
		return response.TopProviders[i].Total.IsGreaterThan(response.TopProviders[j].Total)
	})
	response.TopProviders = herd.Take(response.TopProviders, int(query.TopProviders))
	sort.Slice(response.TopLabels, func(i, j int) bool {
		return response.TopLabels[i].Total.IsGreaterThan(response.TopLabels[j].Total)
	})
	response.TopLabels = herd.Take(response.TopLabels, int(query.TopLabels))
	return result.Success(response)
}
