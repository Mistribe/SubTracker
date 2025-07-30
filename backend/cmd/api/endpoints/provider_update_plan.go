package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"net/http"
	"time"
)

type ProviderPlanUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdatePlanCommand, provider.Plan]
}

func NewProviderPlanUpdateEndpoint(handler core.CommandHandler[command.UpdatePlanCommand, provider.Plan]) *ProviderPlanUpdateEndpoint {
	return &ProviderPlanUpdateEndpoint{handler: handler}
}

type updatePlanModel struct {
	Name        string
	Description *string
	UpdateAt    *time.Time
}

func (m updatePlanModel) Command(providerId, planId uuid.UUID) (command.UpdatePlanCommand, error) {
	return command.UpdatePlanCommand{
		ProviderId:  providerId,
		Id:          planId,
		Name:        m.Name,
		Description: m.Description,
		UpdateAt:    m.UpdateAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Update provider plan by ID
//	@Description	Update an existing provider plan's information
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			providerId	path		string			true	"Provider ID (UUID format)"
//	@Param			planId		path		string			true	"Plan ID (UUID format)"
//	@Param			plan		body		updatePlanModel	true	"Updated plan data"
//	@Success		200			{object}	PlanModel		"Successfully updated plan"
//	@Failure		400			{object}	httpError		"Bad Request - Invalid input data or IDs"
//	@Failure		404			{object}	httpError		"Provider or plan not found"
//	@Failure		500			{object}	httpError		"Internal Server Error"
//	@Router			/providers/{providerId}/plans/{planId} [put]
func (e ProviderPlanUpdateEndpoint) Handle(c *gin.Context) {
	var model updatePlanModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

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

	cmd, err := model.Command(providerId, planId)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withStatus[provider.Plan](http.StatusOK),
		withMapping[provider.Plan](func(p provider.Plan) any {
			return newPlanModel(p)
		}))
}

func (e ProviderPlanUpdateEndpoint) Pattern() []string {
	return []string{
		":providerId/plans/:planId",
	}
}

func (e ProviderPlanUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (e ProviderPlanUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
