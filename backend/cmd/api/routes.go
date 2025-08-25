package main

import (
	"go.uber.org/fx"

	. "github.com/mistribe/subtracker/cmd/api/endpoints"
	"github.com/mistribe/subtracker/cmd/api/ginfx"
	"github.com/mistribe/subtracker/cmd/api/middlewares"
)

func BuildRoutesModule() fx.Option {
	return fx.Module("routes",
		fx.Provide(
			middlewares.NewAuthenticationMiddleware,
			middlewares.NewLanguageMiddleware,

			NewSubscriptionGetEndpoint,
			NewSubscriptionGetAllEndpoint,
			NewSubscriptionCreateEndpoint,
			NewSubscriptionUpdateEndpoint,
			NewSubscriptionDeleteEndpoint,
			NewSubscriptionPatchEndpoint,
			NewSubscriptionSummaryEndpoint,
			ginfx.AsRouteGroup(NewSubscriptionEndpointGroup),

			NewLabelGetAllEndpoint,
			NewLabelGetEndpoint,
			NewLabelCreateEndpoint,
			NewLabelUpdateEndpoint,
			NewLabelDeleteEndpoint,
			NewDefaultLabelEndpoint,
			ginfx.AsRouteGroup(NewLabelEndpointGroup),

			NewFamilyCreateEndpoint,
			NewFamilyUpdateEndpoint,
			NewFamilyPatchEndpoint,
			NewFamilyDeleteEndpoint,
			NewFamilyInviteEndpoint,
			NewFamilyAcceptInvitationEndpoint,
			NewFamilyDeclineEndpoint,
			NewFamilyRevokeEndpoint,
			NewFamilySeeInvitationEndpoint,
			NewFamilyMemberGetEndpoint,
			NewFamilyMemberCreateEndpoint,
			NewFamilyMemberUpdateEndpoint,
			NewFamilyMemberDeleteEndpoint,
			ginfx.AsRouteGroup(NewFamilyEndpointGroup),

			NewProviderGetEndpoint,
			NewProviderGetAllEndpoint,
			NewProviderCreateEndpoint,
			NewProviderPlanCreateEndpoint,
			NewProviderPriceCreateEndpoint,
			NewProviderUpdateEndpoint,
			NewProviderPlanUpdateEndpoint,
			NewProviderPriceUpdateEndpoint,
			NewProviderDeleteEndpoint,
			NewProviderPlanDeleteEndpoint,
			NewProviderPriceDeleteEndpoint,
			NewProviderPatchEndpoint,
			ginfx.AsRouteGroup(NewProviderEndpointGroup),

			NewCurrencySupportedEndpoint,
			NewCurrencyGetRateEndpoint,
			NewCurrencyRefreshRatesEndpoint,
			ginfx.AsRouteGroup(NewCurrencyGroupEndpointGroup),

			NewUserGetPreferredCurrencyEndpoint,
			NewUserUpdatePreferredCurrencyEndpoint,
			NewUserUpdateProfileEndpoint,
			NewUserDeleteEndpoint,
			ginfx.AsRouteGroup(NewUserEndpointGroup),

			ginfx.AsRoute(NewHealthCheckLiveEndpoint),
		),
	)
}
