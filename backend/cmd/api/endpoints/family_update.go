package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/pkg/langext/option"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/ext"
)

type FamilyUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateFamilyCommand, family.Family]
}

type updateFamilyModel struct {
	Name      string     `json:"name" binding:"required"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m updateFamilyModel) Command(familyId uuid.UUID) (command.UpdateFamilyCommand, error) {
	return command.UpdateFamilyCommand{
		Id:        familyId,
		Name:      m.Name,
		UpdatedAt: option.Some[time.Time](ext.ValueOrDefault[time.Time](m.UpdatedAt, time.Now())),
	}, nil

}

// Handle godoc
//
//	@Summary		Update a family
//	@Description	Update family information such as name and other details
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string				true	"Family ID (UUID format)"
//	@Param			family		body		updateFamilyModel	true	"Updated family data"
//	@Success		200			{object}	familyModel			"Successfully updated family"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid input data or family ID"
//	@Failure		401			{object}	HttpErrorResponse	"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse	"Family not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/families/{familyId} [put]
func (f FamilyUpdateEndpoint) Handle(c *gin.Context) {
	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
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

	var model updateFamilyModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	cmd, err := model.Command(familyId)
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
		withStatus[family.Family](http.StatusOK),
		withMapping[family.Family](func(f family.Family) any {
			return newFamilyModel(userId, f)
		}))
}

func (f FamilyUpdateEndpoint) Pattern() []string {
	return []string{
		"/:familyId",
	}
}

func (f FamilyUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (f FamilyUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyUpdateEndpoint(handler core.CommandHandler[command.UpdateFamilyCommand, family.Family]) *FamilyUpdateEndpoint {
	return &FamilyUpdateEndpoint{
		handler: handler,
	}
}
