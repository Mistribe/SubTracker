package label

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/pkg/ginx"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

type UpdateEndpoint struct {
	handler ports.CommandHandler[command.UpdateLabelCommand, label.Label]
}

func updateLabelRequestToCommand(m dto.UpdateLabelRequest, id uuid.UUID) (command.UpdateLabelCommand, error) {
	updatedAt := option.None[time.Time]()

	if m.UpdatedAt != nil {
		updatedAt = option.Some(*m.UpdatedAt)
	}

	return command.UpdateLabelCommand{
		LabelID:   id,
		Name:      m.Name,
		Color:     strings.ToUpper(m.Color),
		UpdatedAt: updatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Update label by LabelID
//	@Description	Update an existing label's name and color by its unique identifier
//	@Tags			labels
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string					true	"Label LabelID (UUID format)"
//	@Param			label	body		dto.UpdateLabelRequest	true	"Updated label data"
//	@Success		200		{object}	dto.LabelModel			"Successfully updated label"
//	@Failure		400		{object}	HttpErrorResponse		"Bad Request - Invalid LabelID format or input data"
//	@Failure		404		{object}	HttpErrorResponse		"Label not found"
//	@Failure		500		{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/labels/{id} [put]
func (l UpdateEndpoint) Handle(c *gin.Context) {
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

	var model dto.UpdateLabelRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	cmd, err := updateLabelRequestToCommand(model, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	r := l.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithMapping[label.Label](func(lbl label.Label) any {
			return dto.NewLabelModel(lbl)
		}))
}

func (l UpdateEndpoint) Pattern() []string {
	return []string{
		"/:id",
	}
}

func (l UpdateEndpoint) Method() string {
	return http.MethodPut
}

func (l UpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewUpdateEndpoint(handler ports.CommandHandler[command.UpdateLabelCommand, label.Label]) *UpdateEndpoint {
	return &UpdateEndpoint{
		handler: handler,
	}
}
