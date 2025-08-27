package family

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type FamilyCreateEndpoint struct {
	handler ports.CommandHandler[command.CreateFamilyCommand, family.Family]
}

// @name	CreateFamily
type createFamilyModel struct {
	Id          *string    `json:"id,omitempty"`
	Name        string     `json:"name" binding:"required"`
	CreatorName string     `json:"creator_name" binding:"required"`
	CreatedAt   *time.Time `json:"created_at,omitempty" format:"date-time"`
}

func (m createFamilyModel) Command() (command.CreateFamilyCommand, error) {
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
//	@Param			family	body		createFamilyModel	true	"Family creation data"
//	@Success		201		{object}	familyModel			"Successfully created family"
//	@Failure		400		{object}	HttpErrorResponse	"Bad Request - Invalid input data"
//	@Failure		401		{object}	HttpErrorResponse	"Unauthorized - Invalid user authentication"
//	@Failure		500		{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/families [post]
func (f FamilyCreateEndpoint) Handle(c *gin.Context) {
	var model createFamilyModel
	if err := c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}
	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		// todo
		FromError(c, errors.New("invalid user id"))
		return
	}

	cmd, err := model.Command()
	if err != nil {
		FromError(c, err)
		return
	}
	r := f.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[family.Family](http.StatusCreated),
		WithMapping[family.Family](func(f family.Family) any {
			return newFamilyModel(userId, f)
		}))
}

func (f FamilyCreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (f FamilyCreateEndpoint) Method() string {
	return http.MethodPost
}

func (f FamilyCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyCreateEndpoint(handler ports.CommandHandler[command.CreateFamilyCommand, family.Family]) *FamilyCreateEndpoint {
	return &FamilyCreateEndpoint{
		handler: handler,
	}
}
