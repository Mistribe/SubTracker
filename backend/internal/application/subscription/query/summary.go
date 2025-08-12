package query

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/currency"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type SummaryQuery struct {
	TopProviders     uint8
	UpcomingRenewals uint8
	TotalMonthly     bool
	TotalYearly      bool
}

type SummaryQueryUpcomingRenewalsResponse struct {
	ProviderId uuid.UUID
	At         time.Time
	Total      float64
}

type SummaryQueryResponse struct {
	Currency         currency.Unit
	TotalMonthly     float64
	TotalYearly      float64
	UpcomingRenewals []SummaryQueryUpcomingRenewalsResponse
	TopProviders     map[uuid.UUID]float64
}

type SummaryQueryHandler struct {
	subscriptionRepository subscription.Repository
	currencyRepository     currency.Repository
	userRepository         user.Repository
	authService            auth.Service
}

func NewSummaryQueryHandler(
	subscriptionRepository subscription.Repository,
	currencyRepository currency.Repository,
	userRepository user.Repository,
	authService auth.Service) *SummaryQueryHandler {
	return &SummaryQueryHandler{
		subscriptionRepository: subscriptionRepository,
		currencyRepository:     currencyRepository,
		userRepository:         userRepository,
		authService:            authService,
	}
}

func (h SummaryQueryHandler) Handle(ctx context.Context, query SummaryQuery) result.Result[SummaryQueryResponse] {
	userId := h.authService.MustGetUserId(ctx)
	userProfile, err := h.userRepository.GetUserProfile(ctx, userId)
	if err != nil {
		return result.Fail[SummaryQueryResponse](err)
	}
	currencyRates, err := h.currencyRepository.GetRatesByDate(ctx, time.Now())
	if err != nil {
		return result.Fail[SummaryQueryResponse](err)
	}
	preferredCurrency := userProfile.Currency()

	response := SummaryQueryResponse{
		Currency: preferredCurrency,
	}
	for sub := range h.subscriptionRepository.GetAllIt(ctx) {
		if query.TotalMonthly {
			monthlyAmount := sub.GetRecurrencyAmount(subscription.MonthlyRecurrency)
			if monthlyAmount.IsValid() {
				response.TotalMonthly += monthlyAmount.ToCurrency(preferredCurrency, currencyRates).Value()
			}
		}
		if query.TotalYearly {
			yearlyAmount := sub.GetRecurrencyAmount(subscription.YearlyRecurrency)
			if yearlyAmount.IsValid() {
				response.TotalYearly += yearlyAmount.ToCurrency(preferredCurrency, currencyRates).Value()
			}
		}

		if query.UpcomingRenewals > 0 {
			renewalDate := sub.GetNextRenewalDate()
			price := sub.GetPrice().ToCurrency(preferredCurrency, currencyRates)
			if renewalDate != nil && price.IsValid() {
				response.UpcomingRenewals = append(response.UpcomingRenewals, SummaryQueryUpcomingRenewalsResponse{
					ProviderId: sub.ProviderId(),
					At:         *renewalDate,
					Total:      price.Value(),
				})
			}
		}

		if query.TopProviders > 0 {
			totalSpent := sub.GetTotalSpent()
			if totalSpent.IsValid() {
				if response.TopProviders == nil {
					response.TopProviders = make(map[uuid.UUID]float64)
				}
				existingValue, ok := response.TopProviders[sub.ProviderId()]
				if ok {
					response.TopProviders[sub.ProviderId()] = existingValue + totalSpent.Value()
				} else {
					response.TopProviders[sub.ProviderId()] = totalSpent.Value()
				}
			}
		}
	}

	return result.Success(response)
}
