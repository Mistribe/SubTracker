package endpoints

import (
	"github.com/google/uuid"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
)

type ProviderDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteProviderCommand, bool]
}

func NewProviderDeleteEndpoint(handler core.CommandHandler[command.DeleteProviderCommand, bool]) *ProviderDeleteEndpoint {
	return &ProviderDeleteEndpoint{handler: handler}
}

func (e ProviderDeleteEndpoint) Handle(c *gin.Context) {
	providerId, err := uuid.Parse(c.Param("providerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	cmd := command.DeleteProviderCommand{
		ProviderId: providerId,
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[bool]())
}

func (e ProviderDeleteEndpoint) Pattern() []string {
	return []string{
		":providerId",
	}
}

func (e ProviderDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (e ProviderDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
