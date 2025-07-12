package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateFamilyMemberCommand struct {
	Id        uuid.UUID
	Name      string
	Email     option.Option[string]
	IsKid     bool
	UpdatedAt option.Option[time.Time]
}

type UpdateFamilyMemberCommandHandler struct {
	repository family.Repository
}

func NewUpdateFamilyMemberCommandHandler(repository family.Repository) *UpdateFamilyMemberCommandHandler {
	return &UpdateFamilyMemberCommandHandler{repository: repository}
}

func (h UpdateFamilyMemberCommandHandler) Handle(
	ctx context.Context,
	command UpdateFamilyMemberCommand) result.Result[family.Member] {
	mbrOpt, err := h.repository.Get(ctx, command.Id)
	if err != nil {
		return result.Fail[family.Member](err)
	}

	return option.Match(mbrOpt, func(in family.Member) result.Result[family.Member] {
		return h.updateFamilyMember(ctx, command, in)
	}, func() result.Result[family.Member] {
		return result.Fail[family.Member](family.ErrFamilyMemberNotFound)
	})
}

func (h UpdateFamilyMemberCommandHandler) updateFamilyMember(
	ctx context.Context,
	command UpdateFamilyMemberCommand,
	mbr family.Member) result.Result[family.Member] {
	mbr.SetName(command.Name)
	if command.IsKid {
		mbr.SetAsKid()
	} else {
		mbr.SetAsAdult()
	}
	mbr.SetEmail(command.Email)

	command.UpdatedAt.IfSome(func(updatedAt time.Time) {
		mbr.SetUpdatedAt(updatedAt)
	})
	command.UpdatedAt.IfNone(func() {
		mbr.SetUpdatedAt(time.Now())
	})

	if err := mbr.Validate(); err != nil {
		return result.Fail[family.Member](err)
	}

	if err := h.repository.Save(ctx, mbr); err != nil {
		return result.Fail[family.Member](err)
	}

	return result.Success(mbr)
}
