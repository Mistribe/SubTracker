package authorization

import (
	"context"
	"fmt"

	"github.com/mistribe/subtracker/internal/domain/account"
	user "github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
)

type service struct {
	userRepository   ports.AccountRepository
	familyRepository ports.FamilyRepository
	cache            ports.Cache
	authentication   ports.Authentication
}

func (s service) Can(ctx context.Context, permission user.Permission) ports.PermissionRequest {
	connectedAccount := s.authentication.MustGetConnectedAccount(ctx)
	userID := connectedAccount.UserID()
	userRole := connectedAccount.Role()
	familyId, err := s.getFamilyId(ctx, userID)
	return &permissionRequest{
		error:      err,
		userRole:   userRole,
		userId:     userID,
		permission: permission,
		userFamily: familyId,
	}
}

func (s service) getFamilyId(ctx context.Context, userId types.UserID) (*types.FamilyID, error) {
	familyUserCacheKey := fmt.Sprintf("family-user-%s", userId)
	fromCache, ok := s.cache.From(ctx, ports.CacheLevelRequest).Get(familyUserCacheKey).(types.FamilyID)
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

func (s service) getUser(ctx context.Context, userId types.UserID) (account.Account, error) {
	userCacheKey := fmt.Sprintf("user-%s", userId)
	fromCache, ok := s.cache.From(ctx, ports.CacheLevelRequest).Get(userCacheKey).(account.Account)
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

func New(
	userRepository ports.AccountRepository,
	cache ports.Cache) ports.Authorization {
	return &service{
		userRepository: userRepository,
		cache:          cache,
	}
}
