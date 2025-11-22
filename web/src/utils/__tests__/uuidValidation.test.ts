import { describe, it, expect } from 'vitest';
import {
  isEmptyValue,
  isValidUUID,
  sanitizeUUID,
  validateAndSanitizeUUID,
  generateValidationErrorMessage
} from '../uuidValidation';

describe('uuidValidation', () => {
  describe('isEmptyValue', () => {
    it('should return true for undefined', () => {
      expect(isEmptyValue(undefined)).toBe(true);
    });

    it('should return true for null', () => {
      expect(isEmptyValue(null)).toBe(true);
    });

    it('should return true for empty string', () => {
      expect(isEmptyValue('')).toBe(true);
    });

    it('should return false for non-empty string', () => {
      expect(isEmptyValue('test')).toBe(false);
    });

    it('should return false for number', () => {
      expect(isEmptyValue(0)).toBe(false);
      expect(isEmptyValue(123)).toBe(false);
    });

    it('should return false for boolean', () => {
      expect(isEmptyValue(false)).toBe(false);
      expect(isEmptyValue(true)).toBe(false);
    });
  });

  describe('isValidUUID', () => {
    describe('valid UUID formats', () => {
      it('should validate UUID v1 format', () => {
        const uuidV1 = '550e8400-e29b-11d4-a716-446655440000';
        expect(isValidUUID(uuidV1)).toBe(true);
      });

      it('should validate UUID v4 format', () => {
        const uuidV4 = '550e8400-e29b-41d4-a716-446655440000';
        expect(isValidUUID(uuidV4)).toBe(true);
      });

      it('should validate UUID v5 format', () => {
        const uuidV5 = '550e8400-e29b-51d4-a716-446655440000';
        expect(isValidUUID(uuidV5)).toBe(true);
      });

      it('should validate UUID with uppercase letters', () => {
        const uuid = '550E8400-E29B-41D4-A716-446655440000';
        expect(isValidUUID(uuid)).toBe(true);
      });

      it('should validate UUID with mixed case', () => {
        const uuid = '550e8400-E29B-41d4-A716-446655440000';
        expect(isValidUUID(uuid)).toBe(true);
      });

      it('should validate UUID with variant 8', () => {
        const uuid = '550e8400-e29b-41d4-8716-446655440000';
        expect(isValidUUID(uuid)).toBe(true);
      });

      it('should validate UUID with variant 9', () => {
        const uuid = '550e8400-e29b-41d4-9716-446655440000';
        expect(isValidUUID(uuid)).toBe(true);
      });

      it('should validate UUID with variant a', () => {
        const uuid = '550e8400-e29b-41d4-a716-446655440000';
        expect(isValidUUID(uuid)).toBe(true);
      });

      it('should validate UUID with variant b', () => {
        const uuid = '550e8400-e29b-41d4-b716-446655440000';
        expect(isValidUUID(uuid)).toBe(true);
      });
    });

    describe('invalid UUID formats', () => {
      it('should reject UUID with wrong length (too short)', () => {
        const uuid = '550e8400-e29b-41d4-a716-44665544000';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject UUID with wrong length (too long)', () => {
        const uuid = '550e8400-e29b-41d4-a716-4466554400000';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject UUID with invalid characters', () => {
        const uuid = '550e8400-e29b-41d4-a716-44665544000g';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject UUID with missing hyphens', () => {
        const uuid = '550e8400e29b41d4a716446655440000';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject UUID with wrong hyphen positions', () => {
        const uuid = '550e8400-e29-b41d4-a716-446655440000';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject UUID with invalid version (0)', () => {
        const uuid = '550e8400-e29b-01d4-a716-446655440000';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject UUID with invalid version (6)', () => {
        const uuid = '550e8400-e29b-61d4-a716-446655440000';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject UUID with invalid variant (c)', () => {
        const uuid = '550e8400-e29b-41d4-c716-446655440000';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject UUID with invalid variant (7)', () => {
        const uuid = '550e8400-e29b-41d4-7716-446655440000';
        expect(isValidUUID(uuid)).toBe(false);
      });

      it('should reject empty string', () => {
        expect(isValidUUID('')).toBe(false);
      });

      it('should reject random string', () => {
        expect(isValidUUID('not-a-uuid')).toBe(false);
      });

      it('should reject UUID with spaces', () => {
        const uuid = '550e8400 e29b 41d4 a716 446655440000';
        expect(isValidUUID(uuid)).toBe(false);
      });
    });
  });

  describe('sanitizeUUID', () => {
    it('should trim leading whitespace', () => {
      const uuid = '  550e8400-e29b-41d4-a716-446655440000';
      expect(sanitizeUUID(uuid)).toBe('550e8400-e29b-41d4-a716-446655440000');
    });

    it('should trim trailing whitespace', () => {
      const uuid = '550e8400-e29b-41d4-a716-446655440000  ';
      expect(sanitizeUUID(uuid)).toBe('550e8400-e29b-41d4-a716-446655440000');
    });

    it('should trim both leading and trailing whitespace', () => {
      const uuid = '  550e8400-e29b-41d4-a716-446655440000  ';
      expect(sanitizeUUID(uuid)).toBe('550e8400-e29b-41d4-a716-446655440000');
    });

    it('should trim tabs and newlines', () => {
      const uuid = '\t550e8400-e29b-41d4-a716-446655440000\n';
      expect(sanitizeUUID(uuid)).toBe('550e8400-e29b-41d4-a716-446655440000');
    });

    it('should not modify UUID without whitespace', () => {
      const uuid = '550e8400-e29b-41d4-a716-446655440000';
      expect(sanitizeUUID(uuid)).toBe('550e8400-e29b-41d4-a716-446655440000');
    });
  });

  describe('validateAndSanitizeUUID', () => {
    describe('empty/null/undefined values', () => {
      it('should return valid for undefined', () => {
        const result = validateAndSanitizeUUID(undefined);
        expect(result.isValid).toBe(true);
        expect(result.uuid).toBeUndefined();
        expect(result.error).toBeUndefined();
      });

      it('should return valid for null', () => {
        const result = validateAndSanitizeUUID(null);
        expect(result.isValid).toBe(true);
        expect(result.uuid).toBeUndefined();
        expect(result.error).toBeUndefined();
      });

      it('should return valid for empty string', () => {
        const result = validateAndSanitizeUUID('');
        expect(result.isValid).toBe(true);
        expect(result.uuid).toBeUndefined();
        expect(result.error).toBeUndefined();
      });
    });

    describe('valid UUIDs', () => {
      it('should return valid result with sanitized UUID', () => {
        const uuid = '550e8400-e29b-41d4-a716-446655440000';
        const result = validateAndSanitizeUUID(uuid);
        expect(result.isValid).toBe(true);
        expect(result.uuid).toBe(uuid);
        expect(result.error).toBeUndefined();
      });

      it('should sanitize and validate UUID with whitespace', () => {
        const uuid = '  550e8400-e29b-41d4-a716-446655440000  ';
        const result = validateAndSanitizeUUID(uuid);
        expect(result.isValid).toBe(true);
        expect(result.uuid).toBe('550e8400-e29b-41d4-a716-446655440000');
        expect(result.error).toBeUndefined();
      });

      it('should handle UUID passed as non-string type', () => {
        const uuid = '550e8400-e29b-41d4-a716-446655440000';
        const result = validateAndSanitizeUUID(uuid as any);
        expect(result.isValid).toBe(true);
        expect(result.uuid).toBe(uuid);
      });
    });

    describe('invalid UUIDs', () => {
      it('should return invalid result with error message for wrong format', () => {
        const invalidUuid = 'not-a-uuid';
        const result = validateAndSanitizeUUID(invalidUuid);
        expect(result.isValid).toBe(false);
        expect(result.uuid).toBeUndefined();
        expect(result.error).toBeDefined();
        expect(result.error).toContain('Invalid UUID format');
        expect(result.error).toContain(invalidUuid);
      });

      it('should return invalid result for UUID with wrong length', () => {
        const invalidUuid = '550e8400-e29b-41d4-a716-44665544000';
        const result = validateAndSanitizeUUID(invalidUuid);
        expect(result.isValid).toBe(false);
        expect(result.error).toBeDefined();
      });

      it('should return invalid result for UUID with invalid characters', () => {
        const invalidUuid = '550e8400-e29b-41d4-a716-44665544000g';
        const result = validateAndSanitizeUUID(invalidUuid);
        expect(result.isValid).toBe(false);
        expect(result.error).toBeDefined();
      });

      it('should sanitize before validating (whitespace should not help invalid UUID)', () => {
        const invalidUuid = '  not-a-uuid  ';
        const result = validateAndSanitizeUUID(invalidUuid);
        expect(result.isValid).toBe(false);
        expect(result.error).toContain('not-a-uuid');
      });
    });
  });

  describe('generateValidationErrorMessage', () => {
    it('should generate error message with invalid UUID', () => {
      const invalidUuid = 'not-a-uuid';
      const message = generateValidationErrorMessage(invalidUuid);
      expect(message).toContain('Invalid UUID format');
      expect(message).toContain(invalidUuid);
      expect(message).toContain('xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx');
    });

    it('should include example UUID in error message', () => {
      const invalidUuid = '123';
      const message = generateValidationErrorMessage(invalidUuid);
      expect(message).toContain('550e8400-e29b-41d4-a716-446655440000');
    });

    it('should show the provided invalid value in quotes', () => {
      const invalidUuid = 'bad-uuid';
      const message = generateValidationErrorMessage(invalidUuid);
      expect(message).toContain('"bad-uuid"');
    });
  });
});
