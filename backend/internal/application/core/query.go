package core

import (
	"context"

	"github.com/oleexo/subtracker/pkg/langext/result"
)

type Query interface{}

type QueryHandler[TQuery Query, TResult any] interface {
	Handle(context.Context, TQuery) result.Result[TResult]
}
