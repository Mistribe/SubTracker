package authorization

import (
	"context"
	"fmt"

	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
)

type service struct {
	userRepository ports.UserRepository
	cache          ports.Cache
}

func (s service) getUser(ctx context.Context, userId string) (user.User, error) {
	userCacheKey := fmt.Sprintf("user-%s", userId)
	fromCache, ok := s.cache.From(ctx, ports.CacheLevelRequest).Get(userCacheKey).(user.User)
	if ok {
		return fromCache, nil
	}
	fromDatabase, err := s.userRepository.GetUser(ctx, userId)
	if err != nil {
		return nil, err
	}
	if fromDatabase != nil {
		s.cache.From(ctx, ports.CacheLevelRequest).Set(userCacheKey, fromDatabase)
		return fromDatabase, nil
	} else {
		return nil, nil
	}
}

func (s service) EnsureWithinLimit(ctx context.Context, userId string, feature user.Feature, delta int) error {
	currentUser, err := s.getUser(ctx, userId)
	if err != nil {
		return err
	}

	if !currentUser.Plan().HasFeature(feature) {
		return user.NewLimitExceededErr(feature, 0, 0)
	}

	info := currentUser.Plan().GetFeatureDetail(feature)
	if info == nil {
		return user.NewLimitExceededErr(feature, 0, 0)
	}

	if !info.HasLimit() {
		return nil
	}

	currentFeatureCount := currentUser.Stats().GetCountFromFeature(feature)
	newFeatureCount := currentFeatureCount + int64(delta)
	if newFeatureCount > *info.Limit() {
		return user.NewLimitExceededErr(feature, *info.Limit(), currentFeatureCount)
	}

	return nil
}

func New(userRepository ports.UserRepository,
	cache ports.Cache) ports.Authorization {
	return &service{
		userRepository: userRepository,
		cache:          cache,
	}
}
