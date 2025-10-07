import { createApiTestHelpers, ApiTestHelpers } from './api-client';
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
  getSubscriptionIdByName(name: string): Promise<string | null>;
  getProviderIdByName(name: string): Promise<string | null>;
  getLabelIdByName(name: string): Promise<string | null>;
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

  async getSubscriptionIdByName(name: string): Promise<string | null> {
    return this.apiHelpers.getSubscriptionIdByName(name);
  }

  async getProviderIdByName(name: string): Promise<string | null> {
    return this.apiHelpers.getProviderIdByName(name);
  }

  async getLabelIdByName(name: string): Promise<string | null> {
    return this.apiHelpers.getLabelIdByName(name);
  }

  isUsingMockApi(): boolean {
    return false;
  }
}

/**
 * Simple mock test helpers for when API is not available
 */
class MockTestHelpers implements TestHelpers {
  private idCounter = 1;

  async createTestProvider(data: ProviderData): Promise<{ id: string; name: string }> {
    const id = `mock-provider-${this.idCounter++}`;
    return { id, name: data.name };
  }

  async createTestSubscription(data: SubscriptionData): Promise<string> {
    return `mock-subscription-${this.idCounter++}`;
  }

  async createTestLabel(data: LabelData): Promise<string> {
    return `mock-label-${this.idCounter++}`;
  }

  async cleanupTestSubscription(id: string): Promise<void> {
    // Mock cleanup - no-op
  }

  async cleanupTestProvider(id: string): Promise<void> {
    // Mock cleanup - no-op
  }

  async cleanupTestLabel(id: string): Promise<void> {
    // Mock cleanup - no-op
  }

  async getSubscriptionIdByName(name: string): Promise<string | null> {
    return `mock-subscription-${this.idCounter++}`;
  }

  async getProviderIdByName(name: string): Promise<string | null> {
    return `mock-provider-${this.idCounter++}`;
  }

  async getLabelIdByName(name: string): Promise<string | null> {
    return `mock-label-${this.idCounter++}`;
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
      await page.waitForLoadState('networkidle', { timeout: 10000 });
      await page.waitForTimeout(1000); // Reduced timeout

      // Try to get the Clerk session token from the page
      authToken = await page.evaluate(async () => {
        // Wait for Clerk to be available with shorter timeout
        let attempts = 0;
        while (!window.Clerk && attempts < 5) {
          await new Promise(resolve => setTimeout(resolve, 200));
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
        console.warn('⚠️ Could not extract authentication token from page. Proceeding without token.');
      }

    } catch (error) {
      console.warn('⚠️ Failed to extract auth token from page:', error);
      // Don't throw error, just proceed without token
    }
  }

  try {
    // Create real API helpers with auth token
    const apiHelpers = await createApiTestHelpers(undefined, authToken);
    console.log('✅ Using real API for tests');
    return new RealApiTestHelpers(apiHelpers);
  } catch (error) {
    console.warn('⚠️ Failed to connect to API:', error);
    console.log('🔄 Falling back to mock API for tests');
    return new MockTestHelpers();
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