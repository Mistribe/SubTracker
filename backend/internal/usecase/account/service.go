package account

import (
	"context"
	"log/slog"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared/i18n"
)

type service struct {
	userRepository ports.AccountRepository
	logger         *slog.Logger
}

func (s service) GetPreferredCurrency(ctx context.Context, userId types.UserID) currency.Unit {
	profile, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		s.logger.Error("cannot retrieve user profile", "err", err, "user_id", userId, "ctx", ctx)
	}

	if profile == nil {
		info := i18n.FromContext(ctx)
		return info.MostPreferred().PreferredCurrency()
	}

	return profile.Currency().ValueOrDefault(currency.USD)

}

func NewService(
	userRepository ports.AccountRepository,
	logger *slog.Logger) ports.AccountService {
	return service{
		userRepository: userRepository,
		logger:         logger,
	}
}
