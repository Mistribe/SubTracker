package provider

import (
	"net/http"
	"time"

	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x/collection"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
	"github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/provider"
)

type ProviderUpdateEndpoint struct {
	handler ports.CommandHandler[command.UpdateProviderCommand, provider.Provider]
}

func NewProviderUpdateEndpoint(handler ports.CommandHandler[command.UpdateProviderCommand, provider.Provider]) *ProviderUpdateEndpoint {
	return &ProviderUpdateEndpoint{handler: handler}
}

type updateProviderModel struct {
	Name           string     `json:"name" binding:"required"`
	Description    *string    `json:"description,omitempty"`
	IconUrl        *string    `json:"icon_url,omitempty"`
	Url            *string    `json:"url,omitempty"`
	PricingPageUrl *string    `json:"pricing_page_url,omitempty"`
	Labels         []string   `json:"labels" binding:"required"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty" format:"date-time"`
}

func (m updateProviderModel) Command(providerId uuid.UUID) (command.UpdateProviderCommand, error) {
	labels, err := collection.SelectErr(m.Labels, uuid.Parse)
	if err != nil {
		return command.UpdateProviderCommand{}, err
	}

	return command.UpdateProviderCommand{
		Id:             providerId,
		Name:           m.Name,
		Description:    m.Description,
		IconUrl:        m.IconUrl,
		Url:            m.Url,
		PricingPageUrl: m.PricingPageUrl,
		Labels:         labels,
		UpdatedAt:      nil,
	}, nil
}

// Handle godoc
//
//	@Summary		Update provider by ID
//	@Description	Update an existing provider's basic information
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			providerId	path		string				true	"Provider ID (UUID format)"
//	@Param			provider	body		updateProviderModel	true	"Updated provider data"
//	@Success		200			{object}	ProviderModel		"Successfully updated provider"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid input data or provider ID"
//	@Failure		404			{object}	HttpErrorResponse	"Provider not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/{providerId} [put]
func (e ProviderUpdateEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("providerId")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: "invalid id format",
		})
		return
	}

	var model updateProviderModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	cmd, err := model.Command(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		c.Abort()
		return
	}
	r := e.handler.Handle(c, cmd)
	FromResult(c,
		r,
		WithMapping[provider.Provider](func(prov provider.Provider) any {
			return newProviderModel(prov)
		}))
}

func (e ProviderUpdateEndpoint) Pattern() []string {
	return []string{
		":providerId",
	}
}

func (e ProviderUpdateEndpoint) Method() string {
	return http.MethodPut
}

func (e ProviderUpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
