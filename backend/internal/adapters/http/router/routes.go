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

// BuildRoutesModule godoc
//
// @title					SubTracker API
// @version				1.0
// @description			This api provide HTTPRest endpoints for the application SubTracker.
// @termsOfService			http://subtracker.mistribe.com/terms/
// @contact.name			API Support
// @contact.url			http://subtracker.mistribe.com/support
// @contact.email			support@mistribe.com
// @license.name			Apache 2.0
// @license.url			http://www.apache.org/licenses/LICENSE-2.0.html
// @servers.url			https://api.subtracker.mistribe.com
// @servers.description	Production server
// @servers.url			http://localhost:8080
// @servers.description	Development server
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
			user.NewUserUpdateProfileEndpoint,
			user.NewUserDeleteEndpoint,
			ginfx.AsEndpointGroup(user.NewEndpointGroup),

			ginfx.AsEndpoint(NewHealthCheckLiveEndpoint),
		),
	)
}
