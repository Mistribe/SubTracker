package billing

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("billing",
		fx.Provide(
			NewEntitlementResolver,
		),
	)
}
