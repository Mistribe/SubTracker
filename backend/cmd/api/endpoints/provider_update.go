package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateProviderCommand, provider.Provider]
}

func NewProviderUpdateEndpoint(handler core.CommandHandler[command.UpdateProviderCommand, provider.Provider]) *ProviderUpdateEndpoint {
	return &ProviderUpdateEndpoint{handler: handler}
}

func (e ProviderUpdateEndpoint) Handle(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (e ProviderUpdateEndpoint) Pattern() []string {
	return []string{
		":providerId",
	}
}

func (e ProviderUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (e ProviderUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
