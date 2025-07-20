package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FamilyMemberDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteFamilyMemberCommand, result.Unit]
}

// Handle godoc
//
//	@Summary		Delete family member by ID
//	@Description	Delete family member by ID
//	@Tags			family
//	@Param			id	path	string	true	"Family member ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	httpError
//	@Failure		404	{object}	httpError
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
	handleResponse(c, r, withNoContent[result.Unit]())
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

func NewFamilyMemberDeleteEndpoint(handler core.CommandHandler[command.DeleteFamilyMemberCommand, result.Unit]) *FamilyMemberDeleteEndpoint {
	return &FamilyMemberDeleteEndpoint{
		handler: handler,
	}
}
