package family

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x"
	"github.com/mistribe/subtracker/pkg/x/collection"

	"github.com/mistribe/subtracker/internal/domain/family"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
)

type FamilyPatchEndpoint struct {
	handler ports.CommandHandler[command.PatchFamilyCommand, family.Family]
}

func NewFamilyPatchEndpoint(handler ports.CommandHandler[command.PatchFamilyCommand, family.Family]) *FamilyPatchEndpoint {
	return &FamilyPatchEndpoint{handler: handler}
}

// @Description	Model for updating family member details
type patchFamilyMemberModel struct {
	// Optional member ID. If not provided, new member will be created
	Id *string `json:"id,omitempty"`
	// member's name
	Name string `json:"name" binding:"required"`
	// Indicates if the member is a kid
	Type string `json:"type" binding:"required" enums:"owner,adult,kid"`
	// Optional timestamp of the last update
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchFamilyMemberModel) Command(familyId uuid.UUID) (family.Member, error) {
	var id uuid.UUID
	var err error

	id, err = x.ParseOrNewUUID(m.Id)
	if err != nil {
		return nil, err
	}
	updatedAt := x.ValueOrDefault(m.UpdatedAt, time.Now())

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
		updatedAt,
		updatedAt,
	), nil

}

// @Description	Model for updating family details
type patchFamilyModel struct {
	// Family ID
	Id *string `json:"id" binding:"required"`
	// Family name
	Name string `json:"name" binding:"required"`
	// List of family members
	Members []patchFamilyMemberModel `json:"members,omitempty"`
	// Optional timestamp of the last update
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchFamilyModel) Command(ownerId string) (command.PatchFamilyCommand, error) {

	familyId, err := x.ParseOrNewUUID(m.Id)
	if err != nil {
		return command.PatchFamilyCommand{}, err
	}
	updatedAt := x.ValueOrDefault(m.UpdatedAt, time.Now())
	members, err := collection.SelectErr(m.Members, func(member patchFamilyMemberModel) (family.Member, error) {
		return member.Command(familyId)
	})
	if err != nil {
		return command.PatchFamilyCommand{}, err
	}
	fam := family.NewFamily(
		familyId,
		ownerId,
		m.Name,
		members,
		updatedAt,
		updatedAt,
	)
	return command.PatchFamilyCommand{
		Family: fam,
	}, nil
}

// Handle godoc
//
//	@Summary		Patch family with members
//	@Description	Update or create a family with specified members. If family doesn't exist, it will be created.
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			family	body		patchFamilyModel	true	"Family update data with members"
//	@Success		200		{object}	familyModel			"Successfully updated family"
//	@Failure		400		{object}	HttpErrorResponse	"Bad Request - Invalid input data"
//	@Failure		401		{object}	HttpErrorResponse	"Unauthorized - Invalid user authentication"
//	@Failure		404		{object}	HttpErrorResponse	"Family not found"
//	@Failure		500		{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/families [patch]
func (e FamilyPatchEndpoint) Handle(c *gin.Context) {
	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		FromError(c, errors.New("invalid user id"))
		return
	}

	var model patchFamilyModel
	if err := c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}

	cmd, err := model.Command(userId)
	if err != nil {
		FromError(c, err)
		return
	}
	r := e.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[family.Family](http.StatusOK),
		WithMapping[family.Family](func(f family.Family) any {
			return newFamilyModel(userId, f)
		}))
}

func (e FamilyPatchEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e FamilyPatchEndpoint) Method() string {
	return http.MethodPatch
}

func (e FamilyPatchEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
