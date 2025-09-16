package updater

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Oleexo/config-go"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

type systemMemberModel struct {
	Id     string  `json:"id"`
	Name   string  `json:"name"`
	UserId *string `json:"user_id,omitempty"`
	Type   string  `json:"type"`
}

type systemFamilyModel struct {
	Id      string              `json:"id"`
	OwnerId string              `json:"owner_id"`
	Name    string              `json:"name"`
	Members []systemMemberModel `json:"members"`
}

type familyUpdater struct {
	downloader       DataDownloader
	familyRepository ports.FamilyRepository
}

func newFamilyUpdater(
	cfg config.Configuration,
	familyRepository ports.FamilyRepository) *familyUpdater {
	providerPath := cfg.GetStringOrDefault("DATA_FAMILY", "")
	var downloader DataDownloader
	if providerPath != "" {
		downloader = newDataDownloader(providerPath, cfg)
	}
	return &familyUpdater{
		downloader:       downloader,
		familyRepository: familyRepository,
	}
}

func (l familyUpdater) Priority() int {
	return highPriority
}

func (l familyUpdater) Update(ctx context.Context) error {
	if l.downloader == nil {
		return nil
	}
	content, err := l.downloader.Download(ctx)
	if err != nil {
		return fmt.Errorf("failed to download familys: %w", err)
	}
	var familys []systemFamilyModel
	if err = json.Unmarshal(content, &familys); err != nil {
		return fmt.Errorf("failed to parse JSON content from %s: %w", l.downloader.String(), err)
	}

	return l.updateDatabase(ctx, familys)
}

func (l familyUpdater) updateDatabase(ctx context.Context, families []systemFamilyModel) error {
	for _, model := range families {
		id := uuid.MustParse(model.Id)
		fam, err := l.familyRepository.GetById(ctx, id)
		if err != nil {
			return err
		}
		if fam == nil {
			members := herd.Select(model.Members, func(memberModel systemMemberModel) family.Member {
				mbr := family.NewMember(
					uuid.MustParse(memberModel.Id),
					id,
					memberModel.Name,
					family.MustParseMemberType(memberModel.Type),
					nil,
					time.Now(),
					time.Now(),
				)

				if memberModel.UserId != nil {
					mbr.SetUserId(memberModel.UserId)
				}
				return mbr
			})
			fam = family.NewFamily(
				id,
				model.OwnerId,
				model.Name,
				nil,
				time.Now(),
				time.Now(),
			)
			for _, mbr := range members {
				fam.AddMember(mbr)
			}

			if err := l.familyRepository.Save(ctx, fam); err != nil {
				return err
			}
		} else {
			// todo
		}
	}

	return nil
}
