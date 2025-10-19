package family

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type InviteEndpoint struct {
	handler ports.CommandHandler[command.InviteMemberCommand, command.InviteMemberResponse]
}

func NewInviteEndpoint(handler ports.CommandHandler[command.InviteMemberCommand, command.InviteMemberResponse]) *InviteEndpoint {
	return &InviteEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Invite a new member to the family
//	@Description	Creates an invitation for a new member to join the family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string						true	"Family LabelID (UUID format)"
//	@Param			request		body		dto.FamilyInviteRequest		true	"Invitation details including email, name, member LabelID and type (adult/kid)"
//	@Success		200			{object}	dto.FamilyInviteResponse	"Successfully created invitation with code and IDs"
//	@Failure		400			{object}	HttpErrorResponse			"Bad Request - Invalid input data"
//	@Failure		401			{object}	HttpErrorResponse			"Unauthorized - Invalid or missing authentication"
//	@Failure		404			{object}	HttpErrorResponse			"Family not found"
//	@Failure		500			{object}	HttpErrorResponse			"Internal Server Error"
//	@Router			/family/{familyId}/invite [post]
func (e InviteEndpoint) Handle(c *gin.Context) {
	familyId, err := types.ParseFamilyID(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	var model dto.FamilyInviteRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}
	familyMemberId, err := types.ParseFamilyMemberID(model.FamilyMemberId)
	if err != nil {
		FromError(c, err)
		return
	}

	var memberType *family.MemberType
	if model.Type != nil {
		var mt family.MemberType
		mt, err = family.ParseMemberType(*model.Type)
		if err != nil {
			FromError(c, err)
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
		return dto.FamilyInviteResponse{
			Code:           response.Code,
			FamilyId:       response.FamilyId.String(),
			FamilyMemberId: response.FamilyMemberId.String(),
		}
	}))
}

func (e InviteEndpoint) Pattern() []string {
	return []string{
		"/:familyId/invite",
	}
}

func (e InviteEndpoint) Method() string {
	return http.MethodPost
}

func (e InviteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
