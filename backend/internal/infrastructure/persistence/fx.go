package persistence

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/internal/infrastructure/startup"
)

func AsRepository[TRepository any](f any) any {
	return fx.Annotate(f,
		fx.As(new(TRepository)),
	)
}

func BuildPersistenceModule() fx.Option {
	return fx.Module("persistence",
		fx.Provide(
			NewRepository,
			AsRepository[subscription.Repository](NewSubscriptionRepository),
			AsRepository[family.Repository](NewFamilyRepository),
			AsRepository[label.Repository](NewLabelRepository),
			startup.AsStartupTask(newLabelTask),
			startup.AsStartupTask(newRepositoryTask),
		),
	)
}
