package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/user/command"
)

type UserUpdatePreferredCurrencyEndpoint struct {
	handler core.CommandHandler[command.UpdatePreferredCurrencyCommand, bool]
}

func NewUserUpdatePreferredCurrencyEndpoint(handler core.CommandHandler[command.UpdatePreferredCurrencyCommand, bool]) *UserUpdatePreferredCurrencyEndpoint {
	return &UserUpdatePreferredCurrencyEndpoint{
		handler: handler,
	}
}

type updateProfileModel struct {
	Currency string `json:"currency" binding:"required"`
}

// Handle godoc
//
//	@Summary		Update user preferred currency
//	@Description	Updates the preferred currency for the authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Bearer token"
//	@Param			request			body		updateProfileModel	true	"Profile update parameters"
//	@Success		200				{object}	UserProfileModel
//	@Failure		400				{object}	httpError
//	@Failure		401				{object}	httpError
//	@Router			/users/profile [put]
func (e UserUpdatePreferredCurrencyEndpoint) Handle(c *gin.Context) {
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

	cmd := command.NewUpdatePreferredCurrencyCommand(newCurrency)

	r := e.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[bool]())
}

func (e UserUpdatePreferredCurrencyEndpoint) Pattern() []string {
	return []string{
		"/profile",
	}
}

func (e UserUpdatePreferredCurrencyEndpoint) Method() string {
	return http.MethodPut
}

func (e UserUpdatePreferredCurrencyEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
