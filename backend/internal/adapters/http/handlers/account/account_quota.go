package account

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x/herd"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/domain/billing"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/account/query"
)

type GetQuotaUsageEndpoint struct {
	handler ports.QueryHandler[query.GetQuotaUsage, []billing.EffectiveEntitlement]
}

// Handle godoc
//
//	@Summary		Get quota usage
//	@Description	Retrieve the current quota usage and limits for the authenticated user
//	@Tags			accounts
//	@Produce		json
//	@Success		200	{object}	[]dto.QuotaUsageModel	"Successfully retrieved quota usage"
//	@Failure		401	{object}	HttpErrorResponse		"Unauthorized - Invalid user authentication"
//	@Failure		500	{object}	HttpErrorResponse		"Internal Server Error"
//	@Router			/accounts/quota/usage [get]
func (e GetQuotaUsageEndpoint) Handle(c *gin.Context) {
	q := query.GetQuotaUsage{}

	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[[]billing.EffectiveEntitlement](func(r []billing.EffectiveEntitlement) any {
			return herd.Select(r, func(in billing.EffectiveEntitlement) dto.QuotaUsageModel {
				return dto.NewQuotaUsageModel(in)
			})
		}))
}

func (e GetQuotaUsageEndpoint) Pattern() []string {
	return []string{
		"quota/usage",
	}
}

func (e GetQuotaUsageEndpoint) Method() string {
	return http.MethodGet
}

func (e GetQuotaUsageEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

func NewGetQuotaUsageEndpoint(handler ports.QueryHandler[query.GetQuotaUsage, []billing.EffectiveEntitlement]) *GetQuotaUsageEndpoint {
	return &GetQuotaUsageEndpoint{handler: handler}
}
