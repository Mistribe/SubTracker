package account

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/account/command"
	"github.com/mistribe/subtracker/internal/usecase/account/query"
)

func Module() fx.Option {
	return fx.Module("app_account",
		fx.Provide(
			NewService,
			ports.AsQueryHandler[query.FindPreferredCurrencyQuery, currency.Unit](query.NewFindPreferredCurrencyQueryHandler),
			ports.AsCommandHandler[command.UpdatePreferredCurrencyCommand, bool](command.NewUpdatePreferredCurrencyCommandHandler),
			ports.AsCommandHandler[command.DeleteAccountCommand, bool](command.NewDeleteUserCommandHandler),
		),
	)
}
