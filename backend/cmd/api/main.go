package main

import (
	"github.com/Oleexo/config-go/dotenv"
	"github.com/Oleexo/config-go/envs"
	configfx "github.com/Oleexo/config-go/fx"
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/adapters/auth/kinde"
	"github.com/mistribe/subtracker/internal/adapters/cache"
	"github.com/mistribe/subtracker/internal/adapters/exchange"
	"github.com/mistribe/subtracker/internal/adapters/http/router"
	logfx2 "github.com/mistribe/subtracker/internal/adapters/logfx"
	"github.com/mistribe/subtracker/internal/adapters/persistence"
	"github.com/mistribe/subtracker/internal/platform/startup"
	"github.com/mistribe/subtracker/internal/platform/startup/updater"
	"github.com/mistribe/subtracker/internal/usecase"
)

func main() {
	opts := []fx.Option{
		configfx.BuildConfigModule(
			dotenv.WithDotenv(),
			envs.WithEnvironmentVariables(),
		),
		fx.WithLogger(logfx2.NewFxLogger),
		logfx2.BuildLoggerModule(),
		persistence.BuildPersistenceModule(),
		router.BuildRoutesModule(),
		router.BuildHttpServerModule(),
		startup.BuildStartupModule(),
		updater.NewUpdaterModule(),
		cache.FxModule(),
		fx.Provide(
			kinde.NewTokenGenerator,
			exchange.New,
		),
	}
	opts = append(opts, usecase.BuildApplicationModules()...)
	app := fx.New(opts...)

	app.Run()
}
