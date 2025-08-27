package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/currency"
)

type UserService interface {
	GetPreferredCurrency(ctx context.Context, userId string) currency.Unit
}
