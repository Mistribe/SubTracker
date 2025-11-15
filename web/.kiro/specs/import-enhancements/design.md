# Design Document: Import Feature Enhancements

## Overview

This design document outlines enhancements to the existing import functionality, focusing on two key improvements:

1. **Complete Template Coverage**: Providing template files in all three supported formats (CSV, JSON, YAML) for each entity type (Labels, Providers, Subscriptions)
2. **UUID Import Support**: Enabling users to specify custom UUIDs during import to maintain identifier consistency across systems

These enhancements build upon the existing import infrastructure without requiring major architectural changes, primarily involving template file creation and validation logic updates.

## Architecture

### High-Level Enhancement Areas

```
┌─────────────────────────────────────────────────────────────┐
│              Import Page (Existing)                          │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Template Download Section (ENHANCED)                 │  │
│  │  - CSV template link                                  │  │
│  │  - JSON template link                                 │  │
│  │  - YAML template link                                 │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  Field Mapper & Validator (ENHANCED)                  │  │
│  │  - UUID validation added                              │  │
│  │  - ID field mapping added                             │  │
│  └──────────────────────────────────────────────────────┘  │
│                                                              │
│  ┌──────────────────────────────────────────────────────┐  │
│  │  ImportPreviewTable (ENHANCED)                        │  │
│  │  - Display ID column when present                     │  │
│  │  - Show UUID validation errors                        │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Template File Structure

```
public/templates/
├── labels/
│   ├── labels-template.csv
│   ├── labels-template.json
│   └── labels-template.yaml
├── providers/
│   ├── providers-template.csv
│   ├── providers-template.json
│   └── providers-template.yaml
└── subscriptions/
    ├── subscriptions-template.csv
    ├── subscriptions-template.json
    └── subscriptions-template.yaml
```

## Components and Interfaces

### 1. Template Files

#### Labels Templates

**CSV Format** (`public/templates/labels/labels-template.csv`):
```csv
id,name,color,ownerType,ownerFamilyId
550e8400-e29b-41d4-a716-446655440001,Entertainment,#FF5733,personal,
550e8400-e29b-41d4-a716-446655440002,Utilities,#33FF57,personal,
550e8400-e29b-41d4-a716-446655440003,Productivity,#3357FF,personal,
,Health & Fitness,#FF33A1,personal,
,Education,#A133FF,personal,
```

**JSON Format** (`public/templates/labels/labels-template.json`):
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Entertainment",
    "color": "#FF5733",
    "ownerType": "personal"
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "Utilities",
    "color": "#33FF57",
    "ownerType": "personal"
  },
  {
    "name": "Productivity",
    "color": "#3357FF",
    "ownerType": "personal"
  }
]
```

**YAML Format** (`public/templates/labels/labels-template.yaml`):
```yaml
- id: "550e8400-e29b-41d4-a716-446655440001"
  name: "Entertainment"
  color: "#FF5733"
  ownerType: "personal"

- id: "550e8400-e29b-41d4-a716-446655440002"
  name: "Utilities"
  color: "#33FF57"
  ownerType: "personal"

- name: "Productivity"
  color: "#3357FF"
  ownerType: "personal"
```

#### Providers Templates

**CSV Format** (`public/templates/providers/providers-template.csv`):
```csv
id,name,description,url,iconUrl,pricingPageUrl,labels
550e8400-e29b-41d4-a716-446655440010,Netflix,Streaming service for movies and TV shows,https://www.netflix.com,https://www.netflix.com/favicon.ico,https://www.netflix.com/signup/planform,"Entertainment,Video Streaming"
550e8400-e29b-41d4-a716-446655440011,Spotify,Music streaming service,https://www.spotify.com,https://www.spotify.com/favicon.ico,https://www.spotify.com/premium,"Entertainment,Music"
,GitHub,Code hosting platform,https://github.com,https://github.com/favicon.ico,https://github.com/pricing,"Development,Productivity"
```

**JSON Format** (`public/templates/providers/providers-template.json`):
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "name": "Netflix",
    "description": "Streaming service for movies and TV shows",
    "url": "https://www.netflix.com",
    "iconUrl": "https://www.netflix.com/favicon.ico",
    "pricingPageUrl": "https://www.netflix.com/signup/planform",
    "labels": ["Entertainment", "Video Streaming"]
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440011",
    "name": "Spotify",
    "description": "Music streaming service",
    "url": "https://www.spotify.com",
    "iconUrl": "https://www.spotify.com/favicon.ico",
    "pricingPageUrl": "https://www.spotify.com/premium",
    "labels": ["Entertainment", "Music"]
  },
  {
    "name": "GitHub",
    "description": "Code hosting platform",
    "url": "https://github.com",
    "iconUrl": "https://github.com/favicon.ico",
    "pricingPageUrl": "https://github.com/pricing",
    "labels": ["Development", "Productivity"]
  }
]
```

**YAML Format** (`public/templates/providers/providers-template.yaml`):
```yaml
- id: "550e8400-e29b-41d4-a716-446655440010"
  name: "Netflix"
  description: "Streaming service for movies and TV shows"
  url: "https://www.netflix.com"
  iconUrl: "https://www.netflix.com/favicon.ico"
  pricingPageUrl: "https://www.netflix.com/signup/planform"
  labels:
    - "Entertainment"
    - "Video Streaming"

- id: "550e8400-e29b-41d4-a716-446655440011"
  name: "Spotify"
  description: "Music streaming service"
  url: "https://www.spotify.com"
  iconUrl: "https://www.spotify.com/favicon.ico"
  pricingPageUrl: "https://www.spotify.com/premium"
  labels:
    - "Entertainment"
    - "Music"

- name: "GitHub"
  description: "Code hosting platform"
  url: "https://github.com"
  iconUrl: "https://github.com/favicon.ico"
  pricingPageUrl: "https://github.com/pricing"
  labels:
    - "Development"
    - "Productivity"
```

#### Subscriptions Templates

**CSV Format** (`public/templates/subscriptions/subscriptions-template.csv`):
```csv
id,providerId,friendlyName,startDate,endDate,recurrency,customRecurrency,customPriceAmount,customPriceCurrency,ownerType,ownerFamilyId,payerType,payerFamilyId,freeTrialStartDate,freeTrialEndDate,labels
550e8400-e29b-41d4-a716-446655440020,your-provider-id-here,Netflix Premium,2024-01-01,,monthly,,15.99,USD,personal,,,,,,"Entertainment,Video Streaming"
550e8400-e29b-41d4-a716-446655440021,your-provider-id-here,Spotify Family,2024-02-15,,monthly,,16.99,USD,personal,,,,,,"Entertainment,Music"
,your-provider-id-here,GitHub Pro,2024-03-01,,monthly,,4.00,USD,personal,,,,2024-03-01,2024-03-31,"Development,Productivity"
```

**JSON Format** (`public/templates/subscriptions/subscriptions-template.json`):
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440020",
    "providerId": "your-provider-id-here",
    "friendlyName": "Netflix Premium",
    "startDate": "2024-01-01",
    "recurrency": "monthly",
    "customPriceAmount": 15.99,
    "customPriceCurrency": "USD",
    "ownerType": "personal",
    "labels": ["Entertainment", "Video Streaming"]
  },
  {
    "id": "550e8400-e29b-41d4-a716-446655440021",
    "providerId": "your-provider-id-here",
    "friendlyName": "Spotify Family",
    "startDate": "2024-02-15",
    "recurrency": "monthly",
    "customPriceAmount": 16.99,
    "customPriceCurrency": "USD",
    "ownerType": "personal",
    "labels": ["Entertainment", "Music"]
  },
  {
    "providerId": "your-provider-id-here",
    "friendlyName": "GitHub Pro",
    "startDate": "2024-03-01",
    "recurrency": "monthly",
    "customPriceAmount": 4.00,
    "customPriceCurrency": "USD",
    "ownerType": "personal",
    "freeTrialStartDate": "2024-03-01",
    "freeTrialEndDate": "2024-03-31",
    "labels": ["Development", "Productivity"]
  }
]
```

**YAML Format** (`public/templates/subscriptions/subscriptions-template.yaml`):
```yaml
- id: "550e8400-e29b-41d4-a716-446655440020"
  providerId: "your-provider-id-here"
  friendlyName: "Netflix Premium"
  startDate: "2024-01-01"
  recurrency: "monthly"
  customPriceAmount: 15.99
  customPriceCurrency: "USD"
  ownerType: "personal"
  labels:
    - "Entertainment"
    - "Video Streaming"

- id: "550e8400-e29b-41d4-a716-446655440021"
  providerId: "your-provider-id-here"
  friendlyName: "Spotify Family"
  startDate: "2024-02-15"
  recurrency: "monthly"
  customPriceAmount: 16.99
  customPriceCurrency: "USD"
  ownerType: "personal"
  labels:
    - "Entertainment"
    - "Music"

- providerId: "your-provider-id-here"
  friendlyName: "GitHub Pro"
  startDate: "2024-03-01"
  recurrency: "monthly"
  customPriceAmount: 4.00
  customPriceCurrency: "USD"
  ownerType: "personal"
  freeTrialStartDate: "2024-03-01"
  freeTrialEndDate: "2024-03-31"
  labels:
    - "Development"
    - "Productivity"
```

### 2. Template Download Component

**File:** `src/components/import/TemplateDownloadSection.tsx`

**Interface:**
```typescript
interface TemplateDownloadSectionProps {
  entityType: 'labels' | 'providers' | 'subscriptions';
}

interface TemplateLink {
  format: 'csv' | 'json' | 'yaml';
  label: string;
  path: string;
  icon: React.ReactNode;
}
```

**Features:**
- Display three download buttons/links for each format
- Clear labeling with format icons
- Descriptive text explaining template purpose
- Responsive layout
- Accessible keyboard navigation

**Implementation:**
```typescript
const TemplateDownloadSection: React.FC<TemplateDownloadSectionProps> = ({ entityType }) => {
  const templates: TemplateLink[] = [
    {
      format: 'csv',
      label: 'CSV Template',
      path: `/templates/${entityType}/${entityType}-template.csv`,
      icon: <FileSpreadsheet />
    },
    {
      format: 'json',
      label: 'JSON Template',
      path: `/templates/${entityType}/${entityType}-template.json`,
      icon: <FileJson />
    },
    {
      format: 'yaml',
      label: 'YAML Template',
      path: `/templates/${entityType}/${entityType}-template.yaml`,
      icon: <FileCode />
    }
  ];

  return (
    <div className="template-download-section">
      <h3>Download Template Files</h3>
      <p>Choose a template format to get started with your import</p>
      <div className="template-links">
        {templates.map(template => (
          <a
            key={template.format}
            href={template.path}
            download
            className="template-link"
          >
            {template.icon}
            <span>{template.label}</span>
          </a>
        ))}
      </div>
    </div>
  );
};
```

### 3. Enhanced Field Mapper with UUID Support

**File:** `src/services/importMapper.ts` (Enhancement)

**UUID Validation Function:**
```typescript
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

function validateUUID(id: string): boolean {
  return UUID_REGEX.test(id);
}

function validateAndMapId(rawId: any): { id?: string; error?: ValidationError } {
  // If no ID provided, that's fine - backend will generate one
  if (rawId === undefined || rawId === null || rawId === '') {
    return {};
  }

  // Convert to string and trim
  const idString = String(rawId).trim();

  // Validate UUID format
  if (!validateUUID(idString)) {
    return {
      error: {
        field: 'id',
        message: `Invalid UUID format: "${idString}". Expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`,
        severity: 'error'
      }
    };
  }

  return { id: idString };
}
```

**Enhanced Mapper Interface:**
```typescript
interface FieldMapper<T> {
  mapFields(rawRecord: Record<string, any>): {
    data: Partial<T>;
    errors: ValidationError[];
  };
}
```

**Updated Label Mapper:**
```typescript
class LabelFieldMapper implements FieldMapper<DtoCreateLabelRequest> {
  mapFields(rawRecord: Record<string, any>): {
    data: Partial<DtoCreateLabelRequest>;
    errors: ValidationError[];
  } {
    const errors: ValidationError[] = [];
    const data: Partial<DtoCreateLabelRequest> = {};

    // Handle ID field
    const idResult = validateAndMapId(rawRecord.id);
    if (idResult.error) {
      errors.push(idResult.error);
    } else if (idResult.id) {
      data.id = idResult.id;
    }

    // ... existing field mapping logic ...

    return { data, errors };
  }
}
```

### 4. Enhanced Import Preview Table

**File:** `src/components/import/ImportPreviewTable.tsx` (Enhancement)

**Column Configuration Update:**
```typescript
// Add ID column when IDs are present in data
const columns: ImportColumnDef<T>[] = [
  // Show ID column if any record has an ID
  ...(records.some(r => r.data.id) ? [{
    key: 'id' as keyof T,
    label: 'ID',
    render: (value: any, record: ParsedImportRecord<T>) => (
      <div className="id-cell">
        <code className="text-xs">{value || 'Auto-generated'}</code>
      </div>
    )
  }] : []),
  // ... existing columns ...
];
```

**UUID Error Display:**
```typescript
// Enhanced error display for UUID validation errors
const renderValidationErrors = (errors: ValidationError[]) => {
  return (
    <div className="validation-errors">
      {errors.map((error, idx) => (
        <div key={idx} className={`error-item ${error.severity}`}>
          <AlertCircle className="error-icon" />
          <span className="error-field">{error.field}:</span>
          <span className="error-message">{error.message}</span>
        </div>
      ))}
    </div>
  );
};
```

### 5. API Error Handling for Duplicate UUIDs

**File:** `src/hooks/import/useImportManager.ts` (Enhancement)

**Enhanced Error Handling:**
```typescript
const importRecord = async (record: ParsedImportRecord<T>, index: number) => {
  try {
    await createMutation.mutateAsync(record.data);
    updateStatus(index, { status: 'success' });
  } catch (error: any) {
    let errorMessage = 'Import failed';

    // Check for UUID conflict errors
    if (error.response?.status === 409 || error.response?.status === 400) {
      if (error.response?.data?.message?.includes('already exists') ||
          error.response?.data?.message?.includes('duplicate')) {
        errorMessage = `Entity with ID ${record.data.id} already exists`;
      }
    }

    updateStatus(index, {
      status: 'error',
      error: errorMessage
    });
  }
};
```

## Data Models

### Enhanced Import Record Type

```typescript
interface ImportRecord<T> {
  index: number;
  rawData: Record<string, any>;
  mappedData: Partial<T> & { id?: string }; // ID is now explicitly optional
  validationErrors: ValidationError[];
  isValid: boolean;
  importStatus: ImportStatus;
}
```

### UUID Validation Error Type

```typescript
interface UUIDValidationError extends ValidationError {
  field: 'id';
  message: string;
  severity: 'error';
  providedValue?: string;
}
```

## Error Handling

### UUID Validation Errors

1. **Invalid Format**
   - Display clear error message with expected format
   - Highlight the invalid UUID in the preview table
   - Prevent import until corrected
   - Show example of valid UUID format

2. **Duplicate UUID (API Conflict)**
   - Catch 409 Conflict responses from API
   - Display specific error message indicating duplicate
   - Allow user to edit file and re-import
   - Continue with remaining records

3. **Empty/Null ID**
   - Treat as valid (backend will generate)
   - Display "Auto-generated" in preview table
   - No validation error

## Testing Strategy

### Unit Tests

1. **UUID Validation Tests**
   - Test valid UUID formats (v1, v4, v5)
   - Test invalid UUID formats
   - Test empty/null/undefined IDs
   - Test whitespace handling

2. **Field Mapper Tests**
   - Test ID field mapping for each entity type
   - Test validation error generation for invalid UUIDs
   - Test records with and without IDs
   - Test mixed scenarios (some with IDs, some without)

3. **Template File Tests**
   - Validate all template files parse correctly
   - Test that templates contain valid example data
   - Verify UUID examples in templates are valid

### Integration Tests

1. **Template Download Tests**
   - Test download links work for all formats
   - Verify downloaded files have correct content
   - Test file naming conventions

2. **Import with UUID Tests**
   - Test importing records with valid UUIDs
   - Test importing records without UUIDs
   - Test mixed import (some with, some without UUIDs)
   - Test duplicate UUID handling

### E2E Tests

1. **Complete Workflow with UUIDs**
   - Download template file
   - Modify with custom UUIDs
   - Upload and preview
   - Import successfully
   - Verify entities created with specified UUIDs

2. **Error Scenarios**
   - Upload file with invalid UUIDs
   - Verify validation errors displayed
   - Test duplicate UUID conflict handling

## User Experience Enhancements

### Template Download Section Design

```
┌─────────────────────────────────────────────────────────┐
│  Download Template Files                                 │
│  Choose a template format to get started                 │
│                                                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ 📊 CSV   │  │ {} JSON  │  │ 📝 YAML  │             │
│  │ Template │  │ Template │  │ Template │             │
│  └──────────┘  └──────────┘  └──────────┘             │
└─────────────────────────────────────────────────────────┘
```

### UUID Field Display

- Show ID column only when at least one record has an ID
- Display "Auto-generated" for records without IDs
- Use monospace font for UUID display
- Truncate long UUIDs with tooltip showing full value

### Validation Error Messages

- Clear, actionable error messages
- Show expected UUID format
- Highlight specific invalid characters or structure
- Provide link to UUID format documentation

## Performance Considerations

### Template File Size

- Keep template files small (< 50KB each)
- Include 3-5 example records per template
- Balance between comprehensive examples and file size

### UUID Validation Performance

- Use compiled regex for UUID validation
- Validate during field mapping phase (already in memory)
- No additional API calls for UUID validation

## Security Considerations

### UUID Handling

- Validate UUID format before sending to API
- Sanitize UUID strings (trim whitespace)
- Let backend handle duplicate UUID conflicts
- Don't expose internal UUID generation logic

### Template Files

- Serve templates as static files
- No dynamic template generation
- Use example UUIDs that don't exist in production
- Include disclaimer in templates about example data

## Accessibility

### Template Download Links

- Keyboard accessible
- Clear focus indicators
- ARIA labels for screen readers
- Descriptive link text

### UUID Display

- Monospace font for readability
- High contrast for UUID text
- Tooltip support for full UUID display
- Screen reader announces UUID values

## Migration and Backward Compatibility

### Existing Import Functionality

- All existing imports without IDs continue to work
- No breaking changes to existing field mappers
- ID field is purely additive (optional)
- Existing templates remain functional

### Gradual Rollout

1. Add new template files
2. Update field mappers to support ID field
3. Update UI to display ID column
4. Update documentation

## Documentation Updates

### User Documentation

- Add section on UUID import capability
- Explain when to use custom UUIDs
- Document UUID format requirements
- Provide examples of UUID import scenarios

### Template Documentation

- Add README in templates directory
- Explain each template format
- Document all available fields
- Include field descriptions and constraints

## File Organization

```
public/templates/
├── README.md (new)
├── labels/
│   ├── labels-template.csv (new)
│   ├── labels-template.json (new)
│   └── labels-template.yaml (new)
├── providers/
│   ├── providers-template.csv (new)
│   ├── providers-template.json (new)
│   └── providers-template.yaml (new)
└── subscriptions/
    ├── subscriptions-template.csv (new)
    ├── subscriptions-template.json (new)
    └── subscriptions-template.yaml (new)
```

## Future Enhancements

1. **Template Customization**
   - Allow users to customize which fields appear in templates
   - Save custom template preferences

2. **UUID Generation Helper**
   - Provide in-app UUID generator
   - Bulk UUID generation for multiple records

3. **Template Validation**
   - Pre-validate templates before download
   - Automated template testing in CI/CD

4. **Advanced UUID Features**
   - UUID conflict resolution UI
   - Bulk UUID remapping
   - UUID import history tracking
