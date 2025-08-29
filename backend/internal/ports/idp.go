package ports

import (
	"context"
)

type IdentityProvider interface {
	ReadSessionToken(ctx context.Context, sessionToken string) (Identity, error)
	DeleteUser(ctx context.Context, userId string) error
}

type Identity struct {
	IsValid bool
	Id      string
}

func NewInvalidIdentity() Identity {
	return Identity{
		IsValid: false,
	}
}
