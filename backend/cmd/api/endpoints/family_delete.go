package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FamilyDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteFamilyCommand, result.Unit]
}

func NewFamilyDeleteEndpoint(handler core.CommandHandler[command.DeleteFamilyCommand, result.Unit]) *FamilyDeleteEndpoint {
	return &FamilyDeleteEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Delete family by ID
//	@Description	Delete family by ID
//	@Tags			family
//	@Param			familyId	path	string	true	"Family ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	httpError
//	@Failure		404	{object}	httpError
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
	handleResponse(c, r, withNoContent[result.Unit]())
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
