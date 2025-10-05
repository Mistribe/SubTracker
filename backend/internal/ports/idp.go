package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/types"
)

type IdentityProvider interface {
	ReadSessionToken(ctx context.Context, sessionToken string) (Identity, error)
	DeleteUser(ctx context.Context, userId types.UserID) error
}

type Identity struct {
	IsValid bool
	Id      string
	Role    string
	Plan    string
}

func NewInvalidIdentity() Identity {
	return Identity{
		IsValid: false,
	}
}
