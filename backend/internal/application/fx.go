package application

import (
	"go.uber.org/fx"

	fmlyCommand "github.com/oleexo/subtracker/internal/application/family/command"
	fmlyQuery "github.com/oleexo/subtracker/internal/application/family/query"
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
			subCommand.NewUpdatePaymentCommandHandler,
			lblQuery.NewFindAllQueryHandler,
			lblQuery.NewFindOneQueryHandler,
			lblCommand.NewCreateLabelCommandHandler,
			lblCommand.NewUpdateLabelCommandHandler,
			lblCommand.NewDeleteLabelCommandHandler,
			fmlyQuery.NewFindAllQueryHandler,
			fmlyQuery.NewFindOneQueryHandler,
			fmlyCommand.NewCreateFamilyMemberCommandHandler,
			fmlyCommand.NewUpdateFamilyMemberCommandHandler,
			fmlyCommand.NewDeleteFamilyMemberCommandHandler,
		),
	)
}
