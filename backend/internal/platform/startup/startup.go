package startup

import (
	"context"
)

type Task interface {
	Priority() int
	OnStart(context.Context) error
	OnStop(context.Context) error
}
