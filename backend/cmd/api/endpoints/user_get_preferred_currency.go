package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/user/query"
)

type UserGetPreferredCurrencyEndpoint struct {
	handler core.QueryHandler[query.FindPreferredCurrencyQuery, currency.Unit]
}

func NewUserGetPreferredCurrencyEndpoint(handler core.QueryHandler[query.FindPreferredCurrencyQuery, currency.Unit]) *UserGetPreferredCurrencyEndpoint {
	return &UserGetPreferredCurrencyEndpoint{handler: handler}
}

type UserPreferredCurrencyModel struct {
	Currency string `json:"currency"`
}

// Handle godoc
//
//	@Summary		Get user preferred currency
//	@Description	Returns the preferred currency for the authenticated user
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	UserPreferredCurrencyModel
//	@Failure		401	{object}	httpError	"Unauthorized"
//	@Router			/users/preferred/currency [get]
func (e UserGetPreferredCurrencyEndpoint) Handle(c *gin.Context) {
	q := query.NewFindPreferredCurrencyQuery()

	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[currency.Unit](func(up currency.Unit) any {
			return UserPreferredCurrencyModel{
				Currency: up.String(),
			}
		}))
}

func (e UserGetPreferredCurrencyEndpoint) Pattern() []string {
	return []string{
		"/preferred/currency",
	}
}

func (e UserGetPreferredCurrencyEndpoint) Method() string {
	return http.MethodGet
}

func (e UserGetPreferredCurrencyEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
