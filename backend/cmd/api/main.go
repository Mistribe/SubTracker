package main

import (
	"github.com/Oleexo/config-go/dotenv"
	"github.com/Oleexo/config-go/envs"
	configfx "github.com/Oleexo/config-go/fx"
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/application"
	"github.com/mistribe/subtracker/internal/application/system"
	"github.com/mistribe/subtracker/internal/infrastructure/cache"
	"github.com/mistribe/subtracker/internal/infrastructure/exchange"
	"github.com/mistribe/subtracker/internal/infrastructure/kinde"
	"github.com/mistribe/subtracker/internal/infrastructure/logfx"
	"github.com/mistribe/subtracker/internal/infrastructure/persistence"
	"github.com/mistribe/subtracker/internal/infrastructure/startup"
)

// @title					SubTracker API
// @version				1.0
// @description			This api provide HTTPRest endpoints for the application SubTracker.
// @termsOfService			http://subtracker.mistribe.com/terms/
// @contact.name			API Support
// @contact.url			http://subtracker.mistribe.com/support
// @contact.email			support@mistribe.com
// @license.name			Apache 2.0
// @license.url			http://www.apache.org/licenses/LICENSE-2.0.html
// @servers.url			https://api.subtracker.mistribe.com
// @servers.description	Production server
// @servers.url			http://localhost:8080
// @servers.description	Development server
func main() {
	opts := []fx.Option{
		configfx.BuildConfigModule(
			dotenv.WithDotenv(),
			envs.WithEnvironmentVariables(),
		),
		fx.WithLogger(logfx.NewFxLogger),
		logfx.BuildLoggerModule(),
		persistence.BuildPersistenceModule(),
		BuildRoutesModule(),
		BuildHttpServerModule(),
		startup.BuildStartupModule(),
		system.NewUpdaterModule(),
		cache.FxModule(),
		fx.Provide(
			kinde.NewTokenGenerator,
			exchange.NewClient,
		),
	}
	opts = append(opts, application.BuildApplicationModules()...)
	app := fx.New(opts...)

	app.Run()
}
