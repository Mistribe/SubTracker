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
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/pkg/ginx"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x"
)

type LabelCreateEndpoint struct {
	handler ports.CommandHandler[command.CreateLabelCommand, label.Label]
}

type createLabelModel struct {
	Id        *string                `json:"id,omitempty"`
	Name      string                 `json:"name" binding:"required"`
	Color     string                 `json:"color" binding:"required"`
	Owner     dto.EditableOwnerModel `json:"owner" binding:"required"`
	CreatedAt *time.Time             `json:"created_at,omitempty" format:"date-time"`
}

func (m createLabelModel) ToLabel(userId string) (label.Label, error) {
	var id uuid.UUID
	var err error
	var createdAt time.Time

	id, err = x.ParseOrNewUUID(m.Id)
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
//	@Tags			labels
//	@Accept			json
//	@Produce		json
//	@Param			label	body		createLabelModel	true	"Label creation data"
//	@Success		201		{object}	labelModel			"Successfully created label"
//	@Failure		400		{object}	HttpErrorResponse	"Bad Request - Invalid input data"
//	@Failure		401		{object}	HttpErrorResponse	"Unauthorized - Invalid user authentication"
//	@Failure		500		{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/labels [post]
func (l LabelCreateEndpoint) Handle(c *gin.Context) {
	var model createLabelModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ginx.HttpErrorResponse{
			Message: "invalid user id",
		})
		return
	}

	cmd, err := model.Command(userId)
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
		WithStatus[label.Label](http.StatusCreated),
		WithMapping[label.Label](func(lbl label.Label) any {
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

func NewLabelCreateEndpoint(handler ports.CommandHandler[command.CreateLabelCommand, label.Label]) *LabelCreateEndpoint {
	return &LabelCreateEndpoint{
		handler: handler,
	}
}
