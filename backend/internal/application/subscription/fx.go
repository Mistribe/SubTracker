package subscription

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/command"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

func Module() fx.Option {
	return fx.Module("app_subscription",
		fx.Provide(
			core.AsQueryHandler[query.SummaryQuery, query.SummaryQueryResponse](query.NewSummaryQueryHandler),
			core.AsQueryHandler[query.FindOneQuery, subscription.Subscription](query.NewFindOneQueryHandler),
			core.AsQueryHandler[query.FindAllQuery, core.PaginatedResponse[subscription.Subscription]](query.NewFindAllQueryHandler),

			core.AsCommandHandler[command.CreateSubscriptionCommand, subscription.Subscription](command.NewCreateSubscriptionCommandHandler),
			core.AsCommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription](command.NewUpdateSubscriptionCommandHandler),
			core.AsCommandHandler[command.DeleteSubscriptionCommand, bool](command.NewDeleteSubscriptionCommandHandler),
			core.AsCommandHandler[command.PatchSubscriptionCommand, subscription.Subscription](command.NewPatchSubscriptionCommandHandler),
		),
	)
}
