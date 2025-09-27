package label

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/label/command"
	"github.com/mistribe/subtracker/pkg/ginx"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

type CreateEndpoint struct {
	handler        ports.CommandHandler[command.CreateLabelCommand, label.Label]
	authentication ports.Authentication
}

func createLabelRequestToCommand(m dto.CreateLabelRequest, userId types.UserID) (command.CreateLabelCommand, error) {
	owner, err := m.Owner.Owner(userId)
	if err != nil {
		return command.CreateLabelCommand{}, err
	}
	labelID, err := types.ParseLabelIDOrNil(m.Id)
	if err != nil {
		return command.CreateLabelCommand{}, err
	}
	return command.CreateLabelCommand{
		LabelID:   option.New(labelID),
		Name:      m.Name,
		Color:     m.Color,
		Owner:     owner,
		CreatedAt: option.New(m.CreatedAt),
	}, nil
}

// Handle godoc
//
//	@Summary		Create a new label
//	@Description	Create a new label with specified name, color, and owner information
//	@Tags			labels
//	@Accept			json
//	@Produce		json
//	@Param			label	body		dto.CreateLabelRequest	true	"Label creation data"
//	@Success		201		{object}	dto.LabelModel			"Successfully created label"
//	@Failure		400		{object}	HttpErrorResponse		"Bad Request - Invalid input data"
//	@Failure		401		{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		500		{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/labels [post]
func (l CreateEndpoint) Handle(c *gin.Context) {
	var model dto.CreateLabelRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	connectedAccount := l.authentication.MustGetConnectedAccount(c)

	cmd, err := createLabelRequestToCommand(model, connectedAccount.UserID())
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
			return dto.NewLabelModel(lbl)
		}))
}

func (l CreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (l CreateEndpoint) Method() string {
	return http.MethodPost
}

func (l CreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewCreateEndpoint(handler ports.CommandHandler[command.CreateLabelCommand, label.Label]) *CreateEndpoint {
	return &CreateEndpoint{
		handler: handler,
	}
}
