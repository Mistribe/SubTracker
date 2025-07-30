package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateProviderCommand struct {
	Id             uuid.UUID
	Name           string
	Description    *string
	IconUrl        *string
	Url            *string
	PricingPageUrl *string
	Labels         []uuid.UUID
	UpdatedAt      *time.Time
}

type UpdateCommandHandler struct {
	providerRepository provider.Repository
	labelRepository    label.Repository
}

func NewUpdateProviderCommandHandler(providerRepository provider.Repository,
	labelRepository label.Repository) *UpdateCommandHandler {
	return &UpdateCommandHandler{
		providerRepository: providerRepository,
		labelRepository:    labelRepository,
	}
}

func (h UpdateCommandHandler) Handle(ctx context.Context, cmd UpdateProviderCommand) result.Result[provider.Provider] {
	prov, err := h.providerRepository.GetById(ctx, cmd.Id)
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	if prov == nil {
		return result.Fail[provider.Provider](provider.ErrProviderNotFound)

	}

	return h.update(ctx, cmd, prov)
}

func (h UpdateCommandHandler) update(
	ctx context.Context,
	cmd UpdateProviderCommand,
	prov provider.Provider) result.Result[provider.Provider] {

	userId, ok := user.FromContext(ctx)
	if !ok {
		return result.Fail[provider.Provider](user.ErrUnknownUser)
	}

	if prov.Owner().Type() == user.PersonalOwner &&
		prov.Owner().UserId() != userId {
		return result.Fail[provider.Provider](provider.ErrOnlyOwnerCanEdit)
	}

	if err := ensureLabelExists(ctx, h.labelRepository, cmd.Labels); err != nil {
		return result.Fail[provider.Provider](err)
	}
	updatedAt := ext.ValueOrDefault(cmd.UpdatedAt, time.Now())

	prov.SetName(cmd.Name)
	prov.SetDescription(cmd.Description)
	prov.SetIconUrl(cmd.IconUrl)
	prov.SetUrl(cmd.Url)
	prov.SetPricingPageUrl(cmd.PricingPageUrl)
	prov.SetLabels(cmd.Labels)
	prov.SetUpdatedAt(updatedAt)

	validationErrors := prov.GetValidationErrors()
	if validationErrors.HasErrors() {
		return result.Fail[provider.Provider](validationErrors)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[provider.Provider](err)
	}

	return result.Success(prov)
}
