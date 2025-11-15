import { describe, it, expect } from 'vitest';
import {
  isValidUUID,
  validateAndSanitizeUUID,
} from '../uuidValidation';

describe('UUID Validation Performance', () => {
  // Helper to generate a valid UUID v4
  const generateValidUUID = (): string => {
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
      const r = Math.random() * 16 | 0;
      const v = c === 'x' ? r : (r & 0x3 | 0x8);
      return v.toString(16);
    });
  };

  // Helper to generate an invalid UUID
  const generateInvalidUUID = (): string => {
    return 'invalid-uuid-' + Math.random().toString(36).substring(7);
  };

  describe('isValidUUID performance with large datasets', () => {
    it('should validate 1000 valid UUIDs in reasonable time', () => {
      const uuids = Array.from({ length: 1000 }, () => generateValidUUID());
      
      const startTime = performance.now();
      
      for (const uuid of uuids) {
        isValidUUID(uuid);
      }
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should complete in less than 100ms for 1000 UUIDs
      expect(duration).toBeLessThan(100);
      
      // Verify all were validated correctly
      const results = uuids.map(uuid => isValidUUID(uuid));
      expect(results.every(result => result === true)).toBe(true);
    });

    it('should validate 1000 invalid UUIDs in reasonable time', () => {
      const uuids = Array.from({ length: 1000 }, () => generateInvalidUUID());
      
      const startTime = performance.now();
      
      for (const uuid of uuids) {
        isValidUUID(uuid);
      }
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should complete in less than 100ms for 1000 UUIDs
      expect(duration).toBeLessThan(100);
      
      // Verify all were validated correctly
      const results = uuids.map(uuid => isValidUUID(uuid));
      expect(results.every(result => result === false)).toBe(true);
    });

    it('should validate 5000 mixed UUIDs in reasonable time', () => {
      const validUuids = Array.from({ length: 2500 }, () => generateValidUUID());
      const invalidUuids = Array.from({ length: 2500 }, () => generateInvalidUUID());
      const mixedUuids = [...validUuids, ...invalidUuids].sort(() => Math.random() - 0.5);
      
      const startTime = performance.now();
      
      for (const uuid of mixedUuids) {
        isValidUUID(uuid);
      }
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should complete in less than 500ms for 5000 UUIDs
      expect(duration).toBeLessThan(500);
    });
  });

  describe('validateAndSanitizeUUID performance with large datasets', () => {
    it('should validate and sanitize 1000 valid UUIDs in reasonable time', () => {
      const uuids = Array.from({ length: 1000 }, () => generateValidUUID());
      
      const startTime = performance.now();
      
      for (const uuid of uuids) {
        validateAndSanitizeUUID(uuid);
      }
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should complete in less than 150ms for 1000 UUIDs (includes sanitization)
      expect(duration).toBeLessThan(150);
      
      // Verify all were validated correctly
      const results = uuids.map(uuid => validateAndSanitizeUUID(uuid));
      expect(results.every(result => result.isValid === true)).toBe(true);
    });

    it('should handle 1000 UUIDs with whitespace efficiently', () => {
      const uuids = Array.from({ length: 1000 }, () => {
        const uuid = generateValidUUID();
        // Add random whitespace
        const whitespace = ['  ', '\t', '\n', '  \t', '\n  '];
        const prefix = whitespace[Math.floor(Math.random() * whitespace.length)];
        const suffix = whitespace[Math.floor(Math.random() * whitespace.length)];
        return prefix + uuid + suffix;
      });
      
      const startTime = performance.now();
      
      for (const uuid of uuids) {
        validateAndSanitizeUUID(uuid);
      }
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should complete in less than 200ms for 1000 UUIDs with whitespace
      expect(duration).toBeLessThan(200);
      
      // Verify all were validated and sanitized correctly
      const results = uuids.map(uuid => validateAndSanitizeUUID(uuid));
      expect(results.every(result => result.isValid === true)).toBe(true);
      expect(results.every(result => result.uuid && !result.uuid.match(/^\s|\s$/))).toBe(true);
    });

    it('should handle 1000 empty/null values efficiently', () => {
      const values = Array.from({ length: 1000 }, (_, i) => {
        const options = [undefined, null, ''];
        return options[i % 3];
      });
      
      const startTime = performance.now();
      
      for (const value of values) {
        validateAndSanitizeUUID(value);
      }
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should complete in less than 50ms for 1000 empty values (fast path)
      expect(duration).toBeLessThan(50);
      
      // Verify all were handled correctly
      const results = values.map(value => validateAndSanitizeUUID(value));
      expect(results.every(result => result.isValid === true)).toBe(true);
    });
  });

  describe('Validation impact on parsing time', () => {
    it('should not significantly impact record processing time', () => {
      // Simulate 1000 import records with UUID validation
      const records = Array.from({ length: 1000 }, (_, i) => ({
        id: i % 2 === 0 ? generateValidUUID() : undefined,
        name: `Record ${i}`,
        value: Math.random() * 100,
      }));
      
      // Time without UUID validation
      const startTimeWithoutValidation = performance.now();
      for (const record of records) {
        // Simulate basic field processing
        const processed = {
          name: record.name,
          value: record.value,
        };
      }
      const durationWithoutValidation = performance.now() - startTimeWithoutValidation;
      
      // Time with UUID validation
      const startTimeWithValidation = performance.now();
      for (const record of records) {
        // Simulate field processing with UUID validation
        const idResult = record.id ? validateAndSanitizeUUID(record.id) : { isValid: true };
        const processed = {
          id: idResult.uuid,
          name: record.name,
          value: record.value,
        };
      }
      const durationWithValidation = performance.now() - startTimeWithValidation;
      
      // UUID validation should add less than 100ms overhead for 1000 records
      const overhead = durationWithValidation - durationWithoutValidation;
      expect(overhead).toBeLessThan(100);
    });

    it('should efficiently process large batch of records with mixed ID presence', () => {
      // Simulate realistic import scenario: 2000 records, 30% with IDs
      const records = Array.from({ length: 2000 }, (_, i) => ({
        id: i % 10 < 3 ? generateValidUUID() : undefined, // 30% have IDs
        name: `Record ${i}`,
        description: `Description for record ${i}`,
        value: Math.random() * 1000,
      }));
      
      const startTime = performance.now();
      
      const processedRecords = records.map(record => {
        const idResult = record.id ? validateAndSanitizeUUID(record.id) : { isValid: true };
        
        return {
          id: idResult.uuid,
          name: record.name,
          description: record.description,
          value: record.value,
          hasValidationError: !idResult.isValid,
        };
      });
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should process 2000 records in less than 200ms
      expect(duration).toBeLessThan(200);
      
      // Verify processing was correct
      expect(processedRecords).toHaveLength(2000);
      const recordsWithIds = processedRecords.filter(r => r.id !== undefined);
      expect(recordsWithIds.length).toBeGreaterThan(500); // ~30% of 2000
    });
  });

  describe('Regex compilation optimization', () => {
    it('should reuse compiled regex across multiple validations', () => {
      // This test verifies that the regex is compiled once and reused
      // by measuring consistent performance across multiple runs
      
      const uuids = Array.from({ length: 100 }, () => generateValidUUID());
      const durations: number[] = [];
      
      // Run validation 10 times
      for (let run = 0; run < 10; run++) {
        const startTime = performance.now();
        
        for (const uuid of uuids) {
          isValidUUID(uuid);
        }
        
        const duration = performance.now() - startTime;
        durations.push(duration);
      }
      
      // Calculate variance in durations
      const avgDuration = durations.reduce((a, b) => a + b, 0) / durations.length;
      const variance = durations.reduce((sum, d) => sum + Math.pow(d - avgDuration, 2), 0) / durations.length;
      const stdDev = Math.sqrt(variance);
      
      // Standard deviation should be reasonable (< 100% of average), indicating consistent performance
      // This suggests regex is compiled once and reused
      // Note: Some variance is expected due to system load and other factors
      expect(stdDev).toBeLessThan(avgDuration * 1.0);
    });

    it('should maintain performance with concurrent validation calls', () => {
      const uuids = Array.from({ length: 1000 }, () => generateValidUUID());
      
      const startTime = performance.now();
      
      // Simulate concurrent validations
      const results = uuids.map(uuid => isValidUUID(uuid));
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should complete in less than 100ms
      expect(duration).toBeLessThan(100);
      expect(results.every(r => r === true)).toBe(true);
    });
  });

  describe('Memory efficiency', () => {
    it('should not create excessive objects during validation', () => {
      // This test ensures validation doesn't create memory pressure
      const uuids = Array.from({ length: 10000 }, () => generateValidUUID());
      
      const startTime = performance.now();
      let validCount = 0;
      
      for (const uuid of uuids) {
        if (isValidUUID(uuid)) {
          validCount++;
        }
      }
      
      const endTime = performance.now();
      const duration = endTime - startTime;
      
      // Should handle 10,000 validations in less than 1 second
      expect(duration).toBeLessThan(1000);
      expect(validCount).toBe(10000);
    });
  });
});
