package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/command"
	"github.com/oleexo/subtracker/internal/domain/label"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderUpdateEndpoint struct {
	handler core.CommandHandler[command.UpdateProviderCommand, provider.Provider]
}

func NewProviderUpdateEndpoint(handler core.CommandHandler[command.UpdateProviderCommand, provider.Provider]) *ProviderUpdateEndpoint {
	return &ProviderUpdateEndpoint{handler: handler}
}

type updateProviderModel struct {
	Name           string   `json:"name" binding:"required"`
	Description    *string  `json:"description,omitempty"`
	IconUrl        *string  `json:"icon_url,omitempty"`
	Url            *string  `json:"url,omitempty"`
	PricingPageUrl *string  `json:"pricing_page_url,omitempty"`
	Labels         []string `json:"labels" binding:"required"`
}

func (e ProviderUpdateEndpoint) Handle(c *gin.Context) {
	idParam := c.Param("providerId")
	if idParam == "" {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "id parameter is required",
		})
		return
	}

	id, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: "invalid id format",
		})
		return
	}

	var model updateLabelModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, httpError{
			Message: err.Error(),
		})
		return
	}
	cmd, err := model.Command(id)
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
		withMapping[label.Label](func(lbl label.Label) any {
			return newLabelModel(lbl)
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
