package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/user/command"
)

type UserDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteUserCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete user
//	@Description	Deletes the authenticated user's account
//	@Tags			users
//	@Produce		json
//	@Success		204
//	@Failure		401	{object}	HttpErrorResponse
//	@Failure		500	{object}	HttpErrorResponse
//	@Router			/users [delete]
func (e UserDeleteEndpoint) Handle(c *gin.Context) {
	cmd := command.NewDeleteUserCommand()

	r := e.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[bool]())
}

func (e UserDeleteEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e UserDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (e UserDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewUserDeleteEndpoint(handler core.CommandHandler[command.DeleteUserCommand, bool]) *UserDeleteEndpoint {
	return &UserDeleteEndpoint{handler: handler}
}
