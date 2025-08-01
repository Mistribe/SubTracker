package system

import (
	"context"
	"sort"

	"github.com/Oleexo/config-go"
	"go.uber.org/fx"
)

const (
	highPriorty    = 0
	mediumPriority = 100
	lowPriority    = 200
)

type Updater interface {
	Priority() int
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
				return runServices(ctx, params.Services)
			},
		})
	}
	return &UpdaterService{
		services: params.Services,
	}
}

func runServices(ctx context.Context, services []Updater) error {
	sort.Slice(services, func(i, j int) bool {
		return services[i].Priority() < services[j].Priority()
	})

	for _, service := range services {
		if err := service.Update(ctx); err != nil {
			return err
		}
	}
	return nil
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
			AsUpdater(newProviderUpdater),
			AsUpdater(newLabelUpdater),
			NewUpdaterService,
		),
		fx.Invoke(func(s *UpdaterService) {}),
	)
}
