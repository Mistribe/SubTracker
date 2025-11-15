# Requirements Document

## Introduction

This feature enables users to import Labels, Providers, and Subscriptions from external files (CSV, JSON, or YAML formats) directly in the web application. The import process is entirely client-side, with parsed data being sent to the backend one record at a time to avoid overloading the server. Users can preview all imported data and selectively choose which items to import, providing full control over the import process.

## Requirements

### Requirement 1: Import Menu Access

**User Story:** As a user managing labels, providers, or subscriptions, I want to access an import feature from each respective page, so that I can easily bring in data from external files.

#### Acceptance Criteria

1. WHEN the user is on the Labels page THEN the system SHALL display an additional menu option labeled "Import from file"
2. WHEN the user is on the Providers page THEN the system SHALL display an additional menu option labeled "Import from file"
3. WHEN the user is on the Subscriptions page THEN the system SHALL display an additional menu option labeled "Import from file"
4. WHEN the user clicks "Import from file" THEN the system SHALL navigate to a new import page specific to that entity type

### Requirement 2: File Selection and Upload Interface

**User Story:** As a user wanting to import data, I want to select or drag-and-drop a file on the import page, so that I can easily provide the data I want to import.

#### Acceptance Criteria

1. WHEN the user navigates to the import page THEN the system SHALL display a file selection area in the center of the page
2. WHEN the user clicks the file selection area THEN the system SHALL open a file picker dialog
3. WHEN the user drags a file over the selection area THEN the system SHALL provide visual feedback indicating the drop zone is active
4. WHEN the user drops a file on the selection area THEN the system SHALL accept the file for processing
5. WHEN the user selects or drops a file THEN the system SHALL validate that the file format is CSV, JSON, or YAML
6. IF the file format is not CSV, JSON, or YAML THEN the system SHALL display an error message indicating unsupported format
7. WHEN a valid file is selected THEN the system SHALL parse the file entirely on the client side
8. WHEN the user is on the import page THEN the system SHALL provide download links for template files in all three formats (CSV, JSON, and YAML)

### Requirement 3: Client-Side File Parsing

**User Story:** As a user who has selected a file to import, I want the system to parse the file on my browser, so that my data remains secure and the backend is not overloaded with large file uploads.

#### Acceptance Criteria

1. WHEN a CSV file is selected THEN the system SHALL parse the CSV format and extract all records
2. WHEN a JSON file is selected THEN the system SHALL parse the JSON format and extract all records
3. WHEN a YAML file is selected THEN the system SHALL parse the YAML format and extract all records
4. IF the file parsing fails THEN the system SHALL display a clear error message explaining the parsing failure
5. WHEN parsing is successful THEN the system SHALL extract all entity records from the file
6. WHEN parsing is complete THEN the system SHALL NOT send the original file to the backend
7. WHEN parsing is complete THEN the system SHALL display the parsed records to the user

### Requirement 4: Import Preview and Selection

**User Story:** As a user who has uploaded a file, I want to see all the parsed records in a list, so that I can review and select which items to import.

#### Acceptance Criteria

1. WHEN file parsing is complete THEN the system SHALL display all parsed records in a table or list format
2. WHEN displaying records THEN the system SHALL show relevant fields for each entity type (name, description, etc.)
3. WHEN displaying the list THEN the system SHALL provide a checkbox for each individual record
4. WHEN displaying the list THEN the system SHALL provide a "Select All" checkbox to select all records at once
5. WHEN the user clicks an individual checkbox THEN the system SHALL toggle the selection state for that record
6. WHEN the user clicks "Select All" THEN the system SHALL select or deselect all records
7. WHEN records are displayed THEN the system SHALL show the total count of records found in the file

### Requirement 5: Individual Record Import

**User Story:** As a user reviewing imported records, I want to import individual records one at a time, so that I have precise control over what gets added to my account.

#### Acceptance Criteria

1. WHEN the user selects a single record and clicks an import action THEN the system SHALL make an API call to the backend to create that record
2. WHEN the API call is in progress THEN the system SHALL display a loading indicator for that specific record
3. IF the API call succeeds THEN the system SHALL mark the record as successfully imported with a visual indicator
4. IF the API call fails THEN the system SHALL display an error message specific to that record
5. IF the API call fails THEN the system SHALL allow the user to retry importing that record
6. WHEN a record is successfully imported THEN the system SHALL keep the record visible in the list with its success status

### Requirement 6: Bulk Import All Records

**User Story:** As a user with many records to import, I want to import all selected records at once, so that I can efficiently add multiple items without manual intervention for each one.

#### Acceptance Criteria

1. WHEN the user has selected multiple records and clicks "Import All" THEN the system SHALL import each selected record sequentially
2. WHEN importing all records THEN the system SHALL make individual API calls to the backend for each record
3. WHEN importing all records THEN the system SHALL NOT send multiple records in a single API call
4. WHEN importing all records THEN the system SHALL display progress indicating how many records have been imported
5. WHEN importing all records THEN the system SHALL update each record's status individually as it is processed
6. IF any record fails during bulk import THEN the system SHALL continue importing remaining records
7. IF any record fails during bulk import THEN the system SHALL mark the failed record with an error indicator
8. WHEN bulk import is complete THEN the system SHALL display a summary showing successful and failed imports
9. WHEN bulk import is in progress THEN the system SHALL allow the user to cancel the remaining imports

### Requirement 7: Entity-Specific Field Mapping

**User Story:** As a user importing different types of entities, I want the system to correctly map file fields to the appropriate entity properties, so that my data is imported accurately.

#### Acceptance Criteria

1. WHEN importing Labels THEN the system SHALL map file fields to Label entity properties (name, color, etc.)
2. WHEN importing Providers THEN the system SHALL map file fields to Provider entity properties (name, description, etc.)
3. WHEN importing Subscriptions THEN the system SHALL map file fields to Subscription entity properties (name, provider, price, billing cycle, etc.)
4. IF required fields are missing from a record THEN the system SHALL mark that record as invalid and display which fields are missing
5. IF field data types do not match expected types THEN the system SHALL display a validation error for that record
6. WHEN displaying records THEN the system SHALL highlight any validation errors before import attempts
7. WHEN a record includes an id field with a valid UUID THEN the system SHALL use that id when creating the entity
8. IF a record includes an id field with an invalid UUID format THEN the system SHALL display a validation error for that record
9. WHEN a record does not include an id field THEN the system SHALL allow the backend to generate a new UUID for that entity

### Requirement 8: Import Page Navigation and State Management

**User Story:** As a user on the import page, I want clear navigation options and state feedback, so that I understand where I am in the import process and can navigate back if needed.

#### Acceptance Criteria

1. WHEN the user is on the import page THEN the system SHALL display a clear page title indicating which entity type is being imported
2. WHEN the user is on the import page THEN the system SHALL provide a way to navigate back to the original entity list page
3. WHEN the user navigates away from the import page THEN the system SHALL warn the user if there are unimported records
4. WHEN file parsing is in progress THEN the system SHALL display a loading indicator
5. WHEN import operations are in progress THEN the system SHALL prevent the user from selecting a new file
6. WHEN all selected records are successfully imported THEN the system SHALL provide an option to import another file or return to the entity list
