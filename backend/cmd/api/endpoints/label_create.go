package endpoints

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/label/command"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/label"
)

type LabelCreateEndpoint struct {
	handler core.CommandHandler[command.CreateLabelCommand, label.Label]
}

type createLabelModel struct {
	Id        *string            `json:"id,omitempty"`
	Name      string             `json:"name" binding:"required"`
	Color     string             `json:"color" binding:"required"`
	Owner     EditableOwnerModel `json:"owner" binding:"required"`
	CreatedAt *time.Time         `json:"created_at,omitempty" format:"date-time"`
}

func (m createLabelModel) ToLabel(userId string) (label.Label, error) {
	var id uuid.UUID
	var err error
	var createdAt time.Time

	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}

	owner, err := m.Owner.Owner(userId)
	if err != nil {
		return nil, err
	}
	return label.NewLabel(
		id,
		owner,
		m.Name,
		nil,
		strings.ToUpper(m.Color),
		createdAt,
		createdAt,
	), nil
}

func (m createLabelModel) Command(userId string) (command.CreateLabelCommand, error) {
	lbl, err := m.ToLabel(userId)
	if err != nil {
		return command.CreateLabelCommand{}, err
	}
	return command.CreateLabelCommand{
		Label: lbl,
	}, nil
}

// Handle godoc
//
//	@Summary		Create a new label
//	@Description	Create a new label with specified name, color, and owner information
//	@Tags			label
//	@Accept			json
//	@Produce		json
//	@Param			label	body		createLabelModel	true	"Label creation data"
//	@Success		201		{object}	labelModel			"Successfully created label"
//	@Failure		400		{object}	httpError			"Bad Request - Invalid input data"
//	@Failure		401		{object}	httpError			"Unauthorized - Invalid user authentication"
//	@Failure		500		{object}	httpError			"Internal Server Error"
//	@Router			/labels [post]
func (l LabelCreateEndpoint) Handle(c *gin.Context) {
	var model createLabelModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, httpError{
			Message: "invalid user id",
		})
		return
	}

	cmd, err := model.Command(userId)
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
		withStatus[label.Label](http.StatusCreated),
		withMapping[label.Label](func(lbl label.Label) any {
			return newLabelModel(lbl)
		}))
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
