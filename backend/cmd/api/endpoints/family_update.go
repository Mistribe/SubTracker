package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FamilyMemberUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateFamilyMemberCommand, family.Member]
}

type updateFamilyMemberModel struct {
	Name      string     `json:"name"`
	Email     *string    `json:"email,omitempty"`
	IsKid     bool       `json:"id_kid"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}

func (m updateFamilyMemberModel) Command(id uuid.UUID) result.Result[command.UpdateFamilyMemberCommand] {
	email := option.None[string]()
	if m.Email != nil {
		email = option.Some(*m.Email)
	}

	updatedAt := option.None[time.Time]()
	if m.UpdatedAt != nil {
		updatedAt = option.Some(*m.UpdatedAt)
	}

	return result.Success(command.UpdateFamilyMemberCommand{
		Id:        id,
		Name:      m.Name,
		Email:     email,
		IsKid:     m.IsKid,
		UpdatedAt: updatedAt,
	})
}

// Handle godoc
// @Summary		Update family member by ID
// @Description	Update family member by ID
// @Tags			family
// @Accept			json
// @Produce		json
// @Param			id		path		uuid.UUID				true	"Family member ID"
// @Param			member	body		updateFamilyMemberModel	true	"Family member data"
// @Success		200		{object}	familyMemberModel
// @Failure		400		{object}	httpError
// @Failure		404		{object}	httpError
// @Router			/families/members/{id} [put]
func (f FamilyMemberUpdateEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "invalid id format",
		})
		return
	}

	var model updateFamilyMemberModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	result.Match[command.UpdateFamilyMemberCommand, result.Unit](model.Command(id),
		func(cmd command.UpdateFamilyMemberCommand) result.Unit {
			r := f.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withMapping[family.Member](func(mbr family.Member) any {
					return newFamilyMemberModel(mbr)
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

func (f FamilyMemberUpdateEndpoint) Pattern() []string {
	return []string{
		"/members/:id",
	}
}

func (f FamilyMemberUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (f FamilyMemberUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberUpdateEndpoint(handler core.CommandHandler[command.UpdateFamilyMemberCommand, family.Member]) *FamilyMemberUpdateEndpoint {
	return &FamilyMemberUpdateEndpoint{
		handler: handler,
	}
}
