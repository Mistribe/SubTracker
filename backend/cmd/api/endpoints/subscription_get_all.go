package endpoints

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/subscription"
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
//	@Description	Get all subscriptions
//	@Tags			subscription
//	@Produce		json
//	@Param			size	query		integer	false	"Number of items per page"
//	@Param			page	query		integer	false	"Page number"
//	@Success		200		{object}	paginatedResponseModel[subscriptionModel]
//	@Failure		400		{object}	httpError
//	@Router			/subscriptions [get]
func (s SubscriptionGetAllEndpoint) Handle(c *gin.Context) {
	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		size = 10
	}
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	q := query.NewFindAllQuery(size, page)
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
