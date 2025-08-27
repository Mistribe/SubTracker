package currency

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/ports"
	currencyQuery "github.com/mistribe/subtracker/internal/usecase/currency/query"
)

func Module() fx.Option {
	return fx.Module("app_currency",
		fx.Provide(
			ports.AsQueryHandler[currencyQuery.CurrencyRateQuery, currencyQuery.CurrencyRateResponse](currencyQuery.NewCurrencyRateQueryHandler),
		),
	)
}
