package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type UpdateLabelCommand struct {
	Id        uuid.UUID
	Name      string
	Color     string
	UpdatedAt option.Option[time.Time]
}

type UpdateLabelCommandHandler struct {
	repository label.Repository
}

func NewUpdateLabelCommandHandler(repository label.Repository) *UpdateLabelCommandHandler {
	return &UpdateLabelCommandHandler{repository: repository}
}

func (h UpdateLabelCommandHandler) Handle(ctx context.Context, command UpdateLabelCommand) result.Result[label.Label] {
	labelOpt, err := h.repository.Get(ctx, command.Id)
	if err != nil {
		return result.Fail[label.Label](err)
	}
	return option.Match(labelOpt, func(in label.Label) result.Result[label.Label] {
		return h.updateLabel(ctx, command, in)
	}, func() result.Result[label.Label] {
		return result.Fail[label.Label](label.ErrLabelNotFound)
	})
}

func (h UpdateLabelCommandHandler) updateLabel(ctx context.Context,
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

	if err := lbl.Validate(); err != nil {
		return result.Fail[label.Label](err)
	}

	if err := h.repository.Save(ctx, lbl); err != nil {
		return result.Fail[label.Label](err)
	}

	return result.Success(lbl)
}
