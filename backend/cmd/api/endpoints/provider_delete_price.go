package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"net/http"
)

type ProviderPriceDeleteEndpoint struct {
	handler core.CommandHandler[command.DeletePriceCommand, bool]
}

func NewProviderPriceDeleteEndpoint(handler core.CommandHandler[command.DeletePriceCommand, bool]) *ProviderPriceDeleteEndpoint {
	return &ProviderPriceDeleteEndpoint{handler: handler}
}

func (e ProviderPriceDeleteEndpoint) Handle(c *gin.Context) {
	providerId, err := uuid.Parse(c.Param("providerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}
	planId, err := uuid.Parse(c.Param("planId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	priceId, err := uuid.Parse(c.Param("priceId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
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
