/**
 * Combined Test Fixtures
 * 
 * Exports all test fixtures in a single, convenient interface
 * Combines authentication, test data, and utility fixtures
 */

import { test as authTest } from './auth';
import { test as dataTest } from './test-data';
import type { AuthFixtures } from './auth';
import type { TestDataFixtures } from './test-data';
import { ClerkTestUserManager } from '../utils/clerk-user-manager';
import { SessionManager } from '../utils/session-manager';
import { TestDataCleanup } from './test-data';

/**
 * Combined fixture interface
 */
export interface CombinedFixtures extends AuthFixtures, TestDataFixtures {
  // Additional combined utilities can be added here
}

/**
 * Combined test fixture with authentication and test data
 * 
 * This combines both authentication fixtures and test data fixtures
 * into a single test instance that can be used throughout the test suite
 */
export const test = authTest.extend<TestDataFixtures>({
  // Unique ID fixture for test isolation
  uniqueId: async ({}, use) => {
    const uniqueId = `test-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    await use(uniqueId);
  },

  // Test subscription fixture
  testSubscription: async ({ uniqueId }, use) => {
    const subscription = {
      name: `Test Subscription ${uniqueId}`,
      providerId: `provider-${uniqueId}`,
      amount: Math.floor(Math.random() * 100) + 10,
      currency: 'USD',
      billingCycle: Math.random() > 0.5 ? 'monthly' : 'yearly' as 'monthly' | 'yearly',
      nextBillingDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(),
      labels: [`label-${uniqueId}`],
    };
    await use(subscription);
  },

  // Test provider fixture
  testProvider: async ({ uniqueId }, use) => {
    const providers = ['Netflix', 'Spotify', 'Adobe', 'Microsoft', 'Google'];
    const randomProvider = providers[Math.floor(Math.random() * providers.length)];
    
    const provider = {
      name: `${randomProvider} Test ${uniqueId}`,
      description: `Test provider for ${randomProvider} - ${uniqueId}`,
      website: `https://test-${uniqueId}.example.com`,
      logoUrl: `https://test-${uniqueId}.example.com/logo.png`,
    };
    await use(provider);
  },

  // Test label fixture
  testLabel: async ({ uniqueId }, use) => {
    const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD'];
    const categories = ['Entertainment', 'Productivity', 'Development', 'Design', 'Business'];
    
    const label = {
      name: `${categories[Math.floor(Math.random() * categories.length)]} ${uniqueId}`,
      color: colors[Math.floor(Math.random() * colors.length)],
      description: `Test label for categorization - ${uniqueId}`,
    };
    await use(label);
  },

  // Test family fixture
  testFamily: async ({ uniqueId }, use) => {
    const family = {
      name: `Test Family ${uniqueId}`,
      ownerId: 'default-test-user',
      members: [
        {
          email: `member1-${uniqueId}@example.com`,
          role: 'member' as const,
          permissions: ['view_subscriptions'],
          status: 'accepted' as const,
        },
        {
          email: `member2-${uniqueId}@example.com`,
          role: 'member' as const,
          permissions: ['view_subscriptions', 'edit_subscriptions'],
          status: 'pending' as const,
        },
      ],
    };
    await use(family);
  },

  // Test data cleanup fixture
  testDataCleanup: async ({ uniqueId }, use) => {
    const cleanup = async () => {
      console.log(`Cleaning up test data for ${uniqueId}`);
      // In a real implementation, this would call API endpoints to delete data
    };

    await use(cleanup);
    
    // Always run cleanup after test, even if test fails
    await cleanup();
  },
});

/**
 * Export expect for convenience
 */
export { expect } from '@playwright/test';

/**
 * Export individual fixture components
 */
export type { AuthFixtures } from './auth';
export type { TestDataFixtures } from './test-data';
export { ClerkTestUserManager } from '../utils/clerk-user-manager';
export { SessionManager } from '../utils/session-manager';
export { TestDataCleanup } from './test-data';

/**
 * Global test setup and teardown utilities
 */
export class GlobalTestSetup {
  /**
   * Initialize all test utilities
   */
  static async initialize(): Promise<void> {
    await ClerkTestUserManager.initialize();
    SessionManager.clearAllSessions();
  }

  /**
   * Cleanup all test resources
   */
  static async cleanup(): Promise<void> {
    await TestDataCleanup.cleanupAllTestData();
    await ClerkTestUserManager.cleanupAllTestUsers();
    SessionManager.clearAllSessions();
  }

  /**
   * Get test status summary
   */
  static getTestSummary() {
    return {
      testData: TestDataCleanup.getCleanupSummary(),
      activeSessions: SessionManager.getActiveSessions().length,
    };
  }
}

/**
 * Test utilities for common operations
 */
export class TestUtils {
  /**
   * Wait for element with retry logic
   */
  static async waitForElementWithRetry(
    page: any,
    selector: string,
    options: { timeout?: number; retries?: number } = {}
  ): Promise<void> {
    const { timeout = 10000, retries = 3 } = options;
    
    for (let i = 0; i < retries; i++) {
      try {
        await page.waitForSelector(selector, { timeout });
        return;
      } catch (error) {
        if (i === retries - 1) throw error;
        await page.waitForTimeout(1000); // Wait 1 second before retry
      }
    }
  }

  /**
   * Fill form with retry logic
   */
  static async fillFormWithRetry(
    page: any,
    formData: Record<string, string>,
    options: { timeout?: number; retries?: number } = {}
  ): Promise<void> {
    const { timeout = 5000, retries = 3 } = options;
    
    for (const [selector, value] of Object.entries(formData)) {
      for (let i = 0; i < retries; i++) {
        try {
          await page.fill(selector, value, { timeout });
          break;
        } catch (error) {
          if (i === retries - 1) throw error;
          await page.waitForTimeout(500);
        }
      }
    }
  }

  /**
   * Click element with retry logic
   */
  static async clickWithRetry(
    page: any,
    selector: string,
    options: { timeout?: number; retries?: number } = {}
  ): Promise<void> {
    const { timeout = 5000, retries = 3 } = options;
    
    for (let i = 0; i < retries; i++) {
      try {
        await page.click(selector, { timeout });
        return;
      } catch (error) {
        if (i === retries - 1) throw error;
        await page.waitForTimeout(500);
      }
    }
  }

  /**
   * Navigate with loading state verification
   */
  static async navigateAndWait(
    page: any,
    url: string,
    options: { waitUntil?: string; timeout?: number } = {}
  ): Promise<void> {
    const { waitUntil = 'networkidle', timeout = 30000 } = options;
    
    await page.goto(url, { waitUntil, timeout });
    await page.waitForLoadState('domcontentloaded');
  }

  /**
   * Take screenshot with timestamp
   */
  static async takeTimestampedScreenshot(
    page: any,
    name: string,
    options: { fullPage?: boolean } = {}
  ): Promise<void> {
    const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
    const filename = `${name}-${timestamp}.png`;
    
    await page.screenshot({
      path: `test-results/${filename}`,
      fullPage: options.fullPage || false,
    });
  }

  /**
   * Generate test email
   */
  static generateTestEmail(prefix: string = 'test'): string {
    const timestamp = Date.now();
    const random = Math.random().toString(36).substr(2, 9);
    return `${prefix}-${timestamp}-${random}@example.com`;
  }

  /**
   * Generate test phone number
   */
  static generateTestPhone(): string {
    const areaCode = Math.floor(Math.random() * 900) + 100;
    const exchange = Math.floor(Math.random() * 900) + 100;
    const number = Math.floor(Math.random() * 9000) + 1000;
    return `+1${areaCode}${exchange}${number}`;
  }

  /**
   * Format currency for testing
   */
  static formatCurrency(amount: number, currency: string = 'USD'): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(amount);
  }

  /**
   * Generate future date
   */
  static generateFutureDate(daysFromNow: number = 30): string {
    const date = new Date();
    date.setDate(date.getDate() + daysFromNow);
    return date.toISOString().split('T')[0]; // YYYY-MM-DD format
  }

  /**
   * Wait for API response
   */
  static async waitForApiResponse(
    page: any,
    urlPattern: string | RegExp,
    options: { timeout?: number; method?: string } = {}
  ): Promise<any> {
    const { timeout = 10000, method = 'GET' } = options;
    
    return await page.waitForResponse(
      (response: any) => {
        const url = response.url();
        const matchesUrl = typeof urlPattern === 'string' 
          ? url.includes(urlPattern)
          : urlPattern.test(url);
        
        return matchesUrl && response.request().method() === method;
      },
      { timeout }
    );
  }
}