package family

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/ginx"
)

type FamilyInviteEndpoint struct {
	handler ports.CommandHandler[command.InviteMemberCommand, command.InviteMemberResponse]
}

func NewFamilyInviteEndpoint(handler ports.CommandHandler[command.InviteMemberCommand, command.InviteMemberResponse]) *FamilyInviteEndpoint {
	return &FamilyInviteEndpoint{handler: handler}
}

// FamilyInviteRequest represents the request body for inviting a family member
type FamilyInviteRequest struct {
	// Email of the invited member
	Email *string `json:"email,omitempty"`
	// ID of the family member to be invited
	FamilyMemberId string `json:"family_member_id" binding:"required"`
	// Name of the invited member
	Name *string `json:"name,omitempty"`
	// Type of the member (adult or kid)
	Type *string `json:"type,omitempty" enums:"adult,kid"`
}

type FamilyInviteResponse struct {
	Code           string `json:"code" binding:"required" example:"123456"`
	FamilyId       string `json:"family_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	FamilyMemberId string `json:"family_member_id" binding:"required" example:"123e4567-e89b-12d3-a456-426614174001"`
}

// Handle godoc
//
//	@Summary		Invite a new member to the family
//	@Description	Creates an invitation for a new member to join the family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string					true	"Family ID (UUID format)"
//	@Param			request		body		FamilyInviteRequest		true	"Invitation details including email, name, member ID and type (adult/kid)"
//	@Success		200			{object}	FamilyInviteResponse	"Successfully created invitation with code and IDs"
//	@Failure		400			{object}	HttpErrorResponse		"Bad Request - Invalid input data"
//	@Failure		401			{object}	HttpErrorResponse		"Unauthorized - Invalid or missing authentication"
//	@Failure		404			{object}	HttpErrorResponse		"Family not found"
//	@Failure		500			{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/families/{familyId}/invite [post]
func (e FamilyInviteEndpoint) Handle(c *gin.Context) {
	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	var model FamilyInviteRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}
	familyMemberId, err := uuid.Parse(model.FamilyMemberId)
	if err != nil {
		FromError(c, err)
		return
	}

	var memberType *family.MemberType
	if model.Type != nil {
		var mt family.MemberType
		mt, err = family.ParseMemberType(*model.Type)
		if err != nil {
			c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
				Message: err.Error(),
			})
			return
		}
		memberType = &mt
	}

	r := e.handler.Handle(c, command.InviteMemberCommand{
		Email:          model.Email,
		Name:           model.Name,
		Type:           memberType,
		FamilyId:       familyId,
		FamilyMemberId: familyMemberId,
	})
	FromResult(c, r, WithMapping(func(response command.InviteMemberResponse) any {
		return FamilyInviteResponse{
			Code:           response.Code,
			FamilyId:       response.FamilyId.String(),
			FamilyMemberId: response.FamilyMemberId.String(),
		}
	}))
}

func (e FamilyInviteEndpoint) Pattern() []string {
	return []string{
		"/:familyId/invite",
	}
}

func (e FamilyInviteEndpoint) Method() string {
	return http.MethodPost
}

func (e FamilyInviteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
