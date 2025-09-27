package account

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/account/command"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type DeleteEndpoint struct {
	handler ports.CommandHandler[command.DeleteAccountCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete user
//	@Description	Deletes the authenticated user's account
//	@Tags			accounts
//	@Produce		json
//	@Success		204
//	@Failure		401	{object}	HttpErrorResponse
//	@Failure		500	{object}	HttpErrorResponse
//	@Router			/accounts [delete]
func (e DeleteEndpoint) Handle(c *gin.Context) {
	cmd := command.NewDeleteAccountCommand()

	r := e.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (e DeleteEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e DeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (e DeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewDeleteEndpoint(handler ports.CommandHandler[command.DeleteAccountCommand, bool]) *DeleteEndpoint {
	return &DeleteEndpoint{handler: handler}
}
