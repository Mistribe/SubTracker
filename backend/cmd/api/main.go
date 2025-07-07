package main

import (
	"github.com/Oleexo/config-go/dotenv"
	"github.com/Oleexo/config-go/envs"
	configfx "github.com/Oleexo/config-go/fx"
	"github.com/oleexo/subtracker/internal/infrastructure/api"
	"github.com/oleexo/subtracker/internal/infrastructure/logfx"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		configfx.BuildConfigModule(
			dotenv.WithDotenv(),
			envs.WithEnvironmentVariables(),
		),
		fx.WithLogger(logfx.NewFxLogger),
		logfx.BuildLoggerModule(),
		api.BuildRoutesModule(),
		api.BuildHttpServerModule(),
	)

	app.Run()
}
