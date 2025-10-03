import { randomBytes } from 'crypto';

/**
 * Test data interfaces matching the application models
 */
export interface SubscriptionData {
  name: string;
  providerId: string;
  amount: number;
  currency: string;
  billingCycle: 'monthly' | 'yearly' | 'weekly' | 'daily';
  nextBillingDate: string;
  labels?: string[];
  payerId?: string;
  description?: string;
  isActive?: boolean;
  freeTrialEndDate?: string;
}

export interface ProviderData {
  name: string;
  description?: string;
  website?: string;
  logoUrl?: string;
  category?: string;
  isActive?: boolean;
}

export interface LabelData {
  name: string;
  color: string;
  description?: string;
  isActive?: boolean;
}

export interface FamilyMemberData {
  email: string;
  role: 'owner' | 'member';
  permissions: string[];
  firstName?: string;
  lastName?: string;
}

export interface FamilyData {
  name: string;
  description?: string;
  ownerId: string;
  members?: FamilyMemberData[];
}

export interface UserProfileData {
  firstName: string;
  lastName: string;
  email: string;
  preferredCurrency: string;
  timezone?: string;
  dateFormat?: string;
  notifications?: {
    email: boolean;
    push: boolean;
    renewalReminders: boolean;
  };
}

export interface TestUser {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  clerkUserId: string;
  sessionToken?: string;
}

/**
 * Utility class for generating unique test identifiers
 */
export class TestIdGenerator {
  private static counter = 0;
  private static readonly prefix = 'test';
  
  /**
   * Generate unique test ID with timestamp and counter
   */
  static generate(type?: string): string {
    this.counter++;
    const timestamp = Date.now();
    const random = randomBytes(4).toString('hex');
    const typePrefix = type ? `${type}-` : '';
    return `${this.prefix}-${typePrefix}${timestamp}-${this.counter}-${random}`;
  }

  /**
   * Generate unique email for testing
   */
  static generateEmail(domain = 'test.example.com'): string {
    return `${this.generate('user')}@${domain}`;
  }

  /**
   * Generate unique name with test prefix
   */
  static generateName(type: string): string {
    return `Test ${type} ${this.generate()}`;
  }

  /**
   * Reset counter (useful for test isolation)
   */
  static reset(): void {
    this.counter = 0;
  }
}

/**
 * Test data generators for creating realistic test data
 */
export class TestDataGenerators {
  private static readonly currencies = ['USD', 'EUR', 'GBP', 'CAD', 'AUD'];
  private static readonly billingCycles: Array<'monthly' | 'yearly' | 'weekly' | 'daily'> = ['monthly', 'yearly', 'weekly', 'daily'];
  private static readonly colors = ['#FF6B6B', '#4ECDC4', '#45B7D1', '#96CEB4', '#FFEAA7', '#DDA0DD', '#98D8C8', '#F7DC6F'];
  private static readonly providerCategories = ['Streaming', 'Software', 'Cloud Storage', 'Productivity', 'Entertainment', 'Security', 'Communication'];

  /**
   * Generate a cryptographically secure random integer between min and max, inclusive
   */
  private static getCryptoRandomInt(min: number, max: number): number {
    // Inclusive of min and max
    const range = max - min + 1;
    if (range <= 0) throw new Error('Invalid random range');
    const maxUint32 = 0xFFFFFFFF;
    const buckets = Math.floor(maxUint32 / range);
    const limit = buckets * range;
    let rand;
    do {
      rand = randomBytes(4).readUInt32BE(0);
    } while (rand >= limit);
    return min + (rand % range);
  }
  /**
   * Generate realistic subscription test data
   */
  static generateSubscription(overrides?: Partial<SubscriptionData>): SubscriptionData {
    const baseData: SubscriptionData = {
      name: TestIdGenerator.generateName('Subscription'),
      providerId: TestIdGenerator.generate('provider'),
      amount: TestDataGenerators.getCryptoRandomInt(5, 105), // $5-$105
      currency: this.getRandomItem(this.currencies),
      billingCycle: this.getRandomItem(this.billingCycles),
      nextBillingDate: this.generateFutureDate(1, 365), // 1-365 days from now
      description: `Test subscription for ${TestIdGenerator.generate()}`,
      isActive: true,
    };

    // Add optional fields randomly
    // random true with probability ~0.3
    if (TestDataGenerators.getCryptoRandomInt(1, 100) > 70) {
      baseData.labels = [TestIdGenerator.generate('label'), TestIdGenerator.generate('label')];
    }

    // random true with probability ~0.2
    if (TestDataGenerators.getCryptoRandomInt(1, 100) > 80) {
      baseData.freeTrialEndDate = this.generateFutureDate(7, 30);
    }

    return { ...baseData, ...overrides };
  }

  /**
   * Generate realistic provider test data
   */
  static generateProvider(overrides?: Partial<ProviderData>): ProviderData {
    const baseData: ProviderData = {
      name: TestIdGenerator.generateName('Provider'),
      description: `Test provider for ${TestIdGenerator.generate()}`,
      website: `https://${TestIdGenerator.generate()}.example.com`,
      category: this.getRandomItem(this.providerCategories),
      isActive: true,
    };

    return { ...baseData, ...overrides };
  }

  /**
   * Generate realistic label test data
   */
  static generateLabel(overrides?: Partial<LabelData>): LabelData {
    const baseData: LabelData = {
      name: TestIdGenerator.generateName('Label'),
      color: this.getRandomItem(this.colors),
      description: `Test label for ${TestIdGenerator.generate()}`,
      isActive: true,
    };

    return { ...baseData, ...overrides };
  }

  /**
   * Generate realistic family member test data
   */
  static generateFamilyMember(overrides?: Partial<FamilyMemberData>): FamilyMemberData {
    const baseData: FamilyMemberData = {
      email: TestIdGenerator.generateEmail(),
      role: Math.random() > 0.8 ? 'owner' : 'member',
      permissions: ['view_subscriptions', 'manage_own_subscriptions'],
      firstName: `Test${TestIdGenerator.generate()}`,
      lastName: `User${TestIdGenerator.generate()}`,
    };

    // Add additional permissions for owners
    if (baseData.role === 'owner') {
      baseData.permissions.push('manage_family', 'manage_all_subscriptions', 'invite_members');
    }

    return { ...baseData, ...overrides };
  }

  /**
   * Generate realistic family test data
   */
  static generateFamily(overrides?: Partial<FamilyData>): FamilyData {
    const baseData: FamilyData = {
      name: TestIdGenerator.generateName('Family'),
      description: `Test family for ${TestIdGenerator.generate()}`,
      ownerId: TestIdGenerator.generate('user'),
    };

    return { ...baseData, ...overrides };
  }

  /**
   * Generate realistic user profile test data
   */
  static generateUserProfile(overrides?: Partial<UserProfileData>): UserProfileData {
    const baseData: UserProfileData = {
      firstName: `Test${TestIdGenerator.generate()}`,
      lastName: `User${TestIdGenerator.generate()}`,
      email: TestIdGenerator.generateEmail(),
      preferredCurrency: this.getRandomItem(this.currencies),
      timezone: 'America/New_York',
      dateFormat: 'MM/DD/YYYY',
      notifications: {
        email: true,
        push: Math.random() > 0.5,
        renewalReminders: true,
      },
    };

    return { ...baseData, ...overrides };
  }

  /**
   * Generate test user data for authentication
   */
  static generateTestUser(overrides?: Partial<TestUser>): TestUser {
    const baseData: TestUser = {
      id: TestIdGenerator.generate('user'),
      email: TestIdGenerator.generateEmail(),
      firstName: `Test${TestIdGenerator.generate()}`,
      lastName: `User${TestIdGenerator.generate()}`,
      clerkUserId: TestIdGenerator.generate('clerk'),
    };

    return { ...baseData, ...overrides };
  }

  /**
   * Generate multiple subscriptions for testing lists and pagination
   */
  static generateSubscriptions(count: number, overrides?: Partial<SubscriptionData>): SubscriptionData[] {
    return Array.from({ length: count }, () => this.generateSubscription(overrides));
  }

  /**
   * Generate multiple providers for testing lists and pagination
   */
  static generateProviders(count: number, overrides?: Partial<ProviderData>): ProviderData[] {
    return Array.from({ length: count }, () => this.generateProvider(overrides));
  }

  /**
   * Generate multiple labels for testing
   */
  static generateLabels(count: number, overrides?: Partial<LabelData>): LabelData[] {
    return Array.from({ length: count }, () => this.generateLabel(overrides));
  }

  /**
   * Generate family with multiple members
   */
  static generateFamilyWithMembers(memberCount: number, overrides?: Partial<FamilyData>): FamilyData {
    const family = this.generateFamily(overrides);
    family.members = Array.from({ length: memberCount }, () => this.generateFamilyMember());
    
    // Ensure at least one owner
    if (family.members.length > 0) {
      family.members[0].role = 'owner';
      family.members[0].permissions = ['view_subscriptions', 'manage_own_subscriptions', 'manage_family', 'manage_all_subscriptions', 'invite_members'];
    }

    return family;
  }

  /**
   * Generate subscription with related entities (provider, labels)
   */
  static generateSubscriptionWithRelations(): {
    subscription: SubscriptionData;
    provider: ProviderData;
    labels: LabelData[];
  } {
    const provider = this.generateProvider();
    const labels = this.generateLabels(2);
    const subscription = this.generateSubscription({
      providerId: provider.name, // Use provider name as ID for testing
      labels: labels.map(label => label.name),
    });

    return { subscription, provider, labels };
  }

  /**
   * Utility method to get random item from array
   */
  private static getRandomItem<T>(array: T[]): T {
    return array[TestDataGenerators.getCryptoRandomInt(0, array.length - 1)];
  }

  /**
   * Generate future date within specified range
   */
  private static generateFutureDate(minDays: number, maxDays: number): string {
    const today = new Date();
    const daysToAdd = TestDataGenerators.getCryptoRandomInt(minDays, maxDays);
    const futureDate = new Date(today.getTime() + daysToAdd * 24 * 60 * 60 * 1000);
    return futureDate.toISOString().split('T')[0]; // Return YYYY-MM-DD format
  }

  /**
   * Generate past date within specified range
   */
  private static generatePastDate(minDays: number, maxDays: number): string {
    const today = new Date();
    const daysToSubtract = TestDataGenerators.getCryptoRandomInt(minDays, maxDays);
    const pastDate = new Date(today.getTime() - daysToSubtract * 24 * 60 * 60 * 1000);
    return pastDate.toISOString().split('T')[0]; // Return YYYY-MM-DD format
  }
}

/**
 * Data validation utilities for test data
 */
export class TestDataValidator {
  /**
   * Validate subscription data structure
   */
  static validateSubscription(data: SubscriptionData): { isValid: boolean; errors: string[] } {
    const errors: string[] = [];

    if (!data.name || data.name.trim().length === 0) {
      errors.push('Subscription name is required');
    }

    if (!data.providerId || data.providerId.trim().length === 0) {
      errors.push('Provider ID is required');
    }

    if (data.amount <= 0) {
      errors.push('Amount must be greater than 0');
    }

    if (!this.currencies.includes(data.currency)) {
      errors.push(`Invalid currency: ${data.currency}`);
    }

    if (!this.billingCycles.includes(data.billingCycle)) {
      errors.push(`Invalid billing cycle: ${data.billingCycle}`);
    }

    if (!this.isValidDate(data.nextBillingDate)) {
      errors.push('Invalid next billing date format');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  /**
   * Validate provider data structure
   */
  static validateProvider(data: ProviderData): { isValid: boolean; errors: string[] } {
    const errors: string[] = [];

    if (!data.name || data.name.trim().length === 0) {
      errors.push('Provider name is required');
    }

    if (data.website && !this.isValidUrl(data.website)) {
      errors.push('Invalid website URL format');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  /**
   * Validate label data structure
   */
  static validateLabel(data: LabelData): { isValid: boolean; errors: string[] } {
    const errors: string[] = [];

    if (!data.name || data.name.trim().length === 0) {
      errors.push('Label name is required');
    }

    if (!this.isValidHexColor(data.color)) {
      errors.push('Invalid color format (must be hex color)');
    }

    return {
      isValid: errors.length === 0,
      errors,
    };
  }

  /**
   * Validate email format
   */
  private static isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  /**
   * Validate URL format
   */
  private static isValidUrl(url: string): boolean {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Validate hex color format
   */
  private static isValidHexColor(color: string): boolean {
    const hexColorRegex = /^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$/;
    return hexColorRegex.test(color);
  }

  /**
   * Validate date format (YYYY-MM-DD)
   */
  private static isValidDate(dateString: string): boolean {
    const date = new Date(dateString);
    return date instanceof Date && !isNaN(date.getTime());
  }

  private static readonly currencies = ['USD', 'EUR', 'GBP', 'CAD', 'AUD'];
  private static readonly billingCycles = ['monthly', 'yearly', 'weekly', 'daily'];
}

/**
 * Utility functions for test data cleanup and management
 */
export class TestDataManager {
  private static testDataRegistry: Map<string, any[]> = new Map();

  /**
   * Register test data for cleanup
   */
  static registerTestData(testId: string, data: any): void {
    if (!this.testDataRegistry.has(testId)) {
      this.testDataRegistry.set(testId, []);
    }
    this.testDataRegistry.get(testId)!.push(data);
  }

  /**
   * Get all test data for a specific test
   */
  static getTestData(testId: string): any[] {
    return this.testDataRegistry.get(testId) || [];
  }

  /**
   * Clear test data registry for a specific test
   */
  static clearTestData(testId: string): void {
    this.testDataRegistry.delete(testId);
  }

  /**
   * Clear all test data registry
   */
  static clearAllTestData(): void {
    this.testDataRegistry.clear();
  }

  /**
   * Get all registered test IDs
   */
  static getAllTestIds(): string[] {
    return Array.from(this.testDataRegistry.keys());
  }
}