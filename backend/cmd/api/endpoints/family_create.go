package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/domain/user"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FamilyCreateEndpoint struct {
	handler core.CommandHandler[command.CreateFamilyCommand, family.Family]
}

type createFamilyModel struct {
	Id               *string    `json:"id,omitempty"`
	Name             string     `json:"name"`
	HaveJointAccount bool       `json:"have_joint_account"`
	CreatedAt        *time.Time `json:"created_at,omitempty"`
}

func (m createFamilyModel) ToFamily(ownerId string) result.Result[family.Family] {
	var id uuid.UUID
	var err error
	var createdAt time.Time

	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return result.Fail[family.Family](err)
	}

	createdAt = ext.ValueOrDefault(m.CreatedAt, time.Now())

	return family.NewFamily(
		id,
		ownerId,
		m.Name,
		m.HaveJointAccount,
		createdAt,
		createdAt,
	)
}

func (m createFamilyModel) Command(ownerId string) result.Result[command.CreateFamilyCommand] {
	return result.Bind[family.Family, command.CreateFamilyCommand](
		m.ToFamily(ownerId),
		func(f family.Family) result.Result[command.CreateFamilyCommand] {
			return result.Success(command.CreateFamilyCommand{
				Family: f,
			})
		})
}

// Handle godoc
// @Summary		Create a new family
// @Description	Create a new family
// @Tags			family
// @Accept			json
// @Produce			json
// @Param			family	body		createFamilyModel	true	"Family data"
// @Success		201		{object}	familyModel
// @Failure		400		{object}	httpError
// @Router			/families [post]
func (f FamilyCreateEndpoint) Handle(c *gin.Context) {
	var model createFamilyModel
	if err := c.ShouldBindJSON(&model); err != nil {
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

	cmd := model.Command(userId)
	result.Match[command.CreateFamilyCommand, result.Unit](cmd,
		func(cmd command.CreateFamilyCommand) result.Unit {
			r := f.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withStatus[family.Family](http.StatusCreated),
				withMapping[family.Family](func(f family.Family) any {
					return newFamilyModel(f)
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
