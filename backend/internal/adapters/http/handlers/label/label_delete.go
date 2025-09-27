package label

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/pkg/ginx"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type DeleteEndpoint struct {
	handler ports.CommandHandler[command.DeleteLabelCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete label by LabelID
//	@Description	Permanently delete a label by its unique identifier
//	@Tags			labels
//	@Param			labelId	path	string	true	"Label LabelID (UUID format)"
//	@Success		204	"No Content - Label successfully deleted"
//	@Failure		400	{object}	HttpErrorResponse	"Bad Request - Invalid LabelID format"
//	@Failure		404	{object}	HttpErrorResponse	"Label not found"
//	@Failure		500	{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/labels/{labelId} [delete]
func (l DeleteEndpoint) Handle(c *gin.Context) {
	labelID, err := types.ParseLabelID(c.Param("labelId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: "invalid labelID format",
		})
		return
	}

	cmd := command.DeleteLabelCommand{
		LabelID: labelID,
	}

	r := l.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (l DeleteEndpoint) Pattern() []string {
	return []string{
		"/:labelId",
	}
}

func (l DeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (l DeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewDeleteEndpoint(handler ports.CommandHandler[command.DeleteLabelCommand, bool]) *DeleteEndpoint {
	return &DeleteEndpoint{
		handler: handler,
	}
}
