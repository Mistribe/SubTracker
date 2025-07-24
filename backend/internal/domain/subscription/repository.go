package subscription

import (
	"context"

	"github.com/google/uuid"

	"github.com/oleexo/subtracker/pkg/langext/option"
)

type Repository interface {
	Get(ctx context.Context, subscriptionId uuid.UUID) (option.Option[Subscription], error)
	GetAll(ctx context.Context, size, page int) ([]Subscription, error)
	GetAllCount(ctx context.Context) (int64, error)
	Save(ctx context.Context, subscription *Subscription) error
	Delete(ctx context.Context, subscriptionId uuid.UUID) error
	DeletePayment(ctx context.Context, paymentId uuid.UUID) error
}
