package system

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/Oleexo/config-go"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/pkg/slicesx"
)

type labelMap map[string]uuid.UUID

func (m labelMap) getUuids(keys []string) ([]uuid.UUID, bool) {
	var results []uuid.UUID

	for _, key := range keys {
		id, ok := m[key]
		if !ok {
			return nil, false
		}
		results = append(results, id)
	}

	return results, true
}

func (m labelMap) Keys() []string {
	results := make([]string, len(m))

	idx := 0
	for k := range m {
		results[idx] = k
		idx++
	}

	return results
}

type systemProviderModel struct {
	Name        string   `json:"name"`
	Key         string   `json:"key"`
	Description *string  `json:"description,omitempty"`
	Website     string   `json:"website"`
	PricingPage string   `json:"pricing_page"`
	Icon        string   `json:"icon"`
	Categories  []string `json:"categories"`
}

type providerUpdater struct {
	downloader         DataDownloader
	providerRepository provider.Repository
	labelRepository    label.Repository
	logger             *slog.Logger
}

func (l providerUpdater) Priority() int {
	return mediumPriority
}

func (l providerUpdater) Update(ctx context.Context) error {
	if l.downloader == nil {
		return nil
	}
	content, err := l.downloader.Download(ctx)
	if err != nil {
		return fmt.Errorf("failed to download providers: %w", err)
	}
	var providers []systemProviderModel
	if err := json.Unmarshal(content, &providers); err != nil {
		return fmt.Errorf("failed to parse JSON content from %s: %w", l.downloader.String(), err)
	}

	return l.updateDatabase(ctx, providers)
}

func (l providerUpdater) getSystemLabels(ctx context.Context) (labelMap, error) {
	lbls, err := l.labelRepository.GetSystemLabels(ctx)
	if err != nil {
		return nil, err
	}

	return slicesx.ToMap(lbls,
		func(lbl label.Label) string {
			return *lbl.Key()
		}, func(lbl label.Label) uuid.UUID {
			return lbl.Id()
		}), nil
}

func (l providerUpdater) updateDatabase(ctx context.Context, sourceProviders []systemProviderModel) error {
	systemProviders, _, err := l.providerRepository.GetSystemProviders(ctx)
	if err != nil {
		return err
	}

	labels, err := l.getSystemLabels(ctx)
	if err != nil {
		return err
	}

	systemProviderMap := slicesx.ToMap(systemProviders, func(prov provider.Provider) string {
		return *prov.Key()
	}, func(prov provider.Provider) provider.Provider { return prov })
	sourceProviderMap := slicesx.ToMap(sourceProviders, func(prov systemProviderModel) string {
		return prov.Key
	}, func(prov systemProviderModel) systemProviderModel { return prov })
	for key, prov := range sourceProviderMap {
		existing, ok := systemProviderMap[key]
		if ok {
			labelIds, ok := labels.getUuids(prov.Categories)
			if !ok {
				l.logger.Warn("missing or invalid labels",
					slog.Any("expected", prov.Categories),
					slog.Any("categories", labels.Keys()))
				continue
			}

			existing.SetName(prov.Name)
			existing.SetDescription(prov.Description)
			existing.SetIconUrl(&prov.Icon)
			existing.SetUrl(&prov.Website)
			existing.SetPricingPageUrl(&prov.PricingPage)
			existing.SetLabels(labelIds)
			// todo bulk ?
			if existing.IsDirty() {
				existing.SetUpdatedAt(time.Now())
				if err := l.providerRepository.Save(ctx, existing); err != nil {
					return err
				}
			}
		} else {
			labelIds, ok := labels.getUuids(prov.Categories)
			if !ok {
				l.logger.Warn("missing or invalid labels",
					slog.Any("expected", prov.Categories),
					slog.Any("categories", labels.Keys()))
				continue
			}

			newProvider := provider.NewProvider(
				uuid.Must(uuid.NewV7()),
				prov.Name,
				&prov.Key,
				prov.Description,
				&prov.Icon,
				&prov.Website,
				&prov.PricingPage,
				labelIds,
				nil,
				auth.SystemOwner,
				time.Now(),
				time.Now(),
			)

			if err := newProvider.GetValidationErrors(); err != nil {
				return err
			}

			if err := l.providerRepository.Save(ctx, newProvider); err != nil {
				return err
			}
		}
	}

	for key, lbl := range systemProviderMap {
		_, ok := sourceProviderMap[key]
		if !ok {
			_, err := l.providerRepository.Delete(ctx, lbl.Id())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func newProviderUpdater(
	cfg config.Configuration,
	providerRepository provider.Repository,
	logger *slog.Logger,
	labelRepository label.Repository) *providerUpdater {
	providerPath := cfg.GetStringOrDefault("DATA_PROVIDER", "")
	var downloader DataDownloader
	if providerPath != "" {
		downloader = newDataDownloader(providerPath, cfg)
	}
	return &providerUpdater{
		downloader:         downloader,
		providerRepository: providerRepository,
		labelRepository:    labelRepository,
		logger:             logger,
	}
}
