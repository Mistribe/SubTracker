package authentication

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/account"
)

const (
	ContextConnectedAccountKey string = "SUBTRACKER_CONNECTED_ACCOUNT"
)

func GetAccountFromContext(ctx context.Context) (account.ConnectedAccount, bool) {
	connectedAccount, ok := ctx.Value(ContextConnectedAccountKey).(account.ConnectedAccount)
	return connectedAccount, ok
}
