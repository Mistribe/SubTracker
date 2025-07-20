package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/langext/option"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FamilyUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateFamilyCommand, family.Family]
}

type updateFamilyModel struct {
	Name             string     `json:"name"`
	HaveJointAccount bool       `json:"have_joint_account"`
	UpdatedAt        *time.Time `json:"updated_at,omitempty"`
}

func (m updateFamilyModel) ToFamily(id uuid.UUID, ownerId string, createdAt time.Time) result.Result[family.Family] {
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, time.Now())
	return family.NewFamily(
		id,
		ownerId,
		m.Name,
		m.HaveJointAccount,
		createdAt,
		updatedAt,
	)
}

func (m updateFamilyModel) Command(familyId uuid.UUID) result.Result[command.UpdateFamilyCommand] {
	return result.Success(command.UpdateFamilyCommand{
		Id:               familyId,
		Name:             m.Name,
		UpdatedAt:        option.Some[time.Time](ext.ValueOrDefault[time.Time](m.UpdatedAt, time.Now())),
		HaveJointAccount: m.HaveJointAccount,
	})

}

// Handle godoc
// @Summary		Update a family
// @Description	Update a family
// @Tags			family
// @Accept			json
// @Produce			json
// @Param			familyId	path		string	true	"Family member ID"
// @Param			family	body		updateFamilyModel	true	"Family data"
// @Success		200		{object}	familyModel
// @Failure		400		{object}	httpError
// @Router			/families/{familyId} [put]
func (f FamilyUpdateEndpoint) Handle(c *gin.Context) {
	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
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

	var model updateFamilyModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	cmd := model.Command(familyId)
	result.Match[command.UpdateFamilyCommand, result.Unit](cmd,
		func(cmd command.UpdateFamilyCommand) result.Unit {
			r := f.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withStatus[family.Family](http.StatusOK),
				withMapping[family.Family](func(f family.Family) any {
					return newFamilyModel(userId, f)
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
