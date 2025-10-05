package family

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type DeleteEndpoint struct {
	handler ports.CommandHandler[command.DeleteFamilyCommand, bool]
}

func NewDeleteEndpoint(handler ports.CommandHandler[command.DeleteFamilyCommand, bool]) *DeleteEndpoint {
	return &DeleteEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Delete family by LabelID
//	@Description	Permanently delete a family and all its members
//	@Tags			family
//	@Param			familyId	path	string	true	"Family LabelID (UUID format)"
//	@Success		204			"No Content - Family successfully deleted"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid family LabelID format"
//	@Failure		404			{object}	HttpErrorResponse	"Family not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/family/{familyId} [delete]
func (e DeleteEndpoint) Handle(c *gin.Context) {
	familyID, err := types.ParseFamilyID(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	cmd := command.DeleteFamilyCommand{
		FamilyId: familyID,
	}

	r := e.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (e DeleteEndpoint) Pattern() []string {
	return []string{
		"/:familyId",
	}
}

func (e DeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (e DeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
