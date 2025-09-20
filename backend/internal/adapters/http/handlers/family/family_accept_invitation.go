package family

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type AcceptInvitationEndpoint struct {
	handler ports.CommandHandler[command.AcceptInvitationCommand, bool]
}

func NewAcceptInvitationEndpoint(handler ports.CommandHandler[command.AcceptInvitationCommand, bool]) *AcceptInvitationEndpoint {
	return &AcceptInvitationEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Accept a family invitation
//	@Description	Accepts an invitation to join a family using the provided invitation code
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string								true	"Family ID (UUID format)"
//	@Param			request		body		dto.FamilyAcceptInvitationRequest	true	"Invitation acceptance details"
//	@Success		204			{object}	nil									"Successfully accepted invitation"
//	@Failure		400			{object}	HttpErrorResponse					"Bad Request - Invalid input data"
//	@Failure		401			{object}	HttpErrorResponse					"Unauthorized - Invalid or missing authentication"
//	@Failure		404			{object}	HttpErrorResponse					"Family not found"
//	@Failure		500			{object}	HttpErrorResponse					"Internal Server Error"
//	@Router			/family/{familyId}/accept [post]
func (e AcceptInvitationEndpoint) Handle(c *gin.Context) {
	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	var model dto.FamilyAcceptInvitationRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}
	familyMemberId, err := uuid.Parse(model.FamilyMemberId)
	if err != nil {
		FromError(c, err)
		return
	}

	r := e.handler.Handle(c, command.AcceptInvitationCommand{
		InvitationCode: model.InvitationCode,
		FamilyId:       familyId,
		FamilyMemberId: familyMemberId,
	})

	FromResult(c, r, WithNoContent[bool]())
}

func (e AcceptInvitationEndpoint) Pattern() []string {
	return []string{
		"/:familyId/accept",
	}
}

func (e AcceptInvitationEndpoint) Method() string {
	return http.MethodPost
}

func (e AcceptInvitationEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
