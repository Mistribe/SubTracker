package auth

import (
	"context"
)

const ContextKey string = "SUBTRACKER_USER_ID"

func GetUserIdFromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(ContextKey).(string)
	return userId, ok
}
