/**
 * Clerk Test User Manager
 * 
 * Utility class for managing test users with Clerk authentication
 * This class provides methods to create, manage, and cleanup test users
 */

import { getTestEnvironment } from '../config/test-env';

const testEnv = getTestEnvironment();

export interface ClerkTestUser {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  password: string;
  role: 'owner' | 'member';
  clerkUserId?: string;
  sessionToken?: string;
}

export interface CreateTestUserOptions {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  role?: 'owner' | 'member';
}

/**
 * Clerk Test User Manager
 * 
 * Manages test users for Playwright E2E tests using Clerk's testing utilities
 */
export class ClerkTestUserManager {
  private static testUsers: Map<string, ClerkTestUser> = new Map();
  private static initialized = false;

  /**
   * Initialize the test user manager
   */
  static async initialize(): Promise<void> {
    if (this.initialized) return;

    // Setup default test users from environment
    const defaultUser: ClerkTestUser = {
      id: 'default-test-user',
      email: testEnv.testUser.email,
      password: testEnv.testUser.password,
      firstName: 'Test',
      lastName: 'User',
      role: 'owner',
    };

    this.testUsers.set(defaultUser.id, defaultUser);

    // Setup family owner if configured
    if (testEnv.familyOwner) {
      const familyOwner: ClerkTestUser = {
        id: 'family-owner',
        email: testEnv.familyOwner.email,
        password: testEnv.familyOwner.password,
        firstName: 'Family',
        lastName: 'Owner',
        role: 'owner',
      };
      this.testUsers.set(familyOwner.id, familyOwner);
    }

    // Setup family member if configured
    if (testEnv.familyMember) {
      const familyMember: ClerkTestUser = {
        id: 'family-member',
        email: testEnv.familyMember.email,
        password: testEnv.familyMember.password,
        firstName: 'Family',
        lastName: 'Member',
        role: 'member',
      };
      this.testUsers.set(familyMember.id, familyMember);
    }

    this.initialized = true;
  }

  /**
   * Create a new test user
   * 
   * Note: In a real implementation, this would use Clerk's API to create users.
   * For now, it manages test user data locally.
   */
  static async createTestUser(options: CreateTestUserOptions): Promise<ClerkTestUser> {
    await this.initialize();

    const testUser: ClerkTestUser = {
      id: `test-user-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
      email: options.email,
      password: options.password,
      firstName: options.firstName,
      lastName: options.lastName,
      role: options.role || 'member',
    };

    // In a real implementation, you would call Clerk's API here:
    // const clerkUser = await clerkClient.users.createUser({
    //   emailAddress: [options.email],
    //   password: options.password,
    //   firstName: options.firstName,
    //   lastName: options.lastName,
    // });
    // testUser.clerkUserId = clerkUser.id;

    this.testUsers.set(testUser.id, testUser);
    return testUser;
  }

  /**
   * Get a test user by ID
   */
  static async getTestUser(id: string): Promise<ClerkTestUser | null> {
    await this.initialize();
    return this.testUsers.get(id) || null;
  }

  /**
   * Get the default test user
   */
  static async getDefaultTestUser(): Promise<ClerkTestUser> {
    await this.initialize();
    const user = this.testUsers.get('default-test-user');
    if (!user) {
      throw new Error('Default test user not found. Check your environment configuration.');
    }
    return user;
  }

  /**
   * Get family owner test user
   */
  static async getFamilyOwnerUser(): Promise<ClerkTestUser> {
    await this.initialize();
    const user = this.testUsers.get('family-owner');
    if (!user) {
      throw new Error('Family owner test user not configured. Check your environment variables.');
    }
    return user;
  }

  /**
   * Get family member test user
   */
  static async getFamilyMemberUser(): Promise<ClerkTestUser> {
    await this.initialize();
    const user = this.testUsers.get('family-member');
    if (!user) {
      throw new Error('Family member test user not configured. Check your environment variables.');
    }
    return user;
  }

  /**
   * Get all test users
   */
  static async getAllTestUsers(): Promise<ClerkTestUser[]> {
    await this.initialize();
    return Array.from(this.testUsers.values());
  }

  /**
   * Update test user session token
   */
  static async updateUserSessionToken(userId: string, sessionToken: string): Promise<void> {
    const user = this.testUsers.get(userId);
    if (user) {
      user.sessionToken = sessionToken;
    }
  }

  /**
   * Get test user session token
   */
  static async getUserSessionToken(userId: string): Promise<string | null> {
    const user = this.testUsers.get(userId);
    return user?.sessionToken || null;
  }

  /**
   * Delete a test user
   */
  static async deleteTestUser(id: string): Promise<void> {
    const user = this.testUsers.get(id);
    if (!user) return;

    // In a real implementation, you would call Clerk's API here:
    // if (user.clerkUserId) {
    //   await clerkClient.users.deleteUser(user.clerkUserId);
    // }

    this.testUsers.delete(id);
  }

  /**
   * Cleanup all test users
   */
  static async cleanupAllTestUsers(): Promise<void> {
    const users = Array.from(this.testUsers.values());
    
    for (const user of users) {
      // Skip default users from environment
      if (['default-test-user', 'family-owner', 'family-member'].includes(user.id)) {
        continue;
      }
      
      await this.deleteTestUser(user.id);
    }
  }

  /**
   * Generate unique test user data
   */
  static generateTestUserData(role: 'owner' | 'member' = 'member'): CreateTestUserOptions {
    const timestamp = Date.now();
    const randomId = Math.random().toString(36).substr(2, 9);
    
    return {
      email: `test-user-${timestamp}-${randomId}@example.com`,
      password: 'TestPassword123!',
      firstName: `Test${randomId}`,
      lastName: `User${timestamp}`,
      role,
    };
  }

  /**
   * Create test user for specific role
   */
  static async createOwnerUser(): Promise<ClerkTestUser> {
    const userData = this.generateTestUserData('owner');
    return await this.createTestUser(userData);
  }

  /**
   * Create test user for member role
   */
  static async createMemberUser(): Promise<ClerkTestUser> {
    const userData = this.generateTestUserData('member');
    return await this.createTestUser(userData);
  }

  /**
   * Validate test user credentials
   */
  static validateTestUser(user: ClerkTestUser): boolean {
    return !!(
      user.email &&
      user.password &&
      user.firstName &&
      user.lastName &&
      user.role
    );
  }

  /**
   * Get test users by role
   */
  static async getTestUsersByRole(role: 'owner' | 'member'): Promise<ClerkTestUser[]> {
    await this.initialize();
    return Array.from(this.testUsers.values()).filter(user => user.role === role);
  }

  /**
   * Reset all test user sessions
   */
  static async resetAllSessions(): Promise<void> {
    for (const user of this.testUsers.values()) {
      user.sessionToken = undefined;
    }
  }
}