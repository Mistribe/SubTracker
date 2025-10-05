package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateFamilyCommand struct {
	FamilyId    option.Option[types.FamilyID]
	Name        string
	CreatorName string
	CreatedAt   option.Option[time.Time]
}

type CreateFamilyCommandHandler struct {
	familyRepository ports.FamilyRepository
	authentication   ports.Authentication
}

func NewCreateFamilyCommandHandler(
	familyRepository ports.FamilyRepository,
	authService ports.Authentication) *CreateFamilyCommandHandler {
	return &CreateFamilyCommandHandler{
		familyRepository: familyRepository,
		authentication:   authService,
	}
}

func (h CreateFamilyCommandHandler) Handle(ctx context.Context, cmd CreateFamilyCommand) result.Result[family.Family] {
	var familyID types.FamilyID
	// Guard against nil Option (zero-value interface) before calling its methods
	if cmd.FamilyId != nil && cmd.FamilyId.IsSome() { // provided family id option
		familyID = *cmd.FamilyId.Value()
		fam, err := h.familyRepository.GetById(ctx, familyID)
		if err != nil {
			return result.Fail[family.Family](err)
		}

		if fam != nil {
			return result.Fail[family.Family](family.ErrFamilyAlreadyExists)
		}
	} else { // no family id provided, generate new one
		familyID = types.NewFamilyID()
	}
	return h.createFamily(ctx, familyID, cmd)
}

func (h CreateFamilyCommandHandler) createFamily(
	ctx context.Context,
	familyID types.FamilyID,
	cmd CreateFamilyCommand) result.Result[family.Family] {

	connectedAccount := h.authentication.MustGetConnectedAccount(ctx)
	userID := connectedAccount.UserID()
	// todo check if the user already have a family
	memberId := types.NewFamilyMemberID()

	// Guard nil option before calling methods
	var createdAt time.Time
	if cmd.CreatedAt != nil && cmd.CreatedAt.IsSome() {
		createdAt = *cmd.CreatedAt.Value()
	} else {
		createdAt = time.Now()
	}

	fam := family.NewFamily(
		familyID,
		userID,
		cmd.Name,
		[]family.Member{},
		createdAt,
		createdAt,
	)

	creator := family.NewMember(
		memberId,
		familyID,
		cmd.CreatorName,
		family.OwnerMemberType,
		nil,
		createdAt,
		createdAt,
	)
	creator.SetUserId(&userID)

	if err := fam.AddMember(creator); err != nil {
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
