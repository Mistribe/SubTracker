package label

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/command"
	"github.com/oleexo/subtracker/internal/application/label/query"
	"github.com/oleexo/subtracker/internal/domain/label"
)

func Module() fx.Option {
	return fx.Module("app_label",
		fx.Provide(
			core.AsQueryHandler[query.FindAllQuery, core.PaginatedResponse[label.Label]](query.NewFindAllQueryHandler),
			core.AsQueryHandler[query.FindOneQuery, label.Label](query.NewFindOneQueryHandler),
			core.AsQueryHandler[query.DefaultLabelQuery, []label.Label](query.NewDefaultLabelQueryHandler),

			core.AsCommandHandler[command.CreateLabelCommand, label.Label](command.NewCreateLabelCommandHandler),
			core.AsCommandHandler[command.UpdateLabelCommand, label.Label](command.NewUpdateLabelCommandHandler),
			core.AsCommandHandler[command.DeleteLabelCommand, bool](command.NewDeleteLabelCommandHandler),
		))
}
