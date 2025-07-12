package application

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/core/result"
	fmlyCommand "github.com/oleexo/subtracker/internal/application/family/command"
	fmlyQuery "github.com/oleexo/subtracker/internal/application/family/query"
	lblCommand "github.com/oleexo/subtracker/internal/application/label/command"
	lblQuery "github.com/oleexo/subtracker/internal/application/label/query"
	subCommand "github.com/oleexo/subtracker/internal/application/subscription/command"
	subQuery "github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

func AsCommandHandler[TCommand core.Command, TResult any](f any) any {
	return fx.Annotate(f,
		fx.As(new(core.CommandHandler[TCommand, TResult])),
	)
}

func AsQueryHandler[TQuery core.Query, TResult any](f any) any {
	return fx.Annotate(f,
		fx.As(new(core.QueryHandler[TQuery, TResult])),
	)
}

func BuildApplicationModule() fx.Option {
	return fx.Module("application",
		fx.Provide(
			AsQueryHandler[subQuery.FindOneQuery, subscription.Subscription](subQuery.NewFindOneQueryHandler),
			AsQueryHandler[subQuery.FindAllQuery, []subscription.Subscription](subQuery.NewFindAllQueryHandler),

			AsCommandHandler[subCommand.CreateSubscriptionCommand, subscription.Subscription](subCommand.NewCreateSubscriptionCommandHandler),
			AsCommandHandler[subCommand.UpdateSubscriptionCommand, subscription.Subscription](subCommand.NewUpdateSubscriptionCommandHandler),
			AsCommandHandler[subCommand.DeleteSubscriptionCommand, result.Unit](subCommand.NewDeleteSubscriptionCommandHandler),
			AsCommandHandler[subCommand.CreatePaymentCommand, subscription.Subscription](subCommand.NewCreatePaymentCommandHandler),
			AsCommandHandler[subCommand.UpdatePaymentCommand, subscription.Subscription](subCommand.NewUpdatePaymentCommandHandler),
			AsCommandHandler[subCommand.DeletePaymentCommand, result.Unit](subCommand.NewDeletePaymentCommandHandler),

			AsQueryHandler[lblQuery.FindAllQuery, []label.Label](lblQuery.NewFindAllQueryHandler),
			AsQueryHandler[lblQuery.FindOneQuery, label.Label](lblQuery.NewFindOneQueryHandler),
			AsQueryHandler[lblQuery.DefaultLabelQuery, []label.Label](lblQuery.NewDefaultLabelQueryHandler),

			AsCommandHandler[lblCommand.CreateLabelCommand, label.Label](lblCommand.NewCreateLabelCommandHandler),
			AsCommandHandler[lblCommand.UpdateLabelCommand, label.Label](lblCommand.NewUpdateLabelCommandHandler),
			AsCommandHandler[lblCommand.DeleteLabelCommand, result.Unit](lblCommand.NewDeleteLabelCommandHandler),

			AsQueryHandler[fmlyQuery.FindAllQuery, []family.Member](fmlyQuery.NewFindAllQueryHandler),
			AsQueryHandler[fmlyQuery.FindOneQuery, family.Member](fmlyQuery.NewFindOneQueryHandler),

			AsCommandHandler[fmlyCommand.CreateFamilyMemberCommand, family.Member](fmlyCommand.NewCreateFamilyMemberCommandHandler),
			AsCommandHandler[fmlyCommand.UpdateFamilyMemberCommand, family.Member](fmlyCommand.NewUpdateFamilyMemberCommandHandler),
			AsCommandHandler[fmlyCommand.DeleteFamilyMemberCommand, result.Unit](fmlyCommand.NewDeleteFamilyMemberCommandHandler),
		),
	)
}
