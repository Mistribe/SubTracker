package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type FamilyPatchEndpoint struct {
	handler core.CommandHandler[command.PatchFamilyCommand, family.Family]
}

func NewFamilyPatchEndpoint(handler core.CommandHandler[command.PatchFamilyCommand, family.Family]) *FamilyPatchEndpoint {
	return &FamilyPatchEndpoint{handler: handler}
}

// @Description	Model for updating family member details
type patchFamilyMemberModel struct {
	// Optional member ID. If not provided, new member will be created
	Id *string `json:"id,omitempty"`
	// member's name
	Name string `json:"name" binding:"required"`
	// Indicates if the member is a kid
	IsKid bool `json:"is_kid" binding:"required"`
	// Optional timestamp of the last update
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchFamilyMemberModel) Command(familyId uuid.UUID) (family.Member, error) {
	var id uuid.UUID
	var err error

	id, err = parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, time.Now())

	return family.NewMember(
		id,
		familyId,
		m.Name,
		m.IsKid,
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
	familyId, err := parseUuidOrNew(m.Id)
	if err != nil {
		return command.PatchFamilyCommand{}, err
	}
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, time.Now())
	members, err := slicesx.MapErr(m.Members, func(member patchFamilyMemberModel) (family.Member, error) {
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
//	@Failure		400		{object}	httpError			"Bad Request - Invalid input data"
//	@Failure		401		{object}	httpError			"Unauthorized - Invalid user authentication"
//	@Failure		404		{object}	httpError			"Family not found"
//	@Failure		500		{object}	httpError			"Internal Server Error"
//	@Router			/families [patch]
func (e FamilyPatchEndpoint) Handle(c *gin.Context) {
	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, httpError{
			Message: "invalid user id",
		})
		return
	}

	var model patchFamilyModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
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
	r := e.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withStatus[family.Family](http.StatusOK),
		withMapping[family.Family](func(f family.Family) any {
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
