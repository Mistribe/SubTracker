package family

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type FamilyDeclineInvitationEndpoint struct {
	handler ports.CommandHandler[command.DeclineInvitationCommand, bool]
}

func NewFamilyDeclineEndpoint(handler ports.CommandHandler[command.DeclineInvitationCommand, bool]) *FamilyDeclineInvitationEndpoint {
	return &FamilyDeclineInvitationEndpoint{handler: handler}
}

type FamilyDeclineInvitationRequest struct {
	// Code received in the invitation
	InvitationCode string `json:"invitation_code" binding:"required" example:"123456"`
	// ID of the family member accepting the invitation
	FamilyMemberId string `json:"family_member_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// Handle godoc
//
//	@Summary		Decline family invitation
//	@Description	Endpoint to decline an invitation to join a family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path	string							true	"Family ID"
//	@Param			request		body	FamilyDeclineInvitationRequest	true	"Decline invitation request"
//	@Success		204			"No Content"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request"
//	@Router			/families/{familyId}/decline [post]
func (e FamilyDeclineInvitationEndpoint) Handle(c *gin.Context) {
	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	var request FamilyDeclineInvitationRequest
	if err = c.ShouldBindJSON(&request); err != nil {
		FromError(c, err)
		return
	}

	familyMemberId, err := uuid.Parse(request.FamilyMemberId)
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

func (e FamilyDeclineInvitationEndpoint) Pattern() []string {
	return []string{
		"/:familyId/decline",
	}
}

func (e FamilyDeclineInvitationEndpoint) Method() string {
	return http.MethodPost
}

func (e FamilyDeclineInvitationEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
