package family

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type MemberCreateEndpoint struct {
	handler ports.CommandHandler[command.CreateFamilyMemberCommand, family.Family]
}

func createFamilyMemberRequestToFamilyMember(m dto.CreateFamilyMemberRequest, familyId uuid.UUID) (
	family.Member,
	error) {
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

// Handle godoc
//
//	@Summary		Add a new family member
//	@Description	Add a new member to an existing family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string							true	"Family ID (UUID format)"
//	@Param			member		body		dto.CreateFamilyMemberRequest	true	"Family member creation data"
//	@Success		201			{object}	dto.FamilyModel					"Successfully added family member"
//	@Failure		400			{object}	HttpErrorResponse				"Bad Request - Invalid input data or family ID"
//	@Failure		401			{object}	HttpErrorResponse				"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse				"Family not found"
//	@Failure		500			{object}	HttpErrorResponse				"Internal Server Error"
//	@Router			/family/{familyId}/members [post]
func (f MemberCreateEndpoint) Handle(c *gin.Context) {
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

	var model dto.CreateFamilyMemberRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}

	fm, err := createFamilyMemberRequestToFamilyMember(model, familyId)
	if err != nil {
		FromError(c, err)
		return
	}
	cmd := command.CreateFamilyMemberCommand{
		FamilyId: familyId,
		Member:   fm,
	}
	r := f.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[family.Family](http.StatusCreated),
		WithMapping[family.Family](func(fm family.Family) any {
			return dto.NewFamilyModel(userId, fm)
		}))
}

func (f MemberCreateEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members",
	}
}

func (f MemberCreateEndpoint) Method() string {
	return http.MethodPost
}

func (f MemberCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewMemberCreateEndpoint(handler ports.CommandHandler[command.CreateFamilyMemberCommand, family.Family]) *MemberCreateEndpoint {
	return &MemberCreateEndpoint{
		handler: handler,
	}
}
