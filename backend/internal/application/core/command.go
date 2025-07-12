package core

import (
	"context"

	"github.com/oleexo/subtracker/pkg/langext/result"
)

type Command interface{}

type CommandHandler[TCommand Command, TResult any] interface {
	Handle(context.Context, TCommand) result.Result[TResult]
}
