package label

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/mistribe/subtracker/internal/adapters/http/dto"
	"github.com/mistribe/subtracker/internal/adapters/http/export"
	"github.com/mistribe/subtracker/internal/domain/label"
	"github.com/mistribe/subtracker/internal/domain/types"
	"github.com/mistribe/subtracker/internal/ports"
	"github.com/mistribe/subtracker/internal/shared"
	"github.com/mistribe/subtracker/internal/usecase/label/query"
	. "github.com/mistribe/subtracker/pkg/ginx"
)

type ExportEndpoint struct {
	handler       ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[label.Label]]
	exportService export.ExportService
}

func NewExportEndpoint(
	handler ports.QueryHandler[query.FindAllQuery, shared.PaginatedResponse[label.Label]],
	exportService export.ExportService) *ExportEndpoint {
	return &ExportEndpoint{
		handler:       handler,
		exportService: exportService,
	}
}

// Handle godoc
//
//	@Summary		Export labels
//	@Description	Export all labels in CSV, JSON, or YAML format
//	@Tags			labels
//	@Produce		text/csv
//	@Produce		application/json
//	@Produce		application/x-yaml
//	@Param			format	query		string				false	"Export format (csv, json, yaml)"	default(json)
//	@Success		200		{file}		file				"Exported labels file"
//	@Failure		400		{object}	HttpErrorResponse	"Invalid format parameter"
//	@Failure		401		{object}	HttpErrorResponse	"Unauthorized"
//	@Failure		500		{object}	HttpErrorResponse	"Internal Server Error"
//	@Router			/labels/export [get]
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

	// Fetch all labels using high limit
	q := query.NewFindAllQuery("", 10000, 0)
	r := e.handler.Handle(c, q)

	// Handle query errors
	if r.IsFaulted() {
		c.JSON(http.StatusInternalServerError, NewProblem(
			ProblemTypeAboutBlank,
			"Internal Server Error",
			http.StatusInternalServerError,
			"Failed to fetch labels",
			c.FullPath(),
		))
		return
	}

	var paginatedResult shared.PaginatedResponse[label.Label]
	r.IfSuccess(func(value shared.PaginatedResponse[label.Label]) {
		paginatedResult = value
	})
	labels := paginatedResult.Data()

	// Transform domain labels to export models
	exportModels := make([]dto.LabelExportModel, len(labels))
	for i, lbl := range labels {
		exportModels[i] = transformLabelToExportModel(lbl)
	}

	// Set response headers
	c.Header("Content-Type", e.exportService.GetContentType(format))
	timestamp := time.Now().Format("2006-01-02T15-04-05")
	filename := fmt.Sprintf("labels_%s.%s", timestamp, e.exportService.GetFileExtension(format))
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

// transformLabelToExportModel converts a domain label to an export model
func transformLabelToExportModel(lbl label.Label) dto.LabelExportModel {
	var ownerType string
	var ownerFamilyId *string

	switch lbl.Owner().Type() {
	case types.PersonalOwnerType:
		ownerType = "personal"
	case types.FamilyOwnerType:
		ownerType = "family"
		familyId := lbl.Owner().FamilyId().String()
		ownerFamilyId = &familyId
	default:
		ownerType = "unknown"
	}

	return dto.LabelExportModel{
		Id:            lbl.Id().String(),
		Name:          lbl.Name(),
		Color:         lbl.Color(),
		OwnerType:     ownerType,
		OwnerFamilyId: ownerFamilyId,
	}
}
