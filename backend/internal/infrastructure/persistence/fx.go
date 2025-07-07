package persistence

import "go.uber.org/fx"

func BuildPersistenceModule() fx.Option {
	return fx.Module("persistence",
		fx.Provide(
			NewSubsriptionRepository,
		),
	)
}
