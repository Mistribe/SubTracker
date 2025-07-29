package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type UpdateProviderCommand struct {
	Provider provider.Provider
}

func NewUpdateProviderCommand(provider provider.Provider) UpdateProviderCommand {
	return UpdateProviderCommand{
		Provider: provider,
	}
}

type UpdateCommandHandler struct {
	repository provider.Repository
}

func NewUpdateCommandHandler(repository provider.Repository) *UpdateCommandHandler {
	return &UpdateCommandHandler{repository: repository}
}

func (h UpdateCommandHandler) Handle(ctx context.Context, cmd UpdateProviderCommand) result.Result[provider.Provider] {
	prvdr, err := h.repository.GetById(ctx, cmd.Provider.Id())
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	if prvdr == nil {
		return result.Fail[provider.Provider](provider.ErrProviderNotFound)

	}

	return h.update(ctx, cmd, prvdr)
}

func (h UpdateCommandHandler) update(
	ctx context.Context,
	cmd UpdateProviderCommand,
	in provider.Provider) result.Result[provider.Provider] {

	userId, ok := user.FromContext(ctx)
	if !ok {
		return result.Fail[provider.Provider](user.ErrUnknownUser)
	}

	if in.Owner().Type() == user.PersonalOwner &&
		in.Owner().UserId() != userId {
		return result.Fail[provider.Provider](provider.ErrOnlyOwnerCanEdit)
	}

	// todo family edit

	if in.Equal(cmd.Provider) {
		return result.Success(in)
	}

	in.SetName(cmd.Provider.Name())
	in.SetDescription(cmd.Provider.Description())
	in.SetIconUrl(cmd.Provider.IconUrl())
	in.SetUrl(cmd.Provider.Url())
	in.SetPricingPageUrl(cmd.Provider.PricingPageUrl())
	in.SetLabels(cmd.Provider.Labels().Values())
	in.SetPlans(cmd.Provider.Plans().Values())
	in.SetUpdatedAt(cmd.Provider.UpdatedAt())

	validationErrors := in.GetValidationErrors()
	if validationErrors.HasErrors() {
		return result.Fail[provider.Provider](validationErrors)
	}

	if err := h.repository.Save(ctx, in); err != nil {
		return result.Fail[provider.Provider](err)
	}

	return result.Success(in)
}
