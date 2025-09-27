package provider

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	command2 "github.com/mistribe/subtracker/internal/usecase/provider/command"
	query2 "github.com/mistribe/subtracker/internal/usecase/provider/query"
)

func Module() fx.Option {
	return fx.Module("app_provider",
		fx.Provide(
			ports.AsQueryHandler[query2.FindOneQuery, provider.Provider](query2.NewFindOneQueryHandler),
			ports.AsQueryHandler[query2.FindAllQuery, shared.PaginatedResponse[provider.Provider]](query2.NewFindAllQueryHandler),
			ports.AsCommandHandler[command2.CreateProviderCommand, provider.Provider](command2.NewCreateProviderCommandHandler),
			ports.AsCommandHandler[command2.UpdateProviderCommand, provider.Provider](command2.NewUpdateProviderCommandHandler),
			ports.AsCommandHandler[command2.DeleteProviderCommand, bool](command2.NewDeleteProviderCommandHandler),
		),
	)
}
