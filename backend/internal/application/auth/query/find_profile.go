package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindProfileQuery struct {
}

type FindProfileQueryHandler struct {
	authRepository auth.Repository
}

func NewFindProfileQueryHandler(authRepository auth.Repository) *FindProfileQueryHandler {
	return &FindProfileQueryHandler{
		authRepository: authRepository,
	}
}

func (h FindProfileQueryHandler) Handler(ctx context.Context, _ FindProfileQuery) result.Result[auth.UserProfile] {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return result.Fail[auth.UserProfile](auth.ErrNotAuthorized)
	}
	profile, err := h.authRepository.GetUserProfile(ctx, userId)
	if err != nil {
		return result.Fail[auth.UserProfile](err)
	}

	if profile == nil {
		return result.Fail[auth.UserProfile](auth.ErrUnknownUser)
	}

	return result.Success(profile)
}
