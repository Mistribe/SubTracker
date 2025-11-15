# Implementation Plan

- [x] 1. Set up export service infrastructure
  - Create export package at internal/adapters/http/export/
  - Create service.go with ExportService interface and format constants (FormatCSV, FormatJSON, FormatYAML)
  - Implement format validation (ParseFormat) and content type mapping (GetContentType, GetFileExtension)
  - Create export DTO models in internal/adapters/http/export/models.go (LabelExportModel, ProviderExportModel, SubscriptionExportModel)
  - Create label_resolver.go with LabelResolver interface for resolving label IDs to names
  - Implement label resolver using existing LabelRepository
  - Write unit tests for format validation and content type mapping
  - _Requirements: 1.5, 1.6, 2.5, 2.6, 3.5, 3.6, 6.1, 6.2_

- [x] 2. Implement CSV encoder
  - [x] 2.1 Create CSV encoder with struct-to-CSV conversion
    - Create csv_encoder.go in internal/adapters/http/export/
    - Write EncodeCSV function that accepts io.Writer and slice of structs
    - Use reflection to read csv struct tags for headers and field order
    - Handle nil/optional pointer fields by writing empty strings
    - Write unit tests for basic CSV encoding with sample structs
    - _Requirements: 1.2, 2.2, 3.2_

  - [x] 2.2 Add CSV encoder support for slice fields
    - Implement joining of string slice fields with ", " separator in CSV encoder
    - Write unit tests for slice field encoding
    - Test empty data sets produce valid CSV with headers only
    - _Requirements: 1.2, 2.2, 2.11, 3.2, 3.11_

- [x] 3. Implement JSON and YAML encoders
  - [x] 3.1 Create JSON encoder
    - Create json_encoder.go in internal/adapters/http/export/
    - Write EncodeJSON function using encoding/json with streaming to io.Writer
    - Write unit tests for JSON encoding with export DTO models
    - Test empty arrays produce valid JSON
    - _Requirements: 1.3, 2.3, 3.3_

  - [x] 3.2 Create YAML encoder
    - Create yaml_encoder.go in internal/adapters/http/export/
    - Write EncodeYAML function using gopkg.in/yaml.v3 with streaming to io.Writer
    - Write unit tests for YAML encoding with export DTO models
    - Test empty arrays produce valid YAML
    - _Requirements: 1.4, 2.4, 3.4_

- [-] 4. Implement label export endpoint
  - [x] 4.1 Create label export handler
    - Create label_export.go in internal/adapters/http/handlers/label/
    - Create ExportEndpoint struct with FindAllQueryHandler dependency and ExportService
    - Implement Handle method with format query parameter parsing (default to "json")
    - Implement Pattern() returning ["/export"], Method() returning GET, and Middlewares() returning nil
    - Set Content-Type header based on format using ExportService.GetContentType
    - Set Content-Disposition header with timestamp filename (e.g., "labels_2025-11-09T15-30-45.csv")
    - _Requirements: 1.1, 1.5, 1.7, 1.8_

  - [x] 4.2 Integrate label export with query handler and encoders
    - Call FindAllQueryHandler with high limit (e.g., 10000) and offset 0 to fetch all labels
    - Transform domain labels to LabelExportModel (map owner type: PersonalOwnerType → "personal", FamilyOwnerType → "family")
    - Route to appropriate encoder (CSV/JSON/YAML) based on format parameter
    - Stream encoded data directly to gin.Context.Writer
    - Handle invalid format with 400 Bad Request error
    - Handle query errors with 500 Internal Server Error
    - _Requirements: 1.1, 1.6, 1.9, 1.10_

  - [x] 4.3 Register label export endpoint
    - Update NewEndpointGroup in endpoint_group.go to accept ExportEndpoint parameter
    - Add ExportEndpoint to routes slice in NewEndpointGroup
    - Add NewExportEndpoint constructor with fx.Annotate as ginfx.Endpoint
    - _Requirements: 1.1_

  - [x] 4.4 Add Swagger documentation for label export
    - Add godoc comments with @Summary, @Description, @Tags, @Produce annotations
    - Document format query parameter with @Param
    - Document 200, 400, 401, 500 responses with @Success and @Failure
    - Add @Router annotation for /labels/export [get]
    - _Requirements: 1.1_

  - [x] 4.5 Write tests for label export endpoint
    - Create label_export_test.go with unit tests using mocked FindAllQueryHandler
    - Test CSV format returns correct Content-Type and valid CSV output
    - Test JSON format returns correct Content-Type and valid JSON output
    - Test YAML format returns correct Content-Type and valid YAML output
    - Test invalid format returns 400 Bad Request
    - Test empty result sets produce valid output with headers/structure
    - Test Content-Disposition header contains timestamp
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 1.8, 1.9_

- [-] 5. Implement provider export endpoint
  - [x] 5.1 Create provider export handler
    - Create provider_export.go in internal/adapters/http/handlers/provider/
    - Create ExportEndpoint struct with FindAllQueryHandler, LabelResolver, and ExportService dependencies
    - Implement Handle method with format query parameter parsing (default to "json")
    - Implement Pattern() returning ["/export"], Method() returning GET, and Middlewares() returning nil
    - Set Content-Type and Content-Disposition headers
    - _Requirements: 2.1, 2.5, 2.7, 2.8_

  - [x] 5.2 Integrate provider export with query handler and encoders
    - Call FindAllQueryHandler with high limit to fetch all providers
    - Extract all label IDs from providers and call LabelResolver.ResolveLabelNames
    - Transform domain providers to ProviderExportModel with label names (not IDs)
    - Route to appropriate encoder based on format parameter
    - Stream encoded data to HTTP response
    - Handle errors with appropriate status codes
    - _Requirements: 2.1, 2.6, 2.9, 2.10, 2.11_

  - [x] 5.3 Register provider export endpoint
    - Update NewEndpointGroup in endpoint_group.go to accept ExportEndpoint parameter
    - Add ExportEndpoint to routes slice
    - Add NewExportEndpoint constructor with fx.Annotate
    - _Requirements: 2.1_

  - [x] 5.4 Add Swagger documentation for provider export
    - Add godoc comments with OpenAPI annotations
    - Document format query parameter, response types, and status codes
    - Add @Router annotation for /providers/export [get]
    - _Requirements: 2.1_

  - [x] 5.5 Write tests for provider export endpoint
    - Create provider_export_test.go with mocked FindAllQueryHandler and LabelResolver
    - Test each format (CSV, JSON, YAML) with correct Content-Type
    - Test invalid format handling returns 400
    - Test empty result sets
    - Test header values (Content-Type, Content-Disposition)
    - Test labels field contains comma-separated label names
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 2.8, 2.9, 2.11_

- [x] 6. Implement subscription export endpoint
  - [x] 6.1 Create subscription export handler
    - Create subscription_export.go in internal/adapters/http/handlers/subscription/
    - Create ExportEndpoint struct with FindAllQueryHandler, LabelResolver, and ExportService dependencies
    - Implement Handle method with format query parameter parsing (default to "json")
    - Implement Pattern() returning ["/export"], Method() returning GET, and Middlewares() returning nil
    - Set Content-Type and Content-Disposition headers
    - _Requirements: 3.1, 3.5, 3.7, 3.8_

  - [x] 6.2 Integrate subscription export with query handler and encoders
    - Call FindAllQueryHandler with high limit and withInactive=true to fetch all subscriptions
    - Extract label IDs from subscriptions and resolve to names using LabelResolver
    - Transform domain subscriptions to SubscriptionExportModel:
      - Format dates as ISO 8601 (YYYY-MM-DD) using time.Format("2006-01-02")
      - Map owner type (PersonalOwnerType → "personal", FamilyOwnerType → "family")
      - Extract customPriceAmount and customPriceCurrency from Price().Amount()
      - Include freeTrialStartDate and freeTrialEndDate only if FreeTrial() is not nil
      - Convert label references to comma-separated label names
    - Route to appropriate encoder based on format
    - Stream encoded data to HTTP response
    - Handle errors with appropriate status codes
    - _Requirements: 3.1, 3.6, 3.9, 3.10, 3.11, 3.12, 3.13, 3.14, 3.15_

  - [x] 6.3 Register subscription export endpoint
    - Update NewEndpointGroup in endpoint_group.go to accept ExportEndpoint parameter
    - Add ExportEndpoint to routes slice
    - Add NewExportEndpoint constructor with fx.Annotate
    - _Requirements: 3.1_

  - [x] 6.4 Add Swagger documentation for subscription export
    - Add godoc comments with OpenAPI annotations
    - Document format query parameter, response types, and status codes
    - Add @Router annotation for /subscriptions/export [get]
    - _Requirements: 3.1_

  - [x] 6.5 Write tests for subscription export endpoint
    - Create subscription_export_test.go with mocked FindAllQueryHandler and LabelResolver
    - Test each format (CSV, JSON, YAML) with correct Content-Type
    - Test invalid format handling returns 400
    - Test empty result sets
    - Test header values
    - Test labels field contains comma-separated label names
    - Test free trial dates are included when present and omitted when absent
    - Test date formatting is ISO 8601 (YYYY-MM-DD)
    - Test customPriceAmount and customPriceCurrency are correctly extracted
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7, 3.8, 3.9, 3.11, 3.12, 3.13, 3.14, 3.15_

- [ ] 7. Write integration tests
  - [ ] 7.1 Create integration tests for label export
    - Create label_export_integration_test.go in integration/ directory
    - Test label export with real database data and authentication
    - Verify CSV output is valid and parseable by encoding/csv
    - Verify JSON output is valid and parseable by encoding/json
    - Verify YAML output is valid and parseable by gopkg.in/yaml.v3
    - Test authorization ensures only user's labels are exported
    - _Requirements: 1.10, 5.1, 5.2, 5.3_

  - [ ] 7.2 Create integration tests for provider export
    - Create provider_export_integration_test.go in integration/ directory
    - Test provider export with real database data and authentication
    - Verify CSV output is valid and parseable
    - Verify JSON output is valid and parseable
    - Verify YAML output is valid and parseable
    - Test authorization ensures only user's providers are exported
    - Verify label names are correctly resolved in output
    - _Requirements: 2.10, 5.1, 5.2, 5.3_

  - [ ] 7.3 Create integration tests for subscription export
    - Create subscription_export_integration_test.go in integration/ directory
    - Test subscription export with real database data and authentication
    - Verify CSV output is valid and parseable
    - Verify JSON output is valid and parseable
    - Verify YAML output is valid and parseable
    - Test authorization ensures only user's subscriptions are exported
    - Verify label names are correctly resolved in output
    - Verify date formatting and free trial handling
    - _Requirements: 3.10, 5.1, 5.2, 5.3_

- [ ] 8. Update OpenAPI documentation
  - Run swag init or make swagger command to regenerate Swagger documentation files
  - Verify export endpoints appear in swagger.json and swagger.yaml
  - Manually test endpoints using Swagger UI at /swagger/index.html
  - Verify all three export endpoints are accessible and documented correctly
  - _Requirements: 1.1, 2.1, 3.1_
