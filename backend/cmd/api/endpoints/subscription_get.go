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
//
//	@Summary		Get subscription by ID
//	@Description	Retrieve a single subscription with all its details including provider, plan, and pricing information
//	@Tags			subscription
//	@Produce		json
//	@Param			subscriptionId	path		string				true	"Subscription ID (UUID format)"
//	@Success		200				{object}	SubscriptionModel	"Successfully retrieved subscription"
//	@Failure		400				{object}	HttpErrorResponse	"Bad Request - Invalid subscription ID format"
//	@Failure		404				{object}	HttpErrorResponse	"Subscription not found"
//	@Failure		500				{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/subscriptions/{subscriptionId} [get]
func (e SubscriptionGetEndpoint) Handle(c *gin.Context) {
	id, err := paramAsUuid(c, "subscriptionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
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
