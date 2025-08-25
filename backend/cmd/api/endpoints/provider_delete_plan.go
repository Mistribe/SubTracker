package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/application/core"
	"github.com/mistribe/subtracker/internal/application/provider/command"
)

type ProviderPlanDeleteEndpoint struct {
	handler core.CommandHandler[command.DeletePlanCommand, bool]
}

func NewProviderPlanDeleteEndpoint(handler core.CommandHandler[command.DeletePlanCommand, bool]) *ProviderPlanDeleteEndpoint {
	return &ProviderPlanDeleteEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Delete provider plan by ID
//	@Description	Permanently delete a provider plan and all its associated prices
//	@Tags			providers
//	@Param			providerId	path	string	true	"Provider ID (UUID format)"
//	@Param			planId		path	string	true	"Plan ID (UUID format)"
//	@Success		204			"No Content - Plan successfully deleted"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid ID format"
//	@Failure		404			{object}	HttpErrorResponse	"Provider or plan not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/{providerId}/plans/{planId} [delete]
func (e ProviderPlanDeleteEndpoint) Handle(c *gin.Context) {
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
