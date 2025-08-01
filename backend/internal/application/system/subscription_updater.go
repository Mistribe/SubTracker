package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Oleexo/config-go"
	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/slicesx"
	"golang.org/x/text/currency"
)

type systemSubscriptionModel struct {
	Id               string                `json:"id"`
	FriendlyName     string                `json:"friendly_name"`
	ProviderKey      string                `json:"provider_key"`
	CustomPrice      *systemPriceModel     `json:"custom_price"`
	FreeTrial        *systemFreeTrialModel `json:"free_trial"`
	Owner            systemOwnerModel      `json:"owner"`
	Payer            *systemPayerModel     `json:"payer"`
	ServiceUsers     []string              `json:"service_users"`
	StartDate        time.Time             `json:"start_date"`
	EndDate          *time.Time            `json:"end_date"`
	Recurrency       string                `json:"recurrency"`
	CustomRecurrency *uint                 `json:"custom_recurrency"`
}

type systemFreeTrialModel struct {
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}
type systemPriceModel struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
}

type systemOwnerModel struct {
	Type     string  `json:"type"`
	FamilyID *string `json:"family_id,omitempty"`
	UserID   *string `json:"user_id,omitempty"`
}

type systemPayerModel struct {
	Type     string  `json:"type"`
	MemberId *string `json:"member_id"`
}

type subscriptionUpdater struct {
	downloader             DataDownloader
	subscriptionRepository subscription.Repository
	providerRepository     provider.Repository
}

func (l subscriptionUpdater) Priority() int {
	return highPriorty
}

func (l subscriptionUpdater) Update(ctx context.Context) error {
	if l.downloader == nil {
		return nil
	}
	content, err := l.downloader.Download(ctx)
	if err != nil {
		return fmt.Errorf("failed to download subscriptions: %w", err)
	}
	var subscriptions []systemSubscriptionModel
	if err := json.Unmarshal(content, &subscriptions); err != nil {
		return fmt.Errorf("failed to parse JSON content from %s: %w", l.downloader.String(), err)
	}

	return l.updateDatabase(ctx, subscriptions)
}

func (l subscriptionUpdater) getSystemProviders(ctx context.Context) (map[string]uuid.UUID, error) {
	providers, err := l.providerRepository.GetSystemProviders(ctx)
	if err != nil {
		return nil, err
	}

	return slicesx.ToMap(providers, func(prov provider.Provider) string {
		return *prov.Key()
	}, func(prov provider.Provider) uuid.UUID {
		return prov.Id()
	}), nil

}

func (l subscriptionUpdater) updateDatabase(ctx context.Context, subscriptions []systemSubscriptionModel) error {
	providers, err := l.getSystemProviders(ctx)
	if err != nil {
		return err
	}
	for _, model := range subscriptions {
		id := uuid.MustParse(model.Id)
		sub, err := l.subscriptionRepository.GetById(ctx, id)
		if err != nil {
			return err
		}

		if sub == nil {
			var freeTrial subscription.FreeTrial
			if model.FreeTrial != nil {
				freeTrial = subscription.NewFreeTrial(model.FreeTrial.StartDate, model.FreeTrial.EndDate)
			}
			var customPrice subscription.CustomPrice
			if model.CustomPrice != nil {
				cry, err := currency.ParseISO(model.CustomPrice.Currency)
				if err != nil {
					return err
				}
				customPrice = subscription.NewCustomPrice(model.CustomPrice.Amount, cry)
			}
			var familyId *uuid.UUID
			if model.Owner.FamilyID != nil {
				fid, err := uuid.Parse(*model.Owner.FamilyID)
				if err != nil {
					return err
				}
				familyId = &fid
			}
			owner := auth.NewOwner(auth.OwnerType(model.Owner.Type), familyId, model.Owner.UserID)
			var payer subscription.Payer
			if model.Payer != nil {
				var memberId *uuid.UUID
				if model.Payer.MemberId != nil {
					mid, err := uuid.Parse(*model.Payer.MemberId)
					if err != nil {
						return err
					}
					memberId = &mid
				}
				payer = subscription.NewPayer(subscription.PayerType(model.Payer.Type),
					*familyId, memberId)
			}
			var serviceUsers []uuid.UUID
			for _, userId := range model.ServiceUsers {
				serviceUsers = append(serviceUsers, uuid.MustParse(userId))
			}
			providerId, ok := providers[model.ProviderKey]
			if !ok {
				return errors.New("missing provider")
			}
			recurrency, err := subscription.ParseRecurrencyType(model.Recurrency)
			if err != nil {
				return err
			}

			sub = subscription.NewSubscription(
				id,
				&model.FriendlyName,
				freeTrial,
				providerId,
				nil,
				nil,
				customPrice,
				owner,
				payer,
				serviceUsers,
				model.StartDate,
				model.EndDate,
				recurrency,
				model.CustomRecurrency,
				time.Now(),
				time.Now(),
			)

			if err := sub.GetValidationErrors(); err != nil {
				return err
			}
			if err := l.subscriptionRepository.Save(ctx, sub); err != nil {
				return err
			}
		} else {
			// todo
		}
	}

	return nil
}

func newSubscriptionUpdater(cfg config.Configuration,
	providerRepository provider.Repository,
	subscriptionRepository subscription.Repository) *subscriptionUpdater {
	labelPath := cfg.GetStringOrDefault("DATA_SUBSCRIPTION", "")
	var downloader DataDownloader
	if labelPath != "" {
		downloader = newDataDownloader(labelPath)
	}

	return &subscriptionUpdater{
		downloader:             downloader,
		subscriptionRepository: subscriptionRepository,
		providerRepository:     providerRepository,
	}
}
