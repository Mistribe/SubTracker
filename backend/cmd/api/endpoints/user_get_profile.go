package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/auth/query"
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/domain/auth"
)

type UserGetProfileEndpoint struct {
	handler core.QueryHandler[query.FindProfileQuery, auth.UserProfile]
}

func NewUserGetProfileEndpoint(handler core.QueryHandler[query.FindProfileQuery, auth.UserProfile]) *UserGetProfileEndpoint {
	return &UserGetProfileEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get user profile
//	@Description	Returns the profile information for the authenticated user
//	@Tags			Users
//	@Produce		json
//	@Success		200	{object}	UserProfileModel
//	@Failure		401	{object}	httpError	"Unauthorized"
//	@Router			/users/profile [get]
func (e UserGetProfileEndpoint) Handle(c *gin.Context) {
	q := query.FindProfileQuery{}

	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[auth.UserProfile](func(up auth.UserProfile) any {
			return newUserProfileModel(up)
		}))
}

func (e UserGetProfileEndpoint) Pattern() []string {
	return []string{
		"/profile",
	}
}

func (e UserGetProfileEndpoint) Method() string {
	return http.MethodGet
}

func (e UserGetProfileEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
