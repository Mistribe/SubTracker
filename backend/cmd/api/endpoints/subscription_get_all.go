package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/subscription"
)

type SubscriptionGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[subscription.Subscription]]
}

func NewSubscriptionGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[subscription.Subscription]]) *SubscriptionGetAllEndpoint {
	return &SubscriptionGetAllEndpoint{handler: handler}
}

type subscriptionGetAllQueryParams struct {
	Search       string      `form:"search"`
	Recurrencies []string    `form:"recurrencies"`
	FromDate     *time.Time  `form:"from_date"`
	ToDate       *time.Time  `form:"to_date"`
	Users        []uuid.UUID `form:"users"`
	WithInactive bool        `form:"with_inactive"`
	Providers    []uuid.UUID `form:"providers"`
	Limit        int64       `form:"limit"`
	Offset       int64       `form:"offset"`
}

// Handle godoc
//
//	@Summary		Get all subscriptions
//	@Description	Retrieve a paginated list of all subscriptions for the authenticated user
//	@Tags			subscription
//	@Produce		json
//	@Param			search			query		string										false	"Search text"
//	@Param			recurrencies	query		[]string									false	"Filter by recurrency types"
//	@Param			from_date		query		string										false	"Filter by start date (RFC3339)"
//	@Param			to_date			query		string										false	"Filter by end date (RFC3339)"
//	@Param			users			query		[]string									false	"Filter by user IDs"
//	@Param			with_inactive	query		boolean										false	"Include inactive subscriptions"
//	@Param			providers		query		[]string									false	"Filter by provider IDs"
//	@Param			limit			query		integer										false	"Number of items per page (default: 10)"
//	@Param			offset			query		integer										false	"Page number (default: 0)"
//	@Success		200				{object}	PaginatedResponseModel[SubscriptionModel]	"Paginated list of subscriptions"
//	@Failure		400				{object}	HttpErrorResponse							"Bad Request - Invalid query parameters"
//	@Failure		500				{object}	HttpErrorResponse							"Internal Server Error"
//	@Router			/subscriptions [get]
func (s SubscriptionGetAllEndpoint) Handle(c *gin.Context) {
	var params subscriptionGetAllQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		// Ignore binding errors and apply defaults
	}

	// Apply defaults and basic validation
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Offset < 0 {
		params.Offset = 0
	}

	recurrencies := subscription.ParseRecurrencyTypesOrDefault(params.Recurrencies, subscription.UnknownRecurrency)

	q := query.NewFindAllQuery(params.Search,
		recurrencies,
		params.FromDate,
		params.ToDate,
		params.Users,
		params.Providers,
		params.WithInactive,
		params.Limit,
		params.Offset)
	r := s.handler.Handle(c, q)
	handleResponse(c,
		r,
		withMapping[core.PaginatedResponse[subscription.Subscription]](func(paginatedResult core.PaginatedResponse[subscription.Subscription]) any {
			return newPaginatedResponseModel(paginatedResult, newSubscriptionModel)
		}))
}

func (s SubscriptionGetAllEndpoint) Pattern() []string {
	return []string{
		"",
	}
}

func (s SubscriptionGetAllEndpoint) Method() string {
	return http.MethodGet
}

func (s SubscriptionGetAllEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}
