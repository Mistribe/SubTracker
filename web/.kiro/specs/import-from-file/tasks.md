# Implementation Plan

- [x] 1. Set up project dependencies and type definitions
  - Install required npm packages: papaparse, js-yaml, and their TypeScript type definitions
  - Create shared TypeScript type definitions file for import functionality at `src/types/import.ts`
  - _Requirements: 2.7, 3.1, 3.2, 3.3_

- [x] 2. Implement file parser service
  - [x] 2.1 Create base file parser service
    - Implement `src/services/fileParser.ts` with FileReader-based file reading
    - Create CSV parser using papaparse library
    - Create JSON parser using native JSON.parse
    - Create YAML parser using js-yaml library
    - Implement format detection based on file extension
    - Add error handling for malformed files
    - _Requirements: 3.1, 3.2, 3.3, 3.4_

  - [x] 2.2 Write unit tests for file parser
    - Create test file `src/services/__tests__/fileParser.test.ts`
    - Write tests for CSV parsing with various formats
    - Write tests for JSON parsing
    - Write tests for YAML parsing
    - Write tests for error handling with malformed files
    - _Requirements: 3.1, 3.2, 3.3, 3.4_

- [x] 3. Implement field mapping and validation service
  - [x] 3.1 Create field mapper service
    - Implement `src/services/importMapper.ts` with base mapper interface
    - Create LabelFieldMapper with field mapping logic for labels
    - Create ProviderFieldMapper with field mapping logic for providers
    - Create SubscriptionFieldMapper with field mapping logic for subscriptions
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_

  - [x] 3.2 Implement validation logic
    - Add validation rules for required fields in each mapper
    - Implement format validation (hex colors, URLs, dates, enums)
    - Add validation error message generation
    - Implement default value assignment for optional fields
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_

  - [x] 3.3 Write unit tests for field mapping and validation
    - Create test file `src/services/__tests__/importMapper.test.ts`
    - Write tests for field mapping for each entity type
    - Write tests for validation rules (required fields, formats, enums)
    - Write tests for default value assignment
    - Write tests for error message generation
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_

- [x] 4. Create FileUploadZone component
  - [x] 4.1 Implement FileUploadZone component
    - Create `src/components/import/FileUploadZone.tsx` component
    - Implement drag-and-drop functionality with visual feedback
    - Implement click-to-upload file picker
    - Add file format validation (CSV, JSON, YAML)
    - Implement visual states (idle, drag-over, loading, error)
    - Display selected file name and size
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

  - [x] 4.2 Style FileUploadZone component
    - Add CSS styling for drag-and-drop zone
    - Implement hover and active states
    - Add accessibility attributes (ARIA labels, keyboard support)
    - Ensure responsive design
    - _Requirements: 2.1, 2.2, 2.3, 2.4_

- [x] 5. Create ImportPreviewTable component
  - [x] 5.1 Implement ImportPreviewTable component
    - Create `src/components/import/ImportPreviewTable.tsx` component
    - Implement table display with entity-specific columns
    - Add checkbox column for individual record selection
    - Add "Select All" checkbox in table header
    - Implement status column showing import progress/result
    - Display validation errors inline with tooltips
    - Add color-coded row states (valid, invalid, importing, success, error)
    - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7_

  - [x] 5.2 Add action buttons and progress display
    - Implement "Import Selected" button
    - Implement "Import All" button
    - Add progress indicator for bulk import
    - Display import summary (total, completed, failed counts)
    - _Requirements: 4.7, 6.8, 6.9_

  - [x] 5.3 Implement table virtualization for performance
    - Add virtualization for large record sets using react-virtual or similar
    - Implement pagination as alternative to virtualization
    - Optimize rendering performance for 1000+ records
    - _Requirements: 4.1, 4.2_

- [x] 6. Create import manager hook
  - [x] 6.1 Implement useImportManager hook
    - Create `src/hooks/import/useImportManager.ts` hook
    - Implement sequential API call logic for importing records
    - Add status tracking for each record (pending, importing, success, error)
    - Implement progress tracking (total, completed, failed)
    - Add cancellation mechanism for bulk imports
    - Handle API errors and continue with remaining records
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8_

  - [x] 6.2 Add rate limiting and throttling
    - Implement delay between API calls (100-200ms)
    - Add exponential backoff on API errors
    - Ensure sequential processing (no concurrent requests)
    - _Requirements: 6.1, 6.2, 6.3_

  - [x] 6.3 Write unit tests for import manager
    - Create test file `src/hooks/import/__tests__/useImportManager.test.ts`
    - Write tests for sequential import logic
    - Write tests for error handling and continuation
    - Write tests for cancellation mechanism
    - Write tests for progress tracking
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8_

- [x] 7. Create import page for Labels
  - [x] 7.1 Implement ImportLabelsPage component
    - Create `src/pages/ImportLabelsPage.tsx` page component
    - Integrate FileUploadZone component
    - Implement file parsing using fileParser service
    - Integrate field mapping and validation using LabelFieldMapper
    - Integrate ImportPreviewTable component with label-specific columns
    - Integrate useImportManager hook with label creation mutation
    - Add navigation back to labels page
    - Implement unsaved changes warning
    - _Requirements: 1.4, 2.7, 3.7, 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8, 8.1, 8.2, 8.3, 8.4, 8.5, 8.6_

  - [x] 7.2 Add route for ImportLabelsPage
    - Update `src/App.tsx` to add route `/labels/import`
    - Wrap route with ProtectedRoute and AppLayout
    - _Requirements: 1.4, 8.1_

- [x] 8. Create import page for Providers
  - [x] 8.1 Implement ImportProvidersPage component
    - Create `src/pages/ImportProvidersPage.tsx` page component
    - Integrate FileUploadZone component
    - Implement file parsing using fileParser service
    - Integrate field mapping and validation using ProviderFieldMapper
    - Integrate ImportPreviewTable component with provider-specific columns
    - Integrate useImportManager hook with provider creation mutation
    - Add navigation back to providers page
    - Implement unsaved changes warning
    - _Requirements: 1.4, 2.7, 3.7, 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8, 8.1, 8.2, 8.3, 8.4, 8.5, 8.6_

  - [x] 8.2 Add route for ImportProvidersPage
    - Update `src/App.tsx` to add route `/providers/import`
    - Wrap route with ProtectedRoute and AppLayout
    - _Requirements: 1.4, 8.1_

- [x] 9. Create import page for Subscriptions
  - [x] 9.1 Implement ImportSubscriptionsPage component
    - Create `src/pages/ImportSubscriptionsPage.tsx` page component
    - Integrate FileUploadZone component
    - Implement file parsing using fileParser service
    - Integrate field mapping and validation using SubscriptionFieldMapper
    - Integrate ImportPreviewTable component with subscription-specific columns
    - Integrate useImportManager hook with subscription creation mutation
    - Add navigation back to subscriptions page
    - Implement unsaved changes warning
    - _Requirements: 1.4, 2.7, 3.7, 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8, 8.1, 8.2, 8.3, 8.4, 8.5, 8.6_

  - [x] 9.2 Add route for ImportSubscriptionsPage
    - Update `src/App.tsx` to add route `/subscriptions/import`
    - Wrap route with ProtectedRoute and AppLayout
    - _Requirements: 1.4, 8.1_

- [ ] 10. Add import menu button to entity pages
  - [x] 10.1 Add import button to LabelsPage
    - Update `src/pages/LabelsPage.tsx` to add import menu button in PageHeader
    - Implement navigation to `/labels/import` on button click
    - Add appropriate icon and label
    - _Requirements: 1.1, 1.4_

  - [x] 10.2 Add import button to ProvidersPage
    - Update `src/pages/ProvidersPage.tsx` to add import menu button in PageHeader
    - Implement navigation to `/providers/import` on button click
    - Add appropriate icon and label
    - _Requirements: 1.2, 1.4_

  - [x] 10.3 Add import button to SubscriptionsPage
    - Update `src/pages/SubscriptionsPage.tsx` to add import menu button in PageHeader
    - Implement navigation to `/subscriptions/import` on button click
    - Add appropriate icon and label
    - _Requirements: 1.3, 1.4_

- [x] 11. Create helper components for import UI
  - [x] 11.1 Create ImportProgress component
    - Create `src/components/import/ImportProgress.tsx` component
    - Display progress bar with percentage
    - Show counts (total, completed, failed)
    - Display estimated time remaining
    - _Requirements: 6.4, 6.8_

  - [x] 11.2 Create ImportSummary component
    - Create `src/components/import/ImportSummary.tsx` component
    - Display final import results
    - Show success and failure counts
    - List failed records with error messages
    - Provide option to retry failed imports or return to entity list
    - _Requirements: 6.8, 8.6_

- [x] 12. Implement error handling and user feedback
  - [x] 12.1 Add parse error handling
    - Display clear error messages for file format issues
    - Show line numbers for CSV/YAML parse errors
    - Allow user to select different file after error
    - _Requirements: 3.4, 3.5_

  - [x] 12.2 Add validation error display
    - Highlight invalid rows in preview table
    - Show specific field errors in tooltips or expandable rows
    - Prevent import of invalid records
    - Display validation error summary
    - _Requirements: 7.4, 7.5, 7.6_

  - [x] 12.3 Add API error handling
    - Display error messages next to failed records
    - Allow retry for individual failed records
    - Show network error messages
    - Implement retry mechanism with exponential backoff
    - _Requirements: 5.4, 5.5, 6.6, 6.7_

- [x] 13. Add accessibility features
  - [x] 13.1 Implement keyboard navigation
    - Add keyboard support for file upload zone (Enter/Space to open picker)
    - Implement table navigation with arrow keys
    - Add keyboard support for checkbox selection (Space bar)
    - Ensure all interactive elements are keyboard accessible
    - _Requirements: 2.2, 4.3, 4.5_

  - [x] 13.2 Add screen reader support
    - Add ARIA labels for all interactive elements
    - Implement status announcements for import progress
    - Add ARIA live regions for dynamic content updates
    - Ensure error messages are announced to screen readers
    - _Requirements: 2.1, 2.2, 4.1, 4.2, 6.4, 6.8_

  - [x] 13.3 Ensure visual accessibility
    - Use high contrast colors for validation errors
    - Add clear focus indicators for all interactive elements
    - Ensure loading states have both visual and text indicators
    - Test color contrast ratios meet WCAG standards
    - _Requirements: 4.2, 4.5, 4.6, 4.7_

- [x] 14. Write integration tests
  - [x] 14.1 Test file upload flow
    - Write test for drag-and-drop functionality
    - Write test for file picker functionality
    - Write test for format detection and validation
    - Write test for file size limit handling
    - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5, 2.6_

  - [x] 14.2 Test import flow
    - Write test for single record import
    - Write test for bulk import
    - Write test for import cancellation
    - Write test for error recovery and retry
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 5.5, 5.6, 6.1, 6.2, 6.3, 6.4, 6.5, 6.6, 6.7, 6.8, 6.9_

  - [x] 14.3 Test API integration
    - Write test for successful API calls
    - Write test for API error handling
    - Write test for rate limiting behavior
    - Write test for sequential request processing
    - _Requirements: 5.1, 5.2, 5.3, 5.4, 6.1, 6.2, 6.3_

- [ ] 15. Write end-to-end tests
  - [x] 15.1 Test complete import workflow for Labels
    - Write E2E test: navigate to labels import page
    - Upload valid CSV file with labels
    - Select records in preview table
    - Import selected records
    - Verify records appear in labels list
    - _Requirements: 1.1, 1.4, 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 3.1, 3.7, 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 5.1, 5.2, 5.3, 5.6, 6.8, 7.1, 8.1, 8.2, 8.6_

  - [x] 15.2 Test complete import workflow for Providers
    - Write E2E test: navigate to providers import page
    - Upload valid JSON file with providers
    - Select all records
    - Import all records
    - Verify records appear in providers list
    - _Requirements: 1.2, 1.4, 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 3.2, 3.7, 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 6.1, 6.2, 6.3, 6.4, 6.5, 6.8, 7.2, 8.1, 8.2, 8.6_

  - [x] 15.3 Test complete import workflow for Subscriptions
    - Write E2E test: navigate to subscriptions import page
    - Upload valid YAML file with subscriptions
    - Select records in preview table
    - Import selected records
    - Verify records appear in subscriptions list
    - _Requirements: 1.3, 1.4, 2.1, 2.2, 2.3, 2.4, 2.5, 2.6, 2.7, 3.3, 3.7, 4.1, 4.2, 4.3, 4.4, 4.5, 4.6, 4.7, 5.1, 5.2, 5.3, 5.6, 6.8, 7.3, 8.1, 8.2, 8.6_

  - [x] 15.4 Test error scenarios
    - Write E2E test: upload invalid file format
    - Write E2E test: upload file with validation errors
    - Write E2E test: handle API failures during import
    - Write E2E test: cancel import mid-process
    - _Requirements: 2.6, 3.4, 5.4, 5.5, 6.6, 6.7, 6.9, 7.4, 7.5_

- [x] 16. Add documentation and examples
  - [x] 16.1 Create file format templates
    - Create example CSV template for labels at `public/templates/labels-template.csv`
    - Create example JSON template for providers at `public/templates/providers-template.json`
    - Create example YAML template for subscriptions at `public/templates/subscriptions-template.yaml`
    - _Requirements: 7.1, 7.2, 7.3_

  - [x] 16.2 Add help documentation
    - Create help text component explaining import process
    - Add field mapping documentation for each entity type
    - Provide validation error explanations
    - Add download links for template files
    - _Requirements: 7.1, 7.2, 7.3, 7.4, 7.5, 7.6_

- [x] 17. Performance optimization
  - [x] 17.1 Optimize file parsing
    - Implement file size limits (10MB)
    - Add streaming for large CSV files
    - Show progress indicator during parsing
    - Optimize memory usage by releasing file references
    - _Requirements: 2.7, 3.1, 3.7, 8.4_

  - [x] 17.2 Optimize preview table rendering
    - Implement virtualization for tables with 100+ records
    - Optimize re-renders using React.memo and useMemo
    - Lazy load validation for large datasets
    - _Requirements: 4.1, 4.2, 4.7_

  - [x] 17.3 Optimize API calls
    - Implement configurable delay between API calls
    - Add request queuing mechanism
    - Implement exponential backoff for retries
    - Monitor and respect API rate limits
    - _Requirements: 6.1, 6.2, 6.3_
