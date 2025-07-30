package command

import (
	"context"
	"time"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateLabelCommand struct {
	Id        uuid.UUID
	Name      string
	Color     string
	UpdatedAt option.Option[time.Time]
}

type UpdateLabelCommandHandler struct {
	labelRepository  label.Repository
	familyRepository family.Repository
}

func NewUpdateLabelCommandHandler(labelRepository label.Repository,
	familyRepository family.Repository) *UpdateLabelCommandHandler {
	return &UpdateLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
	}
}

func (h UpdateLabelCommandHandler) Handle(ctx context.Context, command UpdateLabelCommand) result.Result[label.Label] {
	existingLabel, err := h.labelRepository.GetById(ctx, command.Id)
	if err != nil {
		return result.Fail[label.Label](err)
	}

	if existingLabel == nil {
		return result.Fail[label.Label](label.ErrLabelNotFound)
	}
	err = auth.EnsureOwnership(ctx, h.familyRepository, existingLabel.Owner())
	if err != nil {
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

	command.UpdatedAt.IfSome(func(updatedAt time.Time) {
		lbl.SetUpdatedAt(updatedAt)
	})
	command.UpdatedAt.IfNone(func() {
		lbl.SetUpdatedAt(time.Now())
	})

	if err := lbl.GetValidationErrors(); err != nil {
		return result.Fail[label.Label](err)
	}

	if err := h.labelRepository.Save(ctx, lbl); err != nil {
		return result.Fail[label.Label](err)
	}

	return result.Success(lbl)
}
