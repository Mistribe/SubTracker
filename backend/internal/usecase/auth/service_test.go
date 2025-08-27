package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	domainauth "github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
)

func TestMustGetUserId(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)

	service := auth.NewAuthenticationService(ports.NewMockUserRepository(t), ports.NewMockFamilyRepository(t))

	assert.Equal(t, userId, service.MustGetUserId(ctx))
}

func TestMustGetFamilies(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)
	expectedFamilies := []uuid.UUID{uuid.New(), uuid.New()}

	mockRepo := ports.NewMockUserRepository(t)
	mockRepo.On("GetUserFamilies", ctx, userId).Return(expectedFamilies, nil)

	service := auth.NewAuthenticationService(mockRepo, ports.NewMockFamilyRepository(t))

	assert.Equal(t, expectedFamilies, service.MustGetFamilies(ctx))
}

func TestMustGetFamilies_Error(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)
	expectedError := errors.New("repository error")

	mockRepo := ports.NewMockUserRepository(t)
	mockRepo.On("GetUserFamilies", ctx, userId).Return(nil, expectedError)

	service := auth.NewAuthenticationService(mockRepo, ports.NewMockFamilyRepository(t))

	assert.Panics(t, func() {
		service.MustGetFamilies(ctx)
	})
}

func TestIsInFamily(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)
	families := []uuid.UUID{uuid.New(), uuid.New()}

	mockRepo := ports.NewMockUserRepository(t)
	mockRepo.On("GetUserFamilies", ctx, userId).Return(families, nil)

	service := auth.NewAuthenticationService(mockRepo, ports.NewMockFamilyRepository(t))

	assert.True(t, service.IsInFamily(ctx, families[0]))
	assert.False(t, service.IsInFamily(ctx, uuid.New()))
}

func TestIsOwner_SystemOwner(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)

	service := auth.NewAuthenticationService(ports.NewMockUserRepository(t), ports.NewMockFamilyRepository(t))

	isOwner, err := service.IsOwner(ctx, domainauth.NewSystemOwner())
	assert.NoError(t, err)
	assert.False(t, isOwner)
}

func TestIsOwner_FamilyOwner_IsMember(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)
	familyId := uuid.New()

	mockRepo := ports.NewMockUserRepository(t)
	mockFamilyRepo := ports.NewMockFamilyRepository(t)
	mockFamilyRepo.On("IsUserMemberOfFamily", ctx, familyId, userId).Return(true, nil)

	service := auth.NewAuthenticationService(mockRepo, mockFamilyRepo)

	isOwner, err := service.IsOwner(ctx, domainauth.NewFamilyOwner(familyId))
	assert.NoError(t, err)
	assert.True(t, isOwner)
}

func TestIsOwner_FamilyOwner_NotMember(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)
	familyId := uuid.New()

	mockRepo := ports.NewMockUserRepository(t)
	mockFamilyRepo := ports.NewMockFamilyRepository(t)
	mockFamilyRepo.On("IsUserMemberOfFamily", ctx, familyId, userId).Return(false, nil)

	service := auth.NewAuthenticationService(mockRepo, mockFamilyRepo)

	isOwner, err := service.IsOwner(ctx, domainauth.NewFamilyOwner(familyId))
	assert.ErrorIs(t, err, domainauth.ErrNotAuthorized)
	assert.False(t, isOwner)
}

func TestIsOwner_PersonalOwner_Match(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)

	service := auth.NewAuthenticationService(ports.NewMockUserRepository(t), ports.NewMockFamilyRepository(t))

	isOwner, err := service.IsOwner(ctx, domainauth.NewPersonalOwner(userId))
	assert.NoError(t, err)
	assert.True(t, isOwner)
}

func TestIsOwner_PersonalOwner_NoMatch(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth.ContextKey, userId)

	service := auth.NewAuthenticationService(ports.NewMockUserRepository(t), ports.NewMockFamilyRepository(t))

	isOwner, err := service.IsOwner(ctx, domainauth.NewPersonalOwner("DifferentUserId"))
	assert.ErrorIs(t, err, domainauth.ErrNotAuthorized)
	assert.False(t, isOwner)
}
