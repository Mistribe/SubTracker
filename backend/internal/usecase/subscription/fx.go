package subscription

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	command2 "github.com/mistribe/subtracker/internal/usecase/subscription/command"
	query2 "github.com/mistribe/subtracker/internal/usecase/subscription/query"
)

func Module() fx.Option {
	return fx.Module("app_subscription",
		fx.Provide(
			ports.AsQueryHandler[query2.SummaryQuery, query2.SummaryQueryResponse](query2.NewSummaryQueryHandler),
			ports.AsQueryHandler[query2.FindOneQuery, subscription.Subscription](query2.NewFindOneQueryHandler),
			ports.AsQueryHandler[query2.FindAllQuery, shared.PaginatedResponse[subscription.Subscription]](query2.NewFindAllQueryHandler),

			ports.AsCommandHandler[command2.CreateSubscriptionCommand, subscription.Subscription](command2.NewCreateSubscriptionCommandHandler),
			ports.AsCommandHandler[command2.UpdateSubscriptionCommand, subscription.Subscription](command2.NewUpdateSubscriptionCommandHandler),
			ports.AsCommandHandler[command2.DeleteSubscriptionCommand, bool](command2.NewDeleteSubscriptionCommandHandler),
			ports.AsCommandHandler[command2.PatchSubscriptionCommand, subscription.Subscription](command2.NewPatchSubscriptionCommandHandler),
		),
	)
}
