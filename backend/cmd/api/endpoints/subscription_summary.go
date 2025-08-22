package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/pkg/slicesx"
)

type SubscriptionSummaryEndpoint struct {
	handler core.QueryHandler[query.SummaryQuery, query.SummaryQueryResponse]
}

func NewSubscriptionSummaryEndpoint(handler core.QueryHandler[query.SummaryQuery, query.SummaryQueryResponse]) *SubscriptionSummaryEndpoint {
	return &SubscriptionSummaryEndpoint{handler: handler}
}

type SubscriptionSummaryRequest struct {
	TopProviders     uint8 `json:"top_providers,omitempty" form:"top_providers" example:"5" description:"Number of top providers to return"`
	UpcomingRenewals uint8 `json:"upcoming_renewals,omitempty" form:"upcoming_renewals" example:"3" description:"Number of upcoming renewals to return"`
	TopLabels        uint8 `json:"top_labels,omitempty" form:"top_labels"  example:"5" description:"Number of top labels to return"`
	TotalMonthly     bool  `json:"total_monthly,omitempty" form:"total_monthly"  example:"true" description:"Include monthly total costs"`
	TotalYearly      bool  `json:"total_yearly,omitempty" form:"total_yearly"  example:"true" description:"Include yearly total costs"`
}

type SubscriptionSummaryTopProviderResponse struct {
	ProviderId string      `json:"provider_id" binding:"required"`
	Total      AmountModel `json:"total"`
	Duration   string      `json:"duration"`
}

type SubscriptionSummaryTopLabelResponse struct {
	LabelId string      `json:"label_id" binding:"required"`
	Total   AmountModel `json:"total"`
}

type SubscriptionSummaryUpcomingRenewalResponse struct {
	ProviderId string       `json:"provider_id" binding:"required"`
	At         time.Time    `json:"at" binding:"required" format:"date-time"`
	Total      AmountModel  `json:"total"`
	Source     *AmountModel `json:"source"`
}

// @Description	Response containing subscription summary information
type SubscriptionSummaryResponse struct {
	Active           uint16                                       `json:"active" example:"10" description:"Number of active subscriptions"`
	TotalMonthly     AmountModel                                  `json:"total_monthly" description:"Total monthly subscription costs"`
	TotalYearly      AmountModel                                  `json:"total_yearly"  description:"Total yearly subscription costs"`
	TotalLastMonth   AmountModel                                  `json:"total_last_month" description:"Total monthly subscription costs for last month"`
	TotalLastYear    AmountModel                                  `json:"total_last_year" description:"Total yearly subscription costs for last year"`
	TopProviders     []SubscriptionSummaryTopProviderResponse     `json:"top_providers" description:"List of top providers by cost"`
	UpcomingRenewals []SubscriptionSummaryUpcomingRenewalResponse `json:"upcoming_renewals" description:"List of upcoming subscription renewals"`
	TopLabels        []SubscriptionSummaryTopLabelResponse        `json:"top_labels" description:"List of top labels by cost"`
}

// Handle godoc
//
//	@Summary		Get subscription summary
//	@Description	Returns summary information about subscriptions including total costs and upcoming renewals
//	@Tags			Subscriptions
//	@Produce		json
//	@Param			top_providers		query		integer	true	"Number of top providers to return"
//	@Param			top_labels			query		integer	true	"Number of top labels to return"
//	@Param			upcoming_renewals	query		integer	true	"Number of upcoming renewals to return"
//	@Param			total_monthly		query		boolean	true	"Include monthly total costs"
//	@Param			total_yearly		query		boolean	true	"Include yearly total costs"
//	@Success		200					{object}	SubscriptionSummaryResponse
//	@Failure		400					{object}	HttpErrorResponse
//	@Router			/subscriptions/summary [get]
func (e SubscriptionSummaryEndpoint) Handle(c *gin.Context) {
	var model SubscriptionSummaryRequest
	if err := c.ShouldBindQuery(&model); err != nil {
		c.JSON(http.StatusBadRequest, HttpErrorResponse{
			Message: err.Error(),
		})
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
	handleResponse(c,
		r,
		withMapping[query.SummaryQueryResponse](func(res query.SummaryQueryResponse) any {
			return SubscriptionSummaryResponse{
				Active:         res.Active,
				TotalMonthly:   newAmount(res.TotalMonthly),
				TotalYearly:    newAmount(res.TotalYearly),
				TotalLastMonth: newAmount(res.TotalLastMonth),
				TotalLastYear:  newAmount(res.TotalLastYear),
				TopProviders: slicesx.Select(res.TopProviders,
					func(topProvider query.SummaryQueryTopProvidersResponse) SubscriptionSummaryTopProviderResponse {
						return SubscriptionSummaryTopProviderResponse{
							ProviderId: topProvider.ProviderId.String(),
							Total:      newAmount(topProvider.Total),
							Duration:   topProvider.Duration.String(),
						}
					}),
				TopLabels: slicesx.Select(res.TopLabels,
					func(topLabel query.SummaryQueryLabelResponse) SubscriptionSummaryTopLabelResponse {
						return SubscriptionSummaryTopLabelResponse{
							LabelId: topLabel.TagId.String(),
							Total:   newAmount(topLabel.Total),
						}
					}),
				UpcomingRenewals: slicesx.Select(res.UpcomingRenewals,
					func(upcomingRenewal query.SummaryQueryUpcomingRenewalsResponse) SubscriptionSummaryUpcomingRenewalResponse {
						m := SubscriptionSummaryUpcomingRenewalResponse{
							ProviderId: upcomingRenewal.ProviderId.String(),
							At:         upcomingRenewal.At,
							Total:      newAmount(upcomingRenewal.Total),
						}
						if upcomingRenewal.Source != nil {
							s := newAmount(*upcomingRenewal.Source)
							m.Source = &s
						}

						return m
					}),
			}
		}))
}

func (e SubscriptionSummaryEndpoint) Pattern() []string {
	return []string{
		"/summary",
	}
}

func (e SubscriptionSummaryEndpoint) Method() string {
	return http.MethodGet
}

func (e SubscriptionSummaryEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
