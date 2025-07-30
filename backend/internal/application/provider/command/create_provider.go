package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type CreateProviderCommand struct {
	Id             *uuid.UUID
	Name           string
	Description    *string
	IconUrl        *string
	Url            *string
	PricingPageUrl *string
	Labels         []uuid.UUID
	Owner          auth.Owner
	CreatedAt      *time.Time
}

type CreateCommandHandler struct {
	providerRepository provider.Repository
	labelRepository    label.Repository
}

func NewCreateProviderCommandHandler(providerRepository provider.Repository,
	labelRepository label.Repository) *CreateCommandHandler {
	return &CreateCommandHandler{
		providerRepository: providerRepository,
		labelRepository:    labelRepository,
	}
}

func (h CreateCommandHandler) Handle(ctx context.Context, cmd CreateProviderCommand) result.Result[provider.Provider] {
	if cmd.Id != nil {
		exists, err := h.providerRepository.Exists(ctx, *cmd.Id)
		if err != nil {
			return result.Fail[provider.Provider](err)
		}
		if exists {
			return result.Fail[provider.Provider](provider.ErrProviderAlreadyExists)
		}
	} else {
		newId, err := uuid.NewV7()
		if err != nil {
			return result.Fail[provider.Provider](err)
		}
		cmd.Id = &newId
	}
	createdAt := ext.ValueOrDefault(cmd.CreatedAt, time.Now())

	if err := ensureLabelExists(ctx, h.labelRepository, cmd.Labels); err != nil {
		return result.Fail[provider.Provider](err)
	}

	prov := provider.NewProvider(
		*cmd.Id,
		cmd.Name,
		cmd.Description,
		cmd.IconUrl,
		cmd.Url,
		cmd.PricingPageUrl,
		cmd.Labels,
		nil,
		cmd.Owner,
		createdAt,
		createdAt,
	)

	if err := prov.GetValidationErrors(); err != nil {
		return result.Fail[provider.Provider](err)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[provider.Provider](err)
	}

	return result.Success(prov)
}
