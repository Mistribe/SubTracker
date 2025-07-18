import axios, { AxiosError, AxiosInstance, AxiosResponse } from 'axios';
import {
  LabelModel,
  CreateLabelModel,
  UpdateLabelModel,
  FamilyModel,
  CreateFamilyModel,
  UpdateFamilyModel,
  CreateFamilyMemberModel,
  UpdateFamilyMemberModel,
  SubscriptionModel,
  CreateSubscriptionModel,
  UpdateSubscriptionModel,
  CreateSubscriptionPaymentModel,
  HttpError
} from '../models';

/**
 * API Client for interacting with the SubTracker API
 */
export class ApiClient {
  private client: AxiosInstance;
  
  /**
   * Creates a new ApiClient instance
   * @param baseURL The base URL for the API (default: http://localhost:8080)
   */
  constructor(baseURL: string = 'http://localhost:8080') {
    this.client = axios.create({
      baseURL,
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        // Handle API errors
        if (error.response) {
          const errorData = error.response.data as HttpError;
          console.error('API Error:', errorData.message || 'Unknown error');
        } else if (error.request) {
          console.error('Network Error:', error.message);
        } else {
          console.error('Error:', error.message);
        }
        return Promise.reject(error);
      }
    );
  }
  
  /**
   * Helper method to handle API responses
   */
  private handleResponse<T>(response: AxiosResponse<T>): T {
    return response.data;
  }
  
  /**
   * Helper method to handle API errors
   */
  private handleError(error: unknown): never {
    if (axios.isAxiosError(error)) {
      if (error.response) {
        const errorData = error.response.data as HttpError;
        throw new Error(errorData.message || `API Error: ${error.response.status}`);
      } else if (error.request) {
        throw new Error(`Network Error: ${error.message}`);
      }
    }
    
    // For non-Axios errors
    if (error instanceof Error) {
      throw new Error(`Error: ${error.message}`);
    }
    
    throw new Error('An unknown error occurred');
  }

  // ==================== LABEL ENDPOINTS ====================
  
  /**
   * Get all labels
   * @param withDefault Include default labels
   * @returns Array of labels
   */
  async getLabels(withDefault?: boolean): Promise<LabelModel[]> {
    try {
      const params = withDefault !== undefined ? { with_default: withDefault } : {};
      const response = await this.client.get<LabelModel[]>('/labels', { params });
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Get default labels
   * @returns Array of default labels
   */
  async getDefaultLabels(): Promise<LabelModel[]> {
    try {
      const response = await this.client.get<LabelModel[]>('/default');
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Get label by ID
   * @param id Label ID
   * @returns Label
   */
  async getLabelById(id: string): Promise<LabelModel> {
    try {
      const response = await this.client.get<LabelModel>(`/labels/${id}`);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Create a new label
   * @param label Label data
   * @returns Created label
   */
  async createLabel(label: CreateLabelModel): Promise<LabelModel> {
    try {
      const response = await this.client.post<LabelModel>('/labels', label);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Update a label
   * @param id Label ID
   * @param label Label data
   * @returns Updated label
   */
  async updateLabel(id: string, label: UpdateLabelModel): Promise<LabelModel> {
    try {
      const response = await this.client.put<LabelModel>(`/labels/${id}`, label);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Delete a label
   * @param id Label ID
   */
  async deleteLabel(id: string): Promise<void> {
    try {
      await this.client.delete(`/labels/${id}`);
    } catch (error) {
      this.handleError(error);
    }
  }

  // ==================== FAMILY ENDPOINTS ====================
  
  /**
   * Get all families
   * @returns Array of families
   */
  async getFamilies(): Promise<FamilyModel[]> {
    try {
      const response = await this.client.get<FamilyModel[]>('/families');
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Get family by ID
   * @param id Family ID
   * @returns Family
   */
  async getFamilyById(id: string): Promise<FamilyModel> {
    try {
      const response = await this.client.get<FamilyModel>(`/families/${id}`);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Create a new family
   * @param family Family data
   * @returns Created family
   */
  async createFamily(family: CreateFamilyModel): Promise<FamilyModel> {
    try {
      const response = await this.client.post<FamilyModel>('/families', family);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Update a family
   * @param id Family ID
   * @param family Family data
   * @returns Updated family
   */
  async updateFamily(id: string, family: UpdateFamilyModel): Promise<FamilyModel> {
    try {
      const response = await this.client.put<FamilyModel>(`/families/${id}`, family);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Add a family member
   * @param familyId Family ID
   * @param member Family member data
   * @returns Updated family
   */
  async addFamilyMember(familyId: string, member: CreateFamilyMemberModel): Promise<FamilyModel> {
    try {
      const response = await this.client.post<FamilyModel>(`/families/${familyId}/members`, member);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Update a family member
   * @param familyId Family ID
   * @param memberId Member ID
   * @param member Family member data
   * @returns Updated family
   */
  async updateFamilyMember(
    familyId: string,
    memberId: string,
    member: UpdateFamilyMemberModel
  ): Promise<FamilyModel> {
    try {
      const response = await this.client.put<FamilyModel>(
        `/families/${familyId}/members/${memberId}`,
        member
      );
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Delete a family member
   * @param familyId Family ID
   * @param memberId Member ID
   */
  async deleteFamilyMember(familyId: string, memberId: string): Promise<void> {
    try {
      await this.client.delete(`/families/${familyId}/members/${memberId}`);
    } catch (error) {
      this.handleError(error);
    }
  }

  // ==================== SUBSCRIPTION ENDPOINTS ====================
  
  /**
   * Get all subscriptions
   * @returns Array of subscriptions
   */
  async getSubscriptions(): Promise<SubscriptionModel[]> {
    try {
      const response = await this.client.get<SubscriptionModel[]>('/subscriptions');
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Get subscription by ID
   * @param id Subscription ID
   * @returns Subscription
   */
  async getSubscriptionById(id: string): Promise<SubscriptionModel> {
    try {
      const response = await this.client.get<SubscriptionModel>(`/subscriptions/${id}`);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Create a new subscription
   * @param subscription Subscription data
   * @returns Created subscription
   */
  async createSubscription(subscription: CreateSubscriptionModel): Promise<SubscriptionModel> {
    try {
      const response = await this.client.post<SubscriptionModel>('/subscriptions', subscription);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Update a subscription
   * @param id Subscription ID
   * @param subscription Subscription data
   * @returns Updated subscription
   */
  async updateSubscription(
    id: string,
    subscription: UpdateSubscriptionModel
  ): Promise<SubscriptionModel> {
    try {
      const response = await this.client.put<SubscriptionModel>(`/subscriptions/${id}`, subscription);
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Delete a subscription
   * @param id Subscription ID
   */
  async deleteSubscription(id: string): Promise<void> {
    try {
      await this.client.delete(`/subscriptions/${id}`);
    } catch (error) {
      this.handleError(error);
    }
  }
  
  /**
   * Add a payment to a subscription
   * @param subscriptionId Subscription ID
   * @param payment Payment data
   * @returns Updated subscription
   */
  async addSubscriptionPayment(
    subscriptionId: string,
    payment: CreateSubscriptionPaymentModel
  ): Promise<SubscriptionModel> {
    try {
      const response = await this.client.post<SubscriptionModel>(
        `/subscriptions/${subscriptionId}/payments`,
        payment
      );
      return this.handleResponse(response);
    } catch (error) {
      return this.handleError(error);
    }
  }
  
  /**
   * Delete a subscription payment
   * @param subscriptionId Subscription ID
   * @param paymentId Payment ID
   */
  async deleteSubscriptionPayment(subscriptionId: string, paymentId: string): Promise<void> {
    try {
      await this.client.delete(`/subscriptions/${subscriptionId}/payments/${paymentId}`);
    } catch (error) {
      this.handleError(error);
    }
  }
}