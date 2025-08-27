package auth_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	domainauth "github.com/mistribe/subtracker/internal/domain/auth"
	auth2 "github.com/mistribe/subtracker/internal/usecase/auth"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/user"
)

func TestMustGetUserId(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)

	service := auth2.NewAuthenticationService(user.NewMockRepository(t), family.NewMockRepository(t))

	assert.Equal(t, userId, service.MustGetUserId(ctx))
}

func TestMustGetFamilies(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)
	expectedFamilies := []uuid.UUID{uuid.New(), uuid.New()}

	mockRepo := user.NewMockRepository(t)
	mockRepo.On("GetUserFamilies", ctx, userId).Return(expectedFamilies, nil)

	service := auth2.NewAuthenticationService(mockRepo, family.NewMockRepository(t))

	assert.Equal(t, expectedFamilies, service.MustGetFamilies(ctx))
}

func TestMustGetFamilies_Error(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)
	expectedError := errors.New("repository error")

	mockRepo := user.NewMockRepository(t)
	mockRepo.On("GetUserFamilies", ctx, userId).Return(nil, expectedError)

	service := auth2.NewAuthenticationService(mockRepo, family.NewMockRepository(t))

	assert.Panics(t, func() {
		service.MustGetFamilies(ctx)
	})
}

func TestIsInFamily(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)
	families := []uuid.UUID{uuid.New(), uuid.New()}

	mockRepo := user.NewMockRepository(t)
	mockRepo.On("GetUserFamilies", ctx, userId).Return(families, nil)

	service := auth2.NewAuthenticationService(mockRepo, family.NewMockRepository(t))

	assert.True(t, service.IsInFamily(ctx, families[0]))
	assert.False(t, service.IsInFamily(ctx, uuid.New()))
}

func TestIsOwner_SystemOwner(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)

	service := auth2.NewAuthenticationService(user.NewMockRepository(t), family.NewMockRepository(t))

	isOwner, err := service.IsOwner(ctx, domainauth.NewSystemOwner())
	assert.NoError(t, err)
	assert.False(t, isOwner)
}

func TestIsOwner_FamilyOwner_IsMember(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)
	familyId := uuid.New()

	mockRepo := user.NewMockRepository(t)
	mockFamilyRepo := family.NewMockRepository(t)
	mockFamilyRepo.On("IsUserMemberOfFamily", ctx, familyId, userId).Return(true, nil)

	service := auth2.NewAuthenticationService(mockRepo, mockFamilyRepo)

	isOwner, err := service.IsOwner(ctx, domainauth.NewFamilyOwner(familyId))
	assert.NoError(t, err)
	assert.True(t, isOwner)
}

func TestIsOwner_FamilyOwner_NotMember(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)
	familyId := uuid.New()

	mockRepo := user.NewMockRepository(t)
	mockFamilyRepo := family.NewMockRepository(t)
	mockFamilyRepo.On("IsUserMemberOfFamily", ctx, familyId, userId).Return(false, nil)

	service := auth2.NewAuthenticationService(mockRepo, mockFamilyRepo)

	isOwner, err := service.IsOwner(ctx, domainauth.NewFamilyOwner(familyId))
	assert.ErrorIs(t, err, domainauth.ErrNotAuthorized)
	assert.False(t, isOwner)
}

func TestIsOwner_PersonalOwner_Match(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)

	service := auth2.NewAuthenticationService(user.NewMockRepository(t), family.NewMockRepository(t))

	isOwner, err := service.IsOwner(ctx, domainauth.NewPersonalOwner(userId))
	assert.NoError(t, err)
	assert.True(t, isOwner)
}

func TestIsOwner_PersonalOwner_NoMatch(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), auth2.ContextKey, userId)

	service := auth2.NewAuthenticationService(user.NewMockRepository(t), family.NewMockRepository(t))

	isOwner, err := service.IsOwner(ctx, domainauth.NewPersonalOwner("DifferentUserId"))
	assert.ErrorIs(t, err, domainauth.ErrNotAuthorized)
	assert.False(t, isOwner)
}
