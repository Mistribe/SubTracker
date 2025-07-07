package api

import (
	"github.com/oleexo/subtracker/internal/infrastructure/api/endpoints"
	"github.com/oleexo/subtracker/internal/infrastructure/api/ginfx"
	"go.uber.org/fx"
)

func BuildRoutesModule() fx.Option {
	return fx.Module("routes",
		fx.Provide(
			endpoints.NewSubscriptionGetEndpoint,
			ginfx.AsRouteGroup(endpoints.NewSubscriptionEndpointGroup),
		),
	)
}
