package subscription

import (
	"net/http"

	"github.com/gin-gonic/gin"

	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/x/herd"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
)

type SummaryEndpoint struct {
	handler ports.QueryHandler[query.SummaryQuery, query.SummaryQueryResponse]
}

func NewSummaryEndpoint(handler ports.QueryHandler[query.SummaryQuery, query.SummaryQueryResponse]) *SummaryEndpoint {
	return &SummaryEndpoint{handler: handler}
}

// Handle godoc
//
//	@Summary		Get subscription summary
//	@Description	Returns summary information about subscriptions including total costs and upcoming renewals
//	@Tags			subscriptions
//	@Produce		json
//	@Param			top_providers		query		integer	true	"Number of top providers to return"
//	@Param			top_labels			query		integer	true	"Number of top labels to return"
//	@Param			upcoming_renewals	query		integer	true	"Number of upcoming renewals to return"
//	@Param			total_monthly		query		boolean	true	"Include monthly total costs"
//	@Param			total_yearly		query		boolean	true	"Include yearly total costs"
//	@Success		200					{object}	dto.SubscriptionSummaryResponse
//	@Failure		400					{object}	HttpErrorResponse
//	@Router			/subscriptions/summary [get]
func (e SummaryEndpoint) Handle(c *gin.Context) {
	var model dto.SubscriptionSummaryRequest
	if err := c.ShouldBindQuery(&model); err != nil {
		FromError(c, err)
		return
	}

	q := query.SummaryQuery{
		TopProviders:     model.TopProviders,
		TopLabels:        model.TopLabels,
		UpcomingRenewals: model.UpcomingRenewals,
		TotalMonthly:     model.TotalMonthly,
		TotalYearly:      model.TotalYearly,
	}

	r := e.handler.Handle(c, q)
	FromResult(c,
		r,
		WithMapping[query.SummaryQueryResponse](func(res query.SummaryQueryResponse) any {
			return dto.SubscriptionSummaryResponse{
				Active:            res.Active,
				ActivePersonal:    res.ActivePersonal,
				ActiveFamily:      res.ActiveFamily,
				TotalMonthly:      dto.NewAmount(res.TotalMonthly),
				TotalYearly:       dto.NewAmount(res.TotalYearly),
				TotalLastMonth:    dto.NewAmount(res.TotalLastMonth),
				TotalLastYear:     dto.NewAmount(res.TotalLastYear),
				PersonalMonthly:   dto.NewAmount(res.PersonalMonthly),
				PersonalYearly:    dto.NewAmount(res.PersonalYearly),
				PersonalLastMonth: dto.NewAmount(res.PersonalLastMonth),
				PersonalLastYear:  dto.NewAmount(res.PersonalLastYear),
				FamilyMonthly:     dto.NewAmount(res.FamilyMonthly),
				FamilyYearly:      dto.NewAmount(res.FamilyYearly),
				FamilyLastMonth:   dto.NewAmount(res.FamilyLastMonth),
				FamilyLastYear:    dto.NewAmount(res.FamilyLastYear),
				TopProviders: herd.Select(res.TopProviders,
					func(topProvider query.SummaryQueryTopProvidersResponse) dto.SubscriptionSummaryTopProviderResponse {
						return dto.SubscriptionSummaryTopProviderResponse{
							ProviderId: topProvider.ProviderID.String(),
							Total:      dto.NewAmount(topProvider.Total),
							Duration:   topProvider.Duration.String(),
						}
					}),
				TopLabels: herd.Select(res.TopLabels,
					func(topLabel query.SummaryQueryLabelResponse) dto.SubscriptionSummaryTopLabelResponse {
						return dto.SubscriptionSummaryTopLabelResponse{
							LabelId: topLabel.LabelID.String(),
							Total:   dto.NewAmount(topLabel.Total),
						}
					}),
				UpcomingRenewals: herd.Select(res.UpcomingRenewals,
					func(upcomingRenewal query.SummaryQueryUpcomingRenewalsResponse) dto.SubscriptionSummaryUpcomingRenewalResponse {
						m := dto.SubscriptionSummaryUpcomingRenewalResponse{
							ProviderId: upcomingRenewal.ProviderId.String(),
							At:         upcomingRenewal.At,
							Total:      dto.NewAmount(upcomingRenewal.Total),
						}
						if upcomingRenewal.Source != nil {
							s := dto.NewAmount(*upcomingRenewal.Source)
							m.Source = &s
						}

						return m
					}),
			}
		}))
}

func (e SummaryEndpoint) Pattern() []string {
	return []string{
		"/summary",
	}
}

func (e SummaryEndpoint) Method() string {
	return http.MethodGet
}

func (e SummaryEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
