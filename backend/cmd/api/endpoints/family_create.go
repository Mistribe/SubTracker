package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/core/option"
	"github.com/oleexo/subtracker/internal/application/core/result"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/ext"
)

type FamilyMemberCreateEndpoint struct {
	handler core.CommandHandler[command.CreateFamilyMemberCommand, family.Member]
}

type createFamilyMemberModel struct {
	Id        *string    `json:"id,omitempty"`
	Name      string     `json:"name"`
	Email     *string    `json:"email,omitempty"`
	IsKid     bool       `json:"is_kid"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
}

func (m createFamilyMemberModel) ToFamilyMember() result.Result[family.Member] {
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
		m.Name,
		email,
		m.IsKid,
		createdAt,
		createdAt,
	)
}

func (m createFamilyMemberModel) Command() result.Result[command.CreateFamilyMemberCommand] {
	return result.Bind[family.Member, command.CreateFamilyMemberCommand](
		m.ToFamilyMember(),
		func(fm family.Member) result.Result[command.CreateFamilyMemberCommand] {
			return result.Success(command.CreateFamilyMemberCommand{
				Member: fm,
			})
		})
}

// Handle godoc
// @Summary		Create a new family member
// @Description	Create a new family member
// @Tags			family
// @Accept			json
// @Produce		json
// @Param			member	body		createFamilyMemberModel	true	"Family member data"
// @Success		201		{object}	familyMemberModel
// @Failure		400		{object}	httpError
// @Router			/families/members [post]
func (f FamilyMemberCreateEndpoint) Handle(c *gin.Context) {
	var model createFamilyMemberModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	result.Match[command.CreateFamilyMemberCommand, result.Unit](model.Command(),
		func(cmd command.CreateFamilyMemberCommand) result.Unit {
			r := f.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withMapping[family.Member](func(fm family.Member) any {
					return newFamilyMemberModel(fm)
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
		"",
	}
}

func (f FamilyMemberCreateEndpoint) Method() string {
	return http.MethodPost
}

func (f FamilyMemberCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberCreateEndpoint(handler core.CommandHandler[command.CreateFamilyMemberCommand, family.Member]) *FamilyMemberCreateEndpoint {
	return &FamilyMemberCreateEndpoint{
		handler: handler,
	}
}
