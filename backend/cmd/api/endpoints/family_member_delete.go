package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
)

type FamilyMemberDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteFamilyMemberCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete family member by ID
//	@Description	Permanently delete a family member from a family
//	@Tags			family
//	@Param			familyId	path	string	true	"Family ID (UUID format)"
//	@Param			id			path	string	true	"Family member ID (UUID format)"
//	@Success		204			"No Content - Family member successfully deleted"
//	@Failure		400			{object}	httpError	"Bad Request - Invalid ID format"
//	@Failure		404			{object}	httpError	"Family or family member not found"
//	@Failure		500			{object}	httpError	"Internal Server Error"
//	@Router			/families/{familyId}/members/{id} [delete]
func (f FamilyMemberDeleteEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "invalid id format",
		})
		return
	}

	familyIdParam := c.Param("familyId")
	if familyIdParam == "" {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "family id parameter is required",
		})
		return
	}

	familyId, err := uuid.Parse(familyIdParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "invalid family id format",
		})
	}

	cmd := command.DeleteFamilyMemberCommand{
		Id:       id,
		FamilyId: familyId,
	}

	r := f.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[bool]())
}

func (f FamilyMemberDeleteEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members/:id",
	}
}

func (f FamilyMemberDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (f FamilyMemberDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberDeleteEndpoint(handler core.CommandHandler[command.DeleteFamilyMemberCommand, bool]) *FamilyMemberDeleteEndpoint {
	return &FamilyMemberDeleteEndpoint{
		handler: handler,
	}
}
