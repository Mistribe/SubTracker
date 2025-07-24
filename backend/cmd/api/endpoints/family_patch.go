package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/family/command"
	"github.com/oleexo/subtracker/internal/domain/family"
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/langext/result"
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
	// Member's name
	Name string `json:"name" binding:"required"`
	// Optional email address
	Email *string `json:"email,omitempty"`
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
		return family.Member{}, err
	}
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, time.Now())

	return family.NewMemberWithoutValidation(
		id,
		familyId,
		m.Name,
		m.Email,
		m.IsKid,
		updatedAt,
		updatedAt,
		false,
	), nil

}

// @Description	Model for updating family details
type patchFamilyModel struct {
	// Family ID
	Id *string `json:"id" binding:"required"`
	// Family name
	Name string `json:"name" binding:"required"`
	// Indicates if family has joint account
	HaveJointAccount bool `json:"have_joint_account,omitempty"`
	// List of family members
	Members []patchFamilyMemberModel `json:"members,omitempty"`
	// Optional timestamp of the last update
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchFamilyModel) Command(ownerId string) result.Result[command.PatchFamilyCommand] {
	familyId, err := parseUuidOrNew(m.Id)
	if err != nil {
		return result.Fail[command.PatchFamilyCommand](err)
	}
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, time.Now())
	members, err := ext.MapErr(m.Members, func(member patchFamilyMemberModel) (family.Member, error) {
		return member.Command(familyId)
	})
	if err != nil {
		return result.Fail[command.PatchFamilyCommand](err)
	}
	fam := family.NewFamilyWithoutValidation(
		familyId,
		ownerId,
		m.Name,
		m.HaveJointAccount,
		members,
		updatedAt,
		updatedAt,
		false,
	)
	return result.Success(command.PatchFamilyCommand{
		Family: fam,
	})
}

// Handle
//
//	@Summary		Patch family with members
//	@Description	Patch family with members
//	@Tags			Family
//	@Accept			json
//	@Produce		json
//	@Param			family	body		patchFamilyModel	true	"Family update data"
//	@Success		200		{object}	familyModel
//	@Failure		400		{object}	httpError
//	@Failure		401		{object}	httpError
//	@Router			/families [patch]
func (e FamilyPatchEndpoint) Handle(c *gin.Context) {
	userId, ok := user.FromContext(c)
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

	cmd := model.Command(userId)
	result.Match[command.PatchFamilyCommand, result.Unit](cmd,
		func(cmd command.PatchFamilyCommand) result.Unit {
			r := e.handler.Handle(c, cmd)
			handleResponse(c,
				r,
				withStatus[family.Family](http.StatusOK),
				withMapping[family.Family](func(f family.Family) any {
					return newFamilyModel(userId, f)
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
