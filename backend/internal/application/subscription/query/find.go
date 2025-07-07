package query

import (
    "context"
    "github.com/google/uuid"
    "github.com/oleexo/subtracker/internal/application/core/result"
    "github.com/oleexo/subtracker/internal/domain/subscription"
)

type FindQuery struct {
    id uuid.UUID
}

func NewFindQuery(id uuid.UUID) FindQuery {
    return FindQuery{id: id}
}

type FindQueryHandler struct {
    respository subscription.Repository
}

func NewFindQueryHandler(respository subscription.Repository) *FindQueryHandler {
    return &FindQueryHandler{respository: respository}
}

func (h FindQueryHandler) Handle(ctx context.Context, query FindQuery) result.Result[subscription.Subscription] {

}
