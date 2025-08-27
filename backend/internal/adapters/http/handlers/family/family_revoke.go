package family

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type FamilyRevokeEndpoint struct {
	handler ports.CommandHandler[command.RevokeMemberCommand, bool]
}

func NewFamilyRevokeEndpoint(handler ports.CommandHandler[command.RevokeMemberCommand, bool]) *FamilyRevokeEndpoint {
	return &FamilyRevokeEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Revoke family member
//	@Description	Revokes a member from the family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId		path		string				true	"Family ID (UUID format)"
//	@Param			familyMemberId	path		string				true	"Family Member ID (UUID format)"
//	@Success		204				{object}	nil					"Successfully revoked member"
//	@Failure		400				{object}	HttpErrorResponse	"Bad Request - Invalid input data"
//	@Failure		401				{object}	HttpErrorResponse	"Unauthorized - Invalid or missing authentication"
//	@Failure		404				{object}	HttpErrorResponse	"Family or member not found"
//	@Failure		500				{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/families/{familyId}/members/{familyMemberId}/revoke [post]
func (e FamilyRevokeEndpoint) Handle(c *gin.Context) {
	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	familyMemberId, err := uuid.Parse(c.Param("familyMemberId"))
	if err != nil {
		FromError(c, err)
		return
	}

	r := e.handler.Handle(c, command.RevokeMemberCommand{
		FamilyId:       familyId,
		FamilyMemberId: familyMemberId,
	})

	FromResult(c, r, WithNoContent[bool]())
}

func (e FamilyRevokeEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members/:familyMemberId/revoke",
	}
}

func (e FamilyRevokeEndpoint) Method() string {
	return http.MethodPost
}

func (e FamilyRevokeEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
