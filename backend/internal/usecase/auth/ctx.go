package auth

import (
	"context"
)

const (
	ContextUserIdKey   string = "SUBTRACKER_USER_ID"
	ContextIdentityKey string = "SUBTRACKER_IDENTITY"
)

func GetUserIdFromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(ContextUserIdKey).(string)
	return userId, ok
}
