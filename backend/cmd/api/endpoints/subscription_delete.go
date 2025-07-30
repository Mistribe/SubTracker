package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/command"
)

type SubscriptionDeleteEndpoint struct {
	handler core.CommandHandler[command.DeleteSubscriptionCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete an existing subscription
//	@Description	Delete an existing subscription
//	@Tags			subscription
//	@Param			subscriptionId	path	string	true	"Subscription ID"
//	@Success		204				"No Content"
//	@Failure		400				{object}	httpError
//	@Router			/subscriptions/{subscriptionId} [delete]
func (s SubscriptionDeleteEndpoint) Handle(c *gin.Context) {
	id, err := uuid.Parse(c.Param("subscriptionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	cmd := command.DeleteSubscriptionCommand{
		Id: id,
	}

	r := s.handler.Handle(c, cmd)
	handleResponse(c, r)
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

func NewSubscriptionDeleteEndpoint(handler core.CommandHandler[command.DeleteSubscriptionCommand, bool]) *SubscriptionDeleteEndpoint {
	return &SubscriptionDeleteEndpoint{
		handler: handler,
	}
}
