# Import Template Files

This directory contains template files to help you import Labels, Providers, and Subscriptions into the application. Each entity type has templates available in three formats: CSV, JSON, and YAML.

## Available Templates

### Labels
- `labels/labels-template.csv` - Comma-separated values format
- `labels/labels-template.json` - JSON array format
- `labels/labels-template.yaml` - YAML format

### Providers
- `providers/providers-template.csv` - Comma-separated values format
- `providers/providers-template.json` - JSON array format
- `providers/providers-template.yaml` - YAML format

### Subscriptions
- `subscriptions/subscriptions-template.csv` - Comma-separated values format
- `subscriptions/subscriptions-template.json` - JSON array format
- `subscriptions/subscriptions-template.yaml` - YAML format

## Using Custom IDs (UUIDs)

All templates support an optional `id` field that allows you to specify custom UUIDs for your entities. This is useful when:

- Migrating data from another system and want to maintain the same identifiers
- Syncing data across multiple environments
- Maintaining referential integrity across related imports

### UUID Format Requirements

The `id` field must be a valid UUID (Universally Unique Identifier) in the standard format:

```
xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

Where each `x` is a hexadecimal digit (0-9, a-f, A-F).

**Valid UUID Examples:**
- `550e8400-e29b-41d4-a716-446655440001`
- `6ba7b810-9dad-11d1-80b4-00c04fd430c8`
- `f47ac10b-58cc-4372-a567-0e02b2c3d479`

**When to Use Custom IDs:**
- ✅ Use custom IDs when you need to maintain specific identifiers
- ✅ Use custom IDs when importing related entities that reference each other
- ❌ Leave the `id` field empty if you want the system to auto-generate UUIDs
- ❌ Don't use duplicate IDs - each entity must have a unique identifier

### Auto-Generated IDs

If you don't specify an `id` field (or leave it empty), the system will automatically generate a unique UUID for each entity. This is the recommended approach for most imports.

## Field Descriptions

### Labels

| Field | Required | Type | Description | Example |
|-------|----------|------|-------------|---------|
| `id` | No | UUID | Custom identifier for the label | `550e8400-e29b-41d4-a716-446655440001` |
| `name` | Yes | String | Display name of the label | `Entertainment` |
| `color` | Yes | Hex Color | Color code for the label | `#FF5733` |
| `ownerType` | Yes | Enum | Owner type: `personal` or `family` | `personal` |
| `ownerFamilyId` | No | UUID | Family ID if ownerType is `family` | `550e8400-e29b-41d4-a716-446655440000` |

### Providers

| Field | Required | Type | Description | Example |
|-------|----------|------|-------------|---------|
| `id` | No | UUID | Custom identifier for the provider | `550e8400-e29b-41d4-a716-446655440010` |
| `name` | Yes | String | Name of the service provider | `Netflix` |
| `description` | No | String | Brief description of the service | `Streaming service for movies and TV shows` |
| `url` | No | URL | Provider's website URL | `https://www.netflix.com` |
| `iconUrl` | No | URL | URL to provider's icon/logo | `https://www.netflix.com/favicon.ico` |
| `pricingPageUrl` | No | URL | URL to provider's pricing page | `https://www.netflix.com/signup/planform` |
| `labels` | No | Array/String | Associated labels (comma-separated in CSV) | `Entertainment,Video Streaming` |

### Subscriptions

| Field | Required | Type | Description | Example |
|-------|----------|------|-------------|---------|
| `id` | No | UUID | Custom identifier for the subscription | `550e8400-e29b-41d4-a716-446655440020` |
| `providerId` | Yes | UUID | ID of the provider for this subscription | `550e8400-e29b-41d4-a716-446655440010` |
| `friendlyName` | Yes | String | Custom name for this subscription | `Netflix Premium` |
| `startDate` | Yes | Date | Start date (YYYY-MM-DD format) | `2024-01-01` |
| `endDate` | No | Date | End date (YYYY-MM-DD format) | `2024-12-31` |
| `recurrency` | Yes | Enum | Billing frequency: `monthly`, `yearly`, `weekly`, `custom` | `monthly` |
| `customRecurrency` | No | Number | Custom recurrency in days (if recurrency is `custom`) | `90` |
| `customPriceAmount` | Yes | Number | Subscription price amount | `15.99` |
| `customPriceCurrency` | Yes | String | Currency code (ISO 4217) | `USD` |
| `ownerType` | Yes | Enum | Owner type: `personal` or `family` | `personal` |
| `ownerFamilyId` | No | UUID | Family ID if ownerType is `family` | `550e8400-e29b-41d4-a716-446655440000` |
| `payerType` | No | Enum | Payer type: `personal` or `family` | `personal` |
| `payerFamilyId` | No | UUID | Family ID if payerType is `family` | `550e8400-e29b-41d4-a716-446655440000` |
| `freeTrialStartDate` | No | Date | Free trial start date (YYYY-MM-DD) | `2024-03-01` |
| `freeTrialEndDate` | No | Date | Free trial end date (YYYY-MM-DD) | `2024-03-31` |
| `labels` | No | Array/String | Associated labels (comma-separated in CSV) | `Entertainment,Video Streaming` |

## Common Import Scenarios

### Scenario 1: Simple Import Without Custom IDs

The easiest way to import data is to omit the `id` field entirely. The system will generate unique IDs automatically.

**Example (JSON):**
```json
[
  {
    "name": "Entertainment",
    "color": "#FF5733",
    "ownerType": "personal"
  }
]
```

### Scenario 2: Import With Custom IDs

When you need to maintain specific identifiers (e.g., migrating from another system):

**Example (JSON):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Entertainment",
    "color": "#FF5733",
    "ownerType": "personal"
  }
]
```

### Scenario 3: Mixed Import (Some With IDs, Some Without)

You can mix records with and without IDs in the same import file:

**Example (JSON):**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Entertainment",
    "color": "#FF5733",
    "ownerType": "personal"
  },
  {
    "name": "Utilities",
    "color": "#33FF57",
    "ownerType": "personal"
  }
]
```

### Scenario 4: Importing Related Entities

When importing subscriptions that reference providers, you can use custom IDs to maintain relationships:

**Step 1: Import Providers with custom IDs**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "name": "Netflix",
    "url": "https://www.netflix.com"
  }
]
```

**Step 2: Import Subscriptions referencing those provider IDs**
```json
[
  {
    "providerId": "550e8400-e29b-41d4-a716-446655440010",
    "friendlyName": "Netflix Premium",
    "startDate": "2024-01-01",
    "recurrency": "monthly",
    "customPriceAmount": 15.99,
    "customPriceCurrency": "USD",
    "ownerType": "personal"
  }
]
```

## Format-Specific Notes

### CSV Format
- First row must contain column headers
- Use commas to separate values
- Wrap values containing commas in double quotes
- For array fields (like `labels`), use comma-separated values within quotes: `"Entertainment,Video Streaming"`
- Empty fields can be left blank (e.g., `,,` for consecutive empty fields)

### JSON Format
- Must be a valid JSON array `[...]`
- Use double quotes for strings
- Array fields should be actual arrays: `["Entertainment", "Video Streaming"]`
- Omit fields that are empty rather than using `null` or empty strings

### YAML Format
- Use proper indentation (2 spaces recommended)
- Each record starts with a dash `-`
- Array fields use YAML array syntax with dashes
- Strings with special characters should be quoted

## Validation and Error Handling

### UUID Validation Errors

If you provide an invalid UUID format, you'll see an error message like:

```
Invalid UUID format: "invalid-id". Expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

**Common UUID Format Mistakes:**
- Missing dashes: `550e8400e29b41d4a716446655440001` ❌
- Wrong length: `550e8400-e29b-41d4-a716` ❌
- Invalid characters: `550e8400-e29b-41d4-a716-44665544000g` ❌
- Correct format: `550e8400-e29b-41d4-a716-446655440001` ✅

### Duplicate UUID Errors

If you try to import an entity with an ID that already exists in the system, you'll receive a conflict error. Make sure all IDs in your import file are unique and don't conflict with existing entities.

## Tips for Successful Imports

1. **Start Small**: Test with a few records first before importing large datasets
2. **Validate Your Data**: Ensure all required fields are present and properly formatted
3. **Use Templates**: Download and modify the provided templates rather than creating files from scratch
4. **Check Date Formats**: Always use `YYYY-MM-DD` format for dates
5. **Verify URLs**: Make sure all URLs are complete and valid
6. **Test UUID Format**: If using custom IDs, verify they match the UUID format before importing
7. **Review Preview**: Always review the import preview before confirming the import

## Getting Help

If you encounter issues during import:
- Check the validation errors displayed in the preview table
- Verify your file format matches one of the supported formats (CSV, JSON, YAML)
- Ensure all required fields are present
- Confirm UUID formats are correct if using custom IDs
- Review this documentation for field requirements and examples

For additional support, please refer to the application's help documentation or contact support.
