/**
 * Script to validate all template files
 * Ensures templates are correctly formatted and contain valid data
 */

import * as fs from 'fs';
import * as path from 'path';
import Papa from 'papaparse';
import yaml from 'js-yaml';

// UUID validation regex (RFC 4122 compliant)
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

interface ValidationResult {
  file: string;
  valid: boolean;
  errors: string[];
  warnings: string[];
}

interface TemplateValidation {
  entityType: 'labels' | 'providers' | 'subscriptions';
  format: 'csv' | 'json' | 'yaml';
  path: string;
}

const TEMPLATES: TemplateValidation[] = [
  // Labels
  { entityType: 'labels', format: 'csv', path: 'public/templates/labels/labels-template.csv' },
  { entityType: 'labels', format: 'json', path: 'public/templates/labels/labels-template.json' },
  { entityType: 'labels', format: 'yaml', path: 'public/templates/labels/labels-template.yaml' },
  // Providers
  { entityType: 'providers', format: 'csv', path: 'public/templates/providers/providers-template.csv' },
  { entityType: 'providers', format: 'json', path: 'public/templates/providers/providers-template.json' },
  { entityType: 'providers', format: 'yaml', path: 'public/templates/providers/providers-template.yaml' },
  // Subscriptions
  { entityType: 'subscriptions', format: 'csv', path: 'public/templates/subscriptions/subscriptions-template.csv' },
  { entityType: 'subscriptions', format: 'json', path: 'public/templates/subscriptions/subscriptions-template.json' },
  { entityType: 'subscriptions', format: 'yaml', path: 'public/templates/subscriptions/subscriptions-template.yaml' },
];

// Required fields for each entity type
const REQUIRED_FIELDS = {
  labels: ['name', 'color', 'ownerType'],
  providers: ['name'],
  subscriptions: ['providerId', 'friendlyName', 'startDate', 'recurrency', 'customPriceAmount', 'customPriceCurrency', 'ownerType'],
};

function validateUUID(id: string): boolean {
  return UUID_REGEX.test(id);
}

function validateCSV(filePath: string, entityType: string): ValidationResult {
  const result: ValidationResult = {
    file: filePath,
    valid: true,
    errors: [],
    warnings: [],
  };

  try {
    const content = fs.readFileSync(filePath, 'utf-8');
    
    const parseResult = Papa.parse(content, {
      header: true,
      skipEmptyLines: true,
      transformHeader: (header) => header.trim(),
    });

    if (parseResult.errors.length > 0) {
      result.valid = false;
      parseResult.errors.forEach((error) => {
        result.errors.push(`Parse error at row ${error.row}: ${error.message}`);
      });
      return result;
    }

    const records = parseResult.data as Record<string, any>[];
    
    if (records.length === 0) {
      result.valid = false;
      result.errors.push('Template contains no records');
      return result;
    }

    // Check for required fields in headers
    const headers = parseResult.meta.fields || [];
    const requiredFields = REQUIRED_FIELDS[entityType as keyof typeof REQUIRED_FIELDS];
    
    requiredFields.forEach((field) => {
      if (!headers.includes(field)) {
        result.valid = false;
        result.errors.push(`Missing required field in headers: ${field}`);
      }
    });

    // Check for id field
    if (!headers.includes('id')) {
      result.warnings.push('Template does not include "id" field in headers');
    }

    // Validate records
    let recordsWithIds = 0;
    let recordsWithoutIds = 0;

    records.forEach((record, index) => {
      const rowNum = index + 2; // +2 for header row and 0-based index

      // Check if record has ID
      if (record.id && record.id.trim() !== '') {
        recordsWithIds++;
        
        // Validate UUID format
        if (!validateUUID(record.id.trim())) {
          result.valid = false;
          result.errors.push(`Row ${rowNum}: Invalid UUID format: "${record.id}"`);
        }
      } else {
        recordsWithoutIds++;
      }

      // Check required fields
      requiredFields.forEach((field) => {
        if (!record[field] || record[field].toString().trim() === '') {
          result.valid = false;
          result.errors.push(`Row ${rowNum}: Missing required field: ${field}`);
        }
      });
    });

    // Verify mix of records with and without IDs
    if (recordsWithIds === 0) {
      result.warnings.push('No records with IDs found (should have examples with IDs)');
    }
    if (recordsWithoutIds === 0) {
      result.warnings.push('No records without IDs found (should have examples without IDs)');
    }

  } catch (error) {
    result.valid = false;
    result.errors.push(`Failed to read or parse file: ${error instanceof Error ? error.message : 'Unknown error'}`);
  }

  return result;
}

function validateJSON(filePath: string, entityType: string): ValidationResult {
  const result: ValidationResult = {
    file: filePath,
    valid: true,
    errors: [],
    warnings: [],
  };

  try {
    const content = fs.readFileSync(filePath, 'utf-8');
    const records = JSON.parse(content);

    if (!Array.isArray(records)) {
      result.valid = false;
      result.errors.push('JSON must be an array of records');
      return result;
    }

    if (records.length === 0) {
      result.valid = false;
      result.errors.push('Template contains no records');
      return result;
    }

    // Validate records
    let recordsWithIds = 0;
    let recordsWithoutIds = 0;
    const requiredFields = REQUIRED_FIELDS[entityType as keyof typeof REQUIRED_FIELDS];

    records.forEach((record, index) => {
      if (typeof record !== 'object' || record === null) {
        result.valid = false;
        result.errors.push(`Record ${index + 1}: Must be an object`);
        return;
      }

      // Check if record has ID
      if (record.id && record.id.trim() !== '') {
        recordsWithIds++;
        
        // Validate UUID format
        if (!validateUUID(record.id.trim())) {
          result.valid = false;
          result.errors.push(`Record ${index + 1}: Invalid UUID format: "${record.id}"`);
        }
      } else {
        recordsWithoutIds++;
      }

      // Check required fields
      requiredFields.forEach((field) => {
        if (!record[field] || (typeof record[field] === 'string' && record[field].trim() === '')) {
          result.valid = false;
          result.errors.push(`Record ${index + 1}: Missing required field: ${field}`);
        }
      });
    });

    // Verify mix of records with and without IDs
    if (recordsWithIds === 0) {
      result.warnings.push('No records with IDs found (should have examples with IDs)');
    }
    if (recordsWithoutIds === 0) {
      result.warnings.push('No records without IDs found (should have examples without IDs)');
    }

  } catch (error) {
    result.valid = false;
    if (error instanceof SyntaxError) {
      result.errors.push(`Invalid JSON syntax: ${error.message}`);
    } else {
      result.errors.push(`Failed to read or parse file: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  return result;
}

function validateYAML(filePath: string, entityType: string): ValidationResult {
  const result: ValidationResult = {
    file: filePath,
    valid: true,
    errors: [],
    warnings: [],
  };

  try {
    const content = fs.readFileSync(filePath, 'utf-8');
    const records = yaml.load(content);

    if (!Array.isArray(records)) {
      result.valid = false;
      result.errors.push('YAML must be an array of records');
      return result;
    }

    if (records.length === 0) {
      result.valid = false;
      result.errors.push('Template contains no records');
      return result;
    }

    // Validate records
    let recordsWithIds = 0;
    let recordsWithoutIds = 0;
    const requiredFields = REQUIRED_FIELDS[entityType as keyof typeof REQUIRED_FIELDS];

    records.forEach((record, index) => {
      if (typeof record !== 'object' || record === null) {
        result.valid = false;
        result.errors.push(`Record ${index + 1}: Must be an object`);
        return;
      }

      // Check if record has ID
      if (record.id && record.id.trim() !== '') {
        recordsWithIds++;
        
        // Validate UUID format
        if (!validateUUID(record.id.trim())) {
          result.valid = false;
          result.errors.push(`Record ${index + 1}: Invalid UUID format: "${record.id}"`);
        }
      } else {
        recordsWithoutIds++;
      }

      // Check required fields
      requiredFields.forEach((field) => {
        if (!record[field] || (typeof record[field] === 'string' && record[field].trim() === '')) {
          result.valid = false;
          result.errors.push(`Record ${index + 1}: Missing required field: ${field}`);
        }
      });
    });

    // Verify mix of records with and without IDs
    if (recordsWithIds === 0) {
      result.warnings.push('No records with IDs found (should have examples with IDs)');
    }
    if (recordsWithoutIds === 0) {
      result.warnings.push('No records without IDs found (should have examples without IDs)');
    }

  } catch (error) {
    result.valid = false;
    if (error instanceof yaml.YAMLException) {
      result.errors.push(`YAML parsing error: ${error.message}`);
    } else {
      result.errors.push(`Failed to read or parse file: ${error instanceof Error ? error.message : 'Unknown error'}`);
    }
  }

  return result;
}

function validateTemplate(template: TemplateValidation): ValidationResult {
  switch (template.format) {
    case 'csv':
      return validateCSV(template.path, template.entityType);
    case 'json':
      return validateJSON(template.path, template.entityType);
    case 'yaml':
      return validateYAML(template.path, template.entityType);
  }
}

function main() {
  console.log('🔍 Validating template files...\n');

  const results: ValidationResult[] = [];
  let allValid = true;

  TEMPLATES.forEach((template) => {
    const result = validateTemplate(template);
    results.push(result);

    if (!result.valid) {
      allValid = false;
    }
  });

  // Print results
  results.forEach((result) => {
    const status = result.valid ? '✅' : '❌';
    console.log(`${status} ${result.file}`);

    if (result.errors.length > 0) {
      console.log('  Errors:');
      result.errors.forEach((error) => {
        console.log(`    - ${error}`);
      });
    }

    if (result.warnings.length > 0) {
      console.log('  Warnings:');
      result.warnings.forEach((warning) => {
        console.log(`    - ${warning}`);
      });
    }

    console.log('');
  });

  // Summary
  const validCount = results.filter((r) => r.valid).length;
  const totalCount = results.length;

  console.log('─'.repeat(60));
  console.log(`Summary: ${validCount}/${totalCount} templates valid`);

  if (allValid) {
    console.log('✅ All templates are valid!');
    process.exit(0);
  } else {
    console.log('❌ Some templates have errors. Please fix them.');
    process.exit(1);
  }
}

main();
