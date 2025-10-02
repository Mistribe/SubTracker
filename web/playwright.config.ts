import { defineConfig, devices } from '@playwright/test';
import * as dotenv from 'dotenv';
import * as path from 'path';

// Set NODE_ENV to test if not already set when running Playwright
if (!process.env.NODE_ENV) {
  process.env.NODE_ENV = 'test';
}

// Load environment-specific configuration
const envFile = process.env.NODE_ENV === 'test' ? '.env.test' : '.env';
dotenv.config({ path: envFile });

// Load local overrides if they exist
dotenv.config({ path: `${envFile}.local`, override: true });

/**
 * Playwright Configuration for SubTracker E2E Tests
 * 
 * This configuration supports multiple environments:
 * - Local development (default)
 * - Test environment (NODE_ENV=test)
 * - CI/CD environments
 */
export default defineConfig({
  // Test directory structure
  testDir: './tests/e2e',
  
  // Global test configuration
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 1,
  workers: process.env.CI ? 1 : undefined,
  
  // Test execution settings
  timeout: 30 * 1000, // 30 seconds per test
  expect: {
    timeout: 10 * 1000, // 10 seconds for assertions
  },
  
  // Global setup and teardown
  globalSetup: './tests/config/global-setup.ts',
  globalTeardown: './tests/config/global-teardown.ts',
  
  // Test output and reporting
  outputDir: 'test-results/artifacts',
  reporter: [
    // Console reporter for development
    ['list'],
    // HTML reporter for detailed results
    ['html', { 
      outputFolder: 'test-results/html-report',
      open: process.env.CI ? 'never' : 'on-failure'
    }],
    // JUnit reporter for CI/CD integration
    ['junit', { 
      outputFile: 'test-results/junit-report.xml'
    }],
    // JSON reporter for custom processing
    ['json', {
      outputFile: 'test-results/test-results.json'
    }]
  ],
  
  // Global test configuration
  use: {
    // Base URL for the application
    baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:5173',
    
    // Browser context settings
    viewport: { width: 1280, height: 720 },
    ignoreHTTPSErrors: true,
    
    // Action timeouts
    actionTimeout: 15 * 1000, // 15 seconds for actions
    navigationTimeout: 30 * 1000, // 30 seconds for navigation
    
    // Test artifacts
    screenshot: 'only-on-failure',
    video: 'retain-on-failure',
    trace: 'retain-on-failure',
    
    // Additional context options
    locale: 'en-US',
    timezoneId: 'America/New_York',
    
    // Extra HTTP headers
    extraHTTPHeaders: {
      'Accept-Language': 'en-US,en;q=0.9',
    },
  },

  // Browser projects configuration - Chromium only for faster execution
  projects: [
    // Setup project for authentication and global setup
    {
      name: 'setup',
      testMatch: /.*\.setup\.ts/,
      teardown: 'cleanup',
    },
    
    // Cleanup project
    {
      name: 'cleanup',
      testMatch: /.*\.teardown\.ts/,
    },

    // Desktop Chrome (Primary browser for all tests)
    {
      name: 'chromium',
      use: { 
        ...devices['Desktop Chrome'],
        // Enable additional Chrome features for testing
        launchOptions: {
          args: [
            '--disable-web-security',
            '--disable-features=VizDisplayCompositor',
            '--disable-dev-shm-usage',
            '--no-sandbox',
          ],
        },
      },
      dependencies: ['setup'],
    },
  ],

  // Development server configuration
  webServer: process.env.CI ? undefined : {
    command: 'npm run dev',
    url: 'http://localhost:5173',
    reuseExistingServer: !process.env.CI,
    timeout: 120 * 1000, // 2 minutes to start dev server
    env: {
      NODE_ENV: 'development',
    },
  },

  // Environment-specific overrides
  ...(process.env.NODE_ENV === 'test' && {
    // Test environment specific settings
    use: {
      baseURL: process.env.PLAYWRIGHT_BASE_URL || 'http://localhost:5173',
      screenshot: 'only-on-failure',
      video: 'retain-on-failure',
      trace: 'retain-on-failure',
    },
    
    // Chromium only for faster test execution in test env
    projects: [
      {
        name: 'setup',
        testMatch: /.*\.setup\.ts/,
        teardown: 'cleanup',
      },
      {
        name: 'cleanup',
        testMatch: /.*\.teardown\.ts/,
      },
      {
        name: 'chromium',
        use: { 
          ...devices['Desktop Chrome'],
          launchOptions: {
            args: [
              '--disable-web-security',
              '--disable-features=VizDisplayCompositor',
              '--disable-dev-shm-usage',
              '--no-sandbox',
            ],
          },
        },
        dependencies: ['setup'],
      },
    ],
  }),

  // CI-specific overrides
  ...(process.env.CI && {
    // CI environment specific settings
    workers: 2, // Limit workers in CI
    retries: 3, // More retries in CI due to potential flakiness
    
    use: {
      screenshot: 'only-on-failure',
      video: 'retain-on-failure',
      trace: 'retain-on-failure',
    },
    
    // Disable dev server in CI (assume app is already running)
    webServer: undefined,
  }),
});