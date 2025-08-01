package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/text/currency"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"github.com/oleexo/subtracker/internal/domain/auth"
	"github.com/oleexo/subtracker/internal/domain/provider"
	"github.com/oleexo/subtracker/pkg/ext"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type ProviderPatchEndpoint struct {
	handler core.CommandHandler[command.PatchProviderCommand, provider.Provider]
}

func NewProviderPatchEndpoint(handler core.CommandHandler[command.PatchProviderCommand, provider.Provider]) *ProviderPatchEndpoint {
	return &ProviderPatchEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Patch provider with plans and prices
//	@Description	Update or create a provider with all its plans and prices. If provider doesn't exist, it will be created.
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			provider	body		patchProviderModel	true	"Complete provider data with plans and prices"
//	@Success		200			{object}	ProviderModel		"Successfully updated provider"
//	@Failure		400			{object}	httpError			"Bad Request - Invalid input data"
//	@Failure		401			{object}	httpError			"Unauthorized - Invalid user authentication"
//	@Failure		404			{object}	httpError			"Provider not found"
//	@Failure		500			{object}	httpError			"Internal Server Error"
//	@Router			/providers [patch]

type patchPriceModel struct {
	Id        *string    `json:"id" `
	Currency  string     `json:"currency" binding:"required"`
	StartDate time.Time  `json:"start_date" binding:"required" format:"date-time"`
	EndDate   *time.Time `json:"end_date,omitempty" format:"date-time"`
	Amount    float64    `json:"amount" binding:"required"`
	CreatedAt *time.Time `json:"created_at,omitempty" format:"date-time"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchPriceModel) Price() (provider.Price, error) {
	id, err := parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
	cry, err := currency.ParseISO(m.Currency)
	if err != nil {
		return nil, err
	}
	createdAt := ext.ValueOrDefault(m.CreatedAt, time.Now())
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, createdAt)
	if updatedAt.Before(createdAt) {
		updatedAt = createdAt
	}
	return provider.NewPrice(
		id,
		m.StartDate,
		m.EndDate,
		cry,
		m.Amount,
		createdAt,
		updatedAt,
	), nil
}

type patchPlanModel struct {
	Id          *string           `json:"id"`
	Name        string            `json:"name" binding:"required"`
	Description *string           `json:"description,omitempty"`
	Prices      []patchPriceModel `json:"prices" binding:"required"`
	CreatedAt   *time.Time        `json:"created_at,omitempty"  format:"date-time"`
	UpdatedAt   *time.Time        `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchPlanModel) Plan() (provider.Plan, error) {
	id, err := parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
	createdAt := ext.ValueOrDefault(m.CreatedAt, time.Now())
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, createdAt)
	if updatedAt.Before(createdAt) {
		updatedAt = createdAt
	}
	prices, err := slicesx.SelectErr(m.Prices, func(prce patchPriceModel) (provider.Price, error) {
		return prce.Price()
	})
	if err != nil {
		return nil, err
	}
	return provider.NewPlan(
		id,
		m.Name,
		m.Description,
		prices,
		createdAt,
		updatedAt,
	), nil
}

type patchProviderModel struct {
	Id             *string            `json:"id" `
	Name           string             `json:"name" binding:"required"`
	Description    *string            `json:"description,omitempty"`
	IconUrl        *string            `json:"icon_url,omitempty"`
	Url            *string            `json:"url,omitempty"`
	PricingPageUrl *string            `json:"pricing_page_url,omitempty"`
	Labels         []string           `json:"labels" binding:"required"`
	Plans          []patchPlanModel   `json:"plans" binding:"required"`
	Owner          EditableOwnerModel `json:"owner" binding:"required"`
	CreatedAt      *time.Time         `json:"created_at,omitempty"  format:"date-time"`
	UpdatedAt      *time.Time         `json:"updated_at,omitempty" format:"date-time"`
}

func (m patchProviderModel) Provider(userId string) (provider.Provider, error) {
	id, err := parseUuidOrNew(m.Id)
	if err != nil {
		return nil, err
	}
	createdAt := ext.ValueOrDefault(m.CreatedAt, time.Now())
	updatedAt := ext.ValueOrDefault(m.UpdatedAt, createdAt)
	if updatedAt.Before(createdAt) {
		updatedAt = createdAt
	}
	owner, err := m.Owner.Owner(userId)
	if err != nil {
		return nil, err
	}

	labels, err := slicesx.SelectErr(m.Labels, uuid.Parse)
	if err != nil {
		return nil, err
	}

	plans, err := slicesx.SelectErr(m.Plans, func(prce patchPlanModel) (provider.Plan, error) {
		return prce.Plan()
	})
	if err != nil {
		return nil, err
	}

	return provider.NewProvider(
		id,
		m.Name,
		nil,
		m.Description,
		m.IconUrl,
		m.Url,
		m.PricingPageUrl,
		labels,
		plans,
		owner,
		createdAt,
		updatedAt,
	), nil
}
func (m patchProviderModel) Command(userId string) (command.PatchProviderCommand, error) {
	prov, err := m.Provider(userId)
	if err != nil {
		return command.PatchProviderCommand{}, err
	}
	return command.PatchProviderCommand{
		Provider: prov,
	}, nil
}

func (e ProviderPatchEndpoint) Handle(c *gin.Context) {
	var model patchProviderModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
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
		withStatus[provider.Provider](http.StatusOK),
		withMapping[provider.Provider](func(f provider.Provider) any {
			return newProviderModel(f)
		}))
}

func (e ProviderPatchEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e ProviderPatchEndpoint) Method() string {
	return http.MethodPatch
}

func (e ProviderPatchEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
