package provider

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	. "github.com/mistribe/subtracker/pkg/ginx"

	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/provider/query"
)

type GetEndpoint struct {
	handler ports.QueryHandler[query.FindOneQuery, provider.Provider]
}

func NewGetEndpoint(handler ports.QueryHandler[query.FindOneQuery, provider.Provider]) *GetEndpoint {
	return &GetEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get provider by LabelID
//	@Description	Retrieve a single provider with all its plans and prices by LabelID
//	@Tags			providers
//	@Produce		json
//	@Param			providerId	path		string				true	"Provider LabelID (UUID format)"
//	@Success		200			{object}	dto.ProviderModel	"Successfully retrieved provider"
//	@Failure		400			{object}	HttpErrorResponse	"Bad Request - Invalid provider LabelID format"
//	@Failure		404			{object}	HttpErrorResponse	"Provider not found"
//	@Failure		500			{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/{providerId} [get]
func (e GetEndpoint) Handle(c *gin.Context) {
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
			return dto.NewProviderModel(prvdr)
		}))
}

func (e GetEndpoint) Pattern() []string {
	return []string{
		":providerId",
	}
}

func (e GetEndpoint) Method() string {
	return http.MethodGet
}

func (e GetEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
