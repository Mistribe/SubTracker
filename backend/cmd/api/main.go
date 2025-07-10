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

// @title           SubTracker API
// @version         1.0
// @description     This api provide HTTPRest endpoints for the application SubTracker.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

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
