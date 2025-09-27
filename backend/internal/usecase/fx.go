package usecase

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/usecase/account"
	"github.com/mistribe/subtracker/internal/usecase/currency"
	"github.com/mistribe/subtracker/internal/usecase/family"
	"github.com/mistribe/subtracker/internal/usecase/label"
	"github.com/mistribe/subtracker/internal/usecase/provider"
	"github.com/mistribe/subtracker/internal/usecase/subscription"
)

func BuildApplicationModules() []fx.Option {
	return []fx.Option{
		currency.Module(),
		family.Module(),
		label.Module(),
		provider.Module(),
		subscription.Module(),
		account.Module(),
	}
}
