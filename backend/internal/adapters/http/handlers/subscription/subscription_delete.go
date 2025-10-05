package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/command"
)

type DeleteEndpoint struct {
	handler ports.CommandHandler[command.DeleteSubscriptionCommand, bool]
}

// Handle godoc
//
//	@Summary		Delete subscription by LabelID
//	@Description	Permanently delete an existing subscription
//	@Tags			subscriptions
//	@Param			subscriptionId	path	string	true	"Subscription LabelID (UUID format)"
//	@Success		204				"No Content - Subscription successfully deleted"
//	@Failure		400				{object}	HttpErrorResponse	"Bad Request - Invalid subscription LabelID format"
//	@Failure		404				{object}	HttpErrorResponse	"Subscription not found"
//	@Failure		500				{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/subscriptions/{subscriptionId} [delete]
func (s DeleteEndpoint) Handle(c *gin.Context) {
	subscriptionID, err := types.ParseSubscriptionID(c.Param("subscriptionId"))
	if err != nil {
		FromError(c, err)
		return
	}

	cmd := command.DeleteSubscriptionCommand{
		SubscriptionID: subscriptionID,
	}

	r := s.handler.Handle(c, cmd)
	FromResult(c, r, WithNoContent[bool]())
}

func (s DeleteEndpoint) Pattern() []string {
	return []string{
		"/:subscriptionId",
	}
}

func (s DeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (s DeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewDeleteEndpoint(handler ports.CommandHandler[command.DeleteSubscriptionCommand, bool]) *DeleteEndpoint {
	return &DeleteEndpoint{
		handler: handler,
	}
}
