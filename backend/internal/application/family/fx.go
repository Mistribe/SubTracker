package family

import (
	"go.uber.org/fx"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/application/family/query"
	"github.com/oleexo/subtracker/internal/domain/family"
)

func Module() fx.Option {
	return fx.Module("app_family",
		fx.Provide(
			core.AsQueryHandler[query.FindUserFamilyQuery, family.Family](query.NewFindOneQueryHandler),

			core.AsCommandHandler[command.CreateFamilyCommand, family.Family](command.NewCreateFamilyCommandHandler),
			core.AsCommandHandler[command.UpdateFamilyCommand, family.Family](command.NewUpdateFamilyCommandHandler),
			core.AsCommandHandler[command.PatchFamilyCommand, family.Family](command.NewPatchFamilyCommandHandler),
			core.AsCommandHandler[command.DeleteFamilyCommand, bool](command.NewDeleteFamilyCommandHandler),
			core.AsCommandHandler[command.InviteMemberCommand, command.InviteMemberResponse](command.NewInviteMemberCommandHandler),
			core.AsCommandHandler[command.CreateFamilyMemberCommand, family.Family](command.NewCreateFamilyMemberCommandHandler),
			core.AsCommandHandler[command.UpdateFamilyMemberCommand, family.Family](command.NewUpdateFamilyMemberCommandHandler),
			core.AsCommandHandler[command.DeleteFamilyMemberCommand, bool](command.NewDeleteFamilyMemberCommandHandler),
		),
	)
}
