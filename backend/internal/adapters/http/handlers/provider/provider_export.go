package provider

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/adapters/http/export"
	"github.com/mistribe/subtracker/internal/domain/provider"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/provider/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type ExportEndpoint struct {
	handler       ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[provider.Provider]]
	labelResolver export.LabelResolver
	exportService export.ExportService
}

func NewExportEndpoint(
	handler ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[provider.Provider]],
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
//	@Summary		Export providers
//	@Description	Export all providers in CSV, JSON, or YAML format
//	@Tags			providers
//	@Produce		text/csv
//	@Produce		application/json
//	@Produce		application/x-yaml
//	@Param			format	query		string				false	"Export format (csv, json, yaml)"	default(json)
//	@Success		200		{file}		file				"Exported providers file"
//	@Failure		400		{object}	HttpErrorResponse	"Invalid format parameter"
//	@Failure		401		{object}	HttpErrorResponse	"Unauthorized"
//	@Failure		500		{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/providers/export [get]
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

	// Fetch all providers using high limit
	q := query.NewFindAllQuery("", 10000, 0)
	r := e.handler.Handle(c, q)

	// Handle query errors
	if r.IsFaulted() {
		c.JSON(http.StatusInternalServerError, NewProblem(
			ProblemTypeAboutBlank,
			"Internal Server Error",
			http.StatusInternalServerError,
			"Failed to fetch providers",
			c.FullPath(),
		))
		return
	}

	var paginatedResult shared.PaginatedResponse[provider.Provider]
	r.IfSuccess(func(value shared.PaginatedResponse[provider.Provider]) {
		paginatedResult = value
	})
	providers := paginatedResult.Data()

	// Extract all unique label IDs from providers
	labelIDsMap := make(map[types.LabelID]bool)
	for _, prov := range providers {
		for labelID := range prov.Labels().It() {
			labelIDsMap[labelID] = true
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
		// We need to resolve each label ID individually to maintain the mapping
		for _, labelID := range labelIDs {
			names, err := e.labelResolver.ResolveLabelNames(c, []types.LabelID{labelID})
			if err == nil && len(names) > 0 {
				labelIDToName[labelID] = names[0]
			}
		}
	}

	// Transform domain providers to export models
	exportModels := make([]dto.ProviderExportModel, len(providers))
	for i, prov := range providers {
		exportModels[i] = transformProviderToExportModel(prov, labelIDToName)
	}

	// Set response headers
	c.Header("Content-Type", e.exportService.GetContentType(format))
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("providers_%s.%s", timestamp, e.exportService.GetFileExtension(format))
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
		return
	}
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

// transformProviderToExportModel converts a domain provider to an export model
func transformProviderToExportModel(prov provider.Provider,
	labelIDToName map[types.LabelID]string) dto.ProviderExportModel {
	// Resolve label IDs to names
	labelNames := make([]string, 0)
	for labelID := range prov.Labels().It() {
		if name, ok := labelIDToName[labelID]; ok {
			labelNames = append(labelNames, name)
		}
	}

	return dto.ProviderExportModel{
		Id:             prov.Id().String(),
		Name:           prov.Name(),
		Description:    prov.Description(),
		Url:            prov.Url(),
		IconUrl:        prov.IconUrl(),
		PricingPageUrl: prov.PricingPageUrl(),
		Labels:         labelNames,
	}
}
