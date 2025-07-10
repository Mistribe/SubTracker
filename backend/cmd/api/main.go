package main

import (
	"github.com/Oleexo/config-go/dotenv"
	"github.com/Oleexo/config-go/envs"
	configfx "github.com/Oleexo/config-go/fx"
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/application"
	"github.com/oleexo/subtracker/internal/infrastructure/api"
	"github.com/oleexo/subtracker/internal/infrastructure/logfx"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence"
)

func main() {
	app := fx.New(
		configfx.BuildConfigModule(
			dotenv.WithDotenv(),
			envs.WithEnvironmentVariables(),
		),
		fx.WithLogger(logfx.NewFxLogger),
		logfx.BuildLoggerModule(),
		application.BuildApplicationModule(),
		persistence.BuildPersistenceModule(),
		api.BuildRoutesModule(),
		api.BuildHttpServerModule(),
	)

	app.Run()
}
