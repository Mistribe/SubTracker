package ports

import (
	"context"

	"github.com/mistribe/subtracker/internal/domain/currency"
	"github.com/mistribe/subtracker/internal/domain/types"
)

type AccountService interface {
	GetPreferredCurrency(ctx context.Context, userId types.UserID) currency.Unit
}
