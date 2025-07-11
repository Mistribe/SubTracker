package main

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/cmd/api/endpoints"
	"github.com/oleexo/subtracker/cmd/api/ginfx"
)

func BuildRoutesModule() fx.Option {
	return fx.Module("routes",
		fx.Provide(
			endpoints.NewSubscriptionGetEndpoint,
			endpoints.NewSubscriptionGetAllEndpoint,
			endpoints.NewSubscriptionCreateEndpoint,
			endpoints.NewSubscriptionUpdateEndpoint,
			endpoints.NewSubscriptionDeleteEndpoint,
			endpoints.NewSubscriptionPaymentCreateEndpoint,
			endpoints.NewSubscriptionPaymentUpdateEndpoint,
			ginfx.AsRouteGroup(endpoints.NewSubscriptionEndpointGroup),
			endpoints.NewLabelGetAllEndpoint,
			endpoints.NewLabelGetEndpoint,
			endpoints.NewLabelCreateEndpoint,
			endpoints.NewLabelUpdateEndpoint,
			endpoints.NewLabelDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints.NewLabelEndpointGroup),
			endpoints.NewFamilyMemberGetAllEndpoint,
			endpoints.NewFamilyMemberGetEndpoint,
			endpoints.NewFamilyMemberCreateEndpoint,
			endpoints.NewFamilyMemberUpdateEndpoint,
			endpoints.NewFamilyMemberDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints.NewFamilyEndpointGroup),
		),
	)
}
