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
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type LabelCreateEndpoint struct {
	handler core.CommandHandler[command.CreateLabelCommand, label.Label]
}

type createLabelModel struct {
	Id        *string    `json:"id,omitempty"`
	Name      string     `json:"name" binding:"required"`
	IsDefault *bool      `json:"is_default,omitempty"`
	Color     string     `json:"color" binding:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty" format:"date-time"`
}

func (m createLabelModel) ToLabel(userId string) result.Result[label.Label] {
	var id uuid.UUID
	var err error
	var createdAt time.Time

	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return result.Fail[label.Label](err)
	}

	createdAt = ext.ValueOrDefault(m.CreatedAt, time.Now())
	isDefault := ext.ValueOrDefault(m.IsDefault, false)

	return label.NewLabel(
		id,
		&userId,
		m.Name,
		isDefault,
		strings.ToUpper(m.Color),
		createdAt,
		createdAt,
	)
}

func (m createLabelModel) Command(userId string) result.Result[command.CreateLabelCommand] {
	return result.Bind[label.Label, command.CreateLabelCommand](
		m.ToLabel(userId),
		func(lbl label.Label) result.Result[command.CreateLabelCommand] {
			return result.Success(command.CreateLabelCommand{
				Label: lbl,
			})
		})
}

// Handle godoc
// @Summary		Create a new label
// @Description	Create a new label
// @Tags			label
// @Accept			json
// @Produce		json
// @Param			label	body		createLabelModel	true	"Label data"
// @Success		201		{object}	labelModel
// @Failure		400		{object}	httpError
// @Router			/labels [post]
func (l LabelCreateEndpoint) Handle(c *gin.Context) {
	var model createLabelModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	userId, ok := user.FromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, httpError{
			Message: "invalid user id",
		})
		return
	}

	cmd := model.Command(userId)

	result.Match[command.CreateLabelCommand, result.Unit](cmd,
		func(cmd command.CreateLabelCommand) result.Unit {
			r := l.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withStatus[label.Label](http.StatusCreated),
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

func (l LabelCreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (l LabelCreateEndpoint) Method() string {
	return http.MethodPost
}

func (l LabelCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewLabelCreateEndpoint(handler core.CommandHandler[command.CreateLabelCommand, label.Label]) *LabelCreateEndpoint {
	return &LabelCreateEndpoint{
		handler: handler,
	}
}
