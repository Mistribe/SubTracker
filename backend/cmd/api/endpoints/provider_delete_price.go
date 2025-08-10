package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
)

type ProviderPriceDeleteEndpoint struct {
	handler core.CommandHandler[command.DeletePriceCommand, bool]
}

func NewProviderPriceDeleteEndpoint(handler core.CommandHandler[command.DeletePriceCommand, bool]) *ProviderPriceDeleteEndpoint {
	return &ProviderPriceDeleteEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Delete provider price by ID
//	@Description	Permanently delete a specific price from a provider plan
//	@Tags			providers
//	@Param			providerId	path	string	true	"Provider ID (UUID format)"
//	@Param			planId		path	string	true	"Plan ID (UUID format)"
//	@Param			priceId		path	string	true	"Price ID (UUID format)"
//	@Success		204			"No Content - Price successfully deleted"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid ID format"
//	@Failure		404			{object}	HttpErrorResponse	"Provider, plan, or price not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/{providerId}/plans/{planId}/prices/{priceId} [delete]
func (e ProviderPriceDeleteEndpoint) Handle(c *gin.Context) {
	providerId, err := uuid.Parse(c.Param("providerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	planId, err := uuid.Parse(c.Param("planId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	priceId, err := uuid.Parse(c.Param("priceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	cmd := command.DeletePriceCommand{
		ProviderId: providerId,
		PlanId:     planId,
		Id:         priceId,
		DeletedAt:  nil,
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[bool]())
}

func (e ProviderPriceDeleteEndpoint) Pattern() []string {
	return []string{
		":providerId/plans/:planId/prices/:priceId",
	}
}

func (e ProviderPriceDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (e ProviderPriceDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
