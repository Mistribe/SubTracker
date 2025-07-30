package endpoints

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/command"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type LabelUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateLabelCommand, label.Label]
}

type updateLabelModel struct {
	Name      string     `json:"name" binding:"required"`
	Color     string     `json:"color" binding:"required"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m updateLabelModel) Command(id uuid.UUID) (command.UpdateLabelCommand, error) {
	updatedAt := option.None[time.Time]()

	if m.UpdatedAt != nil {
		updatedAt = option.Some(*m.UpdatedAt)
	}

	return command.UpdateLabelCommand{
		Id:        id,
		Name:      m.Name,
		Color:     strings.ToUpper(m.Color),
		UpdatedAt: updatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Update label by ID
//	@Description	Update an existing label's name and color by its unique identifier
//	@Tags			label
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Label ID (UUID format)"
//	@Param			label	body		updateLabelModel	true	"Updated label data"
//	@Success		200		{object}	labelModel			"Successfully updated label"
//	@Failure		400		{object}	httpError			"Bad Request - Invalid ID format or input data"
//	@Failure		404		{object}	httpError			"Label not found"
//	@Failure		500		{object}	httpError			"Internal Server Error"
//	@Router			/labels/{id} [put]
func (l LabelUpdateEndpoint) Handle(c *gin.Context) {
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

	var model updateLabelModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}
	cmd, err := model.Command(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	r := l.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withMapping[label.Label](func(lbl label.Label) any {
			return newLabelModel(lbl)
		}))
}

func (l LabelUpdateEndpoint) Pattern() []string {
	return []string{
		"/:id",
	}
}

func (l LabelUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (l LabelUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewLabelUpdateEndpoint(handler core.CommandHandler[command.UpdateLabelCommand, label.Label]) *LabelUpdateEndpoint {
	return &LabelUpdateEndpoint{
		handler: handler,
	}
}
