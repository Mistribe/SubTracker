package authorization

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	user "github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
)

type service struct {
	userRepository   ports.AccountRepository
	familyRepository ports.FamilyRepository
	cache            ports.Cache
	authentication   ports.Authentication
}

func (s service) GetCurrentLimits(ctx context.Context, category user.Category) (shared.Limits, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) GetCurrentLimit(ctx context.Context, feature user.Feature) (shared.Limit, error) {
	userId := s.authentication.MustGetUserId(ctx)
	currentUser, err := s.getUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	if !currentUser.Plan().HasFeature(feature) {
		return nil, user.NewLimitExceededErr(feature, 0, 0)
	}

	info := currentUser.Plan().GetFeatureDetail(feature)
	if info == nil {
		return nil, user.NewLimitExceededErr(feature, 0, 0)
	}
	currentCount := currentUser.Stats().GetCountFromFeature(feature)

	return newLimitInfo(currentCount, info), nil
}

func (s service) Can(ctx context.Context, permission user.Permission) ports.PermissionRequest {
	userId := s.authentication.MustGetUserId(ctx)
	userRole := s.authentication.MustGetUserRole(ctx)
	familyId, err := s.getFamilyId(ctx, userId)
	return &permissionRequest{
		error:      err,
		userRole:   userRole,
		userId:     userId,
		permission: permission,
		userFamily: familyId,
	}
}

func (s service) getFamilyId(ctx context.Context, userId string) (*uuid.UUID, error) {
	familyUserCacheKey := fmt.Sprintf("family-user-%s", userId)
	fromCache, ok := s.cache.From(ctx, ports.CacheLevelRequest).Get(familyUserCacheKey).(uuid.UUID)
	if ok {
		return &fromCache, nil
	}
	fromDatabase, err := s.familyRepository.GetAccountFamily(ctx, userId)
	if err != nil {
		return nil, err
	}
	if fromDatabase != nil {
		familyId := fromDatabase.Id()
		s.cache.From(ctx, ports.CacheLevelRequest).Set(familyUserCacheKey, familyId)
		return &familyId, nil
	}
	return nil, nil
}

func (s service) getUser(ctx context.Context, userId string) (user.User, error) {
	userCacheKey := fmt.Sprintf("user-%s", userId)
	fromCache, ok := s.cache.From(ctx, ports.CacheLevelRequest).Get(userCacheKey).(user.User)
	if ok {
		return fromCache, nil
	}
	fromDatabase, err := s.userRepository.GetById(ctx, userId)
	if err != nil {
		return nil, err
	}
	if fromDatabase != nil {
		s.cache.From(ctx, ports.CacheLevelRequest).Set(userCacheKey, fromDatabase)
		return fromDatabase, nil
	}
	return nil, nil
}

func (s service) EnsureWithinLimit(ctx context.Context, feature user.Feature, delta int) error {
	userId := s.authentication.MustGetUserId(ctx)
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

func New(
	userRepository ports.AccountRepository,
	cache ports.Cache) ports.Authorization {
	return &service{
		userRepository: userRepository,
		cache:          cache,
	}
}
