package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/family"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteLabelCommand struct {
	Id uuid.UUID
}

type DeleteLabelCommandHandler struct {
	labelRepository  label.Repository
	familyRepository family.Repository
}

func NewDeleteLabelCommandHandler(labelRepository label.Repository,
	familyRepository family.Repository) *DeleteLabelCommandHandler {
	return &DeleteLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
	}
}

func (h DeleteLabelCommandHandler) Handle(ctx context.Context, command DeleteLabelCommand) result.Result[bool] {
	existingLabel, err := h.labelRepository.GetById(ctx, command.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if existingLabel == nil {
		return result.Fail[bool](label.ErrLabelNotFound)
	}

	err = auth.EnsureOwnership(ctx, h.familyRepository, existingLabel.Owner())
	if err != nil {
		return result.Fail[bool](err)
	}

	ok, err := h.labelRepository.Delete(ctx, command.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](label.ErrLabelNotFound)
	}

	return result.Success(true)
}
