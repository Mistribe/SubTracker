package subscription

import (
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
)

type ExportEndpoint struct {
	handler       ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[subscription.Subscription]]
	labelResolver export.LabelResolver
	exportService export.ExportService
}

func NewExportEndpoint(
	handler ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[subscription.Subscription]],
	labelResolver export.LabelResolver,
	exportService export.ExportService) *ExportEndpoint {
	return &ExportEndpoint{
		handler:       handler,
		labelResolver: labelResolver,
		exportService: exportService,
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

		// Extract all unique label IDs from subscriptions
		labelIDsMap := make(map[types.LabelID]bool)
		for _, sub := range subscriptions {
			for labelRef := range sub.Labels().It() {
				labelIDsMap[labelRef.LabelId] = true
			}
		}

		// Convert map to slice
		labelIDs := make([]types.LabelID, 0, len(labelIDsMap))
		for labelID := range labelIDsMap {
			labelIDs = append(labelIDs, labelID)
		}

		// Resolve label IDs to names
		labelIDToName := make(map[types.LabelID]string)
		if len(labelIDs) > 0 {
			// Resolve each label ID individually to maintain the mapping
			for _, labelID := range labelIDs {
				names, err := e.labelResolver.ResolveLabelNames(c, []types.LabelID{labelID})
				if err == nil && len(names) > 0 {
					labelIDToName[labelID] = names[0]
				}
			}
		}

		// Transform domain subscriptions to export models
		exportModels := make([]dto.SubscriptionExportModel, len(subscriptions))
		for i, sub := range subscriptions {
			exportModels[i] = transformSubscriptionToExportModel(sub, labelIDToName)
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
	labelIDToName map[types.LabelID]string) dto.SubscriptionExportModel {
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
		ProviderId:         sub.ProviderId().String(),
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
