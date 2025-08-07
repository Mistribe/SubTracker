package query

import (
	"context"

	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/lang"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FindPreferredCurrencyQuery struct {
}

func NewFindPreferredCurrencyQuery() FindPreferredCurrencyQuery {
	return FindPreferredCurrencyQuery{}
}

type FindPreferredCurrencyQueryHandler struct {
	userRepository user.Repository
}

func NewFindPreferredCurrencyQueryHandler(userRepository user.Repository) *FindPreferredCurrencyQueryHandler {
	return &FindPreferredCurrencyQueryHandler{
		userRepository: userRepository,
	}
}

func (h FindPreferredCurrencyQueryHandler) Handle(
	ctx context.Context,
	_ FindPreferredCurrencyQuery) result.Result[currency.Unit] {
	userId, ok := auth.GetUserIdFromContext(ctx)
	if !ok {
		return result.Fail[currency.Unit](auth.ErrNotAuthorized)
	}
	profile, err := h.userRepository.GetUserProfile(ctx, userId)
	if err != nil {
		return result.Fail[currency.Unit](err)
	}

	if profile == nil {
		info := lang.FromContext(ctx)
		return result.Success(info.MostPreferred().PreferredCurrency())
	}

	return result.Success(profile.Currency())
}
