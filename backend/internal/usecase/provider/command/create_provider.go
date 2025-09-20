package command

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/user"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x"
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

type CreateProviderCommandHandler struct {
	providerRepository ports.ProviderRepository
	labelRepository    ports.LabelRepository
	authorization      ports.Authorization
}

func NewCreateProviderCommandHandler(
	providerRepository ports.ProviderRepository,
	labelRepository ports.LabelRepository,
	authorization ports.Authorization) *CreateProviderCommandHandler {
	return &CreateProviderCommandHandler{
		providerRepository: providerRepository,
		labelRepository:    labelRepository,
		authorization:      authorization,
	}
}

func (h CreateProviderCommandHandler) Handle(
	ctx context.Context,
	cmd CreateProviderCommand) result.Result[provider.Provider] {
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
	createdAt := x.ValueOrDefault(cmd.CreatedAt, time.Now())

	if err := ensureLabelExists(ctx, h.labelRepository, cmd.Labels); err != nil {
		return result.Fail[provider.Provider](err)
	}

	prov := provider.NewProvider(
		*cmd.Id,
		cmd.Name,
		nil,
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

	if err := h.authorization.Can(ctx, user.PermissionWrite).For(prov); err != nil {
		return result.Fail[provider.Provider](err)
	}

	if err := h.authorization.EnsureWithinLimit(ctx, user.FeatureCustomProviders, 1); err != nil {
		return result.Fail[provider.Provider](err)
	}

	if err := prov.GetValidationErrors(); err != nil {
		return result.Fail[provider.Provider](err)
	}

	if err := h.providerRepository.Save(ctx, prov); err != nil {
		return result.Fail[provider.Provider](err)
	}

	return result.Success(prov)
}
