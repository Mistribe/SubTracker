package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oleexo/subtracker/internal/domain/entity"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type SubscriptionRepository struct {
	repository *DatabaseContext
}

func NewSubscriptionRepository(repository *DatabaseContext) *SubscriptionRepository {
	return &SubscriptionRepository{
		repository: repository,
	}
}

func (r SubscriptionRepository) GetById(ctx context.Context, id uuid.UUID) (subscription.Subscription, error) {
	var model subscriptionSqlModel
	result := r.repository.db.WithContext(ctx).
		Preload("Owner").
		Preload("Payer").
		First(&model, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return newSubscription(model), nil
}

func (r SubscriptionRepository) GetAll(ctx context.Context,
	parameters entity.QueryParameters) ([]subscription.Subscription, error) {
	var models []subscriptionSqlModel
	query := r.repository.db.WithContext(ctx).
		Preload("Owner").
		Preload("Payer").
		Preload("ServiceUsers")

	if parameters.Offset > 0 {
		query = query.Offset(parameters.Offset)
	}
	if parameters.Limit > 0 {
		query = query.Limit(parameters.Limit)
	}

	result := query.Find(&models)
	if result.Error != nil {
		return nil, result.Error
	}

	subscriptions := make([]subscription.Subscription, len(models))
	for _, model := range models {
		subscriptions = append(subscriptions, newSubscription(model))
	}
	return subscriptions, nil
}

func (r SubscriptionRepository) GetAllCount(ctx context.Context) (int64, error) {
	var count int64
	if result := r.repository.db.WithContext(ctx).
		Model(&subscriptionSqlModel{}).
		Count(&count); result.Error != nil {
		return 0, result.Error
	}

	return count, nil
}

func (r SubscriptionRepository) Save(ctx context.Context, dirtySubscription subscription.Subscription) error {
	if dirtySubscription.IsDirty() == false {
		return nil
	}

	dbSubscription := newSubscriptionSqlModel(dirtySubscription)
	var result *gorm.DB
	if dirtySubscription.IsExists() {
		result = r.repository.db.WithContext(ctx).
			Omit("Owner").
			Omit("Payer").
			Omit("ServiceUsers").
			Save(&dbSubscription)
	} else {
		result = r.repository.db.WithContext(ctx).
			Omit("Owner").
			Omit("Payer").
			Omit("ServiceUsers").
			Create(&dbSubscription)
	}
	if result.Error != nil {
		return result.Error
	}

	if dirtySubscription.ServiceUsers().HasChanges() {
		if err := saveTrackedSlice(ctx, r.repository.db,
			dirtySubscription.ServiceUsers(),
			func(serviceUser uuid.UUID) subscriptionServiceUserModel {
				return newSubscriptionServiceUserModel(serviceUser, dirtySubscription.Id())
			}); err != nil {
			return err
		}
	}

	dirtySubscription.Clean()
	return nil
}

func (r SubscriptionRepository) Delete(ctx context.Context, subscriptionId uuid.UUID) (bool, error) {
	result := r.repository.db.WithContext(ctx).Delete(&subscriptionSqlModel{}, subscriptionId)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (r SubscriptionRepository) Exists(ctx context.Context, ids ...uuid.UUID) (bool, error) {
	var count int64
	if result := r.repository.db.WithContext(ctx).
		Model(&subscriptionSqlModel{}).
		Where("id IN ?", ids).
		Count(&count); result.Error != nil {
		return false, result.Error
	}
	return count == int64(len(ids)), nil
}
