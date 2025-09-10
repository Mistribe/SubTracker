package authentication

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("auth_idp",
		fx.Provide(
			NewClerkIdentityProvider,
		),
	)
}
