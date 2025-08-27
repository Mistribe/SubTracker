package command

import (
	"context"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteLabelCommand struct {
	Id uuid.UUID
}

type DeleteLabelCommandHandler struct {
	labelRepository  ports.LabelRepository
	familyRepository ports.FamilyRepository
	authService      ports.AuthService
}

func NewDeleteLabelCommandHandler(labelRepository ports.LabelRepository,
	familyRepository ports.FamilyRepository,
	authService ports.AuthService) *DeleteLabelCommandHandler {
	return &DeleteLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
		authService:      authService,
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

	ok, err := h.authService.IsOwner(ctx, existingLabel.Owner())
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](family.ErrFamilyNotFound)
	}

	ok, err = h.labelRepository.Delete(ctx, command.Id)
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](label.ErrLabelNotFound)
	}

	return result.Success(true)
}
