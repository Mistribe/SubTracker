package auth

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("app_auth",
		fx.Provide(
			NewAuthenticationService,
		),
	)
}
