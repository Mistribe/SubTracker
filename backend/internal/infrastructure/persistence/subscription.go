package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type subscriptionLabelModel struct {
	LabelId uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
	SubId   uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
}
type subscriptionFamilyMemberModel struct {
	FamilyMemberId uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
	SubId          uuid.UUID `gorm:"primaryKey;type:uuid;not null"`
}

type subscriptionPaymentModel struct {
	BaseModel
	SubId     uuid.UUID  `gorm:"type:uuid;not null"`
	Price     float64    `gorm:"type:numeric(10,2);not null"`
	StartDate time.Time  `gorm:"type:timestamp;not null"`
	EndDate   *time.Time `gorm:"type:timestamp"`
	Months    int        `gorm:"type:integer;not null"`
	Currency  string     `gorm:"type:varchar(10);not null"`
}
type subscriptionModel struct {
	BaseModel
	Name          string                          `gorm:"type:varchar(100);not null"`
	Payments      []subscriptionPaymentModel      `gorm:"foreignKey:SubId"`
	Labels        []subscriptionLabelModel        `gorm:"foreignKey:SubId"`
	FamilyMembers []subscriptionFamilyMemberModel `gorm:"foreignKey:SubId"`
	PayerId       *uuid.UUID                      `gorm:"type:uuid"`
}

type SubscriptionRepository struct {
	repository *Repository
}

func NewSubscriptionRepository(repository *Repository) *SubscriptionRepository {
	return &SubscriptionRepository{
		repository: repository,
	}
}

func (r SubscriptionRepository) toPaymentModel(subId uuid.UUID, source subscription.Payment) subscriptionPaymentModel {
	return subscriptionPaymentModel{
		BaseModel: BaseModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
		},
		SubId:     subId,
		Price:     source.Price(),
		StartDate: source.StartDate(),
		EndDate: option.Match(source.EndDate(), func(in time.Time) *time.Time {
			return &in
		}, func() *time.Time {
			return nil
		}),
		Months:   source.Months(),
		Currency: source.Currency(),
	}
}

func (r SubscriptionRepository) toModel(source *subscription.Subscription) subscriptionModel {
	return subscriptionModel{
		BaseModel: BaseModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
		},
		Name: source.Name(),
		Payments: ext.Map(source.Payments(), func(in subscription.Payment) subscriptionPaymentModel {
			return r.toPaymentModel(source.Id(), in)
		}),
		Labels: ext.Map(source.Labels(), func(in uuid.UUID) subscriptionLabelModel {
			return subscriptionLabelModel{
				LabelId: in,
				SubId:   source.Id(),
			}
		}),
		FamilyMembers: ext.Map(source.FamilyMembers(), func(in uuid.UUID) subscriptionFamilyMemberModel {
			return subscriptionFamilyMemberModel{
				FamilyMemberId: in,
				SubId:          source.Id(),
			}
		}),
		PayerId: option.Match(source.Payer(), func(in uuid.UUID) *uuid.UUID {
			return &in
		}, func() *uuid.UUID {
			return nil
		}),
	}
}

func (r SubscriptionRepository) toEntity(source subscriptionModel) subscription.Subscription {
	payments := ext.Map(source.Payments, func(source subscriptionPaymentModel) subscription.Payment {
		return subscription.NewPaymentWithoutValidation(
			source.Id,
			source.Price,
			source.StartDate,
			option.New(source.EndDate),
			source.Months,
			source.Currency,
			source.CreatedAt,
			source.UpdatedAt,
			true,
		)
	})
	labels := ext.Map(source.Labels, func(source subscriptionLabelModel) uuid.UUID {
		return source.LabelId
	})
	familyMembers := ext.Map(source.FamilyMembers, func(source subscriptionFamilyMemberModel) uuid.UUID {
		return source.FamilyMemberId
	})
	sub := subscription.NewSubscriptionWithoutValidation(
		source.Id,
		source.Name,
		payments,
		labels,
		familyMembers,
		option.New(source.PayerId),
		source.CreatedAt,
		source.UpdatedAt,
		true,
	)
	sub.Clean()
	return sub
}

func (r SubscriptionRepository) Get(ctx context.Context, id uuid.UUID) (
	option.Option[subscription.Subscription],
	error) {
	var model subscriptionModel
	result := r.repository.db.WithContext(ctx).
		InnerJoins("Payments").
		InnerJoins("Labels").
		InnerJoins("FamilyMembers").
		First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return option.None[subscription.Subscription](), nil
		}
		return option.None[subscription.Subscription](), result.Error
	}
	return option.Some(r.toEntity(model)), nil
}

func (r SubscriptionRepository) GetAll(ctx context.Context) ([]subscription.Subscription, error) {
	var models []subscriptionModel
	if result := r.repository.db.WithContext(ctx).
		Preload("Payments").
		Preload("Labels").
		Preload("FamilyMembers").
		Find(&models); result.Error != nil {
		return nil, result.Error
	}

	result := make([]subscription.Subscription, 0, len(models))
	for _, model := range models {
		result = append(result, r.toEntity(model))
	}
	return result, nil
}

func (r SubscriptionRepository) Save(ctx context.Context, subscription subscription.Subscription) error {
	if subscription.IsDirty() == false {
		return nil
	}

	dbSubscription := r.toModel(&subscription)
	var result *gorm.DB
	if subscription.IsExists() {
		result = r.repository.db.WithContext(ctx).Save(&dbSubscription)
	} else {
		result = r.repository.db.WithContext(ctx).Create(&dbSubscription)
	}
	if result.Error != nil {
		return result.Error
	}
	subscription.Clean()
	return nil
}

func (r SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if result := r.repository.db.WithContext(ctx).Delete(&subscriptionModel{}, id); result.Error != nil {
		return result.Error
	}
	return nil
}
