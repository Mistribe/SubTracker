package currency

import (
	"context"
	"time"

	"github.com/oleexo/subtracker/internal/domain/currency"
)

type service struct {
	currencyRepository currency.Repository
	// todo add cache
}

func (s service) ConvertTo(ctx context.Context, from currency.Amount, to currency.Unit) (currency.Amount, error) {
	rates, err := s.currencyRepository.GetRatesByDate(ctx, time.Now())
	if err != nil {
		return currency.NewInvalidAmount(), err
	}
	rates = rates.WithReverse()
	return from.ToCurrency(to, rates), nil
}

func NewService(currencyRepository currency.Repository) currency.Service {
	return service{
		currencyRepository: currencyRepository,
	}
}
