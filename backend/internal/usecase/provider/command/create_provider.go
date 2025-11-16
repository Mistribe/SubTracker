package command

import (
	"context"
	"time"

	"github.com/mistribe/subtracker/internal/domain/authorization"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/shared"
	"github.com/mistribe/subtracker/pkg/langext/option"
	"github.com/mistribe/subtracker/pkg/langext/result"
)

type CreateProviderCommand struct {
	ProviderID     option.Option[types.ProviderID]
	Name           string
	Description    *string
	IconUrl        *string
	Url            *string
	PricingPageUrl *string
	Labels         []types.LabelID
	Owner          types.OwnerType
	CreatedAt      option.Option[time.Time]
}

type CreateProviderCommandHandler struct {
	providerRepository ports.ProviderRepository
	labelRepository    ports.LabelRepository
	authorization      ports.Authorization
	entitlement        ports.EntitlementResolver
	ownerFactory       shared.OwnerFactory
}

func NewCreateProviderCommandHandler(
	providerRepository ports.ProviderRepository,
	labelRepository ports.LabelRepository,
	authorization ports.Authorization,
	ownerFactory shared.OwnerFactory,
	entitlement ports.EntitlementResolver) *CreateProviderCommandHandler {
	return &CreateProviderCommandHandler{
		providerRepository: providerRepository,
		labelRepository:    labelRepository,
		authorization:      authorization,
		entitlement:        entitlement,
		ownerFactory:       ownerFactory,
	}
}

func (h *CreateProviderCommandHandler) Handle(
	ctx context.Context,
	cmd CreateProviderCommand) result.Result[provider.Provider] {
	var providerID types.ProviderID
	if cmd.ProviderID.IsSome() {
		providerID = *cmd.ProviderID.Value()
		exists, err := h.providerRepository.Exists(ctx, providerID)
		if err != nil {
			return result.Fail[provider.Provider](err)
		}
		if exists {
			return result.Fail[provider.Provider](provider.ErrProviderAlreadyExists)
		}
	} else {
		providerID = types.NewProviderID()
	}
	// Handle nil option (zero-value interface) for CreatedAt gracefully
	var createdAt time.Time
	if cmd.CreatedAt == nil {
		createdAt = time.Now()
	} else {
		createdAt = cmd.CreatedAt.ValueOrDefault(time.Now())
	}
	if err := ensureLabelExists(ctx, h.labelRepository, cmd.Labels); err != nil {
		return result.Fail[provider.Provider](err)
	}
	owner, err := h.ownerFactory.Resolve(ctx, cmd.Owner)
	if err != nil {
		return result.Fail[provider.Provider](err)
	}

	prov := provider.NewProvider(
		providerID,
		cmd.Name,
		cmd.Description,
		cmd.IconUrl,
		cmd.Url,
		cmd.PricingPageUrl,
		cmd.Labels,
		owner,
		createdAt,
		createdAt,
	)

	if err := h.authorization.Can(ctx, authorization.PermissionWrite).For(prov); err != nil {
		return result.Fail[provider.Provider](err)
	}

	allowed, _, qErr := h.entitlement.CheckQuota(ctx, billing.FeatureIdCustomProvidersCount, 1)
	if qErr != nil {
		return result.Fail[provider.Provider](qErr)
	}
	if !allowed {
		return result.Fail[provider.Provider](provider.ErrCustomProviderLimitReached)
	}
	if vErr := prov.GetValidationErrors(); vErr != nil {
		return result.Fail[provider.Provider](vErr)
	}

	if sErr := h.providerRepository.Save(ctx, prov); sErr != nil {
		return result.Fail[provider.Provider](sErr)
	}

	return result.Success(prov)
}
