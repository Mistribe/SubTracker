package user

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/currency"
)

type Service interface {
	GetPreferredCurrency(ctx context.Context, userId string) currency.Unit
}
