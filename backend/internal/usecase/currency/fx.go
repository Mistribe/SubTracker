package currency

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/ports"
	currencyCommand "github.com/mistribe/subtracker/internal/usecase/currency/command"
	currencyQuery "github.com/mistribe/subtracker/internal/usecase/currency/query"
)

func Module() fx.Option {
	return fx.Module("app_currency",
		fx.Provide(
			NewService,
			ports.AsQueryHandler[currencyQuery.CurrencyRateQuery, currencyQuery.CurrencyRateResponse](currencyQuery.NewCurrencyRateQueryHandler),
			ports.AsCommandHandler[currencyCommand.RefreshCurrencyRatesCommand, currencyCommand.RefreshCurrencyRatesResponse](currencyCommand.NewRefreshCurrencyRatesCommandHandler),
		),
	)
}
