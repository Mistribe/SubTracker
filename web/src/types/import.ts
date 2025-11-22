import type {
  DtoCreateLabelRequest,
  DtoCreateProviderRequest,
  DtoCreateSubscriptionRequest,
} from '../api';

/**
 * Supported file formats for import
 */
export type ImportFileFormat = 'csv' | 'json' | 'yaml' | 'unknown';

/**
 * Entity types that can be imported
 */
export type ImportEntityType = 'labels' | 'providers' | 'subscriptions';

/**
 * Parse error information
 */
export interface ParseError {
  line?: number;
  message: string;
}

/**
 * Result of file parsing operation
 */
export interface ParseResult {
  records: Record<string, any>[];
  errors: ParseError[];
}

/**
 * Validation error for a specific field
 */
export interface ValidationError {
  field: string;
  message: string;
  severity: 'error' | 'warning';
}

/**
 * UUID-specific validation error
 */
export interface UUIDValidationError extends ValidationError {
  field: 'id';
  providedValue?: string;
}

/**
 * Result of validation operation
 */
export interface ValidationResult {
  isValid: boolean;
  errors: ValidationError[];
}

/**
 * Import status for individual records
 * - 'skipped' is used when the backend returns HTTP 409 (already exists)
 */
export type ImportStatusType = 'pending' | 'importing' | 'success' | 'error' | 'skipped';

/**
 * Import status with optional error message
 */
export interface ImportStatus {
  status: ImportStatusType;
  error?: string;
}

/**
 * Progress tracking for bulk import operations
 */
export interface ImportProgress {
  total: number;
  completed: number;
  failed: number;
  inProgress: boolean;
  /** Number of records that were skipped because they already exist (HTTP 409) */
  skipped?: number;
}

/**
 * Generic import record wrapper with optional ID field support
 */
export interface ImportRecord<T> {
  index: number;
  rawData: Record<string, any>;
  mappedData: Partial<T> & { id?: string };
  validationErrors: ValidationError[];
  isValid: boolean;
  importStatus: ImportStatus;
}

/**
 * Parsed import record for preview table with optional ID field support
 */
export interface ParsedImportRecord<T> {
  index: number;
  data: Partial<T> & { id?: string };
  validationErrors: ValidationError[];
  isValid: boolean;
}

/**
 * Column definition for import preview table
 */
export interface ImportColumnDef<T> {
  key: keyof T | string;
  label: string;
  render?: (value: any, record: ParsedImportRecord<T>) => React.ReactNode;
}

/**
 * Entity-specific import record types
 */
export type LabelImportRecord = ImportRecord<DtoCreateLabelRequest>;
export type ProviderImportRecord = ImportRecord<DtoCreateProviderRequest>;
export type SubscriptionImportRecord = ImportRecord<DtoCreateSubscriptionRequest>;

/**
 * Field mapper interface for entity-specific mapping
 */
export interface FieldMapper<T> {
  mapFields(rawRecord: Record<string, any>): Partial<T>;
  validate(record: Partial<T>): ValidationResult;
}

/**
 * File parser service interface
 */
export interface FileParserService {
  parseCSV(file: File): Promise<Record<string, any>[]>;
  parseJSON(file: File): Promise<Record<string, any>[]>;
  parseYAML(file: File): Promise<Record<string, any>[]>;
  detectFormat(file: File): ImportFileFormat;
}
