package provider

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x/collection"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/auth"
	"github.com/mistribe/subtracker/internal/usecase/provider/command"
	"github.com/mistribe/subtracker/pkg/ginx"
)

type CreateEndpoint struct {
	handler ports.CommandHandler[command.CreateProviderCommand, provider.Provider]
}

func NewCreateEndpoint(handler ports.CommandHandler[command.CreateProviderCommand, provider.Provider]) *CreateEndpoint {
	return &CreateEndpoint{handler: handler}
}

func createProviderRequestToCommand(r dto.CreateProviderRequest, userId string) (command.CreateProviderCommand, error) {
	var providerId *uuid.UUID
	if r.Id != nil {
		id, err := uuid.Parse(*r.Id)
		if err != nil {
			return command.CreateProviderCommand{}, err
		}

		providerId = &id
	}
	labels, err := collection.SelectErr(r.Labels, uuid.Parse)
	if err != nil {
		return command.CreateProviderCommand{}, err
	}

	owner, err := r.Owner.Owner(userId)
	if err != nil {
		return command.CreateProviderCommand{}, err
	}

	return command.CreateProviderCommand{
		Id: providerId,

		Name:           r.Name,
		Description:    r.Description,
		IconUrl:        r.IconUrl,
		Url:            r.Url,
		PricingPageUrl: r.PricingPageUrl,
		Labels:         labels,
		Owner:          owner,
		CreatedAt:      r.CreatedAt,
	}, nil
}

// Handle godoc
//
//	@Summary		Create a new provider
//	@Description	Create a new service provider with labels and owner information
//	@Tags			providers
//	@Accept			json
//	@Produce		json
//	@Param			provider	body		dto.CreateProviderRequest	true	"Provider creation data"
//	@Success		201			{object}	dto.ProviderModel			"Successfully created provider"
//	@Failure		400			{object}	HttpErrorResponse			"Bad Request - Invalid input data"
//	@Failure		401			{object}	HttpErrorResponse			"Unauthorized - Invalid user authentication"
//	@Failure		500			{object}	HttpErrorResponse			"Internal Server Error"
//	@Router			/providers [post]
func (e CreateEndpoint) Handle(c *gin.Context) {
	var model dto.CreateProviderRequest
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, ginx.HttpErrorResponse{
			Message: err.Error(),
		})
		return
	}

	userId, ok := authentication.GetUserIdFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, ginx.HttpErrorResponse{
			Message: "invalid user id",
		})
		return
	}

	cmd, err := createProviderRequestToCommand(model, userId)
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
			return dto.NewProviderModel(prov)
		}))
}

func (e CreateEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (e CreateEndpoint) Method() string {
	return http.MethodPost
}

func (e CreateEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
