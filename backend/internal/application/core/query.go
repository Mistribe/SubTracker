package core

import (
	"context"

	"go.uber.org/fx"

	"github.com/mistribe/subtracker/pkg/langext/result"
)

type Query interface{}

type QueryHandler[TQuery Query, TResult any] interface {
	Handle(context.Context, TQuery) result.Result[TResult]
}

func AsQueryHandler[TQuery Query, TResult any](f any) any {
	return fx.Annotate(f,
		fx.As(new(QueryHandler[TQuery, TResult])),
	)
}
