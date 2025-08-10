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
	"github.com/oleexo/subtracker/pkg/langext/option"
)

type FamilyMemberUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateFamilyMemberCommand, family.Family]
}

type updateFamilyMemberModel struct {
	Name      string     `json:"name" binding:"required"`
	Type      string     `json:"type" binding:"required" enums:"owner,adult,kid"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m updateFamilyMemberModel) Command(familyId, memberId uuid.UUID) (command.UpdateFamilyMemberCommand, error) {

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
//	@Param			familyId	path		string					true	"Family ID (UUID format)"
//	@Param			id			path		string					true	"Family member ID (UUID format)"
//	@Param			member		body		updateFamilyMemberModel	true	"Updated family member data"
//	@Success		200			{object}	familyModel				"Successfully updated family member"
//	@Failure		400			{object}	HttpErrorResponse		"Bad Request - Invalid input data or ID format"
//	@Failure		401			{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse		"Family or family member not found"
//	@Failure		500			{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/families/{familyId}/members/{id} [put]
func (f FamilyMemberUpdateEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("id")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: "invalid id format",
		})
		return
	}

	familyId, err := uuid.Parse(c.Param("familyId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, HttpErrorResponse{
			Message: "invalid user id",
		})
		return
	}

	var model updateFamilyMemberModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	cmd, err := model.Command(familyId, id)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	r := f.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withMapping[family.Family](func(mbr family.Family) any {
			return newFamilyModel(userId, mbr)
		}))
}

func (f FamilyMemberUpdateEndpoint) Pattern() []string {
	return []string{
		"/:familyId/members/:id",
	}
}

func (f FamilyMemberUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (f FamilyMemberUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewFamilyMemberUpdateEndpoint(handler core.CommandHandler[command.UpdateFamilyMemberCommand, family.Family]) *FamilyMemberUpdateEndpoint {
	return &FamilyMemberUpdateEndpoint{
		handler: handler,
	}
}
