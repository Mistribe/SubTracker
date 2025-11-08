import Papa from 'papaparse';
import yaml from 'js-yaml';
import type { ImportFileFormat, ParseError } from '../types/import';

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
 * Service for parsing import files in various formats
 */
export class FileParser {
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
   * Read file contents as text
   */
  private readFileAsText(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = (event) => {
        const result = event.target?.result;
        if (typeof result === 'string') {
          resolve(result);
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
   * Parse CSV file using papaparse
   */
  async parseCSV(file: File): Promise<Record<string, any>[]> {
    try {
      const content = await this.readFileAsText(file);

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

            resolve(results.data as Record<string, any>[]);
          },
          error: (error: Error) => {
            reject(new FileParseError(`CSV parsing error: ${error.message}`));
          },
        });
      });
    } catch (error) {
      if (error instanceof FileParseError) {
        throw error;
      }
      throw new FileParseError(
        `Failed to parse CSV file: ${error instanceof Error ? error.message : 'Unknown error'}`
      );
    }
  }

  /**
   * Parse JSON file
   */
  async parseJSON(file: File): Promise<Record<string, any>[]> {
    try {
      const content = await this.readFileAsText(file);
      const parsed = JSON.parse(content);

      // Ensure the result is an array
      if (!Array.isArray(parsed)) {
        throw new FileParseError('JSON file must contain an array of records');
      }

      // Ensure all items are objects
      if (!parsed.every((item) => typeof item === 'object' && item !== null)) {
        throw new FileParseError('All items in JSON array must be objects');
      }

      return parsed;
    } catch (error) {
      if (error instanceof FileParseError) {
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
  async parseYAML(file: File): Promise<Record<string, any>[]> {
    try {
      const content = await this.readFileAsText(file);
      const parsed = yaml.load(content);

      // Ensure the result is an array
      if (!Array.isArray(parsed)) {
        throw new FileParseError('YAML file must contain an array of records');
      }

      // Ensure all items are objects
      if (!parsed.every((item) => typeof item === 'object' && item !== null)) {
        throw new FileParseError('All items in YAML array must be objects');
      }

      return parsed;
    } catch (error) {
      if (error instanceof FileParseError) {
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
  async parse(file: File): Promise<Record<string, any>[]> {
    const format = this.detectFormat(file);

    switch (format) {
      case 'csv':
        return this.parseCSV(file);
      case 'json':
        return this.parseJSON(file);
      case 'yaml':
        return this.parseYAML(file);
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
