package family

import (
	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/router/ginfx"
	"github.com/mistribe/subtracker/internal/adapters/http/router/middlewares"
)

type EndpointGroup struct {
	routes      []ginfx.Endpoint
	middlewares []gin.HandlerFunc
}

func NewEndpointGroup(
	familyCreateEndpoint *CreateEndpoint,
	familyUpdateEndpoint *UpdateEndpoint,
	familyDeleteEndpoint *DeleteEndpoint,
	familyInviteEndpoint *InviteEndpoint,
	familyAcceptInvitationEndpoint *AcceptInvitationEndpoint,
	familyDeclineEndpoint *DeclineInvitationEndpoint,
	familyRevokeEndpoint *RevokeEndpoint,
	familySeeInvitationEndpoint *SeeInvitationEndpoint,
	familyMemberCreateEndpoint *MemberCreateEndpoint,
	familyMemberUpdateEndpoint *MemberUpdateEndpoint,
	familyMemberDeleteEndpoint *MemberDeleteEndpoint,
	familyGetEndpoint *GetEndpoint,
	authenticationMiddleware *middlewares.AuthenticationMiddleware) *EndpointGroup {
	return &EndpointGroup{
		routes: []ginfx.Endpoint{
			familyCreateEndpoint,
			familyUpdateEndpoint,
			familyDeleteEndpoint,
			familyInviteEndpoint,
			familyAcceptInvitationEndpoint,
			familyDeclineEndpoint,
			familyRevokeEndpoint,
			familySeeInvitationEndpoint,
			familyMemberCreateEndpoint,
			familyMemberUpdateEndpoint,
			familyMemberDeleteEndpoint,
			familyGetEndpoint,
		},
		middlewares: []gin.HandlerFunc{
			authenticationMiddleware.Middleware(),
		},
	}
}

func (g EndpointGroup) Prefix() string {
	return "/family"
}

func (g EndpointGroup) Routes() []ginfx.Endpoint {
	return g.routes
}

func (g EndpointGroup) Middlewares() []gin.HandlerFunc {
	return g.middlewares
}
