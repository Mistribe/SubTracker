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
	"github.com/oleexo/subtracker/pkg/langext/option"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type FamilyMemberCreateEndpoint struct {
	handler core.CommandHandler[command.CreateFamilyMemberCommand, family.Family]
}

type createFamilyMemberModel struct {
	Id        *string    `json:"id,omitempty"`
	Name      string     `json:"name" binding:"required"`
	Email     *string    `json:"email,omitempty"`
	IsKid     bool       `json:"is_kid" binding:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (m createFamilyMemberModel) ToFamilyMember(familyId uuid.UUID) result.Result[family.Member] {
	var id uuid.UUID
	var err error
	var createdAt time.Time
	email := option.None[string]()
	if m.Email != nil {
		email = option.Some(*m.Email)
	}
	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return result.Fail[family.Member](err)
	}

	createdAt = ext.ValueOrDefault(m.CreatedAt, time.Now())

	return family.NewMember(
		id,
		familyId,
		m.Name,
		email,
		m.IsKid,
		createdAt,
		createdAt,
	)
}

func (m createFamilyMemberModel) Command(familyId uuid.UUID) result.Result[command.CreateFamilyMemberCommand] {
	return result.Bind[family.Member, command.CreateFamilyMemberCommand](
		m.ToFamilyMember(familyId),
		func(fm family.Member) result.Result[command.CreateFamilyMemberCommand] {
			return result.Success(command.CreateFamilyMemberCommand{
				Member: fm,
			})
		})
}

// Handle godoc
//
//	@Summary		Add a new family member
//	@Description	Add a new family member
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			member	body		createFamilyMemberModel	true	"Family member data"
//	@Success		201		{object}	familyModel
//	@Failure		400		{object}	httpError
//	@Router			/families/{familyId}/members [post]
func (f FamilyMemberCreateEndpoint) Handle(c *gin.Context) {
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

	var model createFamilyMemberModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	cmd := model.Command(familyId)
	result.Match[command.CreateFamilyMemberCommand, result.Unit](cmd,
		func(cmd command.CreateFamilyMemberCommand) result.Unit {
			r := f.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withStatus[family.Family](http.StatusCreated),
				withMapping[family.Family](func(fm family.Family) any {
					return newFamilyModel(userId, fm)
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

func (f FamilyMemberCreateEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members",
	}
}

func (f FamilyMemberCreateEndpoint) Method() string {
	return http.MethodPost
}

func (f FamilyMemberCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberCreateEndpoint(handler core.CommandHandler[command.CreateFamilyMemberCommand, family.Family]) *FamilyMemberCreateEndpoint {
	return &FamilyMemberCreateEndpoint{
		handler: handler,
	}
}
