package core

import (
	"context"

	"go.uber.org/fx"

	"github.com/oleexo/subtracker/pkg/langext/result"
)

type Command interface{}

type CommandHandler[TCommand Command, TResult any] interface {
	Handle(context.Context, TCommand) result.Result[TResult]
}

func AsCommandHandler[TCommand Command, TResult any](f any) any {
	return fx.Annotate(f,
		fx.As(new(CommandHandler[TCommand, TResult])),
	)
}
