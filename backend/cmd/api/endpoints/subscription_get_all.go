package endpoints

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/oleexo/subtracker/internal/application/core"
	"github.com/oleexo/subtracker/internal/application/subscription/query"
	"github.com/oleexo/subtracker/internal/domain/subscription"
	"github.com/oleexo/subtracker/pkg/types"
)

type SubscriptionGetAllEndpoint struct {
	handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[subscription.Subscription]]
}

func NewSubscriptionGetAllEndpoint(handler core.QueryHandler[query.FindAllQuery, core.PaginatedResponse[subscription.Subscription]]) *SubscriptionGetAllEndpoint {
	return &SubscriptionGetAllEndpoint{handler: handler}
}

type subscriptionGetAllQueryParams struct {
	Search       string      `form:"search"`
	SortBy       string      `form:"sort_by"`
	SortOrder    string      `form:"sort_order"`
	Recurrencies []string    `form:"recurrencies"`
	FromDate     *time.Time  `form:"fromDate"`
	ToDate       *time.Time  `form:"toDate"`
	Users        []uuid.UUID `form:"users"`
	WithInactive bool        `form:"with_inactive"`
	Providers    []uuid.UUID `form:"providers"`
	Limit        int32       `form:"limit"`
	Offset       int32       `form:"offset"`
}

// Handle godoc
//
//	@Summary		Get all subscriptions
//	@Description	Retrieve a paginated list of all subscriptions for the authenticated user
//	@Tags			subscription
//	@Produce		json
//	@Param			search			query		string										false	"Search text"
//	@Param			sortBy			query		string										false	"Sort by field"
//	@Param			sortOrder		query		string										false	"Sort order (asc, desc)"
//	@Param			recurrencies	query		[]string									false	"Filter by recurrency types"
//	@Param			fromDate		query		string										false	"Filter by start date (RFC3339)"
//	@Param			toDate			query		string										false	"Filter by end date (RFC3339)"
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
	if params.SortBy == "" {
		params.SortBy = "name"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	sortBy := subscription.ParseSortableFieldOrDefault(params.SortBy, subscription.FriendlyNameSortableField)
	sortOrder := types.ParseSortOrderOrDefault(params.SortOrder, types.SortOrderAsc)
	recurrencies := subscription.ParseRecurrencyTypesOrDefault(params.Recurrencies, subscription.UnknownRecurrency)

	q := query.NewFindAllQuery(params.Search,
		recurrencies,
		params.FromDate,
		params.ToDate,
		params.Users,
		params.Providers,
		params.WithInactive,
		sortBy,
		sortOrder,
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
