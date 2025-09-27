package family

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/types"
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
//	@Router			/family/{familyId}/members/{familyMemberId} [delete]
func (f MemberDeleteEndpoint) Handle(c *gin.Context) {

	familyMemberID, err := types.ParseFamilyMemberID(c.Param("familyMemberId"))
	if err != nil {
		FromError(c, err)
		return
	}

	familyID, err := types.ParseFamilyID(c.Param("familyID"))
	if err != nil {
		FromError(c, err)
		return
	}

	cmd := command.DeleteFamilyMemberCommand{
		FamilyMemberID: familyMemberID,
		FamilyID:       familyID,
	}

	r := f.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (f MemberDeleteEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members/:familyMemberId",
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
