package export

import (
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module("export",
		fx.Provide(
			NewExportService,
			NewLabelResolver,
			NewProviderResolver,
		),
	)
}
