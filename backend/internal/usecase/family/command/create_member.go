package command

import (
	"context"
	"errors"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateFamilyMemberCommand struct {
	FamilyID       types.FamilyID
	FamilyMemberID option.Option[types.FamilyMemberID]
	Name           string
	Type           family.MemberType
	CreatedAt      option.Option[time.Time]
}

type CreateFamilyMemberCommandHandler struct {
	familyRepository    ports.FamilyRepository
	authorization       ports.Authorization
	entitlementResolver ports.EntitlementResolver
}

func NewCreateFamilyMemberCommandHandler(
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization,
	entitlementResolver ports.EntitlementResolver) *CreateFamilyMemberCommandHandler {
	return &CreateFamilyMemberCommandHandler{
		familyRepository:    familyRepository,
		authorization:       authorization,
		entitlementResolver: entitlementResolver,
	}
}

func (h CreateFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command CreateFamilyMemberCommand) result.Result[family.Family] {
	fam, err := h.familyRepository.GetById(ctx, command.FamilyID)
	if err != nil {
		return result.Fail[family.Family](err)
	}
	if fam == nil {
		return result.Fail[family.Family](family.ErrFamilyNotFound)

	}

	if err = h.authorization.Can(ctx, authorization.PermissionWrite).For(fam); err != nil {
		return result.Fail[family.Family](err)
	}

	ok, _, err := h.entitlementResolver.CheckQuota(ctx, billing.FeatureIdFamilyMembersCount, 1)
	if err != nil {
		return result.Fail[family.Family](err)
	}
	if !ok {
		return result.Fail[family.Family](family.ErrFamilyMembersLimitReached)
	}

	return h.addFamilyMemberToFamily(ctx, command, fam)
}

func (h CreateFamilyMemberCommandHandler) addFamilyMemberToFamily(
	ctx context.Context,
	command CreateFamilyMemberCommand, fam family.Family) result.Result[family.Family] {

	var familyMemberID types.FamilyMemberID
	if command.FamilyMemberID != nil && command.FamilyMemberID.IsSome() {
		familyMemberID = *command.FamilyMemberID.Value()
	} else {
		familyMemberID = types.NewFamilyMemberID()
	}

	var createdAt time.Time
	if command.CreatedAt != nil && command.CreatedAt.IsSome() {
		createdAt = *command.CreatedAt.Value()
	} else {
		createdAt = time.Now()
	}
	familyMember := family.NewMember(
		familyMemberID,
		command.FamilyID,
		command.Name,
		command.Type,
		nil,
		createdAt,
		createdAt,
	)

	if err := fam.AddMember(familyMember); err != nil {
		// If duplicate or other add error we return success without persisting (idempotent behaviour)
		if errors.Is(err, family.ErrDuplicateMember) {
			return result.Success(fam)
		}
		return result.Fail[family.Family](err)
	}

	if err := fam.GetValidationErrors(); err != nil {
		return result.Fail[family.Family](err)
	}

	if err := h.familyRepository.Save(ctx, fam); err != nil {
		return result.Fail[family.Family](err)
	}

	return result.Success(fam)

}
