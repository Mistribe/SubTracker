package application

import (
	"go.uber.org/fx"

	lblCommand "github.com/oleexo/subtracker/internal/application/label/command"
	lblQuery "github.com/oleexo/subtracker/internal/application/label/query"
	subCommand "github.com/oleexo/subtracker/internal/application/subscription/command"
	subQuery "github.com/oleexo/subtracker/internal/application/subscription/query"
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
			lblQuery.NewFindAllQueryHandler,
			lblCommand.NewCreateLabelCommandHandler,
			lblCommand.NewUpdateLabelCommandHandler,
			lblCommand.NewDeleteLabelCommandHandler,
		),
	)
}
