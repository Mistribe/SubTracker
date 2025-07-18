package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/user"
)

func ensureOwnerIsEditor(ctx context.Context, ownerId string) error {
	userId, ok := user.FromContext(ctx)
	if !ok {
		return family.ErrFamilyNotFound
	}
	if ownerId != userId {
		return family.ErrOnlyOwnerCanEditFamily
	}
	return nil
}
