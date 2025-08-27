package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/usecase/auth"
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
