package user

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("user not found")
)

const ContextKey string = "SUBTRACKER_USER_ID"

func FromContext(ctx context.Context) (string, bool) {
	userId, ok := ctx.Value(ContextKey).(string)
	return userId, ok
}
