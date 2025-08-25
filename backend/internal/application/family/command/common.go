package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/family"
)

func ensureOwnerIsEditor(ctx context.Context, ownerId string) error {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return family.ErrFamilyNotFound
	}
	if ownerId != userId {
		return family.ErrOnlyOwnerCanEditFamily
	}
	return nil
}
