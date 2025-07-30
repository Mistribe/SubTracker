package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"net/http"
)

type ProviderPlanDeleteEndpoint struct {
	handler core.CommandHandler[command.DeletePlanCommand, bool]
}

func NewProviderPlanDeleteEndpoint(handler core.CommandHandler[command.DeletePlanCommand, bool]) *ProviderPlanDeleteEndpoint {
	return &ProviderPlanDeleteEndpoint{handler: handler}
}

func (e ProviderPlanDeleteEndpoint) Handle(c *gin.Context) {
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

	cmd := command.DeletePlanCommand{
		ProviderId: providerId,
		Id:         planId,
		DeletedAt:  nil,
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c, r, withNoContent[bool]())
}

func (e ProviderPlanDeleteEndpoint) Pattern() []string {
	return []string{
		":providerId/plans/:planId",
	}
}

func (e ProviderPlanDeleteEndpoint) Method() string {
	return http.MethodDelete
}

func (e ProviderPlanDeleteEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
