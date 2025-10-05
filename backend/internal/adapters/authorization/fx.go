package authorization

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("authorization",
		fx.Provide(
			New,
		),
	)
}
