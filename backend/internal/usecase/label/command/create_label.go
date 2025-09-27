package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateLabelCommand struct {
	LabelID   option.Option[types.LabelID]
	Name      string
	Color     string
	Owner     types.Owner
	CreatedAt option.Option[time.Time]
}

type CreateLabelCommandHandler struct {
	labelRepository  ports.LabelRepository
	familyRepository ports.FamilyRepository
	authorization    ports.Authorization
	entitlement      ports.EntitlementResolver
}

func NewCreateLabelCommandHandler(labelRepository ports.LabelRepository,
	familyRepository ports.FamilyRepository,
	authorization ports.Authorization,
	entitlement ports.EntitlementResolver) *CreateLabelCommandHandler {
	return &CreateLabelCommandHandler{
		labelRepository:  labelRepository,
		familyRepository: familyRepository,
		authorization:    authorization,
		entitlement:      entitlement,
	}
}

func (h CreateLabelCommandHandler) Handle(ctx context.Context, command CreateLabelCommand) result.Result[label.Label] {
	var labelID types.LabelID
	if command.LabelID.IsSome() {
		labelID = *command.LabelID.Value()
		existingLabel, err := h.labelRepository.GetById(ctx, labelID)
		if err != nil {
			return result.Fail[label.Label](err)
		}
		if existingLabel != nil {
			return result.Fail[label.Label](label.ErrLabelAlreadyExists)
		}
	} else {
		labelID = types.NewLabelID()
	}
	createdAt := command.CreatedAt.ValueOrDefault(time.Now())

	newLabel := label.NewLabel(labelID,
		command.Owner,
		command.Name,
		nil,
		command.Color,
		createdAt,
		createdAt,
	)

	if err := h.authorization.Can(ctx, authorization.PermissionWrite).For(newLabel); err != nil {
		return result.Fail[label.Label](err)
	}

	allowed, _, err := h.entitlement.CheckQuota(ctx, billing.FeatureIdCustomLabelsCount, 1)
	if err != nil {
		return result.Fail[label.Label](err)
	}
	if !allowed {
		return result.Fail[label.Label](label.ErrCustomLabelLimitReached)
	}

	return h.createLabel(ctx, newLabel)
}

func (h CreateLabelCommandHandler) createLabel(
	ctx context.Context,
	lbl label.Label) result.Result[label.Label] {

	if err := lbl.GetValidationErrors(); err != nil {
		return result.Fail[label.Label](err)
	}

	if err := h.labelRepository.Save(ctx, lbl); err != nil {
		return result.Fail[label.Label](err)
	}
	return result.Success(lbl)
}
