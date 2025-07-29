package endpoints

import (
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
	//TODO implement me
	panic("implement me")
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
