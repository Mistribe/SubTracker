package startup

import (
	"context"
)

type Task interface {
	OnStart(context.Context) error
	OnStop(context.Context) error
}
