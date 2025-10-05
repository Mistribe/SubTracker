package cache

import (
	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return fx.Module("cache",
		fx.Provide(
			New,
			NewRequest,
			NewLocal,
			NewDistributed,
		),
	)
}
