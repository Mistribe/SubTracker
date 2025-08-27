package label

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/pkg/ginx"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type LabelDeleteEndpoint struct {
	handler ports.CommandHandler[command.DeleteLabelCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete label by ID
//	@Description	Permanently delete a label by its unique identifier
//	@Tags			label
//	@Param			id	path	string	true	"Label ID (UUID format)"
//	@Success		204	"No Content - Label successfully deleted"
//	@Failure		400	{object}	HttpErrorResponse	"Bad Request - Invalid ID format"
//	@Failure		404	{object}	HttpErrorResponse	"Label not found"
//	@Failure		500	{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/labels/{id} [delete]
func (l LabelDeleteEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: "invalid id format",
		})
		return
	}

	cmd := command.DeleteLabelCommand{
		Id: id,
	}

	r := l.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
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

func NewLabelDeleteEndpoint(handler ports.CommandHandler[command.DeleteLabelCommand, bool]) *LabelDeleteEndpoint {
	return &LabelDeleteEndpoint{
		handler: handler,
	}
}
