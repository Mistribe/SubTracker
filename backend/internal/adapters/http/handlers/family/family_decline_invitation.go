package family

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type DeclineInvitationEndpoint struct {
	handler ports.CommandHandler[command.DeclineInvitationCommand, bool]
}

func NewDeclineEndpoint(handler ports.CommandHandler[command.DeclineInvitationCommand, bool]) *DeclineInvitationEndpoint {
	return &DeclineInvitationEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Decline family invitation
//	@Description	Endpoint to decline an invitation to join a family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path	string								true	"Family LabelID"
//	@Param			request		body	dto.FamilyDeclineInvitationRequest	true	"Decline invitation request"
//	@Success		204			"No Content"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request"
//	@Router			/family/{familyId}/decline [post]
func (e DeclineInvitationEndpoint) Handle(c *gin.Context) {
	familyId, err := types.ParseFamilyID(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	var request dto.FamilyDeclineInvitationRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		FromError(c, err)
		return
	}

	familyMemberId, err := types.ParseFamilyMemberID(request.FamilyMemberId)
	if err != nil {
		FromError(c, err)
		return
	}

	r := e.handler.Handle(c, command.DeclineInvitationCommand{
		InvitationCode: request.InvitationCode,
		FamilyId:       familyId,
		FamilyMemberId: familyMemberId,
	})

	FromResult(c, r, WithNoContent[bool]())
}

func (e DeclineInvitationEndpoint) Pattern() []string {
	return []string{
		"/:familyId/decline",
	}
}

func (e DeclineInvitationEndpoint) Method() string {
	return http.MethodPost
}

func (e DeclineInvitationEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
