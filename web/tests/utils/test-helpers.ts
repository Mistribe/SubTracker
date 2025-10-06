import { createApiTestHelpers, ApiTestHelpers } from './api-client';
import { createMockApiTestHelpers, MockApiTestHelpers } from './mock-api-client';
import { SubscriptionData, ProviderData, LabelData } from './data-generators';

/**
 * Unified test helpers interface
 */
export interface TestHelpers {
  createTestProvider(data: ProviderData): Promise<{ id: string; name: string }>;
  createTestSubscription(data: SubscriptionData): Promise<string>;
  createTestLabel(data: LabelData): Promise<string>;
  cleanupTestSubscription(id: string): Promise<void>;
  cleanupTestProvider(id: string): Promise<void>;
  cleanupTestLabel(id: string): Promise<void>;
  isUsingMockApi(): boolean;
}

/**
 * Wrapper for real API helpers
 */
class RealApiTestHelpers implements TestHelpers {
  constructor(private apiHelpers: ApiTestHelpers) { }

  async createTestProvider(data: ProviderData): Promise<{ id: string; name: string }> {
    return this.apiHelpers.createTestProvider(data);
  }

  async createTestSubscription(data: SubscriptionData): Promise<string> {
    return this.apiHelpers.createTestSubscription(data);
  }

  async createTestLabel(data: LabelData): Promise<string> {
    return this.apiHelpers.createTestLabel(data);
  }

  async cleanupTestSubscription(id: string): Promise<void> {
    return this.apiHelpers.cleanupTestSubscription(id);
  }

  async cleanupTestProvider(id: string): Promise<void> {
    return this.apiHelpers.cleanupTestProvider(id);
  }

  async cleanupTestLabel(id: string): Promise<void> {
    return this.apiHelpers.cleanupTestLabel(id);
  }

  isUsingMockApi(): boolean {
    return false;
  }
}

/**
 * Wrapper for mock API helpers
 */
class MockTestHelpers implements TestHelpers {
  constructor(private mockHelpers: MockApiTestHelpers) { }

  async createTestProvider(data: ProviderData): Promise<{ id: string; name: string }> {
    return this.mockHelpers.createTestProvider(data);
  }

  async createTestSubscription(data: SubscriptionData): Promise<string> {
    return this.mockHelpers.createTestSubscription(data);
  }

  async createTestLabel(data: LabelData): Promise<string> {
    return this.mockHelpers.createTestLabel(data);
  }

  async cleanupTestSubscription(id: string): Promise<void> {
    return this.mockHelpers.cleanupTestSubscription(id);
  }

  async cleanupTestProvider(id: string): Promise<void> {
    return this.mockHelpers.cleanupTestProvider(id);
  }

  async cleanupTestLabel(id: string): Promise<void> {
    return this.mockHelpers.cleanupTestLabel(id);
  }

  isUsingMockApi(): boolean {
    return true;
  }
}

/**
 * Factory function to create test helpers - requires real API
 */
export async function createTestHelpers(page?: any): Promise<TestHelpers> {
  let authToken: string | undefined;

  // Extract auth token from page if available
  if (page) {
    try {
      // Wait for page to be fully loaded and authenticated
      await page.waitForLoadState('networkidle');
      await page.waitForTimeout(2000); // Give time for Clerk to initialize

      // Try to get the Clerk session token from the page
      authToken = await page.evaluate(async () => {
        // Wait for Clerk to be available
        let attempts = 0;
        while (!window.Clerk && attempts < 10) {
          await new Promise(resolve => setTimeout(resolve, 500));
          attempts++;
        }

        // Method 1: Get token from Clerk session
        if (window.Clerk && window.Clerk.session) {
          try {
            const token = await window.Clerk.session.getToken();
            if (token) return token;
          } catch (e) {
            console.log('Failed to get token from Clerk.session.getToken():', e);
          }
        }

        // Method 2: Check cookies for session token
        const cookies = document.cookie.split(';');
        for (const cookie of cookies) {
          const [name, value] = cookie.trim().split('=');
          if (name && name.includes('__session') && value) {
            const decodedValue = decodeURIComponent(value);
            if (decodedValue.includes('.') && decodedValue.split('.').length === 3) {
              return decodedValue;
            }
          }
        }

        // Method 3: Check for any JWT-like tokens in cookies
        for (const cookie of cookies) {
          const [name, value] = cookie.trim().split('=');
          if (value && value.includes('.') && value.split('.').length === 3) {
            return decodeURIComponent(value);
          }
        }

        return null;
      });

      if (!authToken) {
        throw new Error('Could not extract authentication token from page. User may not be properly authenticated.');
      }

    } catch (error) {
      console.error('❌ Failed to extract auth token from page:', error);
      throw new Error(`Authentication token extraction failed: ${error}. Please ensure user is properly authenticated.`);
    }
  }

  try {
    // Create real API helpers with auth token
    const apiHelpers = await createApiTestHelpers(undefined, authToken);
    console.log('✅ Using real API for tests');
    return new RealApiTestHelpers(apiHelpers);
  } catch (error) {
    console.error('❌ Failed to connect to API:', error);
    throw new Error(`API connection failed: ${error}. Please ensure the backend server is running on ${process.env.PLAYWRIGHT_API_URL || 'http://localhost:8080'} and is accessible.`);
  }
}

/**
 * Test configuration helper
 */
export class TestConfig {
  private static instance: TestConfig;
  private usingMockApi = false;

  static getInstance(): TestConfig {
    if (!TestConfig.instance) {
      TestConfig.instance = new TestConfig();
    }
    return TestConfig.instance;
  }

  setUsingMockApi(value: boolean): void {
    this.usingMockApi = value;
  }

  isUsingMockApi(): boolean {
    return this.usingMockApi;
  }

  /**
   * Skip test if real API is required but not available
   */
  skipIfMockApi(testFn: any): void {
    if (this.usingMockApi) {
      testFn.skip();
    }
  }

  /**
   * Only run test if using mock API
   */
  onlyIfMockApi(testFn: any): void {
    if (!this.usingMockApi) {
      testFn.skip();
    }
  }
}