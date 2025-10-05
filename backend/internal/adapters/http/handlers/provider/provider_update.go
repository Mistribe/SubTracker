package provider

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/types"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x/herd"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
)

type UpdateEndpoint struct {
	handler ports.CommandHandler[command.UpdateProviderCommand, provider.Provider]
}

func NewUpdateEndpoint(handler ports.CommandHandler[command.UpdateProviderCommand, provider.Provider]) *UpdateEndpoint {
	return &UpdateEndpoint{handler: handler}
}

func updateProviderRequestToCommand(r dto.UpdateProviderRequest,
	providerId types.ProviderID) (command.UpdateProviderCommand, error) {
	labels, err := herd.SelectErr(r.Labels, types.ParseLabelID)
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
	id, err := types.ParseProviderID(c.Param("providerId"))
	if err != nil {
		FromError(c, err)
		return
	}

	var model dto.UpdateProviderRequest
	if err = c.ShouldBindJSON(&model); err != nil {
		FromError(c, err)
		return
	}
	cmd, err := updateProviderRequestToCommand(model, id)
	if err != nil {
		FromError(c, err)
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
