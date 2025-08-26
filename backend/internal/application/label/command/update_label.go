package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/family"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
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
	authService      auth.Service
}

func NewUpdateLabelCommandHandler(labelRepository label.Repository,
	familyRepository family.Repository,
	authService auth.Service) *UpdateLabelCommandHandler {
	return &UpdateLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
		authService:      authService,
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

	ok, err := h.authService.IsOwner(ctx, existingLabel.Owner())
	if err != nil {
		return result.Fail[label.Label](err)
	}
	if !ok {
		return result.Fail[label.Label](family.ErrFamilyNotFound)
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
