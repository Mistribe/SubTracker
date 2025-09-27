package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
)

type GetEndpoint struct {
	handler ports.QueryHandler[query.FindOneQuery, subscription.Subscription]
}

// Handle godoc
//
//	@Summary		Get subscription by LabelID
//	@Description	Retrieve a single subscription with all its details including provider, plan, and pricing information
//	@Tags			subscriptions
//	@Produce		json
//	@Param			subscriptionId	path		string					true	"Subscription LabelID (UUID format)"
//	@Success		200				{object}	dto.SubscriptionModel	"Successfully retrieved subscription"
//	@Failure		400				{object}	HttpErrorResponse		"Bad Request - Invalid subscription LabelID format"
//	@Failure		404				{object}	HttpErrorResponse		"Subscription not found"
//	@Failure		500				{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/subscriptions/{subscriptionId} [get]
func (e GetEndpoint) Handle(c *gin.Context) {
	id, err := QueryParamAsUUID(c, "subscriptionId")
	if err != nil {
		FromError(c, err)
		return
	}
	q := query.NewFindOneQuery(id)
	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[subscription.Subscription](func(sub subscription.Subscription) any {
			return dto.NewSubscriptionModel(sub)
		}))
}

func (e GetEndpoint) Pattern() []string {
	return []string{
		"/:subscriptionId",
	}
}

func (e GetEndpoint) Method() string {
	return http.MethodGet
}

func (e GetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewGetEndpoint() *GetEndpoint {
	return &GetEndpoint{}
}
