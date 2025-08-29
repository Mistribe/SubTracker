package auth

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/ports"
)

type AuthenticationService struct {
	userRepository   ports.UserRepository
	familyRepository ports.FamilyRepository
}

func NewAuthenticationService(userRepository ports.UserRepository,
	familyRepository ports.FamilyRepository) ports.AuthService {
	return &AuthenticationService{
		userRepository:   userRepository,
		familyRepository: familyRepository,
	}
}

func (s AuthenticationService) MustGetUserId(ctx context.Context) string {
	userId, ok := GetUserIdFromContext(ctx)
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

func (s AuthenticationService) IsOwner(ctx context.Context,
	owner auth.Owner) (bool, error) {
	userId := s.MustGetUserId(ctx)
	switch owner.Type() {
	case auth.SystemOwnerType:
		return false, nil
	case auth.FamilyOwnerType:
		isMember, err := s.familyRepository.IsUserMemberOfFamily(ctx, owner.FamilyId(), userId)
		if err != nil {
			return false, err
		}
		if !isMember {
			return false, auth.ErrUnauthorized
		}
	case auth.PersonalOwnerType:
		if owner.UserId() != userId {
			return false, auth.ErrUnauthorized
		}
	}

	return true, nil
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
