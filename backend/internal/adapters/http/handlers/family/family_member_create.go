package family

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/langext/option"
)

type MemberCreateEndpoint struct {
	handler        ports.CommandHandler[command.CreateFamilyMemberCommand, family.Family]
	authentication ports.Authentication
}

func createFamilyMemberRequestToCommand(m dto.CreateFamilyMemberRequest, familyId types.FamilyID) (
	command.CreateFamilyMemberCommand,
	error) {
	familyMemberID, err := types.ParseFamilyMemberIDOrNil(m.Id)
	if err != nil {
		return command.CreateFamilyMemberCommand{}, err
	}

	memberType, err := family.ParseMemberType(m.Type)
	if err != nil {
		return command.CreateFamilyMemberCommand{}, err
	}

	return command.CreateFamilyMemberCommand{
		FamilyID:       familyId,
		FamilyMemberID: option.New(familyMemberID),
		Name:           m.Name,
		Type:           memberType,
		CreatedAt:      option.New(m.CreatedAt),
	}, nil
}

// Handle godoc
//
//	@Summary		Add a new family member
//	@Description	Add a new member to an existing family
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string							true	"Family LabelID (UUID format)"
//	@Param			member		body		dto.CreateFamilyMemberRequest	true	"Family member creation data"
//	@Success		201			{object}	dto.FamilyModel					"Successfully added family member"
//	@Failure		400			{object}	HttpErrorResponse				"Bad Request - Invalid input data or family LabelID"
//	@Failure		401			{object}	HttpErrorResponse				"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse				"Family not found"
//	@Failure		500			{object}	HttpErrorResponse				"Internal Server Error"
//	@Router			/family/{familyId}/members [post]
func (e MemberCreateEndpoint) Handle(c *gin.Context) {
	familyId, err := types.ParseFamilyID(c.Param("familyId"))
	if err != nil {
		FromError(c, err)
		return
	}

	connectedAccount := e.authentication.MustGetConnectedAccount(c)

	var model dto.CreateFamilyMemberRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}

	cmd, err := createFamilyMemberRequestToCommand(model, familyId)
	if err != nil {
		FromError(c, err)
		return
	}
	r := e.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[family.Family](http.StatusCreated),
		WithMapping[family.Family](func(fm family.Family) any {
			return dto.NewFamilyModel(connectedAccount.UserID(), fm)
		}))
}

func (e MemberCreateEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members",
	}
}

func (e MemberCreateEndpoint) Method() string {
	return http.MethodPost
}

func (e MemberCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewMemberCreateEndpoint(handler ports.CommandHandler[command.CreateFamilyMemberCommand, family.Family],
	authentication ports.Authentication) *MemberCreateEndpoint {
	return &MemberCreateEndpoint{
		handler:        handler,
		authentication: authentication,
	}
}
