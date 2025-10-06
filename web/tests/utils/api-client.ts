import { APIRequestContext, request } from '@playwright/test';
import { 
  SubscriptionData, 
  ProviderData, 
  LabelData, 
  FamilyData, 
  FamilyMemberData,
  UserProfileData,
  TestUser 
} from './data-generators';

/**
 * Configuration for API client
 */
export interface ApiClientConfig {
  baseUrl: string;
  timeout?: number;
  headers?: Record<string, string>;
  authToken?: string;
}

/**
 * Response wrapper for API calls
 */
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  error?: string;
  status: number;
}

/**
 * API client for test data setup and cleanup operations
 */
export class TestApiClient {
  private context: APIRequestContext | null = null;
  private config: ApiClientConfig;

  constructor(config: ApiClientConfig) {
    this.config = {
      timeout: 30000,
      headers: {
        'Content-Type': 'application/json',
        'Accept': 'application/json',
      },
      ...config,
    };
  }

  /**
   * Initialize the API request context
   */
  async initialize(): Promise<void> {
    if (this.context) {
      await this.context.dispose();
    }
    
    this.context = await request.newContext({
      baseURL: this.config.baseUrl,
      timeout: this.config.timeout,
      extraHTTPHeaders: {
        ...this.config.headers,
        ...(this.config.authToken && { 'Authorization': `Bearer ${this.config.authToken}` }),
      },
    });
  }

  /**
   * Set authentication token
   */
  setAuthToken(token: string): void {
    this.config.authToken = token;
    // Need to reinitialize context with new auth token
    if (this.context) {
      this.context.dispose();
      this.context = null;
    }
  }

  /**
   * Dispose of the API request context
   */
  async dispose(): Promise<void> {
    if (this.context) {
      await this.context.dispose();
      this.context = null;
    }
  }

  /**
   * Generic API request method
   */
  private async makeRequest<T>(
    method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH',
    endpoint: string,
    data?: any,
    options?: { timeout?: number; headers?: Record<string, string> }
  ): Promise<ApiResponse<T>> {
    if (!this.context) {
      throw new Error('API client not initialized. Call initialize() first.');
    }

    try {
      const response = await this.context.fetch(endpoint, {
        method,
        data: data ? JSON.stringify(data) : undefined,
        timeout: options?.timeout || this.config.timeout,
        headers: {
          ...this.config.headers,
          ...(this.config.authToken && { 'Authorization': `Bearer ${this.config.authToken}` }),
          ...options?.headers,
        },
      });

      const responseData = response.ok() ? await response.json().catch(() => null) : null;

      return {
        success: response.ok(),
        data: responseData,
        status: response.status(),
        error: response.ok() ? undefined : `HTTP ${response.status()}: ${response.statusText()}`,
      };
    } catch (error) {
      return {
        success: false,
        status: 0,
        error: `Request failed: ${error}`,
      };
    }
  }

  // Subscription API methods
  async createSubscription(data: SubscriptionData): Promise<ApiResponse<{ id: string }>> {
    return this.makeRequest<{ id: string }>('POST', '/subscriptions', data);
  }

  async getSubscription(id: string): Promise<ApiResponse<SubscriptionData>> {
    return this.makeRequest<SubscriptionData>('GET', `/subscriptions/${id}`);
  }

  async updateSubscription(id: string, data: Partial<SubscriptionData>): Promise<ApiResponse<SubscriptionData>> {
    return this.makeRequest<SubscriptionData>('PUT', `/subscriptions/${id}`, data);
  }

  async deleteSubscription(id: string): Promise<ApiResponse<void>> {
    return this.makeRequest<void>('DELETE', `/subscriptions/${id}`);
  }

  async getSubscriptions(params?: { 
    page?: number; 
    limit?: number; 
    search?: string; 
    labels?: string[];
    providerId?: string;
  }): Promise<ApiResponse<{ subscriptions: SubscriptionData[]; total: number }>> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.search) queryParams.append('search', params.search);
    if (params?.labels) params.labels.forEach(label => queryParams.append('labels', label));
    if (params?.providerId) queryParams.append('providerId', params.providerId);

    const endpoint = `/subscriptions${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return this.makeRequest<{ subscriptions: SubscriptionData[]; total: number }>('GET', endpoint);
  }

  // Provider API methods
  async createProvider(data: ProviderData): Promise<ApiResponse<{ id: string }>> {
    return this.makeRequest<{ id: string }>('POST', '/providers', data);
  }

  async getProvider(id: string): Promise<ApiResponse<ProviderData>> {
    return this.makeRequest<ProviderData>('GET', `/providers/${id}`);
  }

  async updateProvider(id: string, data: Partial<ProviderData>): Promise<ApiResponse<ProviderData>> {
    return this.makeRequest<ProviderData>('PUT', `/providers/${id}`, data);
  }

  async deleteProvider(id: string): Promise<ApiResponse<void>> {
    return this.makeRequest<void>('DELETE', `/providers/${id}`);
  }

  async getProviders(params?: { 
    page?: number; 
    limit?: number; 
    search?: string; 
  }): Promise<ApiResponse<{ providers: ProviderData[]; total: number }>> {
    const queryParams = new URLSearchParams();
    if (params?.page) queryParams.append('page', params.page.toString());
    if (params?.limit) queryParams.append('limit', params.limit.toString());
    if (params?.search) queryParams.append('search', params.search);

    const endpoint = `/providers${queryParams.toString() ? `?${queryParams.toString()}` : ''}`;
    return this.makeRequest<{ providers: ProviderData[]; total: number }>('GET', endpoint);
  }

  // Label API methods
  async createLabel(data: LabelData): Promise<ApiResponse<{ id: string }>> {
    return this.makeRequest<{ id: string }>('POST', '/labels', data);
  }

  async getLabel(id: string): Promise<ApiResponse<LabelData>> {
    return this.makeRequest<LabelData>('GET', `/labels/${id}`);
  }

  async updateLabel(id: string, data: Partial<LabelData>): Promise<ApiResponse<LabelData>> {
    return this.makeRequest<LabelData>('PUT', `/labels/${id}`, data);
  }

  async deleteLabel(id: string): Promise<ApiResponse<void>> {
    return this.makeRequest<void>('DELETE', `/labels/${id}`);
  }

  async getLabels(): Promise<ApiResponse<LabelData[]>> {
    return this.makeRequest<LabelData[]>('GET', '/labels');
  }

  // Family API methods
  async createFamily(data: FamilyData): Promise<ApiResponse<{ id: string }>> {
    return this.makeRequest<{ id: string }>('POST', '/families', data);
  }

  async getFamily(id: string): Promise<ApiResponse<FamilyData>> {
    return this.makeRequest<FamilyData>('GET', `/families/${id}`);
  }

  async updateFamily(id: string, data: Partial<FamilyData>): Promise<ApiResponse<FamilyData>> {
    return this.makeRequest<FamilyData>('PUT', `/families/${id}`, data);
  }

  async deleteFamily(id: string): Promise<ApiResponse<void>> {
    return this.makeRequest<void>('DELETE', `/families/${id}`);
  }

  async inviteFamilyMember(familyId: string, memberData: { email: string; role: string }): Promise<ApiResponse<{ invitationId: string }>> {
    return this.makeRequest<{ invitationId: string }>('POST', `/families/${familyId}/invite`, memberData);
  }

  async acceptFamilyInvitation(invitationId: string): Promise<ApiResponse<void>> {
    return this.makeRequest<void>('POST', `/families/invitations/${invitationId}/accept`);
  }

  async removeFamilyMember(familyId: string, memberId: string): Promise<ApiResponse<void>> {
    return this.makeRequest<void>('DELETE', `/families/${familyId}/members/${memberId}`);
  }

  // User Profile API methods
  async getUserProfile(): Promise<ApiResponse<UserProfileData>> {
    return this.makeRequest<UserProfileData>('GET', '/profile');
  }

  async updateUserProfile(data: Partial<UserProfileData>): Promise<ApiResponse<UserProfileData>> {
    return this.makeRequest<UserProfileData>('PUT', '/profile', data);
  }

  async updatePreferredCurrency(currency: string): Promise<ApiResponse<void>> {
    return this.makeRequest<void>('PUT', '/profile/currency', { currency });
  }

  // Dashboard API methods
  async getDashboardSummary(): Promise<ApiResponse<{
    totalSpending: number;
    activeSubscriptions: number;
    upcomingRenewals: any[];
    topProviders: any[];
    topLabels: any[];
  }>> {
    return this.makeRequest('GET', '/dashboard/summary');
  }

  // Health check
  async healthCheck(): Promise<ApiResponse<{ status: string }>> {
    // Try multiple health check endpoints
    const endpoints = ['/healthz/live', '/health', '/api/health', '/healthcheck'];
    
    for (const endpoint of endpoints) {
      try {
        const response = await this.makeRequest<{ status: string }>('GET', endpoint);
        if (response.success) {
          return response;
        }
      } catch {
        continue;
      }
    }
    
    // If all health checks fail, return a generic error
    return {
      success: false,
      status: 0,
      error: 'No health check endpoint responded successfully',
    };
  }
}

/**
 * API helper utilities for test data management
 */
export class ApiTestHelpers {
  private apiClient: TestApiClient;
  private createdEntities: Map<string, string[]> = new Map();

  constructor(apiClient: TestApiClient) {
    this.apiClient = apiClient;
  }

  /**
   * Create test subscription and track for cleanup
   */
  async createTestSubscription(data: SubscriptionData, testId?: string): Promise<string> {
    const response = await this.apiClient.createSubscription(data);
    
    if (!response.success || !response.data?.id) {
      throw new Error(`Failed to create test subscription: ${response.error}`);
    }

    if (testId) {
      this.trackEntity(testId, 'subscriptions', response.data.id);
    }

    return response.data.id;
  }

  /**
   * Create test provider and track for cleanup
   */
  async createTestProvider(data: ProviderData, testId?: string): Promise<{ id: string; name: string }> {
    const response = await this.apiClient.createProvider(data);
    
    if (!response.success || !response.data?.id) {
      throw new Error(`Failed to create test provider: ${response.error}`);
    }

    const provider = { id: response.data.id, name: data.name };

    if (testId) {
      this.trackEntity(testId, 'providers', response.data.id);
    }

    return provider;
  }

  /**
   * Create test label and track for cleanup
   */
  async createTestLabel(data: LabelData, testId?: string): Promise<string> {
    const response = await this.apiClient.createLabel(data);
    
    if (!response.success || !response.data?.id) {
      throw new Error(`Failed to create test label: ${response.error}`);
    }

    if (testId) {
      this.trackEntity(testId, 'labels', response.data.id);
    }

    return response.data.id;
  }

  /**
   * Create test family and track for cleanup
   */
  async createTestFamily(data: FamilyData, testId?: string): Promise<string> {
    const response = await this.apiClient.createFamily(data);
    
    if (!response.success || !response.data?.id) {
      throw new Error(`Failed to create test family: ${response.error}`);
    }

    if (testId) {
      this.trackEntity(testId, 'families', response.data.id);
    }

    return response.data.id;
  }

  /**
   * Setup complete test scenario with related entities
   */
  async setupTestScenario(testId: string): Promise<{
    providerId: string;
    labelIds: string[];
    subscriptionId: string;
    familyId?: string;
  }> {
    // Create provider
    const provider = await this.createTestProvider({
      name: `Test Provider ${testId}`,
      description: 'Test provider for scenario',
    }, testId);

    // Create labels
    const labelIds = await Promise.all([
      this.createTestLabel({
        name: `Test Label 1 ${testId}`,
        color: '#FF6B6B',
      }, testId),
      this.createTestLabel({
        name: `Test Label 2 ${testId}`,
        color: '#4ECDC4',
      }, testId),
    ]);

    // Create subscription with provider and labels
    const subscriptionId = await this.createTestSubscription({
      name: `Test Subscription ${testId}`,
      providerId: provider.id,
      amount: 29.99,
      currency: 'USD',
      billingCycle: 'monthly',
      nextBillingDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
      labels: labelIds,
    }, testId);

    return {
      providerId: provider.id,
      labelIds,
      subscriptionId,
    };
  }

  /**
   * Bulk create subscriptions for testing pagination and filtering
   */
  async createBulkSubscriptions(count: number, testId: string, baseData?: Partial<SubscriptionData>): Promise<string[]> {
    const subscriptionIds: string[] = [];

    // Create a provider for bulk subscriptions if not provided
    const provider = await this.createTestProvider({
      name: `Bulk Provider ${testId}`,
      description: 'Provider for bulk subscription testing',
    }, testId);

    for (let i = 0; i < count; i++) {
      const subscriptionData: SubscriptionData = {
        name: `Bulk Subscription ${i + 1} ${testId}`,
        providerId: provider.id,
        amount: Math.floor(Math.random() * 100) + 10,
        currency: 'USD',
        billingCycle: 'monthly',
        nextBillingDate: new Date(Date.now() + (i + 1) * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
        ...baseData,
      };

      const id = await this.createTestSubscription(subscriptionData, testId);
      subscriptionIds.push(id);
    }

    return subscriptionIds;
  }

  /**
   * Clean up all test data for a specific test
   */
  async cleanupTestData(testId: string): Promise<void> {
    const entities = this.createdEntities.get(testId);
    if (!entities) return;

    const cleanupPromises: Promise<any>[] = [];

    // Parse entities and create cleanup promises
    entities.forEach(entity => {
      const [type, id] = entity.split(':');
      
      switch (type) {
        case 'subscriptions':
          cleanupPromises.push(this.apiClient.deleteSubscription(id));
          break;
        case 'providers':
          cleanupPromises.push(this.apiClient.deleteProvider(id));
          break;
        case 'labels':
          cleanupPromises.push(this.apiClient.deleteLabel(id));
          break;
        case 'families':
          cleanupPromises.push(this.apiClient.deleteFamily(id));
          break;
      }
    });

    // Execute all cleanup operations
    await Promise.allSettled(cleanupPromises);
    
    // Remove from tracking
    this.createdEntities.delete(testId);
  }

  /**
   * Clean up all tracked test data
   */
  async cleanupAllTestData(): Promise<void> {
    const testIds = Array.from(this.createdEntities.keys());
    
    await Promise.allSettled(
      testIds.map(testId => this.cleanupTestData(testId))
    );
  }

  /**
   * Track created entity for cleanup
   */
  private trackEntity(testId: string, type: string, id: string): void {
    if (!this.createdEntities.has(testId)) {
      this.createdEntities.set(testId, []);
    }
    
    this.createdEntities.get(testId)!.push(`${type}:${id}`);
  }

  /**
   * Get all tracked entities for a test
   */
  getTrackedEntities(testId: string): string[] {
    return this.createdEntities.get(testId) || [];
  }

  /**
   * Check if API is healthy and accessible
   */
  async verifyApiHealth(): Promise<boolean> {
    try {
      const response = await this.apiClient.healthCheck();
      return response.success && (response.data?.status === 'ok' || response.data?.status === 'HEALTHY');
    } catch {
      return false;
    }
  }

  /**
   * Convenience method to clean up a test subscription
   */
  async cleanupTestSubscription(id: string): Promise<void> {
    await this.apiClient.deleteSubscription(id);
  }

  /**
   * Convenience method to clean up a test provider
   */
  async cleanupTestProvider(id: string): Promise<void> {
    await this.apiClient.deleteProvider(id);
  }

  /**
   * Convenience method to clean up a test label
   */
  async cleanupTestLabel(id: string): Promise<void> {
    await this.apiClient.deleteLabel(id);
  }

  /**
   * Convenience method to clean up a test family
   */
  async cleanupTestFamily(id: string): Promise<void> {
    await this.apiClient.deleteFamily(id);
  }
}

/**
 * Factory function to create configured API client for tests
 */
export function createTestApiClient(config?: Partial<ApiClientConfig>): TestApiClient {
  const defaultConfig: ApiClientConfig = {
    baseUrl: process.env.PLAYWRIGHT_API_URL || process.env.API_BASE_URL || 'http://localhost:8080',
    timeout: 30000,
  };

  return new TestApiClient({ ...defaultConfig, ...config });
}

/**
 * Factory function to create API helpers with configured client
 */
export async function createApiTestHelpers(config?: Partial<ApiClientConfig>, authToken?: string): Promise<ApiTestHelpers> {
  const apiClient = createTestApiClient(config);
  
  try {
    // Set auth token before initializing if provided
    if (authToken) {
      apiClient.setAuthToken(authToken);
    }
    
    await apiClient.initialize();
    
    // Test if API is available
    const testResponse = await apiClient.healthCheck();
    if (!testResponse.success) {
      throw new Error(`API health check failed: ${testResponse.error}`);
    }
    
    // Verify the response indicates healthy status (be more flexible with status values)
    const validStatuses = ['HEALTHY', 'ok', 'OK', 'healthy', 'up', 'UP'];
    if (testResponse.data?.status && !validStatuses.includes(testResponse.data.status)) {
      throw new Error(`API is not healthy: ${testResponse.data?.status}`);
    }
    
    return new ApiTestHelpers(apiClient);
  } catch (error) {
    // If API is not available, dispose the client and throw
    await apiClient.dispose();
    throw error;
  }
}