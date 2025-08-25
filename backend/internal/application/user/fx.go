package user

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/application/core"
	usrCommand "github.com/mistribe/subtracker/internal/application/user/command"
	usrQuery "github.com/mistribe/subtracker/internal/application/user/query"
	"github.com/mistribe/subtracker/internal/domain/currency"
)

func Module() fx.Option {
	return fx.Module("app_user",
		fx.Provide(
			NewService,
			core.AsQueryHandler[usrQuery.FindPreferredCurrencyQuery, currency.Unit](usrQuery.NewFindPreferredCurrencyQueryHandler),
			core.AsCommandHandler[usrCommand.UpdatePreferredCurrencyCommand, bool](usrCommand.NewUpdatePreferredCurrencyCommandHandler),
			core.AsCommandHandler[usrCommand.UpdateProfileCommand, bool](usrCommand.NewUpdateProfileCommandHandler),
			core.AsCommandHandler[usrCommand.DeleteUserCommand, bool](usrCommand.NewDeleteUserCommandHandler),
		),
	)
}
