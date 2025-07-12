package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateLabelCommand struct {
	Label label.Label
}

type CreateLabelCommandHandler struct {
	repository label.Repository
}

func NewCreateLabelCommandHandler(repository label.Repository) *CreateLabelCommandHandler {
	return &CreateLabelCommandHandler{repository: repository}
}

func (h CreateLabelCommandHandler) Handle(ctx context.Context, command CreateLabelCommand) result.Result[label.Label] {
	existingLabelOpt, err := h.repository.Get(ctx, command.Label.Id())
	if err != nil {
		return result.Fail[label.Label](err)
	}

	return option.Match(existingLabelOpt, func(existingLabel label.Label) result.Result[label.Label] {
		if existingLabel.Equal(command.Label) {
			return result.Success(command.Label)
		}
		return result.Fail[label.Label](label.ErrLabelAlreadyExists)
	}, func() result.Result[label.Label] {
		r := label.NewLabel(command.Label.Id(),
			command.Label.Name(),
			false,
			command.Label.Color(),
			command.Label.CreatedAt(),
			command.Label.CreatedAt(),
		)
		return h.createLabel(ctx, command, r)
	})
}

func (h CreateLabelCommandHandler) createLabel(ctx context.Context, command CreateLabelCommand,
	r result.Result[label.Label]) result.Result[label.Label] {
	return result.Bind[label.Label, label.Label](r,
		func(value label.Label) result.Result[label.Label] {
			if err := h.repository.Save(ctx, command.Label); err != nil {
				return result.Fail[label.Label](err)
			}
			return result.Success(command.Label)
		},
	)
}
