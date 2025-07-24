package persistence

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/currency"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type subscriptionLabelModel struct {
	LabelId      uuid.UUID         `gorm:"primaryKey;type:uuid;not null"`
	Label        labelModel        `gorm:"foreignKey:LabelId;references:Id"`
	SubId        uuid.UUID         `gorm:"primaryKey;type:uuid;not null"`
	Subscription subscriptionModel `gorm:"foreignKey:SubId;references:Id"`
}

func (s subscriptionLabelModel) TableName() string {
	return "subscription_labels"
}

type subscriptionFamilyMemberModel struct {
	FamilyMemberId uuid.UUID         `gorm:"primaryKey;type:uuid;not null"`
	FamilyMember   familyMemberModel `gorm:"foreignKey:FamilyMemberId;references:Id"`
	SubId          uuid.UUID         `gorm:"primaryKey;type:uuid;not null"`
	Subscription   subscriptionModel `gorm:"foreignKey:SubId;references:Id"`
}

func (s subscriptionFamilyMemberModel) TableName() string {
	return "subscription_family_members"
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

func (s subscriptionPaymentModel) TableName() string {
	return "subscription_payments"
}

type subscriptionModel struct {
	BaseModel
	Name                string                          `gorm:"type:varchar(100);not null"`
	Payments            []subscriptionPaymentModel      `gorm:"foreignKey:SubId;references:Id;constraint:OnDelete:CASCADE"`
	Labels              []subscriptionLabelModel        `gorm:"foreignKey:SubId;references:Id;constraint:OnDelete:CASCADE"`
	FamilyMembers       []subscriptionFamilyMemberModel `gorm:"foreignKey:SubId;references:Id;constraint:OnDelete:CASCADE"`
	PayerId             *uuid.UUID                      `gorm:"type:uuid"`
	Payer               *familyMemberModel              `gorm:"foreignKey:PayerId;references:Id"`
	PayedByJointAccount bool                            `gorm:"type:boolean;not null;default:false"`
	FamilyId            *uuid.UUID                      `gorm:"type:uuid"`
	Family              *familyModel                    `gorm:"foreignKey:FamilyId;references:Id"`
}

func (s subscriptionModel) TableName() string {
	return "subscriptions"
}

type SubscriptionRepository struct {
	repository *Repository
}

func NewSubscriptionRepository(repository *Repository) *SubscriptionRepository {
	return &SubscriptionRepository{
		repository: repository,
	}
}

func (r SubscriptionRepository) toPaymentModel(source subscription.Payment) subscriptionPaymentModel {
	return subscriptionPaymentModel{
		BaseModel: BaseModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
		},
		SubId:     source.SubscriptionId(),
		Price:     source.Price(),
		StartDate: source.StartDate(),
		EndDate: option.Match(source.EndDate(), func(in time.Time) *time.Time {
			return &in
		}, func() *time.Time {
			return nil
		}),
		Months:   source.Months(),
		Currency: source.Currency().String(),
	}
}

func (r SubscriptionRepository) toModel(source *subscription.Subscription) subscriptionModel {
	return subscriptionModel{
		BaseModel: BaseModel{
			Id:        source.Id(),
			CreatedAt: source.CreatedAt(),
			UpdatedAt: source.UpdatedAt(),
			Etag:      source.ETag(),
		},
		Name:          source.Name(),
		Payments:      nil,
		Labels:        nil,
		FamilyMembers: nil,
		PayerId: option.Match(source.Payer(), func(in uuid.UUID) *uuid.UUID {
			return &in
		}, func() *uuid.UUID {
			return nil
		}),
		PayedByJointAccount: source.PayedByJointAccount(),
		FamilyId: option.Match(source.FamilyId(), func(in uuid.UUID) *uuid.UUID {
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
			currency.MustParseISO(source.Currency),
			source.SubId,
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
		source.FamilyId,
		source.Name,
		payments,
		labels,
		familyMembers,
		source.PayerId,
		source.PayedByJointAccount,
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
		Preload("Payments").
		Preload("Labels").
		Preload("FamilyMembers").
		First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return option.None[subscription.Subscription](), nil
		}
		return option.None[subscription.Subscription](), result.Error
	}
	return option.Some(r.toEntity(model)), nil
}

func (r SubscriptionRepository) GetAll(ctx context.Context, size, page int) ([]subscription.Subscription, error) {
	var models []subscriptionModel
	if result := r.repository.db.WithContext(ctx).
		Preload("Payments").
		Preload("Labels").
		Preload("FamilyMembers").
		Offset((page - 1) * size).Limit(size).
		Find(&models); result.Error != nil {
		return nil, result.Error
	}

	result := make([]subscription.Subscription, 0, len(models))
	for _, model := range models {
		result = append(result, r.toEntity(model))
	}
	return result, nil
}

func (r SubscriptionRepository) GetAllCount(ctx context.Context) (int64, error) {
	var count int64
	if result := r.repository.db.WithContext(ctx).
		Model(&subscriptionModel{}).
		Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (r SubscriptionRepository) Save(ctx context.Context, subscription *subscription.Subscription) error {
	if subscription.IsDirty() == false {
		return nil
	}

	dbSubscription := r.toModel(subscription)
	var result *gorm.DB
	if subscription.IsExists() {
		result = r.repository.db.WithContext(ctx).
			Omit("Payments").
			Omit("Labels").
			Omit("FamilyMembers").
			Save(&dbSubscription)
	} else {
		result = r.repository.db.WithContext(ctx).Create(&dbSubscription)
	}
	if result.Error != nil {
		return result.Error
	}

	if err := saveTrackedSlice(ctx,
		subscription.Payments(),
		r.createPayment,
		r.updatePayment,
		r.deletePayment); err != nil {
		return err
	}

	if err := saveTrackedSlice(ctx,
		subscription.Labels(),
		func(ctx context.Context, labelId uuid.UUID) error {
			return r.createSubscriptionLabel(ctx, subscription.Id(), labelId)
		},
		func(ctx context.Context, labelId uuid.UUID) error {
			return r.updateSubscriptionLabel(ctx, subscription.Id(), labelId)
		},
		func(ctx context.Context, labelId uuid.UUID) error {
			return r.deleteSubscriptionLabel(ctx, subscription.Id(), labelId)
		}); err != nil {
		return err
	}

	if err := saveTrackedSlice(ctx,
		subscription.FamilyMembers(),
		func(ctx context.Context, memberId uuid.UUID) error {
			return r.createSubscriptionFamilyMember(ctx, subscription.Id(), memberId)
		},
		func(ctx context.Context, memberId uuid.UUID) error {
			return r.updateSubscriptionFamilyMember(ctx, subscription.Id(), memberId)
		},
		func(ctx context.Context, memberId uuid.UUID) error {
			return r.deleteSubscriptionFamilyMember(ctx, subscription.Id(), memberId)
		}); err != nil {
		return err
	}

	subscription.Clean()
	return nil
}

func (r SubscriptionRepository) createPayment(
	ctx context.Context,
	payment subscription.Payment) error {
	if payment.IsDirty() == false {
		return nil
	}

	dbPayment := r.toPaymentModel(payment)
	if payment.IsExists() {
		return subscription.ErrPaymentAlreadyExists
	}
	result := r.repository.db.WithContext(ctx).Create(&dbPayment)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r SubscriptionRepository) updatePayment(
	ctx context.Context,
	payment subscription.Payment) error {
	if payment.IsDirty() == false {
		return nil
	}

	dbPayment := r.toPaymentModel(payment)
	if !payment.IsExists() {
		return subscription.ErrPaymentNotAlreadyExists
	}
	result := r.repository.db.WithContext(ctx).Save(&dbPayment)

	if result.Error != nil {
		return result.Error
	}

	payment.Clean()
	return nil
}

func (r SubscriptionRepository) deletePayment(ctx context.Context, payment subscription.Payment) error {
	result := r.repository.db.WithContext(ctx).Delete(&subscriptionPaymentModel{}, payment.Id())
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r SubscriptionRepository) createSubscriptionLabel(
	ctx context.Context,
	subscriptionId uuid.UUID,
	labelId uuid.UUID) error {
	dbLabelLink := subscriptionLabelModel{
		LabelId: labelId,
		SubId:   subscriptionId,
	}
	if result := r.repository.db.WithContext(ctx).Create(&dbLabelLink); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r SubscriptionRepository) updateSubscriptionLabel(
	ctx context.Context,
	subscriptionId uuid.UUID,
	labelId uuid.UUID) error {
	dbLabelLink := subscriptionLabelModel{
		LabelId: labelId,
		SubId:   subscriptionId,
	}
	if result := r.repository.db.WithContext(ctx).Save(&dbLabelLink); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r SubscriptionRepository) deleteSubscriptionLabel(
	ctx context.Context,
	subscriptionId uuid.UUID,
	labelId uuid.UUID) error {
	dbLabelLink := subscriptionLabelModel{
		LabelId: labelId,
		SubId:   subscriptionId,
	}
	if result := r.repository.db.WithContext(ctx).Delete(&dbLabelLink); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r SubscriptionRepository) createSubscriptionFamilyMember(
	ctx context.Context,
	subscriptionId uuid.UUID,
	familyMemberId uuid.UUID) error {
	dbFamilyMemberLink := subscriptionFamilyMemberModel{
		FamilyMemberId: familyMemberId,
		SubId:          subscriptionId,
	}

	if result := r.repository.db.WithContext(ctx).Create(&dbFamilyMemberLink); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r SubscriptionRepository) updateSubscriptionFamilyMember(
	ctx context.Context,
	subscriptionId uuid.UUID,
	familyMemberId uuid.UUID) error {
	dbFamilyMemberLink := subscriptionFamilyMemberModel{
		FamilyMemberId: familyMemberId,
		SubId:          subscriptionId,
	}

	if result := r.repository.db.WithContext(ctx).Save(&dbFamilyMemberLink); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r SubscriptionRepository) deleteSubscriptionFamilyMember(
	ctx context.Context,
	subscriptionId uuid.UUID,
	familyMemberId uuid.UUID) error {
	dbFamilyMemberLink := subscriptionFamilyMemberModel{
		FamilyMemberId: familyMemberId,
		SubId:          subscriptionId,
	}

	if result := r.repository.db.WithContext(ctx).Delete(&dbFamilyMemberLink); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r SubscriptionRepository) Delete(ctx context.Context, subscriptionId uuid.UUID) error {
	if result := r.repository.db.WithContext(ctx).Delete(&subscriptionModel{}, subscriptionId); result.Error != nil {
		return result.Error
	}
	return nil
}

func (r SubscriptionRepository) DeletePayment(ctx context.Context, paymentId uuid.UUID) error {
	if result := r.repository.db.WithContext(ctx).Delete(&subscriptionPaymentModel{}, paymentId); result.Error != nil {
		return result.Error
	}
	return nil
}
