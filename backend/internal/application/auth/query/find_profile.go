package query

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/lang"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindProfileQuery struct {
}

func NewFindProfileQuery() FindProfileQuery {
	return FindProfileQuery{}
}

type FindProfileQueryHandler struct {
	authRepository auth.Repository
}

func NewFindProfileQueryHandler(authRepository auth.Repository) *FindProfileQueryHandler {
	return &FindProfileQueryHandler{
		authRepository: authRepository,
	}
}

func (h FindProfileQueryHandler) Handle(ctx context.Context, _ FindProfileQuery) result.Result[auth.UserProfile] {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return result.Fail[auth.UserProfile](auth.ErrNotAuthorized)
	}
	profile, err := h.authRepository.GetUserProfile(ctx, userId)
	if err != nil {
		return result.Fail[auth.UserProfile](err)
	}

	if profile == nil {
		info := lang.FromContext(ctx)
		profile = auth.NewUserProfile(userId, info.MostPreferred().PreferredCurrency())
	}

	return result.Success(profile)
}
