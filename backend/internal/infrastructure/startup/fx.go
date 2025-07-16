package startup

import (
	"context"
	"sort"

	"go.uber.org/fx"
)

func AsStartupTask(f any) any {
	return fx.Annotate(f,
		fx.As(new(Task)),
		fx.ResultTags(`group:"startup_tasks"`),
	)
}

type taskParams struct {
	fx.In
	Lifecycle fx.Lifecycle
	Tasks     []Task `group:"startup_tasks"`
}

func invoke(params taskParams) {
	params.Lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			sort.Slice(params.Tasks, func(i, j int) bool {
				return params.Tasks[i].Priority() < params.Tasks[j].Priority()
			})
			for _, task := range params.Tasks {
				if err := task.OnStart(ctx); err != nil {
					return err
				}
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sort.Slice(params.Tasks, func(i, j int) bool {
				return params.Tasks[i].Priority() < params.Tasks[j].Priority()
			})
			for _, task := range params.Tasks {
				if err := task.OnStop(ctx); err != nil {
					return err
				}
			}
			return nil
		},
	})
}

func BuildStartupModule() fx.Option {
	return fx.Module("startup",
		fx.Invoke(invoke),
	)
}
