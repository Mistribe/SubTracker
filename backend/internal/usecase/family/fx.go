package family

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	command2 "github.com/mistribe/subtracker/internal/usecase/family/command"
	query2 "github.com/mistribe/subtracker/internal/usecase/family/query"
)

func Module() fx.Option {
	return fx.Module("app_family",
		fx.Provide(
			ports.AsQueryHandler[query2.FindUserFamilyQuery, family.Family](query2.NewFindOneQueryHandler),
			ports.AsQueryHandler[query2.SeeInvitationQuery, query2.SeeInvitationQueryResponse](query2.NewSeeInvitationQueryHandler),

			ports.AsCommandHandler[command2.CreateFamilyCommand, family.Family](command2.NewCreateFamilyCommandHandler),
			ports.AsCommandHandler[command2.UpdateFamilyCommand, family.Family](command2.NewUpdateFamilyCommandHandler),
			ports.AsCommandHandler[command2.PatchFamilyCommand, family.Family](command2.NewPatchFamilyCommandHandler),
			ports.AsCommandHandler[command2.DeleteFamilyCommand, bool](command2.NewDeleteFamilyCommandHandler),
			ports.AsCommandHandler[command2.InviteMemberCommand, command2.InviteMemberResponse](command2.NewInviteMemberCommandHandler),
			ports.AsCommandHandler[command2.AcceptInvitationCommand, bool](command2.NewAcceptInvitationCommandHandler),
			ports.AsCommandHandler[command2.DeclineInvitationCommand, bool](command2.NewDeclineInvitationCommandHandler),
			ports.AsCommandHandler[command2.RevokeMemberCommand, bool](command2.NewRevokeMemberCommandHandler),
			ports.AsCommandHandler[command2.CreateFamilyMemberCommand, family.Family](command2.NewCreateFamilyMemberCommandHandler),
			ports.AsCommandHandler[command2.UpdateFamilyMemberCommand, family.Family](command2.NewUpdateFamilyMemberCommandHandler),
			ports.AsCommandHandler[command2.DeleteFamilyMemberCommand, bool](command2.NewDeleteFamilyMemberCommandHandler),
		),
	)
}
