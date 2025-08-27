package provider

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
	"github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/slicesx"
)

type ProviderCreateEndpoint struct {
	handler ports.CommandHandler[command.CreateProviderCommand, provider.Provider]
}

func NewProviderCreateEndpoint(handler ports.CommandHandler[command.CreateProviderCommand, provider.Provider]) *ProviderCreateEndpoint {
	return &ProviderCreateEndpoint{handler: handler}
}

type createProviderModel struct {
	Id             *string                `json:"id,omitempty"`
	Name           string                 `json:"name" binding:"required"`
	Description    *string                `json:"description,omitempty"`
	IconUrl        *string                `json:"icon_url,omitempty"`
	Url            *string                `json:"url,omitempty"`
	PricingPageUrl *string                `json:"pricing_page_url,omitempty"`
	Labels         []string               `json:"labels" binding:"required"`
	Owner          dto.EditableOwnerModel `json:"owner" binding:"required"`
	CreatedAt      *time.Time             `json:"created_at,omitempty" format:"date-time"`
}

func (m createProviderModel) Command(userId string) (command.CreateProviderCommand, error) {
	var providerId *uuid.UUID
	if m.Id != nil {
		id, err := uuid.Parse(*m.Id)
		if err != nil {
			return command.CreateProviderCommand{}, err
		}

		providerId = &id
	}
	labels, err := slicesx.SelectErr(m.Labels, uuid.Parse)
	if err != nil {
		return command.CreateProviderCommand{}, err
	}

	owner, err := m.Owner.Owner(userId)
	if err != nil {
		return command.CreateProviderCommand{}, err
	}

	return command.CreateProviderCommand{
		Id: providerId,

		Name:           m.Name,
		Description:    m.Description,
		IconUrl:        m.IconUrl,
		Url:            m.Url,
		PricingPageUrl: m.PricingPageUrl,
		Labels:         labels,
		Owner:          owner,
		CreatedAt:      m.CreatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Create a new provider
//	@Description	Create a new service provider with labels and owner information
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			provider	body		createProviderModel	true	"Provider creation data"
//	@Success		201			{object}	ProviderModel		"Successfully created provider"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid input data"
//	@Failure		401			{object}	HttpErrorResponse	"Unauthorized - Invalid user authentication"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers [post]
func (e ProviderCreateEndpoint) Handle(c *gin.Context) {
	var model createProviderModel
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	userId, ok := auth.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ginx.HttpErrorResponse{
			Message: "invalid user id",
		})
		return
	}

	cmd, err := model.Command(userId)
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
		WithStatus[provider.Provider](http.StatusCreated),
		WithMapping[provider.Provider](func(prov provider.Provider) any {
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
