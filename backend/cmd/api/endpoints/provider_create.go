package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderCreateEndpoint struct {
	handler core.CommandHandler[command.CreateProviderCommand, provider.Provider]
}

func NewProviderCreateEndpoint(handler core.CommandHandler[command.CreateProviderCommand, provider.Provider]) *ProviderCreateEndpoint {
	return &ProviderCreateEndpoint{handler: handler}
}

func (e ProviderCreateEndpoint) Handle(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (e ProviderCreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e ProviderCreateEndpoint) Method() string {
	return http.MethodPost
}

func (e ProviderCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
