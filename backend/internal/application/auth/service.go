package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
)

type AuthenticationService struct {
	authenticationRepository auth.Repository
}

func (s AuthenticationService) MustGetUserId(ctx context.Context) string {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		panic("missing user id from context")
	}
	return userId
}

func (s AuthenticationService) MustGetFamilies(ctx context.Context) []uuid.UUID {
	userId := s.MustGetUserId(ctx)
	// todo cache
	families, err := s.authenticationRepository.GetUserFamilies(ctx, userId)
	if err != nil {
		panic(err)
	}
	return families
}
