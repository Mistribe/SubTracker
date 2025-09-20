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
			middlewares.NewCacheMiddleware,

			subscription.NewGetEndpoint,
			subscription.NewGetAllEndpoint,
			subscription.NewCreateEndpoint,
			subscription.NewUpdateEndpoint,
			subscription.NewDeleteEndpoint,
			subscription.NewSummaryEndpoint,
			ginfx.AsEndpointGroup(subscription.NewEndpointGroup),

			label.NewGetAllEndpoint,
			label.NewGetEndpoint,
			label.NewCreateEndpoint,
			label.NewUpdateEndpoint,
			label.NewDeleteEndpoint,
			ginfx.AsEndpointGroup(label.NewEndpointGroup),

			family.NewCreateEndpoint,
			family.NewUpdateEndpoint,
			family.NewDeleteEndpoint,
			family.NewInviteEndpoint,
			family.NewAcceptInvitationEndpoint,
			family.NewDeclineEndpoint,
			family.NewRevokeEndpoint,
			family.NewSeeInvitationEndpoint,
			family.NewGetEndpoint,
			family.NewMemberCreateEndpoint,
			family.NewMemberUpdateEndpoint,
			family.NewMemberDeleteEndpoint,
			ginfx.AsEndpointGroup(family.NewEndpointGroup),

			provider.NewGetEndpoint,
			provider.NewGetAllEndpoint,
			provider.NewCreateEndpoint,
			provider.NewUpdateEndpoint,
			provider.NewDeleteEndpoint,
			ginfx.AsEndpointGroup(provider.NewEndpointGroup),

			currency.NewSupportedEndpoint,
			currency.NewGetRateEndpoint,
			ginfx.AsEndpointGroup(currency.NewEndpointGroup),

			user.NewUserGetPreferredCurrencyEndpoint,
			user.NewUpdatePreferredCurrencyEndpoint,
			user.NewDeleteEndpoint,
			ginfx.AsEndpointGroup(user.NewEndpointGroup),

			ginfx.AsEndpoint(NewHealthCheckLiveEndpoint),
		),
	)
}
