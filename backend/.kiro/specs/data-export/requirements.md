# Requirements Document

## Introduction

This feature enables users to export their Labels, Providers, and Subscriptions data in multiple formats (CSV, JSON, YAML) through dedicated HTTP streaming endpoints. The export functionality allows users to download their data for backup, analysis, or migration purposes without storing files on the server. Each resource type (Labels, Providers, Subscriptions) will have its own export endpoint that streams the data directly to the client in the requested format.

## Requirements

### Requirement 1: Export Labels

**User Story:** As a user, I want to export all my labels to a file in CSV, JSON, or YAML format, so that I can backup or analyze my label data outside the application.

#### Acceptance Criteria

1. WHEN a user requests `/labels/export` with a valid format parameter THEN the system SHALL stream the labels data in the requested format
2. WHEN the format parameter is "csv" THEN the system SHALL return labels data as CSV with headers (id, name, color, ownerType, ownerFamilyId)
3. WHEN the format parameter is "json" THEN the system SHALL return labels data as a JSON array with fields (id, name, color, ownerType, ownerFamilyId)
4. WHEN the format parameter is "yaml" THEN the system SHALL return labels data as a YAML document with fields (id, name, color, ownerType, ownerFamilyId)
5. WHEN no format parameter is provided THEN the system SHALL default to JSON format
6. WHEN an invalid format parameter is provided THEN the system SHALL return a 400 Bad Request error with a descriptive message
7. WHEN the export is requested THEN the system SHALL set appropriate Content-Type headers based on the format (text/csv, application/json, application/yaml)
8. WHEN the export is requested THEN the system SHALL set Content-Disposition header to suggest a filename with timestamp (e.g., "labels_2025-11-08.csv")
9. WHEN the user has no labels THEN the system SHALL return an empty dataset in the requested format
10. WHEN the export is requested THEN the system SHALL only include labels that the authenticated user has permission to view

### Requirement 2: Export Providers

**User Story:** As a user, I want to export all my providers to a file in CSV, JSON, or YAML format, so that I can backup or share my provider configurations.

#### Acceptance Criteria

1. WHEN a user requests `/providers/export` with a valid format parameter THEN the system SHALL stream the providers data in the requested format
2. WHEN the format parameter is "csv" THEN the system SHALL return providers data as CSV with headers (id, name, description, url, iconUrl, pricingPageUrl, labels)
3. WHEN the format parameter is "json" THEN the system SHALL return providers data as a JSON array with fields (id, name, description, url, iconUrl, pricingPageUrl, labels)
4. WHEN the format parameter is "yaml" THEN the system SHALL return providers data as a YAML document with fields (id, name, description, url, iconUrl, pricingPageUrl, labels)
5. WHEN no format parameter is provided THEN the system SHALL default to JSON format
6. WHEN an invalid format parameter is provided THEN the system SHALL return a 400 Bad Request error with a descriptive message
7. WHEN the export is requested THEN the system SHALL set appropriate Content-Type headers based on the format
8. WHEN the export is requested THEN the system SHALL set Content-Disposition header to suggest a filename with timestamp (e.g., "providers_2025-11-08.json")
9. WHEN the user has no providers THEN the system SHALL return an empty dataset in the requested format
10. WHEN the export is requested THEN the system SHALL only include providers that the authenticated user has permission to view
11. WHEN exporting providers with multiple labels THEN the system SHALL format the labels field as a comma-separated list of label names (not IDs)

### Requirement 3: Export Subscriptions

**User Story:** As a user, I want to export all my subscriptions to a file in CSV, JSON, or YAML format, so that I can analyze my subscription spending and manage my recurring payments.

#### Acceptance Criteria

1. WHEN a user requests `/subscriptions/export` with a valid format parameter THEN the system SHALL stream the subscriptions data in the requested format
2. WHEN the format parameter is "csv" THEN the system SHALL return subscriptions data as CSV with headers (id, providerId, friendlyName, startDate, endDate, recurrency, customRecurrency, customPriceAmount, customPriceCurrency, ownerType, freeTrialStartDate, freeTrialEndDate, labels)
3. WHEN the format parameter is "json" THEN the system SHALL return subscriptions data as a JSON array with fields (id, providerId, friendlyName, startDate, endDate, recurrency, customRecurrency, customPriceAmount, customPriceCurrency, ownerType, freeTrialStartDate, freeTrialEndDate, labels)
4. WHEN the format parameter is "yaml" THEN the system SHALL return subscriptions data as a YAML document with fields (id, providerId, friendlyName, startDate, endDate, recurrency, customRecurrency, customPriceAmount, customPriceCurrency, ownerType, freeTrialStartDate, freeTrialEndDate, labels)
5. WHEN no format parameter is provided THEN the system SHALL default to JSON format
6. WHEN an invalid format parameter is provided THEN the system SHALL return a 400 Bad Request error with a descriptive message
7. WHEN the export is requested THEN the system SHALL set appropriate Content-Type headers based on the format
8. WHEN the export is requested THEN the system SHALL set Content-Disposition header to suggest a filename with timestamp (e.g., "subscriptions_2025-11-08.yaml")
9. WHEN the user has no subscriptions THEN the system SHALL return an empty dataset in the requested format
10. WHEN the export is requested THEN the system SHALL only include subscriptions that the authenticated user has permission to view
11. WHEN exporting subscriptions with label references THEN the system SHALL format the labels field as a comma-separated list of label names (not IDs)
12. WHEN exporting subscriptions THEN the system SHALL use the providerId field to reference the provider (not provider name)
13. WHEN exporting subscriptions with free trial THEN the system SHALL include freeTrialStartDate and freeTrialEndDate fields
14. WHEN exporting subscriptions without free trial THEN the system SHALL omit or leave empty the freeTrialStartDate and freeTrialEndDate fields
15. WHEN exporting subscriptions THEN the system SHALL use customPriceAmount and customPriceCurrency for the price (not monthly/yearly breakdowns)

### Requirement 4: HTTP Streaming and Performance

**User Story:** As a user, I want my export to stream directly to my browser without server-side file storage, so that the export is fast and doesn't consume unnecessary server resources.

#### Acceptance Criteria

1. WHEN an export is requested THEN the system SHALL stream the data directly to the HTTP response without writing to disk
2. WHEN an export is requested THEN the system SHALL fetch all relevant data in a single query or minimal queries
3. WHEN an export contains a large dataset THEN the system SHALL handle the export without timing out or consuming excessive memory
4. WHEN an export is in progress THEN the system SHALL flush data to the client progressively if the dataset is large
5. WHEN an export fails due to a database error THEN the system SHALL return a 500 Internal Server Error with appropriate logging

### Requirement 5: Authentication and Authorization

**User Story:** As a system administrator, I want to ensure that users can only export their own data, so that data privacy and security are maintained.

#### Acceptance Criteria

1. WHEN an unauthenticated user requests an export endpoint THEN the system SHALL return a 401 Unauthorized error
2. WHEN an authenticated user requests an export THEN the system SHALL only include data owned by that user or their family
3. WHEN an authenticated user requests an export THEN the system SHALL apply the same authorization rules as the corresponding GET all endpoints
4. WHEN a user's session expires during export THEN the system SHALL handle the error gracefully

### Requirement 6: Format Query Parameter Validation

**User Story:** As a developer, I want clear validation of the format parameter, so that users receive helpful error messages when they provide invalid input.

#### Acceptance Criteria

1. WHEN the format parameter is provided THEN the system SHALL accept only "csv", "json", or "yaml" (case-insensitive)
2. WHEN an invalid format is provided THEN the system SHALL return a 400 Bad Request with message "Invalid format. Supported formats: csv, json, yaml"
3. WHEN multiple format parameters are provided THEN the system SHALL use the first valid format parameter
4. WHEN the format parameter is empty THEN the system SHALL treat it as not provided and use the default format
