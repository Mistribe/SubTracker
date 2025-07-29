package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/command"
)

type LabelDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteLabelCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete label by ID
//	@Description	Delete label by ID
//	@Tags			label
//	@Param			id	path	string	true	"Label ID"
//	@Success		204	"No Content"
//	@Failure		400	{object}	httpError
//	@Failure		404	{object}	httpError
//	@Router			/labels/{id} [delete]
func (l LabelDeleteEndpoint) Handle(c *gin.Context) {
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

	cmd := command.DeleteLabelCommand{
		Id: id,
	}

	r := l.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[bool]())
}

func (l LabelDeleteEndpoint) Pattern() []string {
	return []string{
		"/:id",
	}
}

func (l LabelDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (l LabelDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewLabelDeleteEndpoint(handler core.CommandHandler[command.DeleteLabelCommand, bool]) *LabelDeleteEndpoint {
	return &LabelDeleteEndpoint{
		handler: handler,
	}
}
