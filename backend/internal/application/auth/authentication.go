package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type AuthenticationService struct {
	userRepository user.Repository
}

func NewAuthenticationService(userRepository user.Repository) auth.Service {
	return &AuthenticationService{
		userRepository: userRepository,
	}
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
	families, err := s.userRepository.GetUserFamilies(ctx, userId)
	if err != nil {
		panic(err)
	}
	return families
}

func (s AuthenticationService) IsInFamily(ctx context.Context, familyId uuid.UUID) bool {
	families := s.MustGetFamilies(ctx)

	for _, fam := range families {
		if fam == familyId {
			return true
		}
	}

	return false
}
