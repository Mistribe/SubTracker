package provider

import (
	"context"
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
	"github.com/mistribe/subtracker/pkg/x/herd"
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

	// Resolve label IDs to names
	labelIDToName, err := e.extractLabelNames(c, providers)
	if err != nil {
		c.Error(fmt.Errorf("failed to resolve label names: %w", err))
		return
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

func (e ExportEndpoint) extractLabelNames(ctx context.Context,
	providers []provider.Provider) (map[types.LabelID]string, error) {
	// Extract all unique label IDs from providers
	labelIDs := herd.NewSet[types.LabelID]()
	for _, sub := range providers {
		for labelID := range sub.Labels().It() {
			labelIDs.Add(labelID)
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
		Name:           prov.Name(),
		Key:            prov.Key(),
		Description:    prov.Description(),
		Url:            prov.Url(),
		IconUrl:        prov.IconUrl(),
		PricingPageUrl: prov.PricingPageUrl(),
		Labels:         labelNames,
	}
}
