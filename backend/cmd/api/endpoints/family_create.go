package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
)

type FamilyCreateEndpoint struct {
	handler core.CommandHandler[command.CreateFamilyCommand, family.Family]
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
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, HttpErrorResponse{
			Message: "invalid user id",
		})
		return
	}

	cmd, err := model.Command()
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	r := f.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withStatus[family.Family](http.StatusCreated),
		withMapping[family.Family](func(f family.Family) any {
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

func NewFamilyCreateEndpoint(handler core.CommandHandler[command.CreateFamilyCommand, family.Family]) *FamilyCreateEndpoint {
	return &FamilyCreateEndpoint{
		handler: handler,
	}
}
