package provider

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/provider/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type GetAllEndpoint struct {
	handler ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[provider.Provider]]
}

func NewGetAllEndpoint(handler ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[provider.Provider]]) *GetAllEndpoint {
	return &GetAllEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get all providers
//	@Description	Retrieve a paginated list of all providers with their plans and prices
//	@Tags			providers
//	@Produce		json
//	@Param			search	query		string										false	"Search term"
//	@Param			offset	query		int											false	"Offset (default: 0)"
//	@Param			limit	query		int											false	"Limit per request (default: 10)"
//	@Success		200		{object}	dto.PaginatedResponseModel[ProviderModel]	"Paginated list of providers"
//	@Failure		400		{object}	HttpErrorResponse							"Bad Request - Invalid query parameters"
//	@Failure		500		{object}	HttpErrorResponse							"Internal Server Error"
//	@Router			/providers [get]
func (e GetAllEndpoint) Handle(c *gin.Context) {
	search := c.DefaultQuery("search", "")
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	if err != nil {
		offset = 0
	}
	q := query.NewFindAllQuery(search, limit, offset)
	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[shared.PaginatedResponse[provider.Provider]](func(paginatedResult shared.PaginatedResponse[provider.Provider]) any {
			return dto.NewPaginatedResponseModel(paginatedResult, dto.NewProviderModel)
		}))
}

func (e GetAllEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e GetAllEndpoint) Method() string {
	return http.MethodGet
}

func (e GetAllEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
