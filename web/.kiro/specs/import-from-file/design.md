# Design Document: Import from File Feature

## Overview

This feature enables users to import Labels, Providers, and Subscriptions from external files (CSV, JSON, or YAML) through a client-side parsing and validation workflow. The design emphasizes user control, data security, and backend protection by parsing files entirely in the browser and sending only validated individual records to the API.

The import workflow consists of:
1. Accessing the import feature from entity-specific pages
2. Selecting or drag-dropping a file
3. Client-side parsing and validation
4. Preview of parsed records with selection controls
5. Individual or bulk import with progress tracking

## Architecture

### High-Level Component Structure

```
┌─────────────────────────────────────────────────────────────┐
│                    Entity Pages                              │
│  (LabelsPage, ProvidersPage, SubscriptionsPage)            │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  PageHeader with Import Menu Button                   │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ Navigate to
                            ▼
┌─────────────────────────────────────────────────────────────┐
│              Import Page (Entity-Specific)                   │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  FileUploadZone Component                             │  │
│  │  - Drag & Drop                                         │  │
│  │  - File Picker                                         │  │
│  │  - Format Validation                                   │  │
│  └──────────────────────────────────────────────────────┘  │
│                            │                                 │
│                            │ File Selected                   │
│                            ▼                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  File Parser Service                                   │  │
│  │  - CSV Parser                                          │  │
│  │  - JSON Parser                                         │  │
│  │  - YAML Parser                                         │  │
│  └──────────────────────────────────────────────────────┘  │
│                            │                                 │
│                            │ Parsed Records                  │
│                            ▼                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Field Mapper & Validator                             │  │
│  │  - Entity-specific mapping                            │  │
│  │  - Validation rules                                    │  │
│  └──────────────────────────────────────────────────────┘  │
│                            │                                 │
│                            │ Validated Records               │
│                            ▼                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  ImportPreviewTable Component                         │  │
│  │  - Record display                                      │  │
│  │  - Selection controls                                  │  │
│  │  - Validation error display                            │  │
│  │  - Import actions                                      │  │
│  └──────────────────────────────────────────────────────┘  │
│                            │                                 │
│                            │ Import Action                   │
│                            ▼                                 │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Import Manager                                        │  │
│  │  - Sequential API calls                                │  │
│  │  - Progress tracking                                   │  │
│  │  - Error handling                                      │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
                            │
                            │ API Calls (one per record)
                            ▼
┌─────────────────────────────────────────────────────────────┐
│                    Backend API                               │
│  - POST /labels                                              │
│  - POST /providers                                           │
│  - POST /subscriptions                                       │
└─────────────────────────────────────────────────────────────┘
```

### Technology Stack

- **File Parsing Libraries:**
  - `papaparse` - CSV parsing
  - Native `JSON.parse()` - JSON parsing
  - `js-yaml` - YAML parsing

- **UI Components:**
  - Existing shadcn/ui components (Button, Table, Dialog, etc.)
  - Custom FileUploadZone component
  - Custom ImportPreviewTable component

- **State Management:**
  - React hooks (useState, useReducer)
  - React Query for API mutations

- **Routing:**
  - React Router for navigation

## Components and Interfaces

### 1. Import Menu Button (Enhancement to Existing Pages)

**Location:** Added to PageHeader component on Labels, Providers, and Subscriptions pages

**Interface:**
```typescript
interface ImportMenuButtonProps {
  entityType: 'labels' | 'providers' | 'subscriptions';
  onImportClick: () => void;
}
```

**Behavior:**
- Renders a dropdown menu or button next to existing action buttons
- On click, navigates to the entity-specific import page

### 2. Import Page Component

**Files:**
- `src/pages/ImportLabelsPage.tsx`
- `src/pages/ImportProvidersPage.tsx`
- `src/pages/ImportSubscriptionsPage.tsx`

**Interface:**
```typescript
interface ImportPageProps {
  entityType: 'labels' | 'providers' | 'subscriptions';
}

interface ImportPageState {
  file: File | null;
  parsedRecords: ParsedRecord[];
  validationErrors: Map<number, ValidationError[]>;
  selectedRecords: Set<number>;
  importStatus: Map<number, ImportStatus>;
  isImporting: boolean;
  importProgress: ImportProgress;
}

interface ImportProgress {
  total: number;
  completed: number;
  failed: number;
  inProgress: boolean;
}

interface ImportStatus {
  status: 'pending' | 'importing' | 'success' | 'error';
  error?: string;
}
```

### 3. FileUploadZone Component

**File:** `src/components/import/FileUploadZone.tsx`

**Interface:**
```typescript
interface FileUploadZoneProps {
  onFileSelected: (file: File) => void;
  acceptedFormats: string[]; // ['.csv', '.json', '.yaml', '.yml']
  isLoading?: boolean;
}
```

**Features:**
- Drag-and-drop area with visual feedback
- Click to open file picker
- File format validation
- Visual states: idle, drag-over, loading, error
- Display selected file name and size

### 4. File Parser Service

**File:** `src/services/fileParser.ts`

**Interface:**
```typescript
interface FileParserService {
  parseCSV(file: File): Promise<Record<string, any>[]>;
  parseJSON(file: File): Promise<Record<string, any>[]>;
  parseYAML(file: File): Promise<Record<string, any>[]>;
  detectFormat(file: File): 'csv' | 'json' | 'yaml' | 'unknown';
}

interface ParseResult {
  records: Record<string, any>[];
  errors: ParseError[];
}

interface ParseError {
  line?: number;
  message: string;
}
```

**Implementation Details:**
- Uses FileReader API to read file contents
- Delegates to format-specific parsers
- Returns normalized array of objects
- Handles parsing errors gracefully

### 5. Field Mapper & Validator

**File:** `src/services/importMapper.ts`

**Interface:**
```typescript
interface FieldMapper<T> {
  mapFields(rawRecord: Record<string, any>): Partial<T>;
  validate(record: Partial<T>): ValidationResult;
}

interface ValidationResult {
  isValid: boolean;
  errors: ValidationError[];
}

interface ValidationError {
  field: string;
  message: string;
  severity: 'error' | 'warning';
}

// Entity-specific mappers
interface LabelFieldMapper extends FieldMapper<DtoCreateLabelRequest> {}
interface ProviderFieldMapper extends FieldMapper<DtoCreateProviderRequest> {}
interface SubscriptionFieldMapper extends FieldMapper<DtoCreateSubscriptionRequest> {}
```

**Field Mapping Rules:**

**Labels:**
- Required: `name`, `color`
- Optional: `owner` (defaults to personal)
- Color validation: hex format (#RRGGBB or #AARRGGBB)

**Providers:**
- Required: `name`
- Optional: `description`, `url`, `iconUrl`, `pricingPageUrl`, `labels`
- URL validation for url fields

**Subscriptions:**
- Required: `providerId`, `startDate`, `recurrency`, `owner`
- Optional: `friendlyName`, `customPrice`, `customRecurrency`, `endDate`, `labels`, `payer`, `freeTrial`, `familyUsers`
- Date validation for date fields
- Recurrency enum validation
- Provider ID must exist (can be validated against existing providers)

### 6. ImportPreviewTable Component

**File:** `src/components/import/ImportPreviewTable.tsx`

**Interface:**
```typescript
interface ImportPreviewTableProps<T> {
  records: ParsedImportRecord<T>[];
  columns: ImportColumnDef<T>[];
  selectedRecords: Set<number>;
  onSelectionChange: (indices: Set<number>) => void;
  onImportSelected: () => void;
  onImportAll: () => void;
  importStatus: Map<number, ImportStatus>;
  isImporting: boolean;
}

interface ParsedImportRecord<T> {
  index: number;
  data: Partial<T>;
  validationErrors: ValidationError[];
  isValid: boolean;
}

interface ImportColumnDef<T> {
  key: keyof T;
  label: string;
  render?: (value: any, record: ParsedImportRecord<T>) => React.ReactNode;
}
```

**Features:**
- Table display with entity-specific columns
- Checkbox column for selection
- Status column showing import progress/result
- Validation error indicators
- "Select All" checkbox in header
- Action buttons: "Import Selected", "Import All"
- Progress indicator during bulk import
- Color-coded rows: valid (default), invalid (red), importing (yellow), success (green), error (red)

### 7. Import Manager Hook

**File:** `src/hooks/import/useImportManager.ts`

**Interface:**
```typescript
interface UseImportManagerProps<T> {
  records: ParsedImportRecord<T>[];
  createMutation: UseMutationResult<any, any, T, any>;
}

interface UseImportManagerReturn {
  importRecords: (indices: number[]) => Promise<void>;
  importStatus: Map<number, ImportStatus>;
  progress: ImportProgress;
  isImporting: boolean;
  cancelImport: () => void;
}
```

**Behavior:**
- Accepts array of record indices to import
- Calls API sequentially for each record
- Updates status map after each API call
- Continues on individual failures
- Provides cancellation mechanism
- Returns summary of results

## Data Models

### Import Record Types

```typescript
// Generic import record wrapper
interface ImportRecord<T> {
  index: number;
  rawData: Record<string, any>;
  mappedData: Partial<T>;
  validationErrors: ValidationError[];
  isValid: boolean;
  importStatus: ImportStatus;
}

// Entity-specific types
type LabelImportRecord = ImportRecord<DtoCreateLabelRequest>;
type ProviderImportRecord = ImportRecord<DtoCreateProviderRequest>;
type SubscriptionImportRecord = ImportRecord<DtoCreateSubscriptionRequest>;
```

### File Format Examples

**CSV Format (Labels):**
```csv
name,color
Entertainment,#FF5733
Utilities,#33FF57
```

**JSON Format (Providers):**
```json
[
  {
    "name": "Netflix",
    "description": "Streaming service",
    "url": "https://netflix.com",
    "labels": ["entertainment"]
  },
  {
    "name": "AWS",
    "description": "Cloud services",
    "url": "https://aws.amazon.com"
  }
]
```

**YAML Format (Subscriptions):**
```yaml
- providerId: "provider-uuid-here"
  friendlyName: "Netflix Premium"
  startDate: "2024-01-01"
  recurrency: "monthly"
  customPrice:
    amount: 15.99
    currency: "USD"
  owner:
    type: "personal"
```

## Error Handling

### Parse Errors
- Display error message with file format issues
- Allow user to select a different file
- Show line numbers for CSV/YAML errors

### Validation Errors
- Display inline in preview table
- Highlight invalid rows
- Show specific field errors in tooltip or expandable row
- Prevent import of invalid records
- Allow user to fix data in external file and re-upload

### API Errors
- Display error message next to failed record
- Allow retry for individual failed records
- Continue importing remaining records
- Show summary of successes and failures

### Network Errors
- Detect network failures
- Provide retry mechanism
- Preserve import state for recovery

## Testing Strategy

### Unit Tests

1. **File Parser Tests**
   - Test CSV parsing with various formats
   - Test JSON parsing with arrays and objects
   - Test YAML parsing with nested structures
   - Test error handling for malformed files

2. **Field Mapper Tests**
   - Test field mapping for each entity type
   - Test validation rules
   - Test default value assignment
   - Test error message generation

3. **Validation Tests**
   - Test required field validation
   - Test format validation (colors, URLs, dates)
   - Test enum validation
   - Test cross-field validation

### Integration Tests

1. **File Upload Flow**
   - Test drag-and-drop functionality
   - Test file picker functionality
   - Test format detection
   - Test file size limits

2. **Import Flow**
   - Test single record import
   - Test bulk import
   - Test import cancellation
   - Test error recovery

3. **API Integration**
   - Test successful API calls
   - Test API error handling
   - Test rate limiting behavior
   - Test concurrent request handling

### E2E Tests

1. **Complete Import Workflow**
   - Navigate to import page
   - Upload valid file
   - Select records
   - Import successfully
   - Verify records appear in entity list

2. **Error Scenarios**
   - Upload invalid file format
   - Upload file with validation errors
   - Handle API failures
   - Cancel import mid-process

## Performance Considerations

### Client-Side Performance

1. **Large File Handling**
   - Implement file size limits (e.g., 10MB)
   - Use streaming for large CSV files
   - Implement pagination for preview table (virtualization)
   - Show progress indicator during parsing

2. **Memory Management**
   - Release file references after parsing
   - Limit number of records displayed at once
   - Clear import state after completion

### Backend Protection

1. **Rate Limiting**
   - Add delay between API calls (e.g., 100-200ms)
   - Implement exponential backoff on errors
   - Respect API rate limits

2. **Request Throttling**
   - Limit concurrent requests to 1 (sequential)
   - Queue requests for processing
   - Provide user feedback on throttling

## Security Considerations

1. **File Validation**
   - Validate file extensions
   - Validate MIME types
   - Scan for malicious content patterns
   - Limit file sizes

2. **Data Sanitization**
   - Sanitize all parsed data
   - Validate against schema
   - Prevent XSS in preview display
   - Escape special characters

3. **API Security**
   - Use existing authentication
   - Validate ownership before import
   - Respect quota limits
   - Log import activities

## User Experience Enhancements

1. **Progress Feedback**
   - Show parsing progress
   - Show import progress with percentage
   - Display estimated time remaining
   - Show success/failure counts in real-time

2. **Helpful Guidance**
   - Provide file format templates for download
   - Show example files
   - Display field mapping documentation
   - Offer validation error explanations

3. **Undo/Rollback**
   - Consider implementing batch deletion for failed imports
   - Provide export functionality to backup before import
   - Show import history

## Routing Structure

New routes to be added to `src/App.tsx`:

```typescript
<Route
  path="/labels/import"
  element={
    <ProtectedRoute>
      <AppLayout>
        <ImportLabelsPage />
      </AppLayout>
    </ProtectedRoute>
  }
/>
<Route
  path="/providers/import"
  element={
    <ProtectedRoute>
      <AppLayout>
        <ImportProvidersPage />
      </AppLayout>
    </ProtectedRoute>
  }
/>
<Route
  path="/subscriptions/import"
  element={
    <ProtectedRoute>
      <AppLayout>
        <ImportSubscriptionsPage />
      </AppLayout>
    </ProtectedRoute>
  }
/>
```

## Dependencies

New packages to install:

```json
{
  "papaparse": "^5.4.1",
  "js-yaml": "^4.1.0",
  "@types/papaparse": "^5.3.14",
  "@types/js-yaml": "^4.0.9"
}
```

## File Structure

```
src/
├── components/
│   └── import/
│       ├── FileUploadZone.tsx
│       ├── ImportPreviewTable.tsx
│       ├── ImportProgress.tsx
│       └── ImportSummary.tsx
├── hooks/
│   └── import/
│       ├── useImportManager.ts
│       ├── useFileParser.ts
│       └── useImportValidation.ts
├── pages/
│   ├── ImportLabelsPage.tsx
│   ├── ImportProvidersPage.tsx
│   └── ImportSubscriptionsPage.tsx
├── services/
│   ├── fileParser.ts
│   └── importMapper.ts
└── types/
    └── import.ts
```

## Accessibility

1. **Keyboard Navigation**
   - File upload zone accessible via keyboard
   - Table navigation with arrow keys
   - Checkbox selection with space bar

2. **Screen Reader Support**
   - ARIA labels for all interactive elements
   - Status announcements for import progress
   - Error messages announced

3. **Visual Indicators**
   - High contrast for validation errors
   - Clear focus indicators
   - Loading states with spinners and text

## Future Enhancements

1. **Advanced Mapping**
   - Custom field mapping UI
   - Save mapping templates
   - Auto-detect field mappings

2. **Batch Operations**
   - Bulk edit before import
   - Transform data during import
   - Merge with existing records

3. **Import History**
   - Track all imports
   - Rollback capability
   - Export import logs

4. **Format Support**
   - Excel files (.xlsx)
   - XML format
   - API import from URLs
