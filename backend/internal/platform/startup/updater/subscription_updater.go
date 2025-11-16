package updater

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Oleexo/config-go"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

type systemSubscriptionModel struct {
	Id               string                `json:"id"`
	FriendlyName     string                `json:"friendly_name"`
	ProviderKey      string                `json:"provider_key"`
	CustomPrice      *systemPriceModel     `json:"custom_price"`
	FreeTrial        *systemFreeTrialModel `json:"free_trial"`
	Owner            systemOwnerModel      `json:"owner"`
	Payer            *systemPayerModel     `json:"payer"`
	FamilyUsers      []string              `json:"family_users"`
	StartDate        time.Time             `json:"start_date"`
	EndDate          *time.Time            `json:"end_date"`
	Recurrency       string                `json:"recurrency"`
	CustomRecurrency *int32                `json:"custom_recurrency"`
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
	subscriptionRepository ports.SubscriptionRepository
	providerRepository     ports.ProviderRepository
}

func (l subscriptionUpdater) Priority() int {
	return lowPriority
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

func (l subscriptionUpdater) getSystemProviders(ctx context.Context) (herd.Dictionary[string, types.ProviderID],
	error) {
	providers, _, err := l.providerRepository.GetSystemProviders(ctx)
	if err != nil {
		return nil, err
	}

	return herd.NewDictionaryFromSlice(
		providers,
		func(prov provider.Provider) string {
			return prov.Key()
		},
		func(prov provider.Provider) types.ProviderID {
			return prov.Id()
		},
	), nil

}

func (l subscriptionUpdater) updateDatabase(ctx context.Context, subscriptions []systemSubscriptionModel) error {
	providers, err := l.getSystemProviders(ctx)
	if err != nil {
		return err
	}
	for _, model := range subscriptions {
		subscriptionID := types.MustParseSubscriptionID(model.Id)
		sub, err := l.subscriptionRepository.GetById(ctx, subscriptionID)
		if err != nil {
			return err
		}

		if sub == nil {
			var freeTrial subscription.FreeTrial
			if model.FreeTrial != nil {
				freeTrial = subscription.NewFreeTrial(model.FreeTrial.StartDate, model.FreeTrial.EndDate)
			}
			var customPrice subscription.Price
			if model.CustomPrice != nil {
				cry, err := currency.ParseISO(model.CustomPrice.Currency)
				if err != nil {
					return err
				}
				customPrice = subscription.NewPrice(currency.NewAmount(model.CustomPrice.Amount, cry))
			}

			familyId, err := types.ParseFamilyIDOrNil(model.Owner.FamilyID)
			if err != nil {
				return err
			}
			userId, err := types.ParseUserIDOrNil(model.Owner.UserID)
			if err != nil {
				return err
			}
			owner := types.NewOwner(
				types.OwnerType(model.Owner.Type),
				familyId,
				userId,
			)
			var payer subscription.Payer
			if model.Payer != nil {
				memberId, err := types.ParseFamilyMemberIDOrNil(model.Payer.MemberId)
				if err != nil {
					return err
				}
				payer = subscription.NewPayer(
					subscription.PayerType(model.Payer.Type),
					*familyId,
					memberId)
			}
			var familyUsers []types.FamilyMemberID
			for _, familyMemberId := range model.FamilyUsers {
				familyUsers = append(familyUsers, types.MustParseFamilyMemberID(familyMemberId))
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
				subscriptionID,
				&model.FriendlyName,
				freeTrial,
				providerId,
				customPrice,
				owner,
				payer,
				familyUsers,
				nil,
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

func newSubscriptionUpdater(
	cfg config.Configuration,
	providerRepository ports.ProviderRepository,
	subscriptionRepository ports.SubscriptionRepository) *subscriptionUpdater {
	labelPath := cfg.GetStringOrDefault("DATA_SUBSCRIPTION", "")
	var downloader DataDownloader
	if labelPath != "" {
		downloader = newDataDownloader(labelPath, cfg)
	}

	return &subscriptionUpdater{
		downloader:             downloader,
		subscriptionRepository: subscriptionRepository,
		providerRepository:     providerRepository,
	}
}
