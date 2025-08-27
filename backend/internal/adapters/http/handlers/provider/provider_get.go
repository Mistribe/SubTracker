package provider

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/query"
)

type ProviderGetEndpoint struct {
	handler ports.QueryHandler[query.FindOneQuery, provider.Provider]
}

func NewProviderGetEndpoint(handler ports.QueryHandler[query.FindOneQuery, provider.Provider]) *ProviderGetEndpoint {
	return &ProviderGetEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get provider by ID
//	@Description	Retrieve a single provider with all its plans and prices by ID
//	@Tags			providers
//	@Produce		json
//	@Param			providerId	path		string				true	"Provider ID (UUID format)"
//	@Success		200			{object}	ProviderModel		"Successfully retrieved provider"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid provider ID format"
//	@Failure		404			{object}	HttpErrorResponse	"Provider not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/{providerId} [get]
func (e ProviderGetEndpoint) Handle(c *gin.Context) {
	id, err := QueryParamAsUUID(c, "providerId")
	if err != nil {
		FromError(c, err)
		return
	}
	q := query.NewFindOneQuery(id)
	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[provider.Provider](func(prvdr provider.Provider) any {
			return newProviderModel(prvdr)
		}))
}

func (e ProviderGetEndpoint) Pattern() []string {
	return []string{
		":providerId",
	}
}

func (e ProviderGetEndpoint) Method() string {
	return http.MethodGet
}

func (e ProviderGetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
