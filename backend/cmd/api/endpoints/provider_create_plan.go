package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderPlanCreateEndpoint struct {
	handler core.CommandHandler[command.CreatePlanCommand, provider.Plan]
}

func NewProviderPlanCreateEndpoint(handler core.CommandHandler[command.CreatePlanCommand, provider.Plan]) *ProviderPlanCreateEndpoint {
	return &ProviderPlanCreateEndpoint{handler: handler}
}

type createPlanModel struct {
	Id          *string    `json:"id,omitempty"`
	Name        string     `json:"name" binding:"required"`
	Description *string    `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty" format:"date-time"`
}

func (m createPlanModel) Command(providerId uuid.UUID) (command.CreatePlanCommand, error) {
	var planId *uuid.UUID
	if m.Id != nil {
		id, err := uuid.Parse(*m.Id)
		if err != nil {
			return command.CreatePlanCommand{}, err
		}

		planId = &id
	}

	return command.CreatePlanCommand{
		ProviderId: providerId,

		Id:          planId,
		Name:        m.Name,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Create a new provider plan
//	@Description	Create a new subscription plan for an existing provider
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			providerId	path		string			true	"Provider ID (UUID format)"
//	@Param			plan		body		createPlanModel	true	"Plan creation data"
//	@Success		201			{object}	PlanModel		"Successfully created plan"
//	@Failure		400			{object}	HttpErrorResponse		"Bad Request - Invalid input data or provider ID"
//	@Failure		404			{object}	HttpErrorResponse		"Provider not found"
//	@Failure		500			{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/providers/{providerId}/plans [post]
func (e ProviderPlanCreateEndpoint) Handle(c *gin.Context) {
	var model createPlanModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	providerId, err := uuid.Parse(c.Param("providerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	cmd, err := model.Command(providerId)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withStatus[provider.Plan](http.StatusCreated),
		withMapping[provider.Plan](func(p provider.Plan) any {
			return newPlanModel(p)
		}))
}

func (e ProviderPlanCreateEndpoint) Pattern() []string {
	return []string{
		":providerId/plans",
	}
}

func (e ProviderPlanCreateEndpoint) Method() string {
	return http.MethodPost
}

func (e ProviderPlanCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
