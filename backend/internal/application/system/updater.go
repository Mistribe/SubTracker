package system

import (
	"context"

	"github.com/Oleexo/config-go"
	"go.uber.org/fx"
)

type Updater interface {
	Update(ctx context.Context) error
}

type UpdaterService struct {
	services []Updater
}

type UpdaterParams struct {
	fx.In

	Services  []Updater `group:"updater_services"`
	Config    config.Configuration
	Lifecycle fx.Lifecycle
}

func NewUpdaterService(params UpdaterParams) *UpdaterService {
	runAtStart := params.Config.GetBoolOrDefault("UPDATER_AT_START", false)
	if runAtStart {
		params.Lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				for _, service := range params.Services {
					if err := service.Update(ctx); err != nil {
						return err
					}
				}
				return nil
			},
		})
	}
	return &UpdaterService{
		services: params.Services,
	}
}

func AsUpdater(f any) any {
	return fx.Annotate(f,
		fx.As(new(Updater)),
		fx.ResultTags(`group:"updater_services"`),
	)
}

func NewUpdaterModule() fx.Option {
	return fx.Module("updater",
		fx.Provide(
			AsUpdater(newLabelUpdater),
			NewUpdaterService,
		),
		fx.Invoke(func(s *UpdaterService) {}),
	)
}
