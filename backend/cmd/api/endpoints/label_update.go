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
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type LabelUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateLabelCommand, label.Label]
}

type updateLabelModel struct {
	Name      string     `json:"name"`
	Color     string     `json:"color"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m updateLabelModel) Command(id uuid.UUID) result.Result[command.UpdateLabelCommand] {
	updatedAt := option.None[time.Time]()

	if m.UpdatedAt != nil {
		updatedAt = option.Some(*m.UpdatedAt)
	}

	return result.Success(command.UpdateLabelCommand{
		Id:        id,
		Name:      m.Name,
		Color:     strings.ToUpper(m.Color),
		UpdatedAt: updatedAt,
	})
}

// Handle godoc
//
//	@Summary		Update label by ID
//	@Description	Update label by ID
//	@Tags			label
//	@Accept			json
//	@Produce		json
//	@Param			id		path		string				true	"Label ID"
//	@Param			label	body		updateLabelModel	true	"Label data"
//	@Success		200		{object}	labelModel
//	@Failure		400		{object}	httpError
//	@Failure		404		{object}	httpError
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

	result.Match[command.UpdateLabelCommand, result.Unit](model.Command(id),
		func(cmd command.UpdateLabelCommand) result.Unit {
			r := l.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withMapping[label.Label](func(lbl label.Label) any {
					return newLabelModel(lbl)
				}))
			return result.Unit{}
		},
		func(err error) result.Unit {
			c.JSON(http.StatusBadRequest, httpError{
				Message: err.Error(),
			})
			return result.Unit{}
		})
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
