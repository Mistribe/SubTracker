package export

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

type ProviderResolver interface {
	ResolveProviderKeys(ctx context.Context, providerIds []types.ProviderID) (map[types.ProviderID]string, error)
}

type providerResolver struct {
	providerRepository ports.ProviderRepository
}

func NewProviderResolver(providerRepository ports.ProviderRepository) ProviderResolver {
	return &providerResolver{
		providerRepository: providerRepository,
	}
}

func (p providerResolver) ResolveProviderKeys(ctx context.Context,
	providerIds []types.ProviderID) (map[types.ProviderID]string, error) {
	if len(providerIds) == 0 {
		return nil, nil
	}

	ids := herd.NewSetFromSlice(providerIds)
	keys := make(map[types.ProviderID]string)

	for providerId := range ids.It() {
		provider, err := p.providerRepository.GetById(ctx, providerId)
		if err != nil {
			return nil, err
		}
		keys[providerId] = provider.Key()
	}

	return keys, nil
}
