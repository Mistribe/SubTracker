package persistence

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/adapters/persistence/db"
	"github.com/mistribe/subtracker/internal/adapters/persistence/repositories"
)

func AsRepository[TRepository any](f any) any {
	return fx.Annotate(f,
		fx.As(new(TRepository)),
	)
}

func BuildPersistenceModule() fx.Option {
	return fx.Module("persistence",
		fx.Provide(
			db.NewContext,
			repositories.NewSubscriptionRepository,
			repositories.NewFamilyRepository,
			repositories.NewLabelRepository,
			repositories.NewProviderRepository,
			repositories.NewAccountRepository,
			repositories.NewCurrencyRateRepository,
		),
	)
}
