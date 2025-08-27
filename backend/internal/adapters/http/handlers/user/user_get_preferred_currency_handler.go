package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/currency"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/user/query"
)

type UserGetPreferredCurrencyEndpoint struct {
	handler ports.QueryHandler[query.FindPreferredCurrencyQuery, currency.Unit]
}

func NewUserGetPreferredCurrencyEndpoint(handler ports.QueryHandler[query.FindPreferredCurrencyQuery, currency.Unit]) *UserGetPreferredCurrencyEndpoint {
	return &UserGetPreferredCurrencyEndpoint{handler: handler}
}

type UserPreferredCurrencyModel struct {
	Currency string `json:"currency"`
}

// Handle godoc
//
//	@Summary		Get user preferred currency
//	@Description	Returns the preferred currency for the authenticated user
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	UserPreferredCurrencyModel
//	@Failure		401	{object}	HttpErrorResponse	"Unauthorized"
//	@Router			/users/preferred/currency [get]
func (e UserGetPreferredCurrencyEndpoint) Handle(c *gin.Context) {
	q := query.NewFindPreferredCurrencyQuery()

	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[currency.Unit](func(up currency.Unit) any {
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
