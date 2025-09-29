package label

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
)

func Module() fx.Option {
	return fx.Module("app_label",
		fx.Provide(
			ports.AsQueryHandler[query.FindAllQuery, shared.PaginatedResponse[label.Label]](query.NewFindAllQueryHandler),
			ports.AsQueryHandler[query.FindOneQuery, label.Label](query.NewFindOneQueryHandler),
			ports.AsQueryHandler[query.DefaultLabelQuery, []label.Label](query.NewDefaultLabelQueryHandler),
			ports.AsQueryHandler[query.GetQuotaUsageHandler, billing.EffectiveEntitlement](query.NewGetQuotaUsageHandler),

			ports.AsCommandHandler[command.CreateLabelCommand, label.Label](command.NewCreateLabelCommandHandler),
			ports.AsCommandHandler[command.UpdateLabelCommand, label.Label](command.NewUpdateLabelCommandHandler),
			ports.AsCommandHandler[command.DeleteLabelCommand, bool](command.NewDeleteLabelCommandHandler),
		))
}
