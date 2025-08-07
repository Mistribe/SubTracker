package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/user/query"
	"github.com/oleexo/subtracker/internal/domain/user"
)

type UserGetPreferredCurrencyEndpoint struct {
	handler core.QueryHandler[query.FindPreferredCurrencyQuery, user.Profile]
}

func NewUserGetPreferredCurrencyEndpoint(handler core.QueryHandler[query.FindPreferredCurrencyQuery, user.Profile]) *UserGetPreferredCurrencyEndpoint {
	return &UserGetPreferredCurrencyEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get user preferred currency
//	@Description	Returns the preferred currency for the authenticated user
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	UserProfileModel
//	@Failure		401	{object}	httpError	"Unauthorized"
//	@Router			/users/profile [get]
func (e UserGetPreferredCurrencyEndpoint) Handle(c *gin.Context) {
	q := query.NewFindPreferredCurrencyQuery()

	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[user.Profile](func(up user.Profile) any {
			return newUserProfileModel(up)
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
