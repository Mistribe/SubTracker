# Design Document

## Overview

This design implements data export functionality for Labels, Providers, and Subscriptions, allowing users to download their data in CSV, JSON, or YAML formats. The solution follows the existing application architecture using the CQRS pattern with query handlers, HTTP handlers, and streaming responses. The export endpoints will reuse existing query logic to fetch data and transform it into the requested format, streaming the output directly to the client without intermediate file storage.

## Architecture

### High-Level Flow

```
Client Request → HTTP Handler → Query Handler → Repository → Database
                      ↓
                Format Encoder (CSV/JSON/YAML)
                      ↓
                HTTP Response Stream
```

### Component Layers

1. **HTTP Handler Layer**: New export endpoints for each resource type
2. **Use Case Layer**: Reuse existing FindAll query handlers
3. **Encoding Layer**: Format-specific encoders for CSV, JSON, and YAML
4. **Response Streaming**: Direct HTTP streaming without file storage

## Components and Interfaces

### 1. Export Endpoints

Create three new endpoint handlers following the existing pattern:

- `internal/adapters/http/handlers/label/label_export.go`
- `internal/adapters/http/handlers/provider/provider_export.go`
- `internal/adapters/http/handlers/subscription/subscription_export.go`

Each endpoint will:
- Implement the `ginfx.Endpoint` interface
- Parse and validate the `format` query parameter
- Reuse existing FindAll query handlers to fetch data
- Transform domain entities to simplified export models
- Resolve label IDs to label names (for providers and subscriptions)
- Set appropriate HTTP headers (Content-Type, Content-Disposition)
- Stream the encoded data to the response

### 2. Export Service

Create a shared export service to handle format encoding:

**Location**: `internal/adapters/http/export/service.go`

```go
type ExportFormat string

const (
    FormatCSV  ExportFormat = "csv"
    FormatJSON ExportFormat = "json"
    FormatYAML ExportFormat = "yaml"
)

type ExportService interface {
    // EncodeCSV writes data as CSV to the writer
    EncodeCSV(w io.Writer, data interface{}, headers []string) error
    
    // EncodeJSON writes data as JSON to the writer
    EncodeJSON(w io.Writer, data interface{}) error
    
    // EncodeYAML writes data as YAML to the writer
    EncodeYAML(w io.Writer, data interface{}) error
    
    // ParseFormat validates and returns the export format
    ParseFormat(format string) (ExportFormat, error)
    
    // GetContentType returns the MIME type for the format
    GetContentType(format ExportFormat) string
    
    // GetFileExtension returns the file extension for the format
    GetFileExtension(format ExportFormat) string
}
```

### 3. CSV Encoder

**Location**: `internal/adapters/http/export/csv_encoder.go`

The CSV encoder will:
- Use Go's `encoding/csv` package
- Accept a slice of structs and convert to CSV rows
- Use `csv` struct tags to determine column headers and order
- Handle optional fields (pointers) by writing empty strings when nil
- Handle string slices by joining with ", " (comma-space separator)
- Write headers as the first row

### 4. JSON Encoder

**Location**: `internal/adapters/http/export/json_encoder.go`

The JSON encoder will:
- Use Go's `encoding/json` package
- Leverage existing DTO models with JSON tags
- Stream JSON array directly to response

### 5. YAML Encoder

**Location**: `internal/adapters/http/export/yaml_encoder.go`

The YAML encoder will:
- Use `gopkg.in/yaml.v3` package
- Convert data to YAML format
- Stream YAML document to response

## Data Models

### Export DTOs

The export will use simplified DTO models specifically designed for export, containing only user-relevant fields:

**Label Export Model**:
```go
type LabelExportModel struct {
    Id             string  `json:"id" csv:"id"`
    Name           string  `json:"name" csv:"name"`
    Color          string  `json:"color" csv:"color"`
    OwnerType      string  `json:"ownerType" csv:"ownerType"`
    OwnerFamilyId  *string `json:"ownerFamilyId,omitempty" csv:"ownerFamilyId"`
}
```

**Provider Export Model**:
```go
type ProviderExportModel struct {
    Id             string   `json:"id" csv:"id"`
    Name           string   `json:"name" csv:"name"`
    Description    *string  `json:"description,omitempty" csv:"description"`
    Url            *string  `json:"url,omitempty" csv:"url"`
    IconUrl        *string  `json:"iconUrl,omitempty" csv:"iconUrl"`
    PricingPageUrl *string  `json:"pricingPageUrl,omitempty" csv:"pricingPageUrl"`
    Labels         []string `json:"labels" csv:"labels"` // Label names, not IDs
}
```

**Subscription Export Model**:
```go
type SubscriptionExportModel struct {
    Id                   string   `json:"id" csv:"id"`
    ProviderId           string   `json:"providerId" csv:"providerId"`
    FriendlyName         *string  `json:"friendlyName,omitempty" csv:"friendlyName"`
    StartDate            string   `json:"startDate" csv:"startDate"` // ISO 8601 date
    EndDate              *string  `json:"endDate,omitempty" csv:"endDate"` // ISO 8601 date
    Recurrency           string   `json:"recurrency" csv:"recurrency"`
    CustomRecurrency     *int32   `json:"customRecurrency,omitempty" csv:"customRecurrency"`
    CustomPriceAmount    float64  `json:"customPriceAmount" csv:"customPriceAmount"`
    CustomPriceCurrency  string   `json:"customPriceCurrency" csv:"customPriceCurrency"`
    OwnerType            string   `json:"ownerType" csv:"ownerType"`
    FreeTrialStartDate   *string  `json:"freeTrialStartDate,omitempty" csv:"freeTrialStartDate"`
    FreeTrialEndDate     *string  `json:"freeTrialEndDate,omitempty" csv:"freeTrialEndDate"`
    Labels               []string `json:"labels" csv:"labels"` // Label names, not IDs
}
```

### Data Transformation

The export endpoints will need to:
1. Fetch domain entities using existing query handlers
2. Resolve label IDs to label names for providers and subscriptions
3. Transform domain entities to simplified export models
4. Handle owner type mapping (user → "personal", family → "family")
5. Format dates as ISO 8601 strings (YYYY-MM-DD)

### Label Resolution Service

**Location**: `internal/adapters/http/export/label_resolver.go`

Since providers and subscriptions reference labels by ID but exports need label names, we need a label resolution service:

```go
type LabelResolver interface {
    // ResolveLabelNames takes a slice of label IDs and returns their names
    ResolveLabelNames(ctx context.Context, labelIds []types.LabelID) ([]string, error)
}
```

This service will:
- Accept label IDs from providers/subscriptions
- Query the label repository to fetch label entities
- Extract and return label names
- Handle missing labels gracefully (skip or log warning)
- Be injected into export endpoints via dependency injection

## Error Handling

### Error Scenarios

1. **Invalid Format Parameter**
   - Return 400 Bad Request
   - Message: "Invalid format. Supported formats: csv, json, yaml"

2. **Database Query Failure**
   - Return 500 Internal Server Error
   - Log error details for debugging

3. **Encoding Failure**
   - Return 500 Internal Server Error
   - Log error details

4. **Authentication Failure**
   - Return 401 Unauthorized
   - Handled by existing authentication middleware

5. **Authorization Failure**
   - Return 403 Forbidden
   - Handled by existing query handlers

### Error Response Format

Errors will use the existing `ginx.HttpErrorResponse` format for consistency.

## HTTP Response Headers

### Content-Type Headers

- CSV: `text/csv; charset=utf-8`
- JSON: `application/json; charset=utf-8`
- YAML: `application/x-yaml; charset=utf-8`

### Content-Disposition Header

Format: `attachment; filename="<resource>_<timestamp>.<extension>"`

Examples:
- `attachment; filename="labels_2025-11-08T15-30-45.csv"`
- `attachment; filename="providers_2025-11-08T15-30-45.json"`
- `attachment; filename="subscriptions_2025-11-08T15-30-45.yaml"`

Timestamp format: ISO 8601 with colons replaced by hyphens for filesystem compatibility

## Testing Strategy

### Unit Tests

1. **Export Service Tests**
   - Test format parsing and validation
   - Test content type and extension mapping
   - Test each encoder with sample data
   - Test error handling for invalid data

2. **CSV Encoder Tests**
   - Test encoding of simple structs
   - Test handling of nil/optional fields
   - Test handling of nested objects
   - Test handling of slice fields
   - Test empty data sets

3. **JSON Encoder Tests**
   - Test encoding of DTO models
   - Test empty arrays
   - Test streaming behavior

4. **YAML Encoder Tests**
   - Test encoding of DTO models
   - Test empty arrays
   - Test YAML formatting

5. **Endpoint Handler Tests**
   - Test format parameter parsing
   - Test default format behavior
   - Test invalid format handling
   - Test header setting
   - Test integration with query handlers
   - Mock query handlers to test encoding logic

### Integration Tests

1. **End-to-End Export Tests**
   - Test each endpoint with real database data
   - Test authentication and authorization
   - Test each format produces valid output
   - Test empty result sets
   - Test large data sets (performance)

2. **Format Validation Tests**
   - Verify CSV output can be parsed by standard CSV readers
   - Verify JSON output is valid JSON
   - Verify YAML output is valid YAML

## Implementation Considerations

### Performance

1. **Streaming**: Use `gin.Context.Stream()` or write directly to `gin.Context.Writer` to avoid buffering entire response in memory

2. **Query Optimization**: Reuse existing FindAll queries which should already be optimized. For exports, we may want to remove pagination limits to fetch all data.

3. **Memory Management**: For large datasets, consider:
   - Encoding records incrementally rather than loading all into memory
   - Using database cursors if available
   - Setting reasonable limits on export size (can be added later)

### Security

1. **Authentication**: Reuse existing authentication middleware
2. **Authorization**: Leverage existing query handler authorization logic
3. **Data Filtering**: Ensure only user-owned or family-shared data is exported
4. **Rate Limiting**: Consider adding rate limiting for export endpoints (future enhancement)

### Extensibility

The design allows for easy addition of new formats:
1. Implement new encoder function
2. Add format constant
3. Update format validation
4. Add content type mapping

### Dependencies

New external dependencies required:
- `gopkg.in/yaml.v3` - YAML encoding support

Existing dependencies used:
- `encoding/csv` - CSV encoding (standard library)
- `encoding/json` - JSON encoding (standard library)
- `github.com/gin-gonic/gin` - HTTP framework

## Endpoint Specifications

### GET /labels/export

**Query Parameters**:
- `format` (optional): Export format - "csv", "json", or "yaml" (default: "json")

**Response Headers**:
- `Content-Type`: Based on format
- `Content-Disposition`: Suggested filename

**Response Body**: Streamed export data

**Status Codes**:
- 200: Success
- 400: Invalid format parameter
- 401: Unauthorized
- 500: Internal server error

### GET /providers/export

**Query Parameters**:
- `format` (optional): Export format - "csv", "json", or "yaml" (default: "json")

**Response Headers**:
- `Content-Type`: Based on format
- `Content-Disposition`: Suggested filename

**Response Body**: Streamed export data

**Status Codes**:
- 200: Success
- 400: Invalid format parameter
- 401: Unauthorized
- 500: Internal server error

### GET /subscriptions/export

**Query Parameters**:
- `format` (optional): Export format - "csv", "json", or "yaml" (default: "json")

**Response Headers**:
- `Content-Type`: Based on format
- `Content-Disposition`: Suggested filename

**Response Body**: Streamed export data

**Status Codes**:
- 200: Success
- 400: Invalid format parameter
- 401: Unauthorized
- 500: Internal server error

## Dependency Injection

The export service and endpoints will be registered using Uber FX:

**Location**: `internal/adapters/http/export/fx.go`

```go
var Module = fx.Module("export",
    fx.Provide(
        NewExportService,
    ),
)
```

Each endpoint will be registered in its respective handler package's endpoint group constructor, following the existing pattern.

## OpenAPI Documentation

Each export endpoint will include Swagger annotations:

```go
// Handle godoc
//
//  @Summary        Export labels
//  @Description    Export all labels in CSV, JSON, or YAML format
//  @Tags           labels
//  @Produce        text/csv
//  @Produce        application/json
//  @Produce        application/x-yaml
//  @Param          format  query   string  false  "Export format (csv, json, yaml)"  default(json)
//  @Success        200     {file}  file    "Exported labels file"
//  @Failure        400     {object} HttpErrorResponse "Invalid format parameter"
//  @Failure        401     {object} HttpErrorResponse "Unauthorized"
//  @Failure        500     {object} HttpErrorResponse "Internal Server Error"
//  @Router         /labels/export [get]
```

## Migration Path

This is a new feature with no existing data migration required. The implementation will:

1. Add new endpoints without modifying existing ones
2. Reuse existing query handlers and authorization logic
3. Not affect existing API contracts
4. Be backward compatible with all existing functionality
