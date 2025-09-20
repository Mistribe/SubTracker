package auth

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
)

type authentication struct {
	userRepository   ports.UserRepository
	familyRepository ports.FamilyRepository
}

func NewAuthentication(
	userRepository ports.UserRepository,
	familyRepository ports.FamilyRepository) ports.Authentication {
	return &authentication{
		userRepository:   userRepository,
		familyRepository: familyRepository,
	}
}

func (s authentication) MustGetUserId(ctx context.Context) string {
	userId, ok := GetUserIdFromContext(ctx)
	if !ok {
		panic("missing user id from context")
	}
	return userId
}

func (s authentication) MustGetUserRole(ctx context.Context) user.Role {
	//TODO implement me
	panic("implement me")
}
