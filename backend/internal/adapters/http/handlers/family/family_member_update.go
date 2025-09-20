package family

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

type MemberUpdateEndpoint struct {
	handler ports.CommandHandler[command.UpdateFamilyMemberCommand, family.Family]
}

func updateFamilyMemberRequestToCommand(
	m dto.UpdateFamilyMemberRequest,
	familyId, memberId uuid.UUID) (command.UpdateFamilyMemberCommand, error) {

	updatedAt := option.None[time.Time]()
	if m.UpdatedAt != nil {
		updatedAt = option.Some(*m.UpdatedAt)
	}

	memberType, err := family.ParseMemberType(m.Type)
	if err != nil {
		return command.UpdateFamilyMemberCommand{}, err
	}

	return command.UpdateFamilyMemberCommand{
		FamilyId:  familyId,
		Id:        memberId,
		Name:      m.Name,
		Type:      memberType,
		UpdatedAt: updatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Update family member by ID
//	@Description	Update an existing family member's information such as name and kid status
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string							true	"Family ID (UUID format)"
//	@Param			id			path		string							true	"Family member ID (UUID format)"
//	@Param			member		body		dto.UpdateFamilyMemberRequest	true	"Updated family member data"
//	@Success		200			{object}	dto.FamilyModel					"Successfully updated family member"
//	@Failure		400			{object}	HttpErrorResponse				"Bad Request - Invalid input data or ID format"
//	@Failure		401			{object}	HttpErrorResponse				"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse				"Family or family member not found"
//	@Failure		500			{object}	HttpErrorResponse				"Internal Server Error"
//	@Router			/family/{familyId}/members/{id} [put]
func (f MemberUpdateEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		FromError(c, errors.New("id parameter is required"))
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		FromError(c, err)
		return
	}

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

	var model dto.UpdateFamilyMemberRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}

	cmd, err := updateFamilyMemberRequestToCommand(model, familyId, id)
	if err != nil {
		FromError(c, err)
		return
	}

	r := f.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithMapping[family.Family](func(mbr family.Family) any {
			return dto.NewFamilyModel(userId, mbr)
		}))
}

func (f MemberUpdateEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members/:id",
	}
}

func (f MemberUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (f MemberUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewMemberUpdateEndpoint(handler ports.CommandHandler[command.UpdateFamilyMemberCommand, family.Family]) *MemberUpdateEndpoint {
	return &MemberUpdateEndpoint{
		handler: handler,
	}
}
