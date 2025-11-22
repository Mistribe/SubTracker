# Implementation Plan

- [x] 1. Create comprehensive template files for all entity types and formats
  - [x] 1.1 Create Labels template files
    - Create `public/templates/labels/labels-template.csv` with example data including ID field
    - Create `public/templates/labels/labels-template.json` with example data including ID field
    - Create `public/templates/labels/labels-template.yaml` with example data including ID field
    - Include examples with and without IDs to demonstrate both use cases
    - Ensure all required and common optional fields are represented
    - _Requirements: 1.1, 1.5, 1.6, 1.7, 3.1, 3.2, 3.3, 3.4, 3.7, 3.8_

  - [x] 1.2 Create Providers template files
    - Create `public/templates/providers/providers-template.csv` with example data including ID field
    - Create `public/templates/providers/providers-template.json` with example data including ID field
    - Create `public/templates/providers/providers-template.yaml` with example data including ID field
    - Include examples with and without IDs to demonstrate both use cases
    - Ensure all required and common optional fields are represented
    - _Requirements: 1.2, 1.5, 1.6, 1.7, 3.1, 3.2, 3.3, 3.5, 3.6, 3.7, 3.8_

  - [x] 1.3 Create Subscriptions template files
    - Create `public/templates/subscriptions/subscriptions-template.csv` with example data including ID field
    - Create `public/templates/subscriptions/subscriptions-template.json` with example data including ID field
    - Create `public/templates/subscriptions/subscriptions-template.yaml` with example data including ID field
    - Include examples with and without IDs to demonstrate both use cases
    - Ensure all required and common optional fields are represented
    - _Requirements: 1.3, 1.5, 1.6, 1.7, 3.1, 3.2, 3.3, 3.5, 3.6, 3.7, 3.8_

  - [x] 1.4 Create template documentation
    - Create `public/templates/README.md` explaining template usage
    - Document UUID format requirements and when to use custom IDs
    - Include field descriptions for each entity type
    - Add examples of common import scenarios
    - _Requirements: 3.1, 3.2, 3.3, 3.7, 3.8_

- [x] 2. Implement UUID validation functionality
  - [x] 2.1 Create UUID validation utility
    - Create `src/utils/uuidValidation.ts` with UUID validation function
    - Implement regex-based UUID format validation (RFC 4122 compliant)
    - Add function to validate and sanitize UUID strings (trim whitespace)
    - Create helper function to check if value is empty/null/undefined
    - Export validation error message generator
    - _Requirements: 2.1, 2.3, 2.4, 2.7_

  - [x] 2.2 Write unit tests for UUID validation
    - Create test file `src/utils/__tests__/uuidValidation.test.ts`
    - Write tests for valid UUID formats (v1, v4, v5)
    - Write tests for invalid UUID formats (wrong length, invalid characters, wrong structure)
    - Write tests for empty/null/undefined values (should be valid)
    - Write tests for whitespace handling
    - Write tests for error message generation
    - _Requirements: 2.1, 2.3, 2.4, 2.5, 2.7_

- [x] 3. Enhance field mappers to support ID field
  - [x] 3.1 Update Label field mapper
    - Modify `src/services/importMapper.ts` LabelFieldMapper class
    - Add ID field mapping logic using UUID validation utility
    - Generate validation errors for invalid UUID formats
    - Handle empty/null ID values (allow backend generation)
    - Update mapper return type to include ID field
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

  - [x] 3.2 Update Provider field mapper
    - Modify `src/services/importMapper.ts` ProviderFieldMapper class
    - Add ID field mapping logic using UUID validation utility
    - Generate validation errors for invalid UUID formats
    - Handle empty/null ID values (allow backend generation)
    - Update mapper return type to include ID field
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

  - [x] 3.3 Update Subscription field mapper
    - Modify `src/services/importMapper.ts` SubscriptionFieldMapper class
    - Add ID field mapping logic using UUID validation utility
    - Generate validation errors for invalid UUID formats
    - Handle empty/null ID values (allow backend generation)
    - Update mapper return type to include ID field
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

  - [x] 3.4 Write unit tests for enhanced field mappers
    - Update test file `src/services/__tests__/importMapper.test.ts`
    - Write tests for ID field mapping with valid UUIDs
    - Write tests for ID field validation errors with invalid UUIDs
    - Write tests for records without ID fields
    - Write tests for mixed scenarios (some records with IDs, some without)
    - Write tests for whitespace handling in ID fields
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7_

- [x] 4. Create TemplateDownloadSection component
  - [x] 4.1 Implement TemplateDownloadSection component
    - Create `src/components/import/TemplateDownloadSection.tsx` component
    - Accept entityType prop ('labels' | 'providers' | 'subscriptions')
    - Render three download links for CSV, JSON, and YAML formats
    - Add appropriate icons for each format
    - Implement descriptive labels and helper text
    - Add proper download attributes to links
    - Style component with responsive layout
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 4.1, 4.2, 4.3, 4.5_

  - [x] 4.2 Add accessibility features to TemplateDownloadSection
    - Add ARIA labels for download links
    - Implement keyboard navigation support
    - Add focus indicators for interactive elements
    - Ensure screen reader compatibility
    - Add hover states for visual feedback
    - _Requirements: 4.1, 4.2, 4.4, 4.5_

  - [x] 4.3 Write unit tests for TemplateDownloadSection
    - Create test file `src/components/import/__tests__/TemplateDownloadSection.test.tsx`
    - Write tests for rendering all three format links
    - Write tests for correct file paths based on entity type
    - Write tests for download attribute presence
    - Write tests for accessibility attributes
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 4.1, 4.2, 4.3_

- [x] 5. Enhance ImportPreviewTable to display ID column
  - [x] 5.1 Update ImportPreviewTable component
    - Modify `src/components/import/ImportPreviewTable.tsx` component
    - Add conditional ID column that appears when any record has an ID
    - Display UUID values in monospace font
    - Show "Auto-generated" text for records without IDs
    - Add tooltip support for full UUID display
    - Update column configuration to include ID field
    - _Requirements: 2.6, 2.7_

  - [x] 5.2 Enhance validation error display for UUID errors
    - Update error rendering in ImportPreviewTable component
    - Add specific styling for UUID validation errors
    - Display clear error messages for invalid UUID formats
    - Show expected UUID format in error messages
    - Highlight invalid UUID values in error display
    - _Requirements: 2.3, 2.4, 2.7_

  - [x] 5.3 Write tests for enhanced ImportPreviewTable
    - Update test file `src/components/import/__tests__/ImportPreviewTable.test.tsx`
    - Write tests for ID column display when IDs present
    - Write tests for ID column hidden when no IDs present
    - Write tests for "Auto-generated" display for empty IDs
    - Write tests for UUID validation error display
    - Write tests for monospace font rendering
    - _Requirements: 2.6, 2.7_

- [x] 6. Integrate TemplateDownloadSection into import pages
  - [x] 6.1 Update ImportLabelsPage
    - Modify `src/pages/ImportLabelsPage.tsx` to import TemplateDownloadSection
    - Add TemplateDownloadSection component above FileUploadZone
    - Pass 'labels' as entityType prop
    - Ensure component remains visible after file upload
    - _Requirements: 1.1, 1.4, 4.1, 4.6, 4.7_

  - [x] 6.2 Update ImportProvidersPage
    - Modify `src/pages/ImportProvidersPage.tsx` to import TemplateDownloadSection
    - Add TemplateDownloadSection component above FileUploadZone
    - Pass 'providers' as entityType prop
    - Ensure component remains visible after file upload
    - _Requirements: 1.2, 1.4, 4.1, 4.6, 4.7_

  - [x] 6.3 Update ImportSubscriptionsPage
    - Modify `src/pages/ImportSubscriptionsPage.tsx` to import TemplateDownloadSection
    - Add TemplateDownloadSection component above FileUploadZone
    - Pass 'subscriptions' as entityType prop
    - Ensure component remains visible after file upload
    - _Requirements: 1.3, 1.4, 4.1, 4.6, 4.7_

- [x] 7. Enhance import manager to handle UUID conflicts
  - [x] 7.1 Update useImportManager hook
    - Modify `src/hooks/import/useImportManager.ts` hook
    - Add error handling for 409 Conflict responses (duplicate UUID)
    - Add error handling for 400 Bad Request responses (invalid UUID from API)
    - Generate user-friendly error messages for UUID conflicts
    - Display specific UUID in conflict error messages
    - Continue importing remaining records after UUID conflict
    - _Requirements: 2.8, 2.9_

  - [x] 7.2 Write tests for UUID conflict handling
    - Update test file `src/hooks/import/__tests__/useImportManager.test.ts`
    - Write tests for handling 409 Conflict responses
    - Write tests for handling 400 Bad Request responses
    - Write tests for error message generation with UUID conflicts
    - Write tests for continuing import after UUID conflict
    - Write tests for retry functionality with failed UUID imports
    - _Requirements: 2.8, 2.9_

- [x] 8. Update TypeScript types for ID field support
  - [x] 8.1 Update import type definitions
    - Modify `src/types/import.ts` to add ID field to import record types
    - Update ImportRecord interface to include optional id field
    - Update entity-specific import record types (Label, Provider, Subscription)
    - Add UUIDValidationError type definition
    - Update ValidationError type to support UUID-specific errors
    - _Requirements: 2.1, 2.2, 2.3, 2.6_

- [x] 9. Write integration tests for UUID import functionality
  - [x] 9.1 Test file upload and parsing with IDs
    - Create or update integration test file for import flow
    - Write test for uploading CSV file with ID fields
    - Write test for uploading JSON file with ID fields
    - Write test for uploading YAML file with ID fields
    - Write test for parsing files with mixed ID presence (some with, some without)
    - Verify ID values are correctly extracted and validated
    - _Requirements: 2.1, 2.2, 2.3, 2.5, 2.6_

  - [x] 9.2 Test import with valid UUIDs
    - Write test for importing records with valid UUIDs
    - Write test for importing records without UUIDs
    - Write test for mixed import (some with UUIDs, some without)
    - Verify API calls include ID field when present
    - Verify API calls omit ID field when not present
    - _Requirements: 2.1, 2.2, 2.5, 2.6_

  - [x] 9.3 Test UUID validation errors
    - Write test for invalid UUID format detection
    - Write test for validation error display in preview table
    - Write test for preventing import of records with invalid UUIDs
    - Write test for allowing import of valid records while blocking invalid ones
    - _Requirements: 2.3, 2.4, 2.7_

  - [x] 9.4 Test UUID conflict handling
    - Write test for duplicate UUID API error (409 Conflict)
    - Write test for error message display for UUID conflicts
    - Write test for continuing import after UUID conflict
    - Write test for retry functionality after UUID conflict resolution
    - _Requirements: 2.8, 2.9_

- [ ] 10. Write end-to-end tests for complete enhancement workflow
  - [ ] 10.1 Test template download and import with UUIDs for Labels
    - Write E2E test: navigate to labels import page
    - Download CSV template file
    - Verify template file contains ID field examples
    - Upload template file with custom UUIDs
    - Verify ID column appears in preview table
    - Import records with UUIDs
    - Verify labels created with specified UUIDs
    - _Requirements: 1.1, 1.4, 1.5, 1.7, 2.1, 2.2, 2.6, 3.4, 3.7, 4.1_

  - [ ] 10.2 Test template download and import with UUIDs for Providers
    - Write E2E test: navigate to providers import page
    - Download JSON template file
    - Verify template file contains ID field examples
    - Upload template file with custom UUIDs
    - Verify ID column appears in preview table
    - Import records with UUIDs
    - Verify providers created with specified UUIDs
    - _Requirements: 1.2, 1.4, 1.5, 1.7, 2.1, 2.2, 2.6, 3.5, 3.7, 4.1_

  - [ ] 10.3 Test template download and import with UUIDs for Subscriptions
    - Write E2E test: navigate to subscriptions import page
    - Download YAML template file
    - Verify template file contains ID field examples
    - Upload template file with custom UUIDs
    - Verify ID column appears in preview table
    - Import records with UUIDs
    - Verify subscriptions created with specified UUIDs
    - _Requirements: 1.3, 1.4, 1.5, 1.7, 2.1, 2.2, 2.6, 3.6, 3.7, 4.1_

  - [ ] 10.4 Test error scenarios with invalid UUIDs
    - Write E2E test: upload file with invalid UUID formats
    - Verify validation errors displayed in preview table
    - Verify invalid records cannot be imported
    - Verify valid records can still be imported
    - Test UUID conflict scenario (duplicate UUID)
    - Verify appropriate error message displayed for conflicts
    - _Requirements: 2.3, 2.4, 2.7, 2.8, 2.9_

  - [ ] 10.5 Test mixed import scenarios
    - Write E2E test: upload file with some records having IDs, some without
    - Verify ID column displays correctly for mixed data
    - Verify "Auto-generated" shown for records without IDs
    - Import all records
    - Verify records with IDs use specified UUIDs
    - Verify records without IDs get auto-generated UUIDs
    - _Requirements: 2.2, 2.5, 2.6_

- [x] 11. Update documentation and help text
  - [x] 11.1 Add UUID import documentation to import pages
    - Update help text in ImportLabelsPage to mention UUID support
    - Update help text in ImportProvidersPage to mention UUID support
    - Update help text in ImportSubscriptionsPage to mention UUID support
    - Add explanation of when to use custom UUIDs
    - Add UUID format requirements to help text
    - Include link to template README for more details
    - _Requirements: 3.7, 3.8, 4.1_

  - [x] 11.2 Create inline help for UUID field
    - Add tooltip or info icon next to ID column header
    - Provide explanation of UUID field purpose
    - Show example of valid UUID format
    - Explain auto-generation behavior when ID is omitted
    - _Requirements: 2.5, 2.6, 2.7, 3.7_

- [x] 12. Verify template file accessibility and correctness
  - [x] 12.1 Validate all template files
    - Write validation script to parse all template files
    - Verify CSV templates have correct headers and parse correctly
    - Verify JSON templates are valid JSON and parse correctly
    - Verify YAML templates are valid YAML and parse correctly
    - Verify all example UUIDs in templates are valid format
    - Verify all templates include both records with and without IDs
    - _Requirements: 1.5, 1.6, 1.7, 3.1, 3.2, 3.3, 3.4, 3.5, 3.6, 3.7_

  - [x] 12.2 Test template file downloads
    - Write test to verify template files are accessible at correct paths
    - Test download functionality for all template files
    - Verify downloaded files have correct filenames
    - Verify downloaded files have correct content-type headers
    - _Requirements: 1.4, 1.5, 1.7, 4.5_

- [x] 13. Performance testing and optimization
  - [x] 13.1 Test UUID validation performance
    - Write performance test for UUID validation with large datasets
    - Verify validation doesn't significantly impact parsing time
    - Optimize regex compilation if needed
    - Test with 1000+ records containing UUIDs
    - _Requirements: 2.1, 2.3_

  - [x] 13.2 Test template file loading
    - Verify template files load quickly
    - Test template download performance
    - Ensure template files are properly cached
    - _Requirements: 1.4, 1.5, 4.5_
