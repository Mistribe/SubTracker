package core

import (
	"context"

	"github.com/oleexo/subtracker/internal/application/core/result"
)

type Query interface{}

type QueryHandler[Q Query, R any] interface {
	Handle(context.Context, Q) result.Result[R]
}
