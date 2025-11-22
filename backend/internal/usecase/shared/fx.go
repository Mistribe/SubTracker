package shared

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("shared",
		fx.Provide(
			NewOwnerFactory,
		),
	)
}
