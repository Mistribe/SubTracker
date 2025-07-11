package core

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core/result"
)

type Command interface{}

type CommandHandler[TCommand Command, TResult any] interface {
	Handle(context.Context, TCommand) result.Result[TResult]
}
