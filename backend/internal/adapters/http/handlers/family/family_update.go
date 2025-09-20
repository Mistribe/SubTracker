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

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/family/command"
	"github.com/mistribe/subtracker/pkg/langext/option"

	"github.com/mistribe/subtracker/internal/domain/family"
)

type UpdateEndpoint struct {
	handler ports.CommandHandler[command.UpdateFamilyCommand, family.Family]
}

func updateFamilyRequestToCommand(m dto.UpdateFamilyRequest, familyId uuid.UUID) (command.UpdateFamilyCommand, error) {
	return command.UpdateFamilyCommand{
		Id:        familyId,
		Name:      m.Name,
		UpdatedAt: option.Some[time.Time](x.ValueOrDefault[time.Time](m.UpdatedAt, time.Now())),
	}, nil

}

// Handle godoc
//
//	@Summary		Update a family
//	@Description	Update family information such as name and other details
//	@Tags			family
//	@Accept			json
//	@Produce		json
//	@Param			familyId	path		string					true	"Family ID (UUID format)"
//	@Param			family		body		dto.UpdateFamilyRequest	true	"Updated family data"
//	@Success		200			{object}	dto.FamilyModel			"Successfully updated family"
//	@Failure		400			{object}	HttpErrorResponse		"Bad Request - Invalid input data or family ID"
//	@Failure		401			{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	HttpErrorResponse		"Family not found"
//	@Failure		500			{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/family/{familyId} [put]
func (f UpdateEndpoint) Handle(c *gin.Context) {
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

	var model dto.UpdateFamilyRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}

	cmd, err := updateFamilyRequestToCommand(model, familyId)
	if err != nil {
		FromError(c, err)
		return
	}
	r := f.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[family.Family](http.StatusOK),
		WithMapping[family.Family](func(f family.Family) any {
			return dto.NewFamilyModel(userId, f)
		}))
}

func (f UpdateEndpoint) Pattern() []string {
	return []string{
		"/:familyId",
	}
}

func (f UpdateEndpoint) Method() string {
	return http.MethodPut
}

func (f UpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewUpdateEndpoint(handler ports.CommandHandler[command.UpdateFamilyCommand, family.Family]) *UpdateEndpoint {
	return &UpdateEndpoint{
		handler: handler,
	}
}
