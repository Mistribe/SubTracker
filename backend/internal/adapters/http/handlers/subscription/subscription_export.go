package subscription

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/adapters/http/export"
	"github.com/mistribe/subtracker/internal/domain/subscription"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/subscription/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
	"github.com/mistribe/subtracker/pkg/langext/result"
	"github.com/mistribe/subtracker/pkg/x"
	"github.com/mistribe/subtracker/pkg/x/herd"
)

type ExportEndpoint struct {
	handler          ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[subscription.Subscription]]
	labelResolver    export.LabelResolver
	providerResolver export.ProviderResolver
	exportService    export.ExportService
}

func NewExportEndpoint(
	handler ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[subscription.Subscription]],
	labelResolver export.LabelResolver,
	providerResolver export.ProviderResolver,
	exportService export.ExportService) *ExportEndpoint {
	return &ExportEndpoint{
		handler:          handler,
		labelResolver:    labelResolver,
		exportService:    exportService,
		providerResolver: providerResolver,
	}
}

// Handle godoc
//
//	@Summary		Export subscriptions
//	@Description	Export all subscriptions in CSV, JSON, or YAML format
//	@Tags			subscriptions
//	@Produce		text/csv
//	@Produce		application/json
//	@Produce		application/x-yaml
//	@Param			format	query		string				false	"Export format (csv, json, yaml)"	default(json)
//	@Success		200		{file}		file				"Exported subscriptions file"
//	@Failure		400		{object}	HttpErrorResponse	"Invalid format parameter"
//	@Failure		401		{object}	HttpErrorResponse	"Unauthorized"
//	@Failure		500		{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/subscriptions/export [get]
func (e ExportEndpoint) Handle(c *gin.Context) {
	// Parse format parameter (default to json)
	formatParam := c.DefaultQuery("format", "json")
	format, err := e.exportService.ParseFormat(formatParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, NewProblem(
			ProblemTypeAboutBlank,
			"Bad Request",
			http.StatusBadRequest,
			export.ErrInvalidFormat.Error(),
			c.FullPath(),
		))
		return
	}

	// Fetch all subscriptions using high limit and withInactive=true
	// todo do better
	q := query.NewFindAllQuery("", nil, nil, nil, nil, nil, true, 10000, 0)
	r := e.handler.Handle(c, q)

	// Handle query errors
	if r.IsFaulted() {
		c.JSON(http.StatusInternalServerError, NewProblem(
			ProblemTypeAboutBlank,
			"Internal Server Error",
			http.StatusInternalServerError,
			"Failed to fetch subscriptions",
			c.FullPath(),
		))
		return
	}

	result.Match(r, func(paginatedResult shared.PaginatedResponse[subscription.Subscription]) any {
		subscriptions := paginatedResult.Data()

		labelIDToName, err := e.extractLabelNames(c, subscriptions)
		if err != nil {
			FromError(c, err)
			return nil
		}
		providerIDToKey, err := e.extractProviderKeys(c, subscriptions)
		if err != nil {
			FromError(c, err)
			return nil
		}

		// Transform domain subscriptions to export models
		exportModels := make([]dto.SubscriptionExportModel, len(subscriptions))
		for i, sub := range subscriptions {
			exportModels[i] = transformSubscriptionToExportModel(sub, labelIDToName, providerIDToKey)
		}

		// Set response headers
		c.Header("Content-Type", e.exportService.GetContentType(format))
		timestamp := time.Now().Format("2006-01-02T15-04-05")
		filename := fmt.Sprintf("subscriptions_%s.%s", timestamp, e.exportService.GetFileExtension(format))
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

		// Encode and stream data based on format
		var encodeErr error
		switch format {
		case export.FormatCSV:
			encodeErr = export.EncodeCSV(c.Writer, exportModels)
		case export.FormatJSON:
			encodeErr = export.EncodeJSON(c.Writer, exportModels)
		case export.FormatYAML:
			encodeErr = export.EncodeYAML(c.Writer, exportModels)
		}

		if encodeErr != nil {
			// If encoding fails after headers are sent, we can't send a proper error response
			// Log the error (in production, use proper logging)
			c.Error(encodeErr)
			return nil
		}
		return nil
	}, func(err error) any {
		c.JSON(http.StatusInternalServerError, NewProblem(
			ProblemTypeAboutBlank,
			"Internal Server Error",
			http.StatusInternalServerError,
			"Failed to fetch subscriptions",
			c.FullPath(),
		))
		return nil
	})
}

func (e ExportEndpoint) extractLabelNames(ctx context.Context,
	subscriptions []subscription.Subscription) (map[types.LabelID]string, error) {
	// Extract all unique label IDs from subscriptions
	labelIDs := herd.NewSet[types.LabelID]()
	for _, sub := range subscriptions {
		for labelRef := range sub.Labels().It() {
			labelIDs.Add(labelRef.LabelId)
		}
	}

	if len(labelIDs) < 0 {
		return nil, nil
	}
	names, err := e.labelResolver.ResolveLabelNames(ctx, labelIDs.ToSlice())
	if err != nil {
		return nil, err
	}

	return names, nil
}

func (e ExportEndpoint) extractProviderKeys(ctx context.Context,
	subscriptions []subscription.Subscription) (map[types.ProviderID]string, error) {
	providerIDs := herd.NewSet[types.ProviderID]()
	for _, sub := range subscriptions {
		providerIDs.Add(sub.ProviderId())
	}

	if len(providerIDs) < 0 {
		return nil, nil
	}
	keys, err := e.providerResolver.ResolveProviderKeys(ctx, providerIDs.ToSlice())
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (e ExportEndpoint) Pattern() []string {
	return []string{"/export"}
}

func (e ExportEndpoint) Method() string {
	return http.MethodGet
}

func (e ExportEndpoint) Middlewares() []gin.HandlerFunc {
	return nil
}

// transformSubscriptionToExportModel converts a domain subscription to an export model
func transformSubscriptionToExportModel(sub subscription.Subscription,
	labelIDToName map[types.LabelID]string,
	providerIDToKey map[types.ProviderID]string) dto.SubscriptionExportModel {
	// Resolve label IDs to names
	labelNames := make([]string, 0)
	for labelRef := range sub.Labels().It() {
		if name, ok := labelIDToName[labelRef.LabelId]; ok {
			labelNames = append(labelNames, name)
		}
	}

	// Format dates as ISO 8601 (YYYY-MM-DD)
	startDate := sub.StartDate().Format("2006-01-02")
	var endDate *string
	if sub.EndDate() != nil {
		formatted := sub.EndDate().Format("2006-01-02")
		endDate = &formatted
	}

	// Map owner type
	ownerType := "personal"
	if sub.Owner().Type() == types.FamilyOwnerType {
		ownerType = "family"
	}

	// Extract custom price amount and currency
	customPriceAmount := 0.0
	customPriceCurrency := ""
	if sub.Price() != nil && sub.Price().Amount().IsValid() {
		customPriceAmount = sub.Price().Amount().Value()
		customPriceCurrency = sub.Price().Amount().Currency().String()
	}

	// Handle free trial dates
	var freeTrialStartDate *string
	var freeTrialEndDate *string
	if sub.FreeTrial() != nil {
		startFormatted := sub.FreeTrial().StartDate().Format("2006-01-02")
		freeTrialStartDate = &startFormatted
		endFormatted := sub.FreeTrial().EndDate().Format("2006-01-02")
		freeTrialEndDate = &endFormatted
	}

	var payer *string
	var payerMemberId *string
	if sub.Payer() != nil {
		payer = x.P(sub.Payer().Type().String())
		if sub.Payer().Type() == subscription.FamilyMemberPayer {
			payerMemberId = x.P(sub.Payer().MemberId().String())
		}
	}

	var familyUsers []string
	for user := range sub.FamilyUsers().It() {
		familyUsers = append(familyUsers, user.String())
	}

	return dto.SubscriptionExportModel{
		Id:                 sub.Id().String(),
		ProviderKey:        providerIDToKey[sub.ProviderId()],
		FriendlyName:       sub.FriendlyName(),
		StartDate:          startDate,
		EndDate:            endDate,
		Recurrency:         sub.Recurrency().String(),
		CustomRecurrency:   sub.CustomRecurrency(),
		Amount:             customPriceAmount,
		Currency:           customPriceCurrency,
		OwnerType:          ownerType,
		FreeTrialStartDate: freeTrialStartDate,
		FreeTrialEndDate:   freeTrialEndDate,
		Payer:              payer,
		PayerMemberId:      payerMemberId,
		FamilyUsers:        familyUsers,
		Labels:             labelNames,
	}
}
