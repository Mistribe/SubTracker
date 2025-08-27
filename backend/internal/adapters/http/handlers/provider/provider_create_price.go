package provider

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/text/currency"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
	"github.com/mistribe/subtracker/pkg/ginx"
)

type ProviderPriceCreateEndpoint struct {
	handler ports.CommandHandler[command.CreatePriceCommand, provider.Price]
}

func NewProviderPriceCreateEndpoint(handler ports.CommandHandler[command.CreatePriceCommand, provider.Price]) *ProviderPriceCreateEndpoint {
	return &ProviderPriceCreateEndpoint{handler: handler}
}

// @Description	Price information for a plan
type createPriceModel struct {
	Id        *string    `json:"id,omitempty"`
	Currency  string     `json:"currency" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required" format:"date-time"`
	EndDate   *time.Time `json:"end_date,omitempty" format:"date-time"`
	Amount    float64    `json:"amount" binding:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty" format:"date-time"`
}

func (m createPriceModel) Command(providerId, planId uuid.UUID) (command.CreatePriceCommand, error) {
	var priceId *uuid.UUID
	if m.Id != nil {
		id, err := uuid.Parse(*m.Id)
		if err != nil {
			return command.CreatePriceCommand{}, err
		}
		priceId = &id
	}
	cry, err := currency.ParseISO(m.Currency)
	if err != nil {
		return command.CreatePriceCommand{}, err
	}

	return command.CreatePriceCommand{
		ProviderId: providerId,
		PlanId:     planId,
		Id:         priceId,
		StartDate:  m.StartDate,
		EndDate:    m.EndDate,
		Currency:   cry,
		Amount:     m.Amount,
		CreatedAt:  m.CreatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Create a new provider price
//	@Description	Create a new pricing option for a specific provider plan
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			providerId	path		string				true	"Provider ID (UUID format)"
//	@Param			planId		path		string				true	"Plan ID (UUID format)"
//	@Param			price		body		createPriceModel	true	"Price creation data"
//	@Success		201			{object}	PriceModel			"Successfully created price"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid input data or IDs"
//	@Failure		404			{object}	HttpErrorResponse	"Provider or plan not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/{providerId}/plans/{planId}/prices [post]
func (e ProviderPriceCreateEndpoint) Handle(c *gin.Context) {
	var model createPriceModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	providerId, err := uuid.Parse(c.Param("providerId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	planId, err := uuid.Parse(c.Param("planId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	cmd, err := model.Command(providerId, planId)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	r := e.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithStatus[provider.Price](http.StatusCreated),
		WithMapping[provider.Price](func(p provider.Price) any {
			return newPriceModel(p)
		}))
}

func (e ProviderPriceCreateEndpoint) Pattern() []string {
	return []string{
		":providerId/plans/:planId/prices",
	}
}

func (e ProviderPriceCreateEndpoint) Method() string {
	return http.MethodPost
}

func (e ProviderPriceCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
