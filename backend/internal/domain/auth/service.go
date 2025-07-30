package auth

import (
	"context"

	"github.com/google/uuid"
)

type Service interface {
	MustGetUserId(ctx context.Context) string
	MustGetFamilies(ctx context.Context) []uuid.UUID
}
