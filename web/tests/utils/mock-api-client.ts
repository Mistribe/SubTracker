import { 
  SubscriptionData, 
  ProviderData, 
  LabelData, 
  FamilyData, 
  TestDataGenerators 
} from './data-generators';

/**
 * Mock API client for testing UI functionality without a real API server
 */
export class MockApiTestHelpers {
  private mockProviders: Map<string, ProviderData & { id: string }> = new Map();
  private mockSubscriptions: Map<string, SubscriptionData & { id: string }> = new Map();
  private mockLabels: Map<string, LabelData & { id: string }> = new Map();
  private idCounter = 1;

  /**
   * Generate a mock ID
   */
  private generateId(): string {
    return `mock-${this.idCounter++}`;
  }

  /**
   * Create mock test provider
   */
  async createTestProvider(data: ProviderData): Promise<{ id: string; name: string }> {
    const id = this.generateId();
    const provider = { ...data, id };
    this.mockProviders.set(id, provider);
    return { id, name: data.name };
  }

  /**
   * Create mock test subscription
   */
  async createTestSubscription(data: SubscriptionData): Promise<string> {
    const id = this.generateId();
    const subscription = { ...data, id };
    this.mockSubscriptions.set(id, subscription);
    return id;
  }

  /**
   * Create mock test label
   */
  async createTestLabel(data: LabelData): Promise<string> {
    const id = this.generateId();
    const label = { ...data, id };
    this.mockLabels.set(id, label);
    return id;
  }

  /**
   * Mock cleanup methods (no-op for mock)
   */
  async cleanupTestSubscription(id: string): Promise<void> {
    this.mockSubscriptions.delete(id);
  }

  async cleanupTestProvider(id: string): Promise<void> {
    this.mockProviders.delete(id);
  }

  async cleanupTestLabel(id: string): Promise<void> {
    this.mockLabels.delete(id);
  }

  /**
   * Get mock data for testing
   */
  getMockProvider(id: string): (ProviderData & { id: string }) | undefined {
    return this.mockProviders.get(id);
  }

  getMockSubscription(id: string): (SubscriptionData & { id: string }) | undefined {
    return this.mockSubscriptions.get(id);
  }

  getMockLabel(id: string): (LabelData & { id: string }) | undefined {
    return this.mockLabels.get(id);
  }

  /**
   * Clear all mock data
   */
  clearAll(): void {
    this.mockProviders.clear();
    this.mockSubscriptions.clear();
    this.mockLabels.clear();
    this.idCounter = 1;
  }
}

/**
 * Factory function to create mock API helpers
 */
export function createMockApiTestHelpers(): MockApiTestHelpers {
  return new MockApiTestHelpers();
}