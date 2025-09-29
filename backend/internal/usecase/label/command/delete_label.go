package command

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type DeleteLabelCommand struct {
	LabelID types.LabelID
}

type DeleteLabelCommandHandler struct {
	labelRepository  ports.LabelRepository
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewDeleteLabelCommandHandler(
	labelRepository ports.LabelRepository,
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *DeleteLabelCommandHandler {
	return &DeleteLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h DeleteLabelCommandHandler) Handle(ctx context.Context, command DeleteLabelCommand) result.Result[bool] {
	existingLabel, err := h.labelRepository.GetById(ctx, command.LabelID)
	if err != nil {
		return result.Fail[bool](err)
	}
	if existingLabel == nil {
		return result.Fail[bool](label.ErrLabelNotFound)
	}

	permReq := h.authorization.Can(ctx, authorization.PermissionDelete)
	if permReq == nil {
		return result.Fail[bool](authorization.ErrUnauthorized)
	}
	if err = permReq.For(existingLabel); err != nil {
		return result.Fail[bool](err)
	}

	ok, err := h.labelRepository.Delete(ctx, command.LabelID)
	if err != nil {
		return result.Fail[bool](err)
	}
	if !ok {
		return result.Fail[bool](label.ErrLabelNotFound)
	}

	return result.Success(true)
}
