package ports

import (
	"context"

	"go.uber.org/fx"

	"github.com/mistribe/subtracker/pkg/langext/result"
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

type Query interface{}

type QueryHandler[TQuery Query, TResult any] interface {
	Handle(context.Context, TQuery) result.Result[TResult]
}

func AsQueryHandler[TQuery Query, TResult any](f any) any {
	return fx.Annotate(f,
		fx.As(new(QueryHandler[TQuery, TResult])),
	)
}
