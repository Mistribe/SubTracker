package user

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/lang"
	"github.com/mistribe/subtracker/internal/ports"
)

type service struct {
	userRepository ports.UserRepository
}

func (s service) GetPreferredCurrency(ctx context.Context, userId string) currency.Unit {
	profile, err := s.userRepository.GetUserProfile(ctx, userId)
	if err != nil {
		// todo better err
		panic(err)
	}

	if profile == nil {
		info := lang.FromContext(ctx)
		return info.MostPreferred().PreferredCurrency()
	}

	return profile.Currency()

}

func NewService(userRepository ports.UserRepository) ports.UserService {
	return service{
		userRepository: userRepository,
	}
}
