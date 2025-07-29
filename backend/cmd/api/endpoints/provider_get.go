package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/query"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderGetEndpoint struct {
	handler core.QueryHandler[query.FindOneQuery, provider.Provider]
}

func NewProviderGetEndpoint(handler core.QueryHandler[query.FindOneQuery, provider.Provider]) *ProviderGetEndpoint {
	return &ProviderGetEndpoint{handler: handler}
}

func (e ProviderGetEndpoint) Handle(c *gin.Context) {
	id, err := paramAsUuid(c, "providerId")
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
	}
	q := query.NewFindOneQuery(id)
	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[provider.Provider](func(prvdr provider.Provider) any {
			return newProviderModel(prvdr)
		}))
}

func (e ProviderGetEndpoint) Pattern() []string {
	return []string{
		":providerId",
	}
}

func (e ProviderGetEndpoint) Method() string {
	return http.MethodGet
}

func (e ProviderGetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
