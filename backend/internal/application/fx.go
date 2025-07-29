package application

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/application/core"
	fmlyCommand "github.com/oleexo/subtracker/internal/application/family/command"
	fmlyQuery "github.com/oleexo/subtracker/internal/application/family/query"
	lblCommand "github.com/oleexo/subtracker/internal/application/label/command"
	lblQuery "github.com/oleexo/subtracker/internal/application/label/query"
	proCommand "github.com/oleexo/subtracker/internal/application/provider/command"
	proQuery "github.com/oleexo/subtracker/internal/application/provider/query"
	subCommand "github.com/oleexo/subtracker/internal/application/subscription/command"
	subQuery "github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/provider"
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
			AsQueryHandler[proQuery.FindOneQuery, provider.Provider](proQuery.NewFindOneQueryHandler),
			AsQueryHandler[proQuery.FindAllQuery, core.PaginatedResponse[provider.Provider]](proQuery.NewFindAllQueryHandler),
			AsCommandHandler[proCommand.CreateCommand, provider.Provider](proCommand.NewCreateCommandHandler),
			AsCommandHandler[proCommand.UpdateCommand, provider.Provider](proCommand.NewUpdateCommandHandler),
			AsCommandHandler[proCommand.DeleteCommand, bool](proCommand.NewDeleteCommandHandler),

			AsQueryHandler[subQuery.FindOneQuery, subscription.Subscription](subQuery.NewFindOneQueryHandler),
			AsQueryHandler[subQuery.FindAllQuery, core.PaginatedResponse[subscription.Subscription]](subQuery.NewFindAllQueryHandler),

			AsCommandHandler[subCommand.CreateSubscriptionCommand, subscription.Subscription](subCommand.NewCreateSubscriptionCommandHandler),
			AsCommandHandler[subCommand.UpdateSubscriptionCommand, subscription.Subscription](subCommand.NewUpdateSubscriptionCommandHandler),
			AsCommandHandler[subCommand.DeleteSubscriptionCommand, bool](subCommand.NewDeleteSubscriptionCommandHandler),
			AsCommandHandler[subCommand.PatchSubscriptionCommand, subscription.Subscription](subCommand.NewPatchSubscriptionCommandHandler),

			AsQueryHandler[lblQuery.FindAllQuery, core.PaginatedResponse[label.Label]](lblQuery.NewFindAllQueryHandler),
			AsQueryHandler[lblQuery.FindOneQuery, label.Label](lblQuery.NewFindOneQueryHandler),
			AsQueryHandler[lblQuery.DefaultLabelQuery, []label.Label](lblQuery.NewDefaultLabelQueryHandler),

			AsCommandHandler[lblCommand.CreateLabelCommand, label.Label](lblCommand.NewCreateLabelCommandHandler),
			AsCommandHandler[lblCommand.UpdateLabelCommand, label.Label](lblCommand.NewUpdateLabelCommandHandler),
			AsCommandHandler[lblCommand.DeleteLabelCommand, bool](lblCommand.NewDeleteLabelCommandHandler),

			AsQueryHandler[fmlyQuery.FindAllQuery, core.PaginatedResponse[family.Family]](fmlyQuery.NewFindAllQueryHandler),
			AsQueryHandler[fmlyQuery.FindOneQuery, family.Family](fmlyQuery.NewFindOneQueryHandler),

			AsCommandHandler[fmlyCommand.CreateFamilyCommand, family.Family](fmlyCommand.NewCreateFamilyCommandHandler),
			AsCommandHandler[fmlyCommand.UpdateFamilyCommand, family.Family](fmlyCommand.NewUpdateFamilyCommandHandler),
			AsCommandHandler[fmlyCommand.PatchFamilyCommand, family.Family](fmlyCommand.NewPatchFamilyCommandHandler),
			AsCommandHandler[fmlyCommand.DeleteFamilyCommand, bool](fmlyCommand.NewDeleteFamilyCommandHandler),
			AsCommandHandler[fmlyCommand.CreateFamilyMemberCommand, family.Family](fmlyCommand.NewCreateFamilyMemberCommandHandler),
			AsCommandHandler[fmlyCommand.UpdateFamilyMemberCommand, family.Family](fmlyCommand.NewUpdateFamilyMemberCommandHandler),
			AsCommandHandler[fmlyCommand.DeleteFamilyMemberCommand, bool](fmlyCommand.NewDeleteFamilyMemberCommandHandler),
		),
	)
}
