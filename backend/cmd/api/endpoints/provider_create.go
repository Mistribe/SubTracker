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
	"github.com/oleexo/subtracker/internal/domain/user"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type ProviderCreateEndpoint struct {
	handler core.CommandHandler[command.CreateProviderCommand, provider.Provider]
}

func NewProviderCreateEndpoint(handler core.CommandHandler[command.CreateProviderCommand, provider.Provider]) *ProviderCreateEndpoint {
	return &ProviderCreateEndpoint{handler: handler}
}

// @Description Price information for a plan
type createPriceModel struct {
	Id        *string    `json:"id,omitempty"`
	Currency  string     `json:"currency" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required" format:"date-time"`
	EndDate   *time.Time `json:"end_date,omitempty" format:"date-time"`
	Amount    float64    `json:"amount" binding:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty" format:"date-time"`
}

func (m createPriceModel) Price() (provider.Price, error) {
	id, err := parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
	cry, err := currency.ParseISO(m.Currency)
	if err != nil {
		return nil, err
	}
	createdAt := ext.ValueOrDefault(m.CreatedAt, time.Now())

	return provider.NewPrice(
		id,
		m.StartDate,
		m.EndDate,
		cry,
		m.Amount,
		createdAt,
		createdAt), nil
}

// @Description Plan information for a provider
type createPlanModel struct {
	Id          *string            `json:"id,omitempty"`
	Name        *string            `json:"name,omitempty"`
	Description *string            `json:"description,omitempty"`
	Prices      []createPriceModel `json:"prices" binding:"required"`
	CreatedAt   *time.Time         `json:"created_at,omitempty" format:"date-time"`
}

func (m createPlanModel) Plan() (provider.Plan, error) {
	id, err := parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
	createdAt := ext.ValueOrDefault(m.CreatedAt, time.Now())
	prices, err := slicesx.MapErr(m.Prices, func(price createPriceModel) (provider.Price, error) {
		return price.Price()
	})

	return provider.NewPlan(
		id,
		m.Name,
		m.Description,
		prices,
		createdAt,
		createdAt), nil
}

// @Description Provider creation request model
type createProviderModel struct {
	Id             *string             `json:"id,omitempty"`
	Name           string              `json:"name" binding:"required"`
	Description    *string             `json:"description,omitempty"`
	IconUrl        *string             `json:"icon_url,omitempty"`
	Url            *string             `json:"url,omitempty"`
	PricingPageUrl *string             `json:"pricing_page_url,omitempty"`
	Labels         []string            `json:"labels" binding:"required"`
	Plans          []createPlanModel   `json:"plans" binding:"required"`
	Owner          *editableOwnerModel `json:"owner,omitempty"`
	CreatedAt      *time.Time          `json:"created_at,omitempty" format:"date-time"`
}

func (m createProviderModel) Provider(userId string) (provider.Provider, error) {
	id, err := parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
	createdAt := ext.ValueOrDefault(m.CreatedAt, time.Now())
	plans, err := slicesx.MapErr(m.Plans, func(plan createPlanModel) (provider.Plan, error) {
		return plan.Plan()
	})
	if err != nil {
		return nil, err
	}
	labels, err := slicesx.MapErr(m.Labels, uuid.Parse)
	if err != nil {
		return nil, err
	}

	var owner user.Owner
	if m.Owner != nil {
		o, err := m.Owner.Owner(userId)
		if err != nil {
			return nil, err
		}
		owner = o
	}

	return provider.NewProvider(
		id,
		m.Name,
		m.Description,
		m.IconUrl,
		m.Url,
		m.PricingPageUrl,
		labels,
		plans,
		owner,
		createdAt,
		createdAt,
	), nil
}

func (m createProviderModel) Command(userId string) (command.CreateProviderCommand, error) {
	prov, err := m.Provider(userId)
	if err != nil {
		return command.CreateProviderCommand{}, err
	}

	return command.CreateProviderCommand{
		Provider: prov,
	}, nil
}

// Handle godoc
//
// @Summary Create a new provider
// @Description Create a new provider with plans and prices
// @Tags providers
// @Accept json
// @Produce json
// @Param provider body createProviderModel true "Provider information"
// @Success 201 {object} providerModel
// @Failure 400 {object} httpError
// @Router /providers [post]
func (e ProviderCreateEndpoint) Handle(c *gin.Context) {
	var model createProviderModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	userId, ok := user.FromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, httpError{
			Message: "invalid user id",
		})
		return
	}

	cmd, err := model.Command(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		c.Abort()
		return
	}

	r := e.handler.Handle(c, cmd)
	handleResponse(c,
		r,
		withStatus[provider.Provider](http.StatusCreated),
		withMapping[provider.Provider](func(prov provider.Provider) any {
			return newProviderModel(prov)
		}))
}

func (e ProviderCreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e ProviderCreateEndpoint) Method() string {
	return http.MethodPost
}

func (e ProviderCreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
