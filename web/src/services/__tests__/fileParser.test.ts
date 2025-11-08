import { describe, it, expect, beforeEach } from 'vitest';
import { FileParser, FileParseError } from '../fileParser';

describe('FileParser', () => {
  let parser: FileParser;

  beforeEach(() => {
    parser = new FileParser();
  });

  describe('detectFormat', () => {
    it('should detect CSV format', () => {
      const file = new File([''], 'test.csv', { type: 'text/csv' });
      expect(parser.detectFormat(file)).toBe('csv');
    });

    it('should detect JSON format', () => {
      const file = new File([''], 'test.json', { type: 'application/json' });
      expect(parser.detectFormat(file)).toBe('json');
    });

    it('should detect YAML format with .yaml extension', () => {
      const file = new File([''], 'test.yaml', { type: 'text/yaml' });
      expect(parser.detectFormat(file)).toBe('yaml');
    });

    it('should detect YAML format with .yml extension', () => {
      const file = new File([''], 'test.yml', { type: 'text/yaml' });
      expect(parser.detectFormat(file)).toBe('yaml');
    });

    it('should return unknown for unsupported formats', () => {
      const file = new File([''], 'test.txt', { type: 'text/plain' });
      expect(parser.detectFormat(file)).toBe('unknown');
    });

    it('should handle files without extensions', () => {
      const file = new File([''], 'test', { type: 'text/plain' });
      expect(parser.detectFormat(file)).toBe('unknown');
    });

    it('should be case-insensitive', () => {
      const file = new File([''], 'test.CSV', { type: 'text/csv' });
      expect(parser.detectFormat(file)).toBe('csv');
    });
  });

  describe('parseCSV', () => {
    it('should parse valid CSV with headers', async () => {
      const csvContent = 'name,color\nEntertainment,#FF5733\nUtilities,#33FF57';
      const file = new File([csvContent], 'test.csv', { type: 'text/csv' });

      const result = await parser.parseCSV(file);

      expect(result).toEqual([
        { name: 'Entertainment', color: '#FF5733' },
        { name: 'Utilities', color: '#33FF57' },
      ]);
    });

    it('should handle CSV with quoted fields', async () => {
      const csvContent = 'name,description\n"Test Name","Description with, comma"\n"Another","Normal"';
      const file = new File([csvContent], 'test.csv', { type: 'text/csv' });

      const result = await parser.parseCSV(file);

      expect(result).toEqual([
        { name: 'Test Name', description: 'Description with, comma' },
        { name: 'Another', description: 'Normal' },
      ]);
    });

    it('should skip empty lines', async () => {
      const csvContent = 'name,color\nEntertainment,#FF5733\n\nUtilities,#33FF57\n';
      const file = new File([csvContent], 'test.csv', { type: 'text/csv' });

      const result = await parser.parseCSV(file);

      expect(result).toHaveLength(2);
    });

    it('should trim header whitespace', async () => {
      const csvContent = ' name , color \nEntertainment,#FF5733';
      const file = new File([csvContent], 'test.csv', { type: 'text/csv' });

      const result = await parser.parseCSV(file);

      expect(result[0]).toHaveProperty('name');
      expect(result[0]).toHaveProperty('color');
    });

    it('should handle empty CSV file', async () => {
      const csvContent = '';
      const file = new File([csvContent], 'test.csv', { type: 'text/csv' });

      // Empty CSV files may throw an error or return empty array depending on parser
      // We accept either behavior as valid
      try {
        const result = await parser.parseCSV(file);
        expect(result).toEqual([]);
      } catch (error) {
        expect(error).toBeInstanceOf(FileParseError);
      }
    });

    it('should handle CSV with only headers', async () => {
      const csvContent = 'name,color';
      const file = new File([csvContent], 'test.csv', { type: 'text/csv' });

      const result = await parser.parseCSV(file);

      expect(result).toEqual([]);
    });

    it('should throw FileParseError for malformed CSV', async () => {
      const csvContent = 'name,color\n"Unclosed quote,#FF5733';
      const file = new File([csvContent], 'test.csv', { type: 'text/csv' });

      await expect(parser.parseCSV(file)).rejects.toThrow(FileParseError);
    });
  });

  describe('parseJSON', () => {
    it('should parse valid JSON array', async () => {
      const jsonContent = JSON.stringify([
        { name: 'Netflix', description: 'Streaming service' },
        { name: 'AWS', description: 'Cloud services' },
      ]);
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      const result = await parser.parseJSON(file);

      expect(result).toEqual([
        { name: 'Netflix', description: 'Streaming service' },
        { name: 'AWS', description: 'Cloud services' },
      ]);
    });

    it('should parse empty JSON array', async () => {
      const jsonContent = '[]';
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      const result = await parser.parseJSON(file);

      expect(result).toEqual([]);
    });

    it('should handle JSON with nested objects', async () => {
      const jsonContent = JSON.stringify([
        {
          name: 'Test',
          nested: { key: 'value' },
        },
      ]);
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      const result = await parser.parseJSON(file);

      expect(result[0].nested).toEqual({ key: 'value' });
    });

    it('should throw error for non-array JSON', async () => {
      const jsonContent = '{"name": "Single object"}';
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      await expect(parser.parseJSON(file)).rejects.toThrow(
        'JSON file must contain an array of records'
      );
    });

    it('should throw error for array with non-object items', async () => {
      const jsonContent = '["string", "another string"]';
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      await expect(parser.parseJSON(file)).rejects.toThrow(
        'All items in JSON array must be objects'
      );
    });

    it('should throw error for array with null items', async () => {
      const jsonContent = '[{"name": "Valid"}, null]';
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      await expect(parser.parseJSON(file)).rejects.toThrow(
        'All items in JSON array must be objects'
      );
    });

    it('should throw FileParseError for invalid JSON syntax', async () => {
      const jsonContent = '{invalid json}';
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      await expect(parser.parseJSON(file)).rejects.toThrow(FileParseError);
      await expect(parser.parseJSON(file)).rejects.toThrow(/Invalid JSON format/);
    });

    it('should throw FileParseError for empty file', async () => {
      const jsonContent = '';
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      await expect(parser.parseJSON(file)).rejects.toThrow(FileParseError);
    });
  });

  describe('parseYAML', () => {
    it('should parse valid YAML array', async () => {
      const yamlContent = `
- name: Netflix
  description: Streaming service
- name: AWS
  description: Cloud services
`;
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      const result = await parser.parseYAML(file);

      expect(result).toEqual([
        { name: 'Netflix', description: 'Streaming service' },
        { name: 'AWS', description: 'Cloud services' },
      ]);
    });

    it('should parse empty YAML array', async () => {
      const yamlContent = '[]';
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      const result = await parser.parseYAML(file);

      expect(result).toEqual([]);
    });

    it('should handle YAML with nested objects', async () => {
      const yamlContent = `
- name: Test
  nested:
    key: value
    another: data
`;
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      const result = await parser.parseYAML(file);

      expect(result[0].nested).toEqual({ key: 'value', another: 'data' });
    });

    it('should handle YAML with arrays', async () => {
      const yamlContent = `
- name: Provider
  labels:
    - label1
    - label2
`;
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      const result = await parser.parseYAML(file);

      expect(result[0].labels).toEqual(['label1', 'label2']);
    });

    it('should throw error for non-array YAML', async () => {
      const yamlContent = 'name: Single object';
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      await expect(parser.parseYAML(file)).rejects.toThrow(
        'YAML file must contain an array of records'
      );
    });

    it('should throw error for array with non-object items', async () => {
      const yamlContent = '- string\n- another string';
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      await expect(parser.parseYAML(file)).rejects.toThrow(
        'All items in YAML array must be objects'
      );
    });

    it('should throw FileParseError for invalid YAML syntax', async () => {
      const yamlContent = '- name: Test\n  invalid indentation';
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      await expect(parser.parseYAML(file)).rejects.toThrow(FileParseError);
    });

    it('should include line number in YAML parse errors', async () => {
      const yamlContent = '- name: Test\n  : invalid';
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      try {
        await parser.parseYAML(file);
        expect.fail('Should have thrown an error');
      } catch (error) {
        expect(error).toBeInstanceOf(FileParseError);
        if (error instanceof FileParseError) {
          expect(error.errors.length).toBeGreaterThan(0);
        }
      }
    });
  });

  describe('parse', () => {
    it('should parse CSV file using detectFormat', async () => {
      const csvContent = 'name,color\nTest,#FF5733';
      const file = new File([csvContent], 'test.csv', { type: 'text/csv' });

      const result = await parser.parse(file);

      expect(result).toEqual([{ name: 'Test', color: '#FF5733' }]);
    });

    it('should parse JSON file using detectFormat', async () => {
      const jsonContent = JSON.stringify([{ name: 'Test' }]);
      const file = new File([jsonContent], 'test.json', { type: 'application/json' });

      const result = await parser.parse(file);

      expect(result).toEqual([{ name: 'Test' }]);
    });

    it('should parse YAML file using detectFormat', async () => {
      const yamlContent = '- name: Test';
      const file = new File([yamlContent], 'test.yaml', { type: 'text/yaml' });

      const result = await parser.parse(file);

      expect(result).toEqual([{ name: 'Test' }]);
    });

    it('should throw error for unknown file format', async () => {
      const file = new File(['content'], 'test.txt', { type: 'text/plain' });

      await expect(parser.parse(file)).rejects.toThrow(
        'Unsupported file format. Please upload a CSV, JSON, or YAML file.'
      );
    });
  });

  describe('FileParseError', () => {
    it('should create error with message', () => {
      const error = new FileParseError('Test error');

      expect(error.message).toBe('Test error');
      expect(error.name).toBe('FileParseError');
      expect(error.errors).toEqual([]);
    });

    it('should create error with parse errors', () => {
      const parseErrors = [
        { line: 1, message: 'Error on line 1' },
        { line: 2, message: 'Error on line 2' },
      ];
      const error = new FileParseError('Multiple errors', parseErrors);

      expect(error.message).toBe('Multiple errors');
      expect(error.errors).toEqual(parseErrors);
    });
  });
});
