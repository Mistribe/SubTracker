package provider

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
)

type DeleteEndpoint struct {
	handler ports.CommandHandler[command.DeleteProviderCommand, bool]
}

func NewDeleteEndpoint(handler ports.CommandHandler[command.DeleteProviderCommand, bool]) *DeleteEndpoint {
	return &DeleteEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Delete provider by LabelID
//	@Description	Permanently delete a provider and all its associated plans and prices
//	@Tags			providers
//	@Param			providerId	path	string	true	"Provider LabelID (UUID format)"
//	@Success		204			"No Content - Provider successfully deleted"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid provider LabelID format"
//	@Failure		404			{object}	HttpErrorResponse	"Provider not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/{providerId} [delete]
func (e DeleteEndpoint) Handle(c *gin.Context) {
	providerId, err := types.ParseProviderID(c.Param("providerId"))
	if err != nil {
		FromError(c, err)
		return
	}

	cmd := command.DeleteProviderCommand{
		ProviderID: providerId,
	}

	r := e.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (e DeleteEndpoint) Pattern() []string {
	return []string{
		":providerId",
	}
}

func (e DeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (e DeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
