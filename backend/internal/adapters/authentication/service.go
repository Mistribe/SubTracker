package authentication

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/account"
	"github.com/mistribe/subtracker/internal/ports"
)

type authentication struct {
}

func NewAuthentication() ports.Authentication {
	return &authentication{}
}

func (s authentication) MustGetConnectedAccount(ctx context.Context) account.ConnectedAccount {
	acc, ok := GetAccountFromContext(ctx)
	if !ok {
		panic("missing account from context")
	}
	return acc
}
