package label

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	command2 "github.com/mistribe/subtracker/internal/usecase/label/command"
	query2 "github.com/mistribe/subtracker/internal/usecase/label/query"
)

func Module() fx.Option {
	return fx.Module("app_label",
		fx.Provide(
			ports.AsQueryHandler[query2.FindAllQuery, shared.PaginatedResponse[label.Label]](query2.NewFindAllQueryHandler),
			ports.AsQueryHandler[query2.FindOneQuery, label.Label](query2.NewFindOneQueryHandler),
			ports.AsQueryHandler[query2.DefaultLabelQuery, []label.Label](query2.NewDefaultLabelQueryHandler),

			ports.AsCommandHandler[command2.CreateLabelCommand, label.Label](command2.NewCreateLabelCommandHandler),
			ports.AsCommandHandler[command2.UpdateLabelCommand, label.Label](command2.NewUpdateLabelCommandHandler),
			ports.AsCommandHandler[command2.DeleteLabelCommand, bool](command2.NewDeleteLabelCommandHandler),
		))
}
