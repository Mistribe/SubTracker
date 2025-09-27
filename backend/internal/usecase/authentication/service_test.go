package authentication_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
)

func TestMustGetUserId(t *testing.T) {
	userId := "RandomUserId1234567890"
	ctx := context.WithValue(context.Background(), authentication.ContextUserIdKey, userId)

	service := authentication.NewAuthentication(ports.NewMockUserRepository(t), ports.NewMockFamilyRepository(t))

	assert.Equal(t, userId, service.MustGetUserId(ctx))
}
