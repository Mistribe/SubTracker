package provider

import (
	"net/http"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x/collection"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
	"github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/provider"
)

type UpdateEndpoint struct {
	handler ports.CommandHandler[command.UpdateProviderCommand, provider.Provider]
}

func NewUpdateEndpoint(handler ports.CommandHandler[command.UpdateProviderCommand, provider.Provider]) *UpdateEndpoint {
	return &UpdateEndpoint{handler: handler}
}

func updateProviderRequestToCommand(r dto.UpdateProviderRequest, providerId uuid.UUID) (
	command.UpdateProviderCommand,
	error) {
	labels, err := collection.SelectErr(r.Labels, uuid.Parse)
	if err != nil {
		return command.UpdateProviderCommand{}, err
	}

	return command.UpdateProviderCommand{
		ProviderID:     providerId,
		Name:           r.Name,
		Description:    r.Description,
		IconUrl:        r.IconUrl,
		Url:            r.Url,
		PricingPageUrl: r.PricingPageUrl,
		Labels:         labels,
		UpdatedAt:      nil,
	}, nil
}

// Handle godoc
//
//	@Summary		Update provider by LabelID
//	@Description	Update an existing provider's basic information
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			providerId	path		string						true	"Provider LabelID (UUID format)"
//	@Param			provider	body		dto.UpdateProviderRequest	true	"Updated provider data"
//	@Success		200			{object}	dto.ProviderModel			"Successfully updated provider"
//	@Failure		400			{object}	HttpErrorResponse			"Bad Request - Invalid input data or provider LabelID"
//	@Failure		404			{object}	HttpErrorResponse			"Provider not found"
//	@Failure		500			{object}	HttpErrorResponse			"Internal Server Error"
//	@Router			/providers/{providerId} [put]
func (e UpdateEndpoint) Handle(c *gin.Context) {
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

	var model dto.UpdateProviderRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}
	cmd, err := updateProviderRequestToCommand(model, id)
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
			return dto.NewProviderModel(prov)
		}))
}

func (e UpdateEndpoint) Pattern() []string {
	return []string{
		":providerId",
	}
}

func (e UpdateEndpoint) Method() string {
	return http.MethodPut
}

func (e UpdateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
