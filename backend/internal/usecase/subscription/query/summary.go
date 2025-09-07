package query

import (
	"context"
	"sort"
	"time"

	"github.com/google/uuid"
	xcur "golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x/collection"
)

type SummaryQuery struct {
	TopProviders     uint8
	TopLabels        uint8
	UpcomingRenewals uint8
	TotalMonthly     bool
	TotalYearly      bool
}

type SummaryQueryUpcomingRenewalsResponse struct {
	ProviderId uuid.UUID
	At         time.Time
	Total      currency.Amount
	Source     *currency.Amount
}

type SummaryQueryTopProvidersResponse struct {
	ProviderId uuid.UUID
	Total      currency.Amount
	Duration   time.Duration
}

type SummaryQueryLabelResponse struct {
	TagId uuid.UUID
	Total currency.Amount
}

type SummaryQueryResponse struct {
	Active           uint16
	TotalMonthly     currency.Amount
	TotalLastMonth   currency.Amount
	TotalYearly      currency.Amount
	TotalLastYear    currency.Amount
	UpcomingRenewals []SummaryQueryUpcomingRenewalsResponse
	TopProviders     []SummaryQueryTopProvidersResponse
	TopLabels        []SummaryQueryLabelResponse
}

type SummaryQueryHandler struct {
	subscriptionRepository ports.SubscriptionRepository
	currencyRepository     ports.CurrencyRepository
	userRepository         ports.UserRepository
	authService            ports.AuthService
	exchange               ports.Exchange
}

func NewSummaryQueryHandler(
	subscriptionRepository ports.SubscriptionRepository,
	currencyRepository ports.CurrencyRepository,
	userRepository ports.UserRepository,
	authService ports.AuthService,
	exchange ports.Exchange) *SummaryQueryHandler {
	return &SummaryQueryHandler{
		subscriptionRepository: subscriptionRepository,
		currencyRepository:     currencyRepository,
		userRepository:         userRepository,
		authService:            authService,
		exchange:               exchange,
	}
}

func (h SummaryQueryHandler) convertToCurrency(
	ctx context.Context,
	amount currency.Amount,
	preferredCurrency xcur.Unit,
	date time.Time) currency.Amount {
	v, err := h.exchange.ToCurrencyAt(ctx, amount, preferredCurrency, date)
	if err != nil {
		return currency.NewInvalidAmount()
	}
	return v
}

func (h SummaryQueryHandler) Handle(ctx context.Context, query SummaryQuery) result.Result[SummaryQueryResponse] {
	userId := h.authService.MustGetUserId(ctx)
	userProfile, err := h.userRepository.GetUser(ctx, userId)
	if err != nil {
		return result.Fail[SummaryQueryResponse](err)
	}
	currencyRates, err := h.currencyRepository.GetRatesByDate(ctx, time.Now())
	if err != nil {
		return result.Fail[SummaryQueryResponse](err)
	}
	preferredCurrency := currency.USD
	if userProfile != nil {
		preferredCurrency = userProfile.Currency()
	}

	currencyRates = currencyRates.WithReverse()

	topProviders := make(map[uuid.UUID]SummaryQueryTopProvidersResponse)
	topLabels := make(map[uuid.UUID]currency.Amount)
	var upcomingRenewals []SummaryQueryUpcomingRenewalsResponse
	totalMonthly := 0.0
	totalYearly := 0.0
	totalLastMonth := 0.0
	totalLastYear := 0.0
	lastMonth := time.Now().AddDate(0, -1, 0)
	lastYear := time.Now().AddDate(-1, 0, 0)
	var active uint16
	for sub := range h.subscriptionRepository.GetAllIt(ctx, userId, "") {
		if sub.IsActive() {
			active++
		}
		if query.TotalMonthly {
			if sub.IsActive() {
				monthlyAmount := h.convertToCurrency(ctx,
					sub.GetRecurrencyAmount(subscription.MonthlyRecurrency),
					preferredCurrency,
					sub.StartDate())
				if monthlyAmount.IsValid() {
					totalMonthly += monthlyAmount.Value()
				}
			}
			if sub.IsActiveAt(lastMonth) {
				lastMonthAmount := h.convertToCurrency(ctx,
					sub.GetRecurrencyAmount(subscription.MonthlyRecurrency),
					preferredCurrency,
					sub.StartDate())
				if lastMonthAmount.IsValid() {
					totalLastMonth += lastMonthAmount.Value()
				}
			}
		}
		if query.TotalYearly {
			if sub.IsActive() {
				yearlyAmount := h.convertToCurrency(ctx,
					sub.GetRecurrencyAmount(subscription.YearlyRecurrency),
					preferredCurrency,
					sub.StartDate())
				if yearlyAmount.IsValid() {
					totalYearly += yearlyAmount.Value()
				}
			}
			if sub.IsActiveAt(lastYear) {
				lastYearAmount := h.convertToCurrency(ctx,
					sub.GetRecurrencyAmount(subscription.YearlyRecurrency),
					preferredCurrency,
					sub.StartDate())
				if lastYearAmount.IsValid() {
					totalLastYear += lastYearAmount.Value()
				}
			}
		}

		if query.UpcomingRenewals > 0 {
			if sub.IsActive() {
				renewalDate := sub.GetNextRenewalDate()
				price := h.convertToCurrency(ctx,
					sub.GetPrice(),
					preferredCurrency,
					sub.StartDate())
				if renewalDate != nil && price.IsValid() {
					upcomingRenewals = append(upcomingRenewals, SummaryQueryUpcomingRenewalsResponse{
						ProviderId: sub.ProviderId(),
						At:         *renewalDate,
						Total:      price,
						Source:     &price,
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
								ProviderId: sub.ProviderId(),
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
		Active:           active,
		TotalMonthly:     currency.NewAmount(totalMonthly, preferredCurrency),
		TotalYearly:      currency.NewAmount(totalYearly, preferredCurrency),
		TotalLastMonth:   currency.NewAmount(totalLastMonth, preferredCurrency),
		TotalLastYear:    currency.NewAmount(totalLastYear, preferredCurrency),
		UpcomingRenewals: upcomingRenewals,
		TopProviders: collection.MapToArr(topProviders,
			func(providerId uuid.UUID, res SummaryQueryTopProvidersResponse) SummaryQueryTopProvidersResponse {
				return res
			}),
		TopLabels: collection.MapToArr(topLabels,
			func(labelId uuid.UUID, res currency.Amount) SummaryQueryLabelResponse {
				return SummaryQueryLabelResponse{
					TagId: labelId,
					Total: res,
				}
			}),
	}
	sort.Slice(response.UpcomingRenewals, func(i, j int) bool {
		return response.UpcomingRenewals[i].At.Before(response.UpcomingRenewals[j].At)
	})
	response.UpcomingRenewals = collection.Take(response.UpcomingRenewals, int(query.UpcomingRenewals))
	sort.Slice(response.TopProviders, func(i, j int) bool {
		return response.TopProviders[i].Total.IsGreaterThan(response.TopProviders[j].Total)
	})
	response.TopProviders = collection.Take(response.TopProviders, int(query.TopProviders))
	sort.Slice(response.TopLabels, func(i, j int) bool {
		return response.TopLabels[i].Total.IsGreaterThan(response.TopLabels[j].Total)
	})
	response.TopLabels = collection.Take(response.TopLabels, int(query.TopLabels))
	return result.Success(response)
}
