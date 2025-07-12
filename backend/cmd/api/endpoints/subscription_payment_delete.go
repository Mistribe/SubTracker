package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/command"
	"github.com/oleexo/subtracker/pkg/langext/result"
)

type SubscriptionPaymentDeleteEndpoint struct {
	handler core.CommandHandler[command.DeletePaymentCommand, result.Unit]
}

// Handle godoc
//
//	@Summary		Delete a subscription payment
//	@Description	Delete a subscription payment
//	@Tags			subscription
//	@Accept			json
//	@Produce		json
//	@Param			id			path	string	true	"Subscription ID"
//	@Param			paymentId	path	string	true	"Payment ID"
//	@Success		204
//	@Failure		400	{object}	httpError
//	@Router			/subscriptions/{id}/payments/{paymentId} [delete]
func (s SubscriptionPaymentDeleteEndpoint) Handle(c *gin.Context) {
	subscriptionId, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	paymentId, err := uuid.Parse(c.Param("paymentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	cmd := command.DeletePaymentCommand{
		SubscriptionId: subscriptionId,
		PaymentId:      paymentId,
	}

	r := s.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[result.Unit]())
}

func (s SubscriptionPaymentDeleteEndpoint) Pattern() []string {
	return []string{
		"/:id/payments/:paymentId",
	}
}

func (s SubscriptionPaymentDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (s SubscriptionPaymentDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewSubscriptionPaymentDeleteEndpoint(
	handler core.CommandHandler[command.DeletePaymentCommand, result.Unit],
) *SubscriptionPaymentDeleteEndpoint {
	return &SubscriptionPaymentDeleteEndpoint{
		handler: handler,
	}
}
