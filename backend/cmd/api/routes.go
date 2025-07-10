package main

import (
	"go.uber.org/fx"

	endpoints2 "github.com/oleexo/subtracker/cmd/api/endpoints"
	"github.com/oleexo/subtracker/cmd/api/ginfx"
)

func BuildRoutesModule() fx.Option {
	return fx.Module("routes",
		fx.Provide(
			endpoints2.NewSubscriptionGetEndpoint,
			endpoints2.NewSubscriptionGetAllEndpoint,
			endpoints2.NewSubscriptionCreateEndpoint,
			ginfx.AsRouteGroup(endpoints2.NewSubscriptionEndpointGroup),
			endpoints2.NewLabelGetAllEndpoint,
			endpoints2.NewLabelGetEndpoint,
			endpoints2.NewLabelCreateEndpoint,
			endpoints2.NewLabelUpdateEndpoint,
			endpoints2.NewLabelDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints2.NewLabelEndpointGroup),
			endpoints2.NewFamilyMemberGetAllEndpoint,
			endpoints2.NewFamilyMemberGetEndpoint,
			endpoints2.NewFamilyMemberCreateEndpoint,
			endpoints2.NewFamilyMemberUpdateEndpoint,
			endpoints2.NewFamilyMemberDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints2.NewFamilyEndpointGroup),
		),
	)
}
