/**
 * UUID Validation Utility
 * 
 * Provides RFC 4122 compliant UUID validation and sanitization functions.
 */

/**
 * RFC 4122 compliant UUID regex pattern
 * Matches UUID versions 1-5 with proper format:
 * xxxxxxxx-xxxx-Mxxx-Nxxx-xxxxxxxxxxxx
 * where M is the version (1-5) and N is the variant (8, 9, a, or b)
 */
const UUID_REGEX = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i;

/**
 * Checks if a value is empty, null, or undefined
 * @param value - The value to check
 * @returns true if the value is empty, null, or undefined
 */
export function isEmptyValue(value: any): boolean {
  return value === undefined || value === null || value === '';
}

/**
 * Validates if a string is a valid UUID format (RFC 4122)
 * @param uuid - The UUID string to validate
 * @returns true if the UUID is valid, false otherwise
 */
export function isValidUUID(uuid: string): boolean {
  return UUID_REGEX.test(uuid);
}

/**
 * Sanitizes a UUID string by trimming whitespace
 * @param uuid - The UUID string to sanitize
 * @returns The sanitized UUID string
 */
export function sanitizeUUID(uuid: string): string {
  return uuid.trim();
}

/**
 * Validates and sanitizes a UUID string
 * @param rawValue - The raw value to validate (can be any type)
 * @returns An object containing the sanitized UUID (if valid) or an error message
 */
export function validateAndSanitizeUUID(rawValue: any): {
  isValid: boolean;
  uuid?: string;
  error?: string;
} {
  // If no value provided, that's valid - backend will generate one
  if (isEmptyValue(rawValue)) {
    return { isValid: true };
  }

  // Convert to string and sanitize
  const uuidString = sanitizeUUID(String(rawValue));

  // Validate UUID format
  if (!isValidUUID(uuidString)) {
    return {
      isValid: false,
      error: generateValidationErrorMessage(uuidString)
    };
  }

  return {
    isValid: true,
    uuid: uuidString
  };
}

/**
 * Generates a user-friendly validation error message for an invalid UUID
 * @param invalidUUID - The invalid UUID string
 * @returns A formatted error message
 */
export function generateValidationErrorMessage(invalidUUID: string): string {
  return `Invalid UUID format: "${invalidUUID}". Expected format: xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx (e.g., 550e8400-e29b-41d4-a716-446655440000)`;
}
