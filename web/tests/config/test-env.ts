/**
 * Test Environment Configuration
 * 
 * Centralizes environment-specific settings and validation
 * for Playwright E2E tests
 */

export interface TestEnvironment {
  baseUrl: string;
  apiUrl: string;
  clerkPublishableKey: string;
  clerkSecretKey?: string;
  clerkTestingToken?: string;
  testUser: {
    email: string;
    password: string;
  };
  familyOwner?: {
    email: string;
    password: string;
  };
  familyMember?: {
    email: string;
    password: string;
  };
  timeouts: {
    test: number;
    expect: number;
    navigation: number;
    action: number;
  };
  retries: number;
  workers?: number;
}

/**
 * Load and validate test environment configuration
 */
export function getTestEnvironment(): TestEnvironment {
  // Validate required environment variables
  const requiredVars = [
    'VITE_CLERK_PUBLISHABLE_KEY',
    'TEST_USER_EMAIL',
    'TEST_USER_PASSWORD',
  ];
  
  const missingVars = requiredVars.filter(varName => !process.env[varName]);
  
  if (missingVars.length > 0) {
    throw new Error(
      `Missing required environment variables: ${missingVars.join(', ')}\n` +
      'Please check your .env.test file or environment configuration.'
    );
  }
  
  return {
    baseUrl: process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:5173',
    apiUrl: process.env.PLAYWRIGHT_API_URL || 'http://localhost:8080',
    clerkPublishableKey: process.env.VITE_CLERK_PUBLISHABLE_KEY!,
    clerkSecretKey: process.env.CLERK_SECRET_KEY,
    clerkTestingToken: process.env.CLERK_TESTING_TOKEN,
    testUser: {
      email: process.env.TEST_USER_EMAIL!,
      password: process.env.TEST_USER_PASSWORD!,
    },
    familyOwner: process.env.FAMILY_OWNER_EMAIL ? {
      email: process.env.FAMILY_OWNER_EMAIL,
      password: process.env.FAMILY_OWNER_PASSWORD || process.env.TEST_USER_PASSWORD!,
    } : undefined,
    familyMember: process.env.FAMILY_MEMBER_EMAIL ? {
      email: process.env.FAMILY_MEMBER_EMAIL,
      password: process.env.FAMILY_MEMBER_PASSWORD || process.env.TEST_USER_PASSWORD!,
    } : undefined,
    timeouts: {
      test: 30 * 1000, // 30 seconds
      expect: 10 * 1000, // 10 seconds
      navigation: 30 * 1000, // 30 seconds
      action: 15 * 1000, // 15 seconds
    },
    retries: process.env.CI ? 3 : 1,
    workers: process.env.CI ? 2 : undefined,
  };
}

/**
 * Environment-specific configuration presets
 */
export const environments = {
  development: {
    baseUrl: 'http://localhost:5173',
    apiUrl: 'http://localhost:8080',
  },
  staging: {
    baseUrl: process.env.STAGING_URL || 'https://staging.subtracker.app',
    apiUrl: process.env.STAGING_API_URL || 'https://api-staging.subtracker.app',
  },
  production: {
    baseUrl: process.env.PRODUCTION_URL || 'https://subtracker.app',
    apiUrl: process.env.PRODUCTION_API_URL || 'https://api.subtracker.app',
  },
} as const;

/**
 * Get environment-specific URLs
 */
export function getEnvironmentUrls(env: keyof typeof environments = 'development') {
  return environments[env];
}

/**
 * Check if running in CI environment
 */
export function isCI(): boolean {
  return !!process.env.CI;
}

/**
 * Check if running in development mode
 */
export function isDevelopment(): boolean {
  return process.env.NODE_ENV !== 'test' && process.env.NODE_ENV !== 'production';
}

/**
 * Get browser-specific configuration
 */
export function getBrowserConfig(browserName: string) {
  const baseConfig = {
    headless: process.env.CI ? true : false,
    slowMo: process.env.CI ? 0 : 100,
  };
  
  switch (browserName) {
    case 'chromium':
      return {
        ...baseConfig,
        args: [
          '--disable-web-security',
          '--disable-features=VizDisplayCompositor',
          '--disable-dev-shm-usage',
          '--no-sandbox',
        ],
      };
    case 'firefox':
      return {
        ...baseConfig,
        firefoxUserPrefs: {
          'security.tls.insecure_fallback_hosts': 'localhost',
        },
      };
    case 'webkit':
      return baseConfig;
    default:
      return baseConfig;
  }
}