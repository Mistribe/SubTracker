package core

import (
	"context"
	"github.com/oleexo/subtracker/internal/application/core/result"
)

type Command interface{}

type CommandHandler[C Command, R any] interface {
	Handle(context.Context, C) result.Result[R]
}
