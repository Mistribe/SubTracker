package endpoints

import (
	"net/http"
	"strconv"

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

// Handle godoc
//
//	@Summary		Get all providers
//	@Description	Retrieve a paginated list of all providers with their plans and prices
//	@Tags			providers
//	@Produce		json
//	@Param			offset	query		int										false	"Offset (default: 0)"
//	@Param			limit	query		int										false	"Limit per request (default: 10)"
//	@Success		200		{object}	PaginatedResponseModel[ProviderModel]	"Paginated list of providers"
//	@Failure		400		{object}	httpError								"Bad Request - Invalid query parameters"
//	@Failure		500		{object}	httpError								"Internal Server Error"
//	@Router			/providers [get]
func (e ProviderGetAllEndpoint) Handle(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil {
		limit = 10
	}
	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil {
		offset = 1
	}
	q := query.NewFindAllQuery(limit, offset)
	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[core.PaginatedResponse[provider.Provider]](func(paginatedResult core.PaginatedResponse[provider.Provider]) any {
			return newPaginatedResponseModel(paginatedResult, newProviderModel)
		}))
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
