package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/application/auth/command"
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/auth"
)

type UserUpdateProfileEndpoint struct {
	handler core.CommandHandler[command.UpdateProfileCommand, auth.UserProfile]
}

func NewUserUpdateProfileEndpoint(handler core.CommandHandler[command.UpdateProfileCommand, auth.UserProfile]) *UserUpdateProfileEndpoint {
	return &UserUpdateProfileEndpoint{
		handler: handler,
	}
}

type updateProfileModel struct {
	Currency string `json:"currency" binding:"required"`
}

// Handle godoc
//
//	@Summary		Update user profile
//	@Description	Updates the currency preference in user's profile
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer token"
//	@Param			request			body		updateProfileModel	true	"Profile update parameters"
//	@Success		200				{object}	UserProfileModel
//	@Failure		400				{object}	httpError
//	@Failure		401				{object}	httpError
//	@Router			/users/profile [put]
func (e UserUpdateProfileEndpoint) Handle(c *gin.Context) {
	var model updateProfileModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{Message: err.Error()})
		return
	}

	newCurrency, err := currency.ParseISO(model.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{Message: err.Error()})
		return
	}

	cmd := command.UpdateProfileCommand{
		Currency: newCurrency,
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c, r, withMapping[auth.UserProfile](func(result auth.UserProfile) any {
		return result
	}))
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
