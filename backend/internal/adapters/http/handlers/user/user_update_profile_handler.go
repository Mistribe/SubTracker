package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/user/command"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type UserUpdateProfileEndpoint struct {
	handler ports.CommandHandler[command.UpdateProfileCommand, bool]
}

type UpdateProfileModel struct {
	GivenName  string `json:"given_name" binding:"required"`
	FamilyName string `json:"family_name" binding:"required"`
}

// Handle godoc
//
//	@Summary		Update user profile
//	@Description	Updates the given name and family name for the authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string				true	"Bearer token"
//	@Param			request			body	UpdateProfileModel	true	"Profile update parameters"
//	@Success		204
//	@Failure		400	{object}	HttpErrorResponse
//	@Failure		401	{object}	HttpErrorResponse
//	@Router			/users/profile [put]
func (e UserUpdateProfileEndpoint) Handle(c *gin.Context) {
	var model UpdateProfileModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cmd := command.NewUpdateProfileCommand(model.GivenName, model.FamilyName)
	r := e.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (e UserUpdateProfileEndpoint) Pattern() []string {
	return []string{
		"/profile",
	}
}

func (e UserUpdateProfileEndpoint) Method() string {
	return http.MethodPut
}

func (e UserUpdateProfileEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewUserUpdateProfileEndpoint(handler ports.CommandHandler[command.UpdateProfileCommand, bool]) *UserUpdateProfileEndpoint {
	return &UserUpdateProfileEndpoint{
		handler: handler,
	}
}
