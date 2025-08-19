package endpoints

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/types"
)

type SubscriptionGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[subscription.Subscription]]
}

func NewSubscriptionGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[subscription.Subscription]]) *SubscriptionGetAllEndpoint {
	return &SubscriptionGetAllEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get all subscriptions
//	@Description	Retrieve a paginated list of all subscriptions for the authenticated user
//	@Tags			subscription
//	@Produce		json
//	@Param			search		query		string										false	"Search text"
//	@Param			sortBy		query		string										false	"Sort by field"
//	@Param			sortOrder	query		string										false	"Sort order (asc, desc)"
//	@Param			limit		query		integer										false	"Number of items per page (default: 10)"
//	@Param			offset		query		integer										false	"Page number (default: 0)"
//	@Success		200			{object}	PaginatedResponseModel[SubscriptionModel]	"Paginated list of subscriptions"
//	@Failure		400			{object}	HttpErrorResponse							"Bad Request - Invalid query parameters"
//	@Failure		500			{object}	HttpErrorResponse							"Internal Server Error"
//	@Router			/subscriptions [get]
func (s SubscriptionGetAllEndpoint) Handle(c *gin.Context) {
	searchText := c.DefaultQuery("search", "")
	limit, err := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 32)
	if err != nil {
		limit = 10
	}
	offset, err := strconv.ParseInt(c.DefaultQuery("offset", "o"), 10, 32)
	if err != nil {
		offset = 1
	}
	sortBy := subscription.ParseSortableField(c.DefaultQuery("sortBy", "name"))
	sortOrder := types.ParseSortOrder(c.DefaultQuery("sortOrder", "asc"))

	q := query.NewFindAllQuery(searchText, sortBy, sortOrder, int32(limit), int32(offset))
	r := s.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[core.PaginatedResponse[subscription.Subscription]](func(paginatedResult core.PaginatedResponse[subscription.Subscription]) any {
			return newPaginatedResponseModel(paginatedResult, newSubscriptionModel)
		}))
}

func (s SubscriptionGetAllEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (s SubscriptionGetAllEndpoint) Method() string {
	return http.MethodGet
}

func (s SubscriptionGetAllEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
