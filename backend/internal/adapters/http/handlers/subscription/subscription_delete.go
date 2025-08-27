package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
)

type SubscriptionDeleteEndpoint struct {
	handler ports.CommandHandler[command.DeleteSubscriptionCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete subscription by ID
//	@Description	Permanently delete an existing subscription
//	@Tags			subscription
//	@Param			subscriptionId	path	string	true	"Subscription ID (UUID format)"
//	@Success		204				"No Content - Subscription successfully deleted"
//	@Failure		400				{object}	HttpErrorResponse	"Bad Request - Invalid subscription ID format"
//	@Failure		404				{object}	HttpErrorResponse	"Subscription not found"
//	@Failure		500				{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/subscriptions/{subscriptionId} [delete]
func (s SubscriptionDeleteEndpoint) Handle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("subscriptionId"))
	if err != nil {
		FromError(c, err)
		return
	}

	cmd := command.DeleteSubscriptionCommand{
		Id: id,
	}

	r := s.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (s SubscriptionDeleteEndpoint) Pattern() []string {
	return []string{
		"/:subscriptionId",
	}
}

func (s SubscriptionDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (s SubscriptionDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionDeleteEndpoint(handler ports.CommandHandler[command.DeleteSubscriptionCommand, bool]) *SubscriptionDeleteEndpoint {
	return &SubscriptionDeleteEndpoint{
		handler: handler,
	}
}
