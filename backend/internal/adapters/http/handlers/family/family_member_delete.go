package family

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type MemberDeleteEndpoint struct {
	handler ports.CommandHandler[command.DeleteFamilyMemberCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete family member by LabelID
//	@Description	Permanently delete a family member from a family
//	@Tags			family
//	@Param			familyId	path	string	true	"Family LabelID (UUID format)"
//	@Param			id			path	string	true	"Family member LabelID (UUID format)"
//	@Success		204			"No Content - Family member successfully deleted"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid LabelID format"
//	@Failure		404			{object}	HttpErrorResponse	"Family or family member not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/family/{familyId}/members/{id} [delete]
func (f MemberDeleteEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		FromError(c, errors.New("id parameter is required"))
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		FromError(c, err)
		return
	}

	familyIdParam := c.Param("familyId")
	if familyIdParam == "" {
		FromError(c, errors.New("family id parameter is required"))
		return
	}

	familyId, err := uuid.Parse(familyIdParam)
	if err != nil {
		FromError(c, err)
		return
	}

	cmd := command.DeleteFamilyMemberCommand{
		FamilyMemberID: id,
		FamilyID:       familyId,
	}

	r := f.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (f MemberDeleteEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members/:id",
	}
}

func (f MemberDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (f MemberDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewMemberDeleteEndpoint(handler ports.CommandHandler[command.DeleteFamilyMemberCommand, bool]) *MemberDeleteEndpoint {
	return &MemberDeleteEndpoint{
		handler: handler,
	}
}
