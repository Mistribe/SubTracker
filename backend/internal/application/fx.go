package application

import (
	subQuery "github.com/oleexo/subtracker/internal/application/subscription/query"
	"go.uber.org/fx"
)

func BuildApplicationModule() fx.Option {
	return fx.Module("application",
		fx.Provide(
			subQuery.NewFindQueryHandler,
		),
	)
}
