package usecase

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/currency"
	"github.com/mistribe/subtracker/internal/usecase/family"
	"github.com/mistribe/subtracker/internal/usecase/label"
	"github.com/mistribe/subtracker/internal/usecase/provider"
	"github.com/mistribe/subtracker/internal/usecase/subscription"
	"github.com/mistribe/subtracker/internal/usecase/user"
)

func BuildApplicationModules() []fx.Option {
	return []fx.Option{
		currency.Module(),
		family.Module(),
		authentication.Module(),
		label.Module(),
		provider.Module(),
		subscription.Module(),
		user.Module(),
	}
}
