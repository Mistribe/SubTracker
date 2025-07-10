package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/ext"
)

type SubscriptionGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, []subscription.Subscription]
}

// Handle godoc
// @Summary		Get all subscriptions
// @Description	Get all subscriptions
// @Tags			subscription
// @Produce		json
// @Success		200	{array}		subscriptionModel
// @Failure		400	{object}	httpError
// @Router			/subscriptions [get]
func (s SubscriptionGetAllEndpoint) Handle(c *gin.Context) {
	q := query.NewFindAllQuery()
	r := s.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[[]subscription.Subscription](func(subs []subscription.Subscription) any {
			return ext.Map[subscription.Subscription, subscriptionModel](
				subs,
				func(sub subscription.Subscription) subscriptionModel {
					return newSubscriptionModel(sub)
				},
			)
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

func NewSubscriptionGetAllEndpoint() *SubscriptionGetAllEndpoint {
	return &SubscriptionGetAllEndpoint{}
}
