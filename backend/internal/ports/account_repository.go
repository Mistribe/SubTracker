package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type AccountRepository interface {
	GetFamily(ctx context.Context, userId types.UserID) (*uuid.UUID, error)
	GetById(ctx context.Context, userId types.UserID) (account.Account, error)
	Save(ctx context.Context, profile account.Account) error
	CreateDefault(ctx context.Context, userId types.UserID) (account.Account, error)
}
