package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type SubscriptionGetEndpoint struct {
	handler core.QueryHandler[query.FindOneQuery, subscription.Subscription]
}

// Handle godoc
// @Summary		Get subscription by ID
// @Description	Get subscription by ID
// @Tags			subscription
// @Produce		json
// @Param			subscriptionId	path		string	true	"Subscription ID"
// @Success		200	{object}	subscriptionModel
// @Failure		400	{object}	httpError
// @Failure		404	{object}	httpError
// @Router			/subscriptions/{subscriptionId} [get]
func (e SubscriptionGetEndpoint) Handle(c *gin.Context) {
	id, err := paramAsUuid(c, "subscriptionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
	}
	q := query.NewFindOneQuery(id)
	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[subscription.Subscription](func(sub subscription.Subscription) any {
			return newSubscriptionModel(sub)
		}))
}

func (e SubscriptionGetEndpoint) Pattern() []string {
	return []string{
		"/:subscriptionId",
	}
}

func (e SubscriptionGetEndpoint) Method() string {
	return http.MethodGet
}

func (e SubscriptionGetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionGetEndpoint() *SubscriptionGetEndpoint {
	return &SubscriptionGetEndpoint{}
}
