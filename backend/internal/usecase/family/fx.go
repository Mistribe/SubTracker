package family

import (
	"go.uber.org/fx"

	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/internal/usecase/family/query"
)

func Module() fx.Option {
	return fx.Module("app_family",
		fx.Provide(
			ports.AsQueryHandler[query.FindUserFamilyQuery, query.FindUserFamilyQueryResponse](query.NewFindOneQueryHandler),
			ports.AsQueryHandler[query.SeeInvitationQuery, query.SeeInvitationQueryResponse](query.NewSeeInvitationQueryHandler),
			ports.AsQueryHandler[query.GetQuotaUsage, []billing.EffectiveEntitlement](query.NewGetQuotaUsageHandler),

			ports.AsCommandHandler[command.CreateFamilyCommand, family.Family](command.NewCreateFamilyCommandHandler),
			ports.AsCommandHandler[command.UpdateFamilyCommand, family.Family](command.NewUpdateFamilyCommandHandler),
			ports.AsCommandHandler[command.DeleteFamilyCommand, bool](command.NewDeleteFamilyCommandHandler),
			ports.AsCommandHandler[command.InviteMemberCommand, command.InviteMemberResponse](command.NewInviteMemberCommandHandler),
			ports.AsCommandHandler[command.AcceptInvitationCommand, bool](command.NewAcceptInvitationCommandHandler),
			ports.AsCommandHandler[command.DeclineInvitationCommand, bool](command.NewDeclineInvitationCommandHandler),
			ports.AsCommandHandler[command.RevokeMemberCommand, bool](command.NewRevokeMemberCommandHandler),
			ports.AsCommandHandler[command.CreateFamilyMemberCommand, family.Family](command.NewCreateFamilyMemberCommandHandler),
			ports.AsCommandHandler[command.UpdateFamilyMemberCommand, family.Family](command.NewUpdateFamilyMemberCommandHandler),
			ports.AsCommandHandler[command.DeleteFamilyMemberCommand, bool](command.NewDeleteFamilyMemberCommandHandler),
		),
	)
}
