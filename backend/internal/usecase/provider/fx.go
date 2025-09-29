package provider

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
	"github.com/mistribe/subtracker/internal/usecase/provider/query"
)

func Module() fx.Option {
	return fx.Module("app_provider",
		fx.Provide(
			ports.AsQueryHandler[query.FindOneQuery, provider.Provider](query.NewFindOneQueryHandler),
			ports.AsQueryHandler[query.FindAllQuery, shared.PaginatedResponse[provider.Provider]](query.NewFindAllQueryHandler),
			ports.AsQueryHandler[query.GetQuotaUsageHandler, billing.EffectiveEntitlement](query.NewGetQuotaUsageHandler),

			ports.AsCommandHandler[command.CreateProviderCommand, provider.Provider](command.NewCreateProviderCommandHandler),
			ports.AsCommandHandler[command.UpdateProviderCommand, provider.Provider](command.NewUpdateProviderCommandHandler),
			ports.AsCommandHandler[command.DeleteProviderCommand, bool](command.NewDeleteProviderCommandHandler),
		),
	)
}
