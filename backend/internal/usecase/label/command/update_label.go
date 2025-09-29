package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type UpdateLabelCommand struct {
	LabelID   types.LabelID
	Name      string
	Color     string
	UpdatedAt option.Option[time.Time]
}

type UpdateLabelCommandHandler struct {
	labelRepository  ports.LabelRepository
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
}

func NewUpdateLabelCommandHandler(
	labelRepository ports.LabelRepository,
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization) *UpdateLabelCommandHandler {
	return &UpdateLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
		authorization:    authorization,
	}
}

func (h UpdateLabelCommandHandler) Handle(ctx context.Context, command UpdateLabelCommand) result.Result[label.Label] {
	existingLabel, err := h.labelRepository.GetById(ctx, command.LabelID)
	if err != nil {
		return result.Fail[label.Label](err)
	}

	if existingLabel == nil {
		return result.Fail[label.Label](label.ErrLabelNotFound)
	}

	permReq := h.authorization.Can(ctx, authorization.PermissionWrite)
	if permReq == nil {
		return result.Fail[label.Label](authorization.ErrUnauthorized)
	}
	if err = permReq.For(existingLabel); err != nil {
		return result.Fail[label.Label](err)
	}

	return h.updateLabel(ctx, command, existingLabel)
}

func (h UpdateLabelCommandHandler) updateLabel(
	ctx context.Context,
	command UpdateLabelCommand,
	lbl label.Label) result.Result[label.Label] {
	lbl.SetName(command.Name)
	lbl.SetColor(command.Color)

	lbl.SetUpdatedAt(command.UpdatedAt.ValueOrDefault(time.Now()))

	if err := lbl.GetValidationErrors(); err != nil {
		return result.Fail[label.Label](err)
	}

	if err := h.labelRepository.Save(ctx, lbl); err != nil {
		return result.Fail[label.Label](err)
	}

	return result.Success(lbl)
}
