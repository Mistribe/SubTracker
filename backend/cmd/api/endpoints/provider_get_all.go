package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/query"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[provider.Provider]]
}

func NewProviderGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[provider.Provider]]) *ProviderGetAllEndpoint {
	return &ProviderGetAllEndpoint{handler: handler}
}

func (e ProviderGetAllEndpoint) Handle(c *gin.Context) {
	//TODO implement me
	panic("implement me")
}

func (e ProviderGetAllEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e ProviderGetAllEndpoint) Method() string {
	return http.MethodGet
}

func (e ProviderGetAllEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
