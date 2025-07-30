package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
)

type FamilyDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteFamilyCommand, bool]
}

func NewFamilyDeleteEndpoint(handler core.CommandHandler[command.DeleteFamilyCommand, bool]) *FamilyDeleteEndpoint {
	return &FamilyDeleteEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Delete family by ID
//	@Description	Permanently delete a family and all its members
//	@Tags			family
//	@Param			familyId	path	string	true	"Family ID (UUID format)"
//	@Success		204			"No Content - Family successfully deleted"
//	@Failure		400			{object}	httpError	"Bad Request - Invalid family ID format"
//	@Failure		404			{object}	httpError	"Family not found"
//	@Failure		500			{object}	httpError	"Internal Server Error"
//	@Router			/families/{familyId} [delete]
func (e FamilyDeleteEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("familyId")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "familyId parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "invalid familyId format",
		})
		return
	}

	cmd := command.DeleteFamilyCommand{
		FamilyId: id,
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[bool]())
}

func (e FamilyDeleteEndpoint) Pattern() []string {
	return []string{
		"/:familyId",
	}
}

func (e FamilyDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (e FamilyDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
