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
)

type FamilyCreateEndpoint struct {
	handler core.CommandHandler[command.CreateFamilyCommand, family.Family]
}

type createFamilyModel struct {
	Id               *string    `json:"id,omitempty"`
	Name             string     `json:"name" binding:"required"`
	HaveJointAccount bool       `json:"have_joint_account,omitempty"`
	CreatedAt        *time.Time `json:"created_at,omitempty" format:"date-time"`
}

func (m createFamilyModel) ToFamily(ownerId string) (family.Family, error) {
	var familyId uuid.UUID
	var err error
	var createdAt time.Time

	familyId, err = parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}

	createdAt = ext.ValueOrDefault(m.CreatedAt, time.Now())

	creatorId, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	creator := family.NewMember(creatorId,
		familyId,
		"Owner", // todo replace with name from token,
		false,
		time.Now(),
		time.Now(),
	)
	creator.SetUserId(&ownerId)
	return family.NewFamily(
		familyId,
		ownerId,
		m.Name,
		[]family.Member{
			creator,
		},
		createdAt,
		createdAt,
	), nil
}

func (m createFamilyModel) Command(ownerId string) (command.CreateFamilyCommand, error) {
	fam, err := m.ToFamily(ownerId)
	if err != nil {
		return command.CreateFamilyCommand{}, err
	}
	return command.CreateFamilyCommand{
		Family: fam,
	}, nil

}

// Handle godoc
//
//	@Summary		Create a new family
//	@Description	Create a new family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			family	body		createFamilyModel	true	"Family data"
//	@Success		201		{object}	familyModel
//	@Failure		400		{object}	httpError
//	@Router			/families [post]
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

	cmd, err := model.Command(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
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
