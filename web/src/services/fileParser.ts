import Papa from 'papaparse';
import yaml from 'js-yaml';
import type { ImportFileFormat, ParseError } from '../types/import';

/**
 * Maximum file size in bytes (10MB)
 */
const MAX_FILE_SIZE = 10 * 1024 * 1024;

/**
 * File size threshold for streaming (5MB)
 */
const STREAMING_THRESHOLD = 5 * 1024 * 1024;

/**
 * Error thrown when file parsing fails
 */
export class FileParseError extends Error {
  constructor(
    message: string,
    public readonly errors: ParseError[] = []
  ) {
    super(message);
    this.name = 'FileParseError';
  }
}

/**
 * Error thrown when file size exceeds limit
 */
export class FileSizeError extends Error {
  constructor(
    public readonly fileSize: number,
    public readonly maxSize: number
  ) {
    super(`File size (${(fileSize / 1024 / 1024).toFixed(2)}MB) exceeds maximum allowed size (${(maxSize / 1024 / 1024).toFixed(2)}MB)`);
    this.name = 'FileSizeError';
  }
}

/**
 * Service for parsing import files in various formats
 */
export class FileParser {
  /**
   * Validate file size
   */
  private validateFileSize(file: File): void {
    if (file.size > MAX_FILE_SIZE) {
      throw new FileSizeError(file.size, MAX_FILE_SIZE);
    }
  }

  /**
   * Detect file format based on file extension
   */
  detectFormat(file: File): ImportFileFormat {
    const extension = file.name.split('.').pop()?.toLowerCase();

    switch (extension) {
      case 'csv':
        return 'csv';
      case 'json':
        return 'json';
      case 'yaml':
      case 'yml':
        return 'yaml';
      default:
        return 'unknown';
    }
  }

  /**
   * Read file contents as text with progress tracking
   */
  private readFileAsText(
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onprogress = (event) => {
        if (event.lengthComputable && onProgress) {
          const progress = (event.loaded / event.total) * 100;
          onProgress(progress);
        }
      };

      reader.onload = (event) => {
        const result = event.target?.result;
        if (typeof result === 'string') {
          resolve(result);
          // Release file reference to optimize memory
          reader.abort();
        } else {
          reject(new Error('Failed to read file as text'));
        }
      };

      reader.onerror = () => {
        reject(new Error(`File reading failed: ${reader.error?.message || 'Unknown error'}`));
      };

      reader.readAsText(file);
    });
  }

  /**
   * Parse CSV file using papaparse with streaming support for large files
   */
  async parseCSV(
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<Record<string, any>[]> {
    try {
      // Validate file size
      this.validateFileSize(file);

      // Use streaming for large files
      if (file.size > STREAMING_THRESHOLD) {
        return this.parseCSVStreaming(file, onProgress);
      }

      // Use standard parsing for smaller files
      const content = await this.readFileAsText(file, onProgress);

      return new Promise((resolve, reject) => {
        Papa.parse(content, {
          header: true,
          skipEmptyLines: true,
          transformHeader: (header) => header.trim(),
          complete: (results) => {
            if (results.errors.length > 0) {
              const parseErrors: ParseError[] = results.errors.map((error) => ({
                line: error.row !== undefined ? error.row + 1 : undefined,
                message: error.message,
              }));

              reject(new FileParseError('CSV parsing failed', parseErrors));
              return;
            }

            if (onProgress) {
              onProgress(100);
            }

            resolve(results.data as Record<string, any>[]);
          },
          error: (error: Error) => {
            reject(new FileParseError(`CSV parsing error: ${error.message}`));
          },
        });
      });
    } catch (error) {
      if (error instanceof FileParseError || error instanceof FileSizeError) {
        throw error;
      }
      throw new FileParseError(
        `Failed to parse CSV file: ${error instanceof Error ? error.message : 'Unknown error'}`
      );
    }
  }

  /**
   * Parse large CSV files using streaming
   */
  private parseCSVStreaming(
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<Record<string, any>[]> {
    return new Promise((resolve, reject) => {
      const records: Record<string, any>[] = [];
      let processedBytes = 0;

      Papa.parse(file, {
        header: true,
        skipEmptyLines: true,
        transformHeader: (header) => header.trim(),
        chunk: (results, parser) => {
          // Process chunk
          records.push(...(results.data as Record<string, any>[]));

          // Update progress
          processedBytes += new Blob([JSON.stringify(results.data)]).size;
          if (onProgress) {
            const progress = Math.min((processedBytes / file.size) * 100, 99);
            onProgress(progress);
          }

          // Check for errors in chunk
          if (results.errors.length > 0) {
            parser.abort();
            const parseErrors: ParseError[] = results.errors.map((error) => ({
              line: error.row !== undefined ? error.row + 1 : undefined,
              message: error.message,
            }));
            reject(new FileParseError('CSV parsing failed', parseErrors));
          }
        },
        complete: () => {
          if (onProgress) {
            onProgress(100);
          }
          resolve(records);
        },
        error: (error: Error) => {
          reject(new FileParseError(`CSV streaming error: ${error.message}`));
        },
      });
    });
  }

  /**
   * Parse JSON file
   */
  async parseJSON(
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<Record<string, any>[]> {
    try {
      // Validate file size
      this.validateFileSize(file);

      const content = await this.readFileAsText(file, onProgress);
      const parsed = JSON.parse(content);

      // Ensure the result is an array
      if (!Array.isArray(parsed)) {
        throw new FileParseError('JSON file must contain an array of records');
      }

      // Ensure all items are objects
      if (!parsed.every((item) => typeof item === 'object' && item !== null)) {
        throw new FileParseError('All items in JSON array must be objects');
      }

      if (onProgress) {
        onProgress(100);
      }

      return parsed;
    } catch (error) {
      if (error instanceof FileParseError || error instanceof FileSizeError) {
        throw error;
      }

      if (error instanceof SyntaxError) {
        throw new FileParseError(`Invalid JSON format: ${error.message}`);
      }

      throw new FileParseError(
        `Failed to parse JSON file: ${error instanceof Error ? error.message : 'Unknown error'}`
      );
    }
  }

  /**
   * Parse YAML file using js-yaml
   */
  async parseYAML(
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<Record<string, any>[]> {
    try {
      // Validate file size
      this.validateFileSize(file);

      const content = await this.readFileAsText(file, onProgress);
      const parsed = yaml.load(content);

      // Ensure the result is an array
      if (!Array.isArray(parsed)) {
        throw new FileParseError('YAML file must contain an array of records');
      }

      // Ensure all items are objects
      if (!parsed.every((item) => typeof item === 'object' && item !== null)) {
        throw new FileParseError('All items in YAML array must be objects');
      }

      if (onProgress) {
        onProgress(100);
      }

      return parsed;
    } catch (error) {
      if (error instanceof FileParseError || error instanceof FileSizeError) {
        throw error;
      }

      if (error instanceof yaml.YAMLException) {
        const parseError: ParseError = {
          line: error.mark?.line !== undefined ? error.mark.line + 1 : undefined,
          message: error.message,
        };
        throw new FileParseError('YAML parsing failed', [parseError]);
      }

      throw new FileParseError(
        `Failed to parse YAML file: ${error instanceof Error ? error.message : 'Unknown error'}`
      );
    }
  }

  /**
   * Parse file based on detected format
   */
  async parse(
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<Record<string, any>[]> {
    const format = this.detectFormat(file);

    switch (format) {
      case 'csv':
        return this.parseCSV(file, onProgress);
      case 'json':
        return this.parseJSON(file, onProgress);
      case 'yaml':
        return this.parseYAML(file, onProgress);
      case 'unknown':
        throw new FileParseError(
          `Unsupported file format. Please upload a CSV, JSON, or YAML file.`
        );
    }
  }
}

/**
 * Singleton instance of FileParser
 */
export const fileParser = new FileParser();

/**
 * Export constants for external use
 */
export { MAX_FILE_SIZE, STREAMING_THRESHOLD };
