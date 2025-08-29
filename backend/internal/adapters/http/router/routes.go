package router

import (
	"go.uber.org/fx"

	. "github.com/mistribe/subtracker/internal/adapters/http/handlers"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/currency"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/family"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/label"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/provider"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/subscription"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/user"
	"github.com/mistribe/subtracker/internal/adapters/http/router/ginfx"
	"github.com/mistribe/subtracker/internal/adapters/http/router/middlewares"
)

func BuildRoutesModule() fx.Option {
	return fx.Module("routes",
		fx.Provide(
			middlewares.NewAuthenticationMiddleware,
			middlewares.NewLanguageMiddleware,

			subscription.NewSubscriptionGetEndpoint,
			subscription.NewSubscriptionGetAllEndpoint,
			subscription.NewSubscriptionCreateEndpoint,
			subscription.NewSubscriptionUpdateEndpoint,
			subscription.NewSubscriptionDeleteEndpoint,
			subscription.NewSubscriptionPatchEndpoint,
			subscription.NewSubscriptionSummaryEndpoint,
			ginfx.AsEndpointGroup(subscription.NewEndpointGroup),

			label.NewLabelGetAllEndpoint,
			label.NewLabelGetEndpoint,
			label.NewLabelCreateEndpoint,
			label.NewLabelUpdateEndpoint,
			label.NewLabelDeleteEndpoint,
			label.NewDefaultLabelEndpoint,
			ginfx.AsEndpointGroup(label.NewEndpointGroup),

			family.NewFamilyCreateEndpoint,
			family.NewFamilyUpdateEndpoint,
			family.NewFamilyPatchEndpoint,
			family.NewFamilyDeleteEndpoint,
			family.NewFamilyInviteEndpoint,
			family.NewFamilyAcceptInvitationEndpoint,
			family.NewFamilyDeclineEndpoint,
			family.NewFamilyRevokeEndpoint,
			family.NewFamilySeeInvitationEndpoint,
			family.NewFamilyMemberGetEndpoint,
			family.NewFamilyMemberCreateEndpoint,
			family.NewFamilyMemberUpdateEndpoint,
			family.NewFamilyMemberDeleteEndpoint,
			ginfx.AsEndpointGroup(family.NewEndpointGroup),

			provider.NewProviderGetEndpoint,
			provider.NewProviderGetAllEndpoint,
			provider.NewProviderCreateEndpoint,
			provider.NewProviderPlanCreateEndpoint,
			provider.NewProviderPriceCreateEndpoint,
			provider.NewProviderUpdateEndpoint,
			provider.NewProviderPlanUpdateEndpoint,
			provider.NewProviderPriceUpdateEndpoint,
			provider.NewProviderDeleteEndpoint,
			provider.NewProviderPlanDeleteEndpoint,
			provider.NewProviderPriceDeleteEndpoint,
			provider.NewProviderPatchEndpoint,
			ginfx.AsEndpointGroup(provider.NewEndpointGroup),

			currency.NewCurrencySupportedEndpoint,
			currency.NewCurrencyGetRateEndpoint,
			ginfx.AsEndpointGroup(currency.NewEndpointGroup),

			user.NewUserGetPreferredCurrencyEndpoint,
			user.NewUserUpdatePreferredCurrencyEndpoint,
			user.NewUserDeleteEndpoint,
			ginfx.AsEndpointGroup(user.NewEndpointGroup),

			ginfx.AsEndpoint(NewHealthCheckLiveEndpoint),
		),
	)
}
