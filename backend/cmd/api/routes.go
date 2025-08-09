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
			middlewares.NewLanguageMiddleware,

			endpoints.NewSubscriptionGetEndpoint,
			endpoints.NewSubscriptionGetAllEndpoint,
			endpoints.NewSubscriptionCreateEndpoint,
			endpoints.NewSubscriptionUpdateEndpoint,
			endpoints.NewSubscriptionDeleteEndpoint,
			endpoints.NewSubscriptionPatchEndpoint,
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
			endpoints.NewFamilyPatchEndpoint,
			endpoints.NewFamilyDeleteEndpoint,
			endpoints.NewFamilyMemberGetAllEndpoint,
			endpoints.NewFamilyMemberGetEndpoint,
			endpoints.NewFamilyMemberCreateEndpoint,
			endpoints.NewFamilyMemberUpdateEndpoint,
			endpoints.NewFamilyMemberDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints.NewFamilyEndpointGroup),

			endpoints.NewProviderGetEndpoint,
			endpoints.NewProviderGetAllEndpoint,
			endpoints.NewProviderCreateEndpoint,
			endpoints.NewProviderPlanCreateEndpoint,
			endpoints.NewProviderPriceCreateEndpoint,
			endpoints.NewProviderUpdateEndpoint,
			endpoints.NewProviderPlanUpdateEndpoint,
			endpoints.NewProviderPriceUpdateEndpoint,
			endpoints.NewProviderDeleteEndpoint,
			endpoints.NewProviderPlanDeleteEndpoint,
			endpoints.NewProviderPriceDeleteEndpoint,
			endpoints.NewProviderPatchEndpoint,
			ginfx.AsRouteGroup(endpoints.NewProviderEndpointGroup),

			endpoints.NewCurrencySupportedEndpoint,
			endpoints.NewCurrencyRateEndpoint,
			ginfx.AsRouteGroup(endpoints.NewCurrencyGroupEndpointGroup),

			endpoints.NewUserGetPreferredCurrencyEndpoint,
			endpoints.NewUserUpdatePreferredCurrencyEndpoint,
			endpoints.NewUserUpdateProfileEndpoint,
			endpoints.NewUserDeleteEndpoint,
			ginfx.AsRouteGroup(endpoints.NewUserEndpointGroup),

			ginfx.AsRoute(endpoints.NewHealthCheckLiveEndpoint),
		),
	)
}
