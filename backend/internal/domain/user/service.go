package user

import (
	"context"

	"github.com/oleexo/subtracker/internal/domain/currency"
)

type Service interface {
	GetPreferredCurrency(ctx context.Context, userId string) currency.Unit
}
