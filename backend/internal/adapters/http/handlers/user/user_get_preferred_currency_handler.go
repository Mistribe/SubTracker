package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/currency"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/user/query"
)

type GetPreferredCurrencyEndpoint struct {
	handler ports.QueryHandler[query.FindPreferredCurrencyQuery, currency.Unit]
}

func NewUserGetPreferredCurrencyEndpoint(handler ports.QueryHandler[query.FindPreferredCurrencyQuery, currency.Unit]) *GetPreferredCurrencyEndpoint {
	return &GetPreferredCurrencyEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get user preferred currency
//	@Description	Returns the preferred currency for the authenticated user
//	@Tags			users
//	@Produce		json
//	@Success		200	{object}	dto.UserPreferredCurrencyModel
//	@Failure		401	{object}	HttpErrorResponse	"Unauthorized"
//	@Router			/users/preferred/currency [get]
func (e GetPreferredCurrencyEndpoint) Handle(c *gin.Context) {
	q := query.NewFindPreferredCurrencyQuery()

	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[currency.Unit](func(up currency.Unit) any {
			return dto.UserPreferredCurrencyModel{
				Currency: up.String(),
			}
		}))
}

func (e GetPreferredCurrencyEndpoint) Pattern() []string {
	return []string{
		"/preferred/currency",
	}
}

func (e GetPreferredCurrencyEndpoint) Method() string {
	return http.MethodGet
}

func (e GetPreferredCurrencyEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
