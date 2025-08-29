package user

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/user/command"
	usrQuery "github.com/mistribe/subtracker/internal/usecase/user/query"
)

func Module() fx.Option {
	return fx.Module("app_user",
		fx.Provide(
			NewService,
			ports.AsQueryHandler[usrQuery.FindPreferredCurrencyQuery, currency.Unit](usrQuery.NewFindPreferredCurrencyQueryHandler),
			ports.AsCommandHandler[command.UpdatePreferredCurrencyCommand, bool](command.NewUpdatePreferredCurrencyCommandHandler),
			ports.AsCommandHandler[command.DeleteUserCommand, bool](command.NewDeleteUserCommandHandler),
		),
	)
}
