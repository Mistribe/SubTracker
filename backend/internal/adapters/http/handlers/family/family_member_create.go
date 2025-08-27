package family

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type FamilyMemberCreateEndpoint struct {
	handler ports.CommandHandler[command.CreateFamilyMemberCommand, family.Family]
}

type createFamilyMemberModel struct {
	Id        *string    `json:"id,omitempty"`
	Name      string     `json:"name" binding:"required"`
	Type      string     `json:"type" binding:"required" enums:"owner,adult,kid"`
	CreatedAt *time.Time `json:"created_at,omitempty" format:"date-time"`
}

func (m createFamilyMemberModel) ToFamilyMember(familyId uuid.UUID) (family.Member, error) {
	var id uuid.UUID
	var err error
	var createdAt time.Time
	id, err = x.ParseOrNewUUID(m.Id)
	if err != nil {
		return nil, err
	}

	createdAt = x.ValueOrDefault(m.CreatedAt, time.Now())

	memberType, err := family.ParseMemberType(m.Type)
	if err != nil {
		return nil, err
	}
	return family.NewMember(
		id,
		familyId,
		m.Name,
		memberType,
		nil,
		createdAt,
		createdAt,
	), nil
}

func (m createFamilyMemberModel) Command(familyId uuid.UUID) (command.CreateFamilyMemberCommand, error) {
	fam, err := m.ToFamilyMember(familyId)
	if err != nil {
		return command.CreateFamilyMemberCommand{}, err
	}
	return command.CreateFamilyMemberCommand{
		FamilyId: familyId,
		Member:   fam,
	}, nil

}

// Handle godoc
//
//	@Summary		Add a new family member
//	@Description	Add a new member to an existing family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string					true	"Family ID (UUID format)"
//	@Param			member		body		createFamilyMemberModel	true	"Family member creation data"
//	@Success		201			{object}	familyModel				"Successfully added family member"
//	@Failure		400			{object}	HttpErrorResponse		"Bad Request - Invalid input data or family ID"
//	@Failure		401			{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse		"Family not found"
//	@Failure		500			{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/families/{familyId}/members [post]
func (f FamilyMemberCreateEndpoint) Handle(c *gin.Context) {
	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		FromError(c, errors.New("invalid user id"))
		return
	}

	var model createFamilyMemberModel
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}

	cmd, err := model.Command(familyId)
	if err != nil {
		FromError(c, err)
		return
	}
	r := f.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[family.Family](http.StatusCreated),
		WithMapping[family.Family](func(fm family.Family) any {
			return newFamilyModel(userId, fm)
		}))
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

func NewFamilyMemberCreateEndpoint(handler ports.CommandHandler[command.CreateFamilyMemberCommand, family.Family]) *FamilyMemberCreateEndpoint {
	return &FamilyMemberCreateEndpoint{
		handler: handler,
	}
}
