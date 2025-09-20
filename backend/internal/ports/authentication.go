package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/user"
)

type Authentication interface {
	MustGetUserId(ctx context.Context) string
	MustGetUserRole(ctx context.Context) user.Role
}
