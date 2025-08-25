package application

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/application/auth"
	"github.com/mistribe/subtracker/internal/application/currency"
	"github.com/mistribe/subtracker/internal/application/family"
	"github.com/mistribe/subtracker/internal/application/label"
	"github.com/mistribe/subtracker/internal/application/provider"
	"github.com/mistribe/subtracker/internal/application/subscription"
	"github.com/mistribe/subtracker/internal/application/user"
)

func BuildApplicationModules() []fx.Option {
	return []fx.Option{
		currency.Module(),
		family.Module(),
		auth.Module(),
		label.Module(),
		provider.Module(),
		subscription.Module(),
		user.Module(),
	}
}
