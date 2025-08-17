package application

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/application/auth"
	"github.com/oleexo/subtracker/internal/application/currency"
	"github.com/oleexo/subtracker/internal/application/family"
	"github.com/oleexo/subtracker/internal/application/label"
	"github.com/oleexo/subtracker/internal/application/provider"
	"github.com/oleexo/subtracker/internal/application/subscription"
	"github.com/oleexo/subtracker/internal/application/user"
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
