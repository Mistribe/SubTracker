package authentication

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("authentication",
		fx.Provide(
			NewClerkIdentityProvider,
			NewAuthentication,
		),
	)
}
