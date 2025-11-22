# Requirements Document

## Introduction

This feature enhances the existing import functionality by adding comprehensive template file support for all three file formats (CSV, JSON, and YAML) for each entity type (Labels, Providers, and Subscriptions), and enabling the ability to import entities with predefined UUIDs. These enhancements provide users with more flexibility and control over the import process, allowing them to maintain consistent identifiers across systems and have access to complete template examples in their preferred format.

## Requirements

### Requirement 1: Complete Template File Coverage

**User Story:** As a user preparing to import data, I want to download template files in all three supported formats (CSV, JSON, YAML) for each entity type, so that I can choose the format that best suits my workflow and data source.

#### Acceptance Criteria

1. WHEN the user is on the Labels import page THEN the system SHALL provide download links for labels templates in CSV, JSON, and YAML formats
2. WHEN the user is on the Providers import page THEN the system SHALL provide download links for providers templates in CSV, JSON, and YAML formats
3. WHEN the user is on the Subscriptions import page THEN the system SHALL provide download links for subscriptions templates in CSV, JSON, and YAML formats
4. WHEN the user clicks a template download link THEN the system SHALL download the corresponding template file
5. WHEN a template file is downloaded THEN the file SHALL contain example data demonstrating all available fields for that entity type
6. WHEN a template file is downloaded THEN the file SHALL include comments or examples showing optional vs required fields
7. WHEN a template file is downloaded THEN the file SHALL be properly formatted and immediately usable for import

### Requirement 2: UUID Support in Import

**User Story:** As a user importing entities from another system, I want to specify custom UUIDs for my entities during import, so that I can maintain consistent identifiers across different systems and environments.

#### Acceptance Criteria

1. WHEN a record in the import file includes an "id" field THEN the system SHALL accept the id value if it is a valid UUID format
2. WHEN a record includes a valid UUID in the "id" field THEN the system SHALL use that UUID when creating the entity via the API
3. IF a record includes an "id" field with an invalid UUID format THEN the system SHALL display a validation error for that record
4. IF a record includes an "id" field with an invalid UUID format THEN the system SHALL prevent import of that record until corrected
5. WHEN a record does not include an "id" field THEN the system SHALL allow the backend to auto-generate a UUID for that entity
6. WHEN displaying records in the preview table THEN the system SHALL show the "id" field if present in the imported data
7. WHEN a validation error occurs for an invalid UUID THEN the system SHALL display a clear message indicating the UUID format is incorrect
8. WHEN importing a record with a custom UUID THEN the system SHALL handle API errors if the UUID already exists in the system
9. IF the API returns a conflict error for a duplicate UUID THEN the system SHALL display an appropriate error message to the user

### Requirement 3: Template File Content Standards

**User Story:** As a user learning how to format import files, I want template files that clearly demonstrate all available fields and their expected formats, so that I can create valid import files without trial and error.

#### Acceptance Criteria

1. WHEN a template file is provided THEN it SHALL include examples of all required fields for that entity type
2. WHEN a template file is provided THEN it SHALL include examples of commonly used optional fields
3. WHEN a template file is provided THEN it SHALL demonstrate proper data formatting (dates, colors, URLs, enums)
4. WHEN a CSV template is provided THEN it SHALL include a header row with all field names
5. WHEN a JSON template is provided THEN it SHALL be a valid JSON array with multiple example records
6. WHEN a YAML template is provided THEN it SHALL be valid YAML with multiple example records
7. WHEN a template includes the "id" field THEN it SHALL show example UUID values in the correct format
8. WHEN a template is provided THEN it SHALL include diverse examples demonstrating different valid configurations

### Requirement 4: Template File Accessibility

**User Story:** As a user on the import page, I want easy access to template downloads, so that I can quickly get started with the import process.

#### Acceptance Criteria

1. WHEN the user is on the import page THEN the system SHALL display template download links in a prominent, easily discoverable location
2. WHEN template download links are displayed THEN each link SHALL clearly indicate the file format (CSV, JSON, or YAML)
3. WHEN template download links are displayed THEN they SHALL be grouped or organized by entity type
4. WHEN the user hovers over a template download link THEN the system SHALL provide visual feedback
5. WHEN template files are downloaded THEN they SHALL have descriptive filenames indicating the entity type and format
6. WHEN the user has not yet uploaded a file THEN the template download section SHALL be visible and accessible
7. WHEN the user has uploaded a file THEN the template download section SHALL remain accessible for reference
