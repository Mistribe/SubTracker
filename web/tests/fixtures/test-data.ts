/**
 * Test Data Fixtures
 * 
 * Provides fixtures for test data creation, management, and cleanup
 * Ensures test isolation and proper cleanup after test execution
 */

import { test as base } from '@playwright/test';
import { ClerkTestUserManager, ClerkTestUser } from '../utils/clerk-user-manager';
import { SessionManager } from '../utils/session-manager';

/**
 * Test data interfaces
 */
export interface SubscriptionTestData {
  id?: string;
  name: string;
  providerId: string;
  amount: number;
  currency: string;
  billingCycle: 'monthly' | 'yearly';
  nextBillingDate: string;
  labels?: string[];
  payerId?: string;
}

export interface ProviderTestData {
  id?: string;
  name: string;
  description?: string;
  website?: string;
  logoUrl?: string;
}

export interface LabelTestData {
  id?: string;
  name: string;
  color: string;
  description?: string;
}

export interface FamilyTestData {
  id?: string;
  name: string;
  ownerId: string;
  members: FamilyMemberTestData[];
}

export interface FamilyMemberTestData {
  id?: string;
  email: string;
  role: 'owner' | 'member';
  permissions: string[];
  status: 'pending' | 'accepted' | 'declined';
}

/**
 * Test data fixture interface
 */
export interface TestDataFixtures {
  testSubscription: SubscriptionTestData;
  testProvider: ProviderTestData;
  testLabel: LabelTestData;
  testFamily: FamilyTestData;
  testDataCleanup: () => Promise<void>;
  uniqueId: string;
}

/**
 * Test data generator utilities
 */
export class TestDataGenerator {
  private static counter = 0;

  /**
   * Generate unique identifier for test isolation
   */
  static generateUniqueId(): string {
    const timestamp = Date.now();
    const counter = ++this.counter;
    const random = Math.random().toString(36).substr(2, 9);
    return `test-${timestamp}-${counter}-${random}`;
  }

  /**
   * Generate test subscription data
   */
  static generateSubscription(uniqueId: string): SubscriptionTestData {
    return {
      name: `Test Subscription ${uniqueId}`,
      providerId: `provider-${uniqueId}`,
      amount: Math.floor(Math.random() * 100) + 10, // $10-$110
      currency: 'USD',
      billingCycle: Math.random() > 0.5 ? 'monthly' : 'yearly',
      nextBillingDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString(), // 30 days from now
      labels: [`label-${uniqueId}`],
    };
  }

  /**
   * Generate test provider data
   */
  static generateProvider(uniqueId: string): ProviderTestData {
    const providers = ['Netflix', 'Spotify', 'Adobe', 'Microsoft', 'Google'];
    const randomProvider = providers[Math.floor(Math.random() * providers.length)];
    
    return {
      name: `${randomProvider} Test ${uniqueId}`,
      description: `Test provider for ${randomProvider} - ${uniqueId}`,
      website: `https://test-${uniqueId}.example.com`,
      logoUrl: `https://test-${uniqueId}.example.com/logo.png`,
    };
  }

  /**
   * Generate test label data
   */
  static generateLabel(uniqueId: string): LabelTestData {
    const colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD'];
    const categories = ['Entertainment', 'Productivity', 'Development', 'Design', 'Business'];
    
    return {
      name: `${categories[Math.floor(Math.random() * categories.length)]} ${uniqueId}`,
      color: colors[Math.floor(Math.random() * colors.length)],
      description: `Test label for categorization - ${uniqueId}`,
    };
  }

  /**
   * Generate test family data
   */
  static generateFamily(uniqueId: string, ownerId: string): FamilyTestData {
    return {
      name: `Test Family ${uniqueId}`,
      ownerId,
      members: [
        {
          email: `member1-${uniqueId}@example.com`,
          role: 'member',
          permissions: ['view_subscriptions'],
          status: 'accepted',
        },
        {
          email: `member2-${uniqueId}@example.com`,
          role: 'member',
          permissions: ['view_subscriptions', 'edit_subscriptions'],
          status: 'pending',
        },
      ],
    };
  }

  /**
   * Generate realistic test data set
   */
  static generateTestDataSet(uniqueId: string) {
    const provider = this.generateProvider(uniqueId);
    const label = this.generateLabel(uniqueId);
    const subscription = this.generateSubscription(uniqueId);
    
    // Link subscription to provider and label
    subscription.providerId = provider.id || `provider-${uniqueId}`;
    subscription.labels = [label.id || `label-${uniqueId}`];
    
    return {
      provider,
      label,
      subscription,
    };
  }
}

/**
 * Test data cleanup manager
 */
export class TestDataCleanup {
  private static createdData: Map<string, {
    subscriptions: string[];
    providers: string[];
    labels: string[];
    families: string[];
  }> = new Map();

  /**
   * Track created test data
   */
  static trackCreatedData(
    uniqueId: string,
    type: 'subscription' | 'provider' | 'label' | 'family',
    id: string
  ): void {
    if (!this.createdData.has(uniqueId)) {
      this.createdData.set(uniqueId, {
        subscriptions: [],
        providers: [],
        labels: [],
        families: [],
      });
    }

    const data = this.createdData.get(uniqueId)!;
    switch (type) {
      case 'subscription':
        data.subscriptions.push(id);
        break;
      case 'provider':
        data.providers.push(id);
        break;
      case 'label':
        data.labels.push(id);
        break;
      case 'family':
        data.families.push(id);
        break;
    }
  }

  /**
   * Cleanup test data for a specific test
   */
  static async cleanupTestData(uniqueId: string): Promise<void> {
    const data = this.createdData.get(uniqueId);
    if (!data) return;

    try {
      // In a real implementation, you would call API endpoints to delete data
      // For now, we'll just log the cleanup actions
      
      console.log(`Cleaning up test data for ${uniqueId}:`);
      
      if (data.subscriptions.length > 0) {
        console.log(`  - Deleting ${data.subscriptions.length} subscriptions`);
        // await apiClient.deleteSubscriptions(data.subscriptions);
      }
      
      if (data.providers.length > 0) {
        console.log(`  - Deleting ${data.providers.length} providers`);
        // await apiClient.deleteProviders(data.providers);
      }
      
      if (data.labels.length > 0) {
        console.log(`  - Deleting ${data.labels.length} labels`);
        // await apiClient.deleteLabels(data.labels);
      }
      
      if (data.families.length > 0) {
        console.log(`  - Deleting ${data.families.length} families`);
        // await apiClient.deleteFamilies(data.families);
      }

      // Remove from tracking
      this.createdData.delete(uniqueId);
      
    } catch (error) {
      console.error(`Failed to cleanup test data for ${uniqueId}:`, error);
    }
  }

  /**
   * Cleanup all test data
   */
  static async cleanupAllTestData(): Promise<void> {
    const uniqueIds = Array.from(this.createdData.keys());
    
    for (const uniqueId of uniqueIds) {
      await this.cleanupTestData(uniqueId);
    }
  }

  /**
   * Get cleanup summary
   */
  static getCleanupSummary(): { [key: string]: number } {
    let totalSubscriptions = 0;
    let totalProviders = 0;
    let totalLabels = 0;
    let totalFamilies = 0;

    for (const data of this.createdData.values()) {
      totalSubscriptions += data.subscriptions.length;
      totalProviders += data.providers.length;
      totalLabels += data.labels.length;
      totalFamilies += data.families.length;
    }

    return {
      subscriptions: totalSubscriptions,
      providers: totalProviders,
      labels: totalLabels,
      families: totalFamilies,
      tests: this.createdData.size,
    };
  }
}

/**
 * Extended test with test data fixtures
 */
export const test = base.extend<TestDataFixtures>({
  /**
   * Unique ID fixture for test isolation
   */
  uniqueId: async ({}, use) => {
    const uniqueId = TestDataGenerator.generateUniqueId();
    await use(uniqueId);
  },

  /**
   * Test subscription fixture
   */
  testSubscription: async ({ uniqueId }, use) => {
    const subscription = TestDataGenerator.generateSubscription(uniqueId);
    await use(subscription);
    
    // Cleanup is handled by testDataCleanup fixture
  },

  /**
   * Test provider fixture
   */
  testProvider: async ({ uniqueId }, use) => {
    const provider = TestDataGenerator.generateProvider(uniqueId);
    await use(provider);
    
    // Cleanup is handled by testDataCleanup fixture
  },

  /**
   * Test label fixture
   */
  testLabel: async ({ uniqueId }, use) => {
    const label = TestDataGenerator.generateLabel(uniqueId);
    await use(label);
    
    // Cleanup is handled by testDataCleanup fixture
  },

  /**
   * Test family fixture
   */
  testFamily: async ({ uniqueId }, use) => {
    const testUser = await ClerkTestUserManager.getDefaultTestUser();
    const family = TestDataGenerator.generateFamily(uniqueId, testUser.id);
    await use(family);
    
    // Cleanup is handled by testDataCleanup fixture
  },

  /**
   * Test data cleanup fixture
   */
  testDataCleanup: async ({ uniqueId }, use) => {
    // Provide cleanup function to test
    const cleanup = async () => {
      await TestDataCleanup.cleanupTestData(uniqueId);
      await SessionManager.cleanupExpiredSessions();
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
 * Export TestDataFixtures interface for type checking
 */
export type { TestDataFixtures };