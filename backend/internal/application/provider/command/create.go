package command

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateProviderCommand struct {
	Provider provider.Provider
}

func NewCreateProviderCommand(provider provider.Provider) CreateProviderCommand {
	return CreateProviderCommand{
		Provider: provider,
	}
}

type CreateCommandHandler struct {
	repository provider.Repository
}

func NewCreateCommandHandler(repository provider.Repository) *CreateCommandHandler {
	return &CreateCommandHandler{repository: repository}
}

func (h CreateCommandHandler) Handle(ctx context.Context, cmd CreateProviderCommand) result.Result[provider.Provider] {
	exists, err := h.repository.Exists(ctx, cmd.Provider.Id())
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	if exists {
		return result.Fail[provider.Provider](provider.ErrProviderAlreadyExists)
	}

	err = cmd.Provider.GetValidationErrors()
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	err = h.repository.Save(ctx, cmd.Provider)
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	return result.Success(cmd.Provider)
}
