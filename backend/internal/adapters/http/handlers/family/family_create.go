package family

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/langext/option"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type CreateEndpoint struct {
	handler        ports.CommandHandler[command.CreateFamilyCommand, family.Family]
	authentication ports.Authentication
}

func createFamilyRequestToCommand(m dto.CreateFamilyRequest) (command.CreateFamilyCommand, error) {
	familyID, err := types.ParseFamilyIDOrNil(m.Id)
	if err != nil {
		return command.CreateFamilyCommand{}, err
	}
	return command.CreateFamilyCommand{
		FamilyId:    option.New(familyID),
		Name:        m.Name,
		CreatorName: m.CreatorName,
		CreatedAt:   option.New(m.CreatedAt),
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
func (e CreateEndpoint) Handle(c *gin.Context) {
	var model dto.CreateFamilyRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}
	connectedAccount := e.authentication.MustGetConnectedAccount(c)

	cmd, err := createFamilyRequestToCommand(model)
	if err != nil {
		FromError(c, err)
		return
	}
	r := e.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[family.Family](http.StatusCreated),
		WithMapping[family.Family](func(f family.Family) any {
			return dto.NewFamilyModel(connectedAccount.UserID(), f)
		}))
}

func (e CreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e CreateEndpoint) Method() string {
	return http.MethodPost
}

func (e CreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewCreateEndpoint(handler ports.CommandHandler[command.CreateFamilyCommand, family.Family],
	authentication ports.Authentication) *CreateEndpoint {
	return &CreateEndpoint{
		handler:        handler,
		authentication: authentication,
	}
}
