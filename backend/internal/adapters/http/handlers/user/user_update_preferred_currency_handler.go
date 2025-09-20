package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/currency"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/user/command"
	"github.com/mistribe/subtracker/pkg/ginx"
)

type UpdatePreferredCurrencyEndpoint struct {
	handler ports.CommandHandler[command.UpdatePreferredCurrencyCommand, bool]
}

func NewUpdatePreferredCurrencyEndpoint(handler ports.CommandHandler[command.UpdatePreferredCurrencyCommand, bool]) *UpdatePreferredCurrencyEndpoint {
	return &UpdatePreferredCurrencyEndpoint{
		handler: handler,
	}
}

// Handle godoc
//
//	@Summary		Update user preferred currency
//	@Description	Updates the preferred currency for the authenticated user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header	string								true	"Bearer token"
//	@Param			request			body	dto.UpdatePreferredCurrencyRequest	true	"Profile update parameters"
//	@Success		204
//	@Failure		400	{object}	HttpErrorResponse
//	@Failure		401	{object}	HttpErrorResponse
//	@Router			/users/preferred/currency [put]
func (e UpdatePreferredCurrencyEndpoint) Handle(c *gin.Context) {
	var model dto.UpdatePreferredCurrencyRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{Message: err.Error()})
		return
	}

	newCurrency, err := currency.ParseISO(model.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{Message: err.Error()})
		return
	}

	cmd := command.NewUpdatePreferredCurrencyCommand(newCurrency)

	r := e.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (e UpdatePreferredCurrencyEndpoint) Pattern() []string {
	return []string{
		"/preferred/currency",
	}
}

func (e UpdatePreferredCurrencyEndpoint) Method() string {
	return http.MethodPut
}

func (e UpdatePreferredCurrencyEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
