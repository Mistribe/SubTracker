package provider

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/application/core"
	"github.com/mistribe/subtracker/internal/application/provider/command"
	"github.com/mistribe/subtracker/internal/application/provider/query"
	"github.com/mistribe/subtracker/internal/domain/provider"
)

func Module() fx.Option {
	return fx.Module("app_provider",
		fx.Provide(
			core.AsQueryHandler[query.FindOneQuery, provider.Provider](query.NewFindOneQueryHandler),
			core.AsQueryHandler[query.FindAllQuery, core.PaginatedResponse[provider.Provider]](query.NewFindAllQueryHandler),
			core.AsCommandHandler[command.CreateProviderCommand, provider.Provider](command.NewCreateProviderCommandHandler),
			core.AsCommandHandler[command.CreatePlanCommand, provider.Plan](command.NewCreatePlanCommandHandler),
			core.AsCommandHandler[command.CreatePriceCommand, provider.Price](command.NewCreatePriceCommandHandler),
			core.AsCommandHandler[command.UpdateProviderCommand, provider.Provider](command.NewUpdateProviderCommandHandler),
			core.AsCommandHandler[command.UpdatePlanCommand, provider.Plan](command.NewUpdatePlanCommandHandler),
			core.AsCommandHandler[command.UpdatePriceCommand, provider.Price](command.NewUpdatePriceCommandHandler),
			core.AsCommandHandler[command.PatchProviderCommand, provider.Provider](command.NewPatchProviderCommandHandler),
			core.AsCommandHandler[command.DeleteProviderCommand, bool](command.NewDeleteProviderCommandHandler),
			core.AsCommandHandler[command.DeletePlanCommand, bool](command.NewDeletePlanCommandHandler),
			core.AsCommandHandler[command.DeletePriceCommand, bool](command.NewDeletePriceCommandHandler),
		),
	)
}
