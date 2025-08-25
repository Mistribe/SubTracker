package currency

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/application/core"
	currencyCommand "github.com/mistribe/subtracker/internal/application/currency/command"
	currencyQuery "github.com/mistribe/subtracker/internal/application/currency/query"
)

func Module() fx.Option {
	return fx.Module("app_currency",
		fx.Provide(
			NewService,
			core.AsQueryHandler[currencyQuery.CurrencyRateQuery, currencyQuery.CurrencyRateResponse](currencyQuery.NewCurrencyRateQueryHandler),
			core.AsCommandHandler[currencyCommand.RefreshCurrencyRatesCommand, currencyCommand.RefreshCurrencyRatesResponse](currencyCommand.NewRefreshCurrencyRatesCommandHandler),
		),
	)
}
