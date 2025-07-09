package application

import (
	subCommand "github.com/oleexo/subtracker/internal/application/subscription/command"
	subQuery "github.com/oleexo/subtracker/internal/application/subscription/query"
	"go.uber.org/fx"
)

func BuildApplicationModule() fx.Option {
	return fx.Module("application",
		fx.Provide(
			subQuery.NewFindOneQueryHandler,
			subQuery.NewFindAllQueryHandler,
			subCommand.NewCreateSubscriptionCommandHandler,
			subCommand.NewUpdateSubscriptionCommandHandler,
			subCommand.NewDeleteSubscriptionCommandHandler,
			subCommand.NewCreatePaymentCommandHandler,
		),
	)
}
