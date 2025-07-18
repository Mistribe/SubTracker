package main

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/cmd/api/endpoints"
	"github.com/oleexo/subtracker/cmd/api/ginfx"
	"github.com/oleexo/subtracker/cmd/api/middlewares"
)

func BuildRoutesModule() fx.Option {
	return fx.Module("routes",
		fx.Provide(
			middlewares.NewAuthenticationMiddleware,
			endpoints.NewSubscriptionGetEndpoint,
			endpoints.NewSubscriptionGetAllEndpoint,
			endpoints.NewSubscriptionCreateEndpoint,
			endpoints.NewSubscriptionUpdateEndpoint,
			endpoints.NewSubscriptionDeleteEndpoint,
			endpoints.NewSubscriptionPaymentCreateEndpoint,
			endpoints.NewSubscriptionPaymentUpdateEndpoint,
			endpoints.NewSubscriptionPaymentDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints.NewSubscriptionEndpointGroup),

			endpoints.NewLabelGetAllEndpoint,
			endpoints.NewLabelGetEndpoint,
			endpoints.NewLabelCreateEndpoint,
			endpoints.NewLabelUpdateEndpoint,
			endpoints.NewLabelDeleteEndpoint,
			endpoints.NewDefaultLabelEndpoint,
			ginfx.AsRouteGroup(endpoints.NewLabelEndpointGroup),

			endpoints.NewFamilyCreateEndpoint,
			endpoints.NewFamilyUpdateEndpoint,
			endpoints.NewFamilyMemberGetAllEndpoint,
			endpoints.NewFamilyMemberGetEndpoint,
			endpoints.NewFamilyMemberCreateEndpoint,
			endpoints.NewFamilyMemberUpdateEndpoint,
			endpoints.NewFamilyMemberDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints.NewFamilyEndpointGroup),

			ginfx.AsRoute(endpoints.NewHealthCheckLiveEndpoint),
		),
	)
}
