package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Oleexo/config-go"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/auth"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

type systemLabelModel struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type labelUpdater struct {
	downloader      DataDownloader
	labelRepository ports.LabelRepository
}

func (l labelUpdater) Priority() int {
	return highPriority
}

func (l labelUpdater) Update(ctx context.Context) error {
	if l.downloader == nil {
		return nil
	}
	content, err := l.downloader.Download(ctx)
	if err != nil {
		return fmt.Errorf("failed to download labels: %w", err)
	}
	var labels []systemLabelModel
	if err := json.Unmarshal(content, &labels); err != nil {
		return fmt.Errorf("failed to parse JSON content from %s: %w", l.downloader.String(), err)
	}

	return l.updateDatabase(ctx, labels)
}

func (l labelUpdater) updateDatabase(ctx context.Context, sourceLabels []systemLabelModel) error {
	systemLabels, err := l.labelRepository.GetSystemLabels(ctx)
	if err != nil {
		return err
	}

	systemLabelMap := herd.ToMap(systemLabels, func(lbl label.Label) string {
		return *lbl.Key()
	}, func(lbl label.Label) label.Label { return lbl })
	sourceLabelMap := herd.ToMap(sourceLabels, func(lbl systemLabelModel) string {
		return lbl.Key
	}, func(lbl systemLabelModel) systemLabelModel { return lbl })
	for key, lbl := range sourceLabelMap {
		existing, ok := systemLabelMap[key]
		if ok {
			existing.SetColor(lbl.Color)
			// todo bulk ?
			if existing.IsDirty() {
				existing.SetUpdatedAt(time.Now())
				if err := l.labelRepository.Save(ctx, existing); err != nil {
					return err
				}
			}
		} else {
			newLabel := label.NewLabel(
				uuid.Must(uuid.NewV7()),
				auth.SystemOwner,
				lbl.Name,
				&lbl.Key,
				lbl.Color,
				time.Now(),
				time.Now(),
			)

			if err := newLabel.GetValidationErrors(); err != nil {
				return err
			}

			if err := l.labelRepository.Save(ctx, newLabel); err != nil {
				return err
			}
		}
	}

	for key, lbl := range systemLabelMap {
		_, ok := sourceLabelMap[key]
		if !ok {
			_, err := l.labelRepository.Delete(ctx, lbl.Id())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func newLabelUpdater(
	cfg config.Configuration,
	labelRepository ports.LabelRepository) *labelUpdater {
	labelPath := cfg.GetStringOrDefault("DATA_LABEL", "")
	var downloader DataDownloader
	if labelPath != "" {
		downloader = newDataDownloader(labelPath, cfg)
	}
	return &labelUpdater{
		downloader:      downloader,
		labelRepository: labelRepository,
	}
}
