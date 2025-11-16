package export

import (
	"errors"
	"strings"
)

type ExportFormat string

const (
	FormatCSV  ExportFormat = "csv"
	FormatJSON ExportFormat = "json"
	FormatYAML ExportFormat = "yaml"
)

var (
	ErrInvalidFormat = errors.New("invalid format. Supported formats: csv, json, yaml")
)

type ExportService interface {
	// ParseFormat validates and returns the export format
	ParseFormat(format string) (ExportFormat, error)

	// GetContentType returns the MIME type for the format
	GetContentType(format ExportFormat) string

	// GetFileExtension returns the file extension for the format
	GetFileExtension(format ExportFormat) string
}

type exportService struct{}

func NewExportService() ExportService {
	return &exportService{}
}

func (s *exportService) ParseFormat(format string) (ExportFormat, error) {
	normalized := strings.ToLower(strings.TrimSpace(format))

	switch normalized {
	case string(FormatCSV):
		return FormatCSV, nil
	case string(FormatJSON):
		return FormatJSON, nil
	case string(FormatYAML):
		return FormatYAML, nil
	default:
		return "", ErrInvalidFormat
	}
}

func (s *exportService) GetContentType(format ExportFormat) string {
	switch format {
	case FormatCSV:
		return "text/csv; charset=utf-8"
	case FormatJSON:
		return "application/json; charset=utf-8"
	case FormatYAML:
		return "application/x-yaml; charset=utf-8"
	default:
		return "application/octet-stream"
	}
}

func (s *exportService) GetFileExtension(format ExportFormat) string {
	switch format {
	case FormatCSV:
		return "csv"
	case FormatJSON:
		return "json"
	case FormatYAML:
		return "yaml"
	default:
		return "bin"
	}
}
