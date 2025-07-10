package api

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/infrastructure/api/endpoints"
	"github.com/oleexo/subtracker/internal/infrastructure/api/ginfx"
)

func BuildRoutesModule() fx.Option {
	return fx.Module("routes",
		fx.Provide(
			endpoints.NewSubscriptionGetEndpoint,
			endpoints.NewSubscriptionGetAllEndpoint,
			endpoints.NewSubscriptionCreateEndpoint,
			ginfx.AsRouteGroup(endpoints.NewSubscriptionEndpointGroup),
			endpoints.NewLabelGetAllEndpoint,
			endpoints.NewLabelGetEndpoint,
			endpoints.NewLabelCreateEndpoint,
			endpoints.NewLabelUpdateEndpoint,
			endpoints.NewLabelDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints.NewLabelEndpointGroup),
		),
	)
}
