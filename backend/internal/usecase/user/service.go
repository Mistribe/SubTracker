package user

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared/i18n"
)

type service struct {
	userRepository ports.AccountRepository
}

func (s service) GetPreferredCurrency(ctx context.Context, userId string) currency.Unit {
	profile, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		// todo better err
		panic(err)
	}

	if profile == nil {
		info := i18n.FromContext(ctx)
		return info.MostPreferred().PreferredCurrency()
	}

	return profile.Currency()

}

func NewService(userRepository ports.AccountRepository) ports.UserService {
	return service{
		userRepository: userRepository,
	}
}
