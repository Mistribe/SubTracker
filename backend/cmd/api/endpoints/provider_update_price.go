package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderPriceUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdatePriceCommand, provider.Price]
}

func NewProviderPriceUpdateEndpoint(handler core.CommandHandler[command.UpdatePriceCommand, provider.Price]) *ProviderPriceUpdateEndpoint {
	return &ProviderPriceUpdateEndpoint{handler: handler}
}

type updatePriceModel struct {
	StartDate time.Time  `json:"start_date" binding:"required" format:"date-time"`
	EndDate   *time.Time `json:"end_date,omitempty" format:"date-time"`
	Currency  string     `json:"currency" binding:"required"`
	Amount    float64    `json:"amount" binding:"required"`
	UpdatedAt *time.Time `json:"created_at,omitempty" format:"date-time"`
}

func (m updatePriceModel) Command(providerId, planId, priceId uuid.UUID) (command.UpdatePriceCommand, error) {
	cry, err := currency.ParseISO(m.Currency)
	if err != nil {
		return command.UpdatePriceCommand{}, err
	}
	return command.UpdatePriceCommand{
		ProviderId: providerId,
		PlanId:     planId,
		Id:         priceId,
		StartDate:  m.StartDate,
		EndDate:    m.EndDate,
		Currency:   cry,
		Amount:     m.Amount,
		UpdatedAt:  m.UpdatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Update provider price by ID
//	@Description	Update an existing price for a specific provider plan
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			providerId	path		string				true	"Provider ID (UUID format)"
//	@Param			planId		path		string				true	"Plan ID (UUID format)"
//	@Param			priceId		path		string				true	"Price ID (UUID format)"
//	@Param			price		body		updatePriceModel	true	"Updated price data"
//	@Success		200			{object}	PriceModel			"Successfully updated price"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid input data or IDs"
//	@Failure		404			{object}	HttpErrorResponse	"Provider, plan, or price not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/{providerId}/plans/{planId}/prices/{priceId} [put]
func (e ProviderPriceUpdateEndpoint) Handle(c *gin.Context) {
	var model updatePriceModel
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

	cmd, err := model.Command(providerId, planId, priceId)
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withStatus[provider.Price](http.StatusOK),
		withMapping[provider.Price](func(p provider.Price) any {
			return newPriceModel(p)
		}))
}

func (e ProviderPriceUpdateEndpoint) Pattern() []string {
	return []string{
		":providerId/plans/:planId/prices/:priceId",
	}
}

func (e ProviderPriceUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (e ProviderPriceUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
