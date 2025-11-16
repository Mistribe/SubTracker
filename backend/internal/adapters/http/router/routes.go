package router

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/adapters/http/export"
	. "github.com/mistribe/subtracker/internal/adapters/http/handlers"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/account"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/currency"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/family"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/label"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/provider"
	"github.com/mistribe/subtracker/internal/adapters/http/handlers/subscription"
	"github.com/mistribe/subtracker/internal/adapters/http/router/ginfx"
	"github.com/mistribe/subtracker/internal/adapters/http/router/middlewares"
)

func BuildRoutesModule() fx.Option {
	return fx.Module("routes",
		export.Module(),
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
			subscription.NewGetQuotaUsageEndpoint,
			subscription.NewExportEndpoint,
			ginfx.AsEndpointGroup(subscription.NewEndpointGroup),

			label.NewGetAllEndpoint,
			label.NewGetEndpoint,
			label.NewCreateEndpoint,
			label.NewUpdateEndpoint,
			label.NewDeleteEndpoint,
			label.NewGetQuotaUsageEndpoint,
			label.NewExportEndpoint,
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
			family.NewGetQuotaUsageEndpoint,
			ginfx.AsEndpointGroup(family.NewEndpointGroup),

			provider.NewGetEndpoint,
			provider.NewGetAllEndpoint,
			provider.NewCreateEndpoint,
			provider.NewUpdateEndpoint,
			provider.NewDeleteEndpoint,
			provider.NewGetQuotaUsageEndpoint,
			provider.NewExportEndpoint,
			ginfx.AsEndpointGroup(provider.NewEndpointGroup),

			currency.NewSupportedEndpoint,
			currency.NewGetRateEndpoint,
			ginfx.AsEndpointGroup(currency.NewEndpointGroup),

			account.NewUserGetPreferredCurrencyEndpoint,
			account.NewUpdatePreferredCurrencyEndpoint,
			account.NewDeleteEndpoint,
			account.NewGetQuotaUsageEndpoint,
			ginfx.AsEndpointGroup(account.NewEndpointGroup),

			ginfx.AsEndpoint(NewHealthCheckLiveEndpoint),
			ginfx.AsEndpoint(NewVersionEndpoint),
		),
	)
}
