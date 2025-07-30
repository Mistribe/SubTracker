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
//	@Param			page	query		int										false	"Page number (default: 1)"
//	@Param			size	query		int										false	"Items per page (default: 10)"
//	@Success		200		{object}	PaginatedResponseModel[ProviderModel]	"Paginated list of providers"
//	@Failure		400		{object}	httpError								"Bad Request - Invalid query parameters"
//	@Failure		500		{object}	httpError								"Internal Server Error"
//	@Router			/providers [get]
func (e ProviderGetAllEndpoint) Handle(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		size = 10
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	q := query.NewFindAllQuery(size, page)
	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[core.PaginatedResponse[provider.Provider]](func(paginatedResult core.PaginatedResponse[provider.Provider]) any {
			return newPaginatedResponseModel(paginatedResult, newProviderModel)
		}))
}

func (e ProviderGetAllEndpoint) Pattern() []string {
	return []string{
		"/api/v1/providers",
	}
}

func (e ProviderGetAllEndpoint) Method() string {
	return http.MethodGet
}

func (e ProviderGetAllEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
