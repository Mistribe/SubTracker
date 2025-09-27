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
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

type MemberUpdateEndpoint struct {
	handler        ports.CommandHandler[command.UpdateFamilyMemberCommand, family.Family]
	authentication ports.Authentication
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
		FamilyID:       familyId,
		FamilyMemberID: memberId,
		Name:           m.Name,
		Type:           memberType,
		UpdatedAt:      updatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Update family member by LabelID
//	@Description	Update an existing family member's information such as name and kid status
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string							true	"Family LabelID (UUID format)"
//	@Param			id			path		string							true	"Family member LabelID (UUID format)"
//	@Param			member		body		dto.UpdateFamilyMemberRequest	true	"Updated family member data"
//	@Success		200			{object}	dto.FamilyModel					"Successfully updated family member"
//	@Failure		400			{object}	HttpErrorResponse				"Bad Request - Invalid input data or LabelID format"
//	@Failure		401			{object}	HttpErrorResponse				"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse				"Family or family member not found"
//	@Failure		500			{object}	HttpErrorResponse				"Internal Server Error"
//	@Router			/family/{familyId}/members/{id} [put]
func (e MemberUpdateEndpoint) Handle(c *gin.Context) {
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

	connectedAccount := e.authentication.MustGetConnectedAccount(c)

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

	r := e.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithMapping[family.Family](func(mbr family.Family) any {
			return dto.NewFamilyModel(connectedAccount.UserID(), mbr)
		}))
}

func (e MemberUpdateEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members/:id",
	}
}

func (e MemberUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (e MemberUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewMemberUpdateEndpoint(handler ports.CommandHandler[command.UpdateFamilyMemberCommand, family.Family],
	authentication ports.Authentication) *MemberUpdateEndpoint {
	return &MemberUpdateEndpoint{
		handler:        handler,
		authentication: authentication,
	}
}
