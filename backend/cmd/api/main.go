package main

import (
	"github.com/Oleexo/config-go/dotenv"
	"github.com/Oleexo/config-go/envs"
	configfx "github.com/Oleexo/config-go/fx"
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/application"
	"github.com/oleexo/subtracker/internal/infrastructure/logfx"
	"github.com/oleexo/subtracker/internal/infrastructure/persistence"
)

// @title			SubTracker API
// @version		1.0
// @description	This api provide HTTPRest endpoints for the application SubTracker.
// @termsOfService	http://subtracker.mistribe.com/terms/
// @contact.name	API Support
// @contact.url	http://subtracker.mistribe.com/support
// @contact.email	support@mistribe.com
// @license.name	Apache 2.0
// @license.url	http://www.apache.org/licenses/LICENSE-2.0.html
// @host			api.subtracker.mistribe.com
// @BasePath		/
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
		BuildRoutesModule(),
		BuildHttpServerModule(),
	)

	app.Run()
}
