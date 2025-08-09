package currency

import (
	"context"
	"time"

	"github.com/oleexo/subtracker/internal/domain/entity"
)

// Repository defines the interface for currency rate operations
type Repository interface {
	entity.Repository[Rate]

	GetRatesByDate(ctx context.Context, date time.Time) ([]Rate, error)
}
