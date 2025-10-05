package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/account"
)

type Authentication interface {
	MustGetConnectedAccount(ctx context.Context) account.ConnectedAccount
}
