package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"net/http"
)

type SubscriptionGetEndpoint struct {
	handler core.QueryHandler[query.FindOneQuery, subscription.Subscription]
}

func (s SubscriptionGetEndpoint) Handle(c *gin.Context) {
	id, err := paramAsUuid(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
	}
	q := query.NewFindOneQuery(id)
	r := s.handler.Handle(c, q)
	handleResponse(c, r)
}

func (s SubscriptionGetEndpoint) Pattern() []string {
	return []string{
		"/:id",
	}
}

func (s SubscriptionGetEndpoint) Method() string {
	return http.MethodGet
}

func (s SubscriptionGetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionGetEndpoint() *SubscriptionGetEndpoint {
	return &SubscriptionGetEndpoint{}
}
