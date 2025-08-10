package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/provider/query"
	"github.com/oleexo/subtracker/internal/domain/provider"
)

type ProviderGetEndpoint struct {
	handler core.QueryHandler[query.FindOneQuery, provider.Provider]
}

func NewProviderGetEndpoint(handler core.QueryHandler[query.FindOneQuery, provider.Provider]) *ProviderGetEndpoint {
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
	id, err := paramAsUuid(c, "providerId")
	if err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
	}
	q := query.NewFindOneQuery(id)
	r := e.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[provider.Provider](func(prvdr provider.Provider) any {
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
