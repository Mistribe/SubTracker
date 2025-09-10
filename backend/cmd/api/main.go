package main

import (
	"github.com/Oleexo/config-go/dotenv"
	"github.com/Oleexo/config-go/envs"
	configfx "github.com/Oleexo/config-go/fx"
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/adapters/authentication"
	"github.com/mistribe/subtracker/internal/adapters/authorization"
	"github.com/mistribe/subtracker/internal/adapters/cache"
	"github.com/mistribe/subtracker/internal/adapters/exchange"
	"github.com/mistribe/subtracker/internal/adapters/http/router"
	logfx2 "github.com/mistribe/subtracker/internal/adapters/logfx"
	"github.com/mistribe/subtracker/internal/adapters/persistence"
	"github.com/mistribe/subtracker/internal/platform/startup"
	"github.com/mistribe/subtracker/internal/platform/startup/updater"
	"github.com/mistribe/subtracker/internal/usecase"
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
		fx.WithLogger(logfx2.NewFxLogger),
		logfx2.BuildLoggerModule(),
		persistence.BuildPersistenceModule(),
		router.BuildRoutesModule(),
		router.BuildHttpServerModule(),
		startup.BuildStartupModule(),
		updater.NewUpdaterModule(),
		cache.FxModule(),
		authentication.Module(),
		authorization.Module(),
		fx.Provide(

			exchange.New,
		),
	}
	opts = append(opts, usecase.BuildApplicationModules()...)
	app := fx.New(opts...)

	app.Run()
}
