package persistence

import (
	"go.uber.org/fx"
)

func AsRepository[TRepository any](f any) any {
	return fx.Annotate(f,
		fx.As(new(TRepository)),
	)
}

func BuildPersistenceModule() fx.Option {
	return fx.Module("persistence",
		fx.Provide(
			NewDatabaseContext,
			NewSubscriptionRepository,
			NewFamilyRepository,
			NewLabelRepository,
			NewProviderRepository,
			NewAuthenticationRepository,
			//startup.AsStartupTask(newLabelTask),
			//startup.AsStartupTask(newRepositoryTask),
		),
	)
}
