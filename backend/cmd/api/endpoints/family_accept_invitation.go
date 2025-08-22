package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
)

type FamilyAcceptInvitationEndpoint struct {
	handler core.CommandHandler[command.AcceptInvitationCommand, bool]
}

func NewFamilyAcceptInvitationEndpoint(handler core.CommandHandler[command.AcceptInvitationCommand, bool]) *FamilyAcceptInvitationEndpoint {
	return &FamilyAcceptInvitationEndpoint{handler: handler}
}

// FamilyAcceptInvitationRequest represents the request body for accepting a family invitation
type FamilyAcceptInvitationRequest struct {
	// Code received in the invitation
	InvitationCode string `json:"invitation_code" binding:"required" example:"123456"`
	// ID of the family member accepting the invitation
	FamilyMemberId string `json:"family_member_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
}

// Handle godoc
//
//	@Summary		Accept a family invitation
//	@Description	Accepts an invitation to join a family using the provided invitation code
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string							true	"Family ID (UUID format)"
//	@Param			request		body		FamilyAcceptInvitationRequest	true	"Invitation acceptance details"
//	@Success		204			{object}	nil								"Successfully accepted invitation"
//	@Failure		400			{object}	HttpErrorResponse				"Bad Request - Invalid input data"
//	@Failure		401			{object}	HttpErrorResponse				"Unauthorized - Invalid or missing authentication"
//	@Failure		404			{object}	HttpErrorResponse				"Family not found"
//	@Failure		500			{object}	HttpErrorResponse				"Internal Server Error"
//	@Router			/families/{familyId}/accept [post]
func (e FamilyAcceptInvitationEndpoint) Handle(c *gin.Context) {
	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	var model FamilyAcceptInvitationRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	familyMemberId, err := uuid.Parse(model.FamilyMemberId)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	r := e.handler.Handle(c, command.AcceptInvitationCommand{
		InvitationCode: model.InvitationCode,
		FamilyId:       familyId,
		FamilyMemberId: familyMemberId,
	})

	handleResponse(c, r, withNoContent[bool]())
}

func (e FamilyAcceptInvitationEndpoint) Pattern() []string {
	return []string{
		"/:familyId/accept",
	}
}

func (e FamilyAcceptInvitationEndpoint) Method() string {
	return http.MethodPost
}

func (e FamilyAcceptInvitationEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
