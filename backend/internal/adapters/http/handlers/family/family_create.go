package family

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type CreateEndpoint struct {
	handler ports.CommandHandler[command.CreateFamilyCommand, family.Family]
}

func createFamilyRequestToCommand(m dto.CreateFamilyRequest) (command.CreateFamilyCommand, error) {
	var familyId *uuid.UUID
	if m.Id != nil {
		id, err := uuid.Parse(*m.Id)
		if err != nil {
			return command.CreateFamilyCommand{}, err
		}
		familyId = &id
	}
	return command.CreateFamilyCommand{
		FamilyId:    familyId,
		Name:        m.Name,
		CreatorName: m.CreatorName,
		CreatedAt:   m.CreatedAt,
	}, nil

}

// Handle godoc
//
//	@Summary		Create a new family
//	@Description	Create a new family with the authenticated user as the owner and initial member
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			family	body		dto.CreateFamilyRequest	true	"Family creation data"
//	@Success		201		{object}	dto.FamilyModel			"Successfully created family"
//	@Failure		400		{object}	HttpErrorResponse		"Bad Request - Invalid input data"
//	@Failure		401		{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		500		{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/family [post]
func (f CreateEndpoint) Handle(c *gin.Context) {
	var model dto.CreateFamilyRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}
	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		FromError(c, errors.New("invalid user id"))
		return
	}

	cmd, err := createFamilyRequestToCommand(model)
	if err != nil {
		FromError(c, err)
		return
	}
	r := f.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[family.Family](http.StatusCreated),
		WithMapping[family.Family](func(f family.Family) any {
			return dto.NewFamilyModel(userId, f)
		}))
}

func (f CreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (f CreateEndpoint) Method() string {
	return http.MethodPost
}

func (f CreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewCreateEndpoint(handler ports.CommandHandler[command.CreateFamilyCommand, family.Family]) *CreateEndpoint {
	return &CreateEndpoint{
		handler: handler,
	}
}
