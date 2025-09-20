package subscription

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
)

func Module() fx.Option {
	return fx.Module("app_subscription",
		fx.Provide(
			ports.AsQueryHandler[query.SummaryQuery, query.SummaryQueryResponse](query.NewSummaryQueryHandler),
			ports.AsQueryHandler[query.FindOneQuery, subscription.Subscription](query.NewFindOneQueryHandler),
			ports.AsQueryHandler[query.FindAllQuery, shared.PaginatedResponse[subscription.Subscription]](query.NewFindAllQueryHandler),

			ports.AsCommandHandler[command.CreateSubscriptionCommand, subscription.Subscription](command.NewCreateSubscriptionCommandHandler),
			ports.AsCommandHandler[command.UpdateSubscriptionCommand, subscription.Subscription](command.NewUpdateSubscriptionCommandHandler),
			ports.AsCommandHandler[command.DeleteSubscriptionCommand, bool](command.NewDeleteSubscriptionCommandHandler),
		),
	)
}
