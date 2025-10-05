/**
 * Authentication Fixtures for Playwright Tests
 * 
 * This module provides Clerk-based authentication utilities and fixtures
 * for E2E testing. It uses only Clerk's official testing library.
 */

import { test as base, expect, Page } from '@playwright/test';
import { clerkSetup } from '@clerk/testing/playwright';
import { getTestEnvironment } from '../config/test-env';

const testEnv = getTestEnvironment();

/**
 * Test user interface for type safety
 */
export interface TestUser {
  email: string;
  password: string;
  role?: 'owner' | 'member';
}

/**
 * Authentication fixture interface
 */
export interface AuthFixtures {
  authenticatedPage: Page;
  unauthenticatedPage: Page;
  familyOwnerPage: Page;
  familyMemberPage: Page;
  testUser: TestUser;
}

/**
 * Clerk authentication utilities
 */
export class ClerkAuthHelper {
  /**
   * Setup Clerk for testing environment
   */
  static async setupClerk() {
    await clerkSetup({
      publishableKey: testEnv.clerkPublishableKey,
      secretKey: testEnv.clerkSecretKey,
    });
  }

  /**
   * Sign in a user using Clerk authentication elements
   */
  static async signInUser(page: Page, user: TestUser): Promise<void> {
    // Navigate to sign-in page
    await page.goto('/sign-in');
    
    // Wait for Clerk to load
    await page.waitForSelector('[data-clerk-component="SignIn"]', { 
      timeout: testEnv.timeouts.navigation 
    });

    // Wait for the form to be fully loaded
    await page.waitForLoadState('networkidle');

    // Fill in email - try multiple possible selectors
    const emailSelectors = [
      'input[name="identifier"]',
      'input[type="email"]',
      'input[placeholder*="email" i]',
      '.cl-formFieldInput[name="identifier"]',
      '.cl-formFieldInput[type="email"]'
    ];
    
    let emailInput = null;
    for (const selector of emailSelectors) {
      try {
        emailInput = page.locator(selector).first();
        await emailInput.waitFor({ timeout: 2000 });
        break;
      } catch {
        continue;
      }
    }
    
    if (!emailInput) {
      throw new Error('Could not find email input field');
    }
    
    await emailInput.fill(user.email);

    // Fill in password - try multiple possible selectors
    const passwordSelectors = [
      'input[name="password"]',
      'input[type="password"]',
      '.cl-formFieldInput[name="password"]',
      '.cl-formFieldInput[type="password"]'
    ];
    
    let passwordInput = null;
    for (const selector of passwordSelectors) {
      try {
        passwordInput = page.locator(selector).first();
        await passwordInput.waitFor({ timeout: 2000 });
        break;
      } catch {
        continue;
      }
    }
    
    if (!passwordInput) {
      throw new Error('Could not find password input field');
    }
    
    await passwordInput.fill(user.password);

    // Click sign in button - try multiple possible selectors
    const buttonSelectors = [
      'button:has-text("Continue")',
      '.cl-formButtonPrimary',
      'button[type="submit"]:has-text("Continue")',
      'button:has-text("Sign in")',
      '[data-localization-key="signIn.start.actionText"]'
    ];
    
    let signInButton = null;
    for (const selector of buttonSelectors) {
      try {
        signInButton = page.locator(selector).first();
        // Wait for the button to be attached and check if it's the right one
        await signInButton.waitFor({ timeout: 2000 });
        
        // For the Continue button, make sure it's not a social login button
        const buttonText = await signInButton.textContent();
        const isVisible = await signInButton.isVisible();
        const isEnabled = await signInButton.isEnabled();
        
        if (isEnabled && buttonText && buttonText.trim().length > 0) {
          break;
        }
      } catch {
        continue;
      }
    }
    
    if (!signInButton) {
      throw new Error('Could not find sign in button');
    }
    
    // Click the button
    await signInButton.click();

    // Wait for successful authentication and redirect
    await page.waitForURL('/dashboard', { 
      timeout: testEnv.timeouts.navigation 
    });

    // Verify authentication was successful
    await expect(page).toHaveURL(/\/dashboard/);
  }

  /**
   * Sign out the current user
   */
  static async signOutUser(page: Page): Promise<void> {
    try {
      // Look for user menu or sign out button - try multiple selectors
      const userButtonSelectors = [
        '.cl-userButton',
        '[data-clerk-component="UserButton"]',
        '.cl-userButtonTrigger'
      ];
      
      let userButton = null;
      for (const selector of userButtonSelectors) {
        try {
          userButton = page.locator(selector).first();
          if (await userButton.isVisible({ timeout: 2000 })) {
            break;
          }
        } catch {
          continue;
        }
      }
      
      if (userButton && await userButton.isVisible()) {
        await userButton.click();
        
        // Click sign out option - try multiple selectors
        const signOutSelectors = [
          'button:has-text("Sign out")',
          '.cl-userButtonPopoverActionButton:has-text("Sign out")',
          '[data-localization-key="userButton.action__signOut"]',
          '.cl-menuItem:has-text("Sign out")'
        ];
        
        let signOutButton = null;
        for (const selector of signOutSelectors) {
          try {
            signOutButton = page.locator(selector).first();
            if (await signOutButton.isVisible({ timeout: 2000 })) {
              break;
            }
          } catch {
            continue;
          }
        }
        
        if (signOutButton) {
          await signOutButton.click();
          // Wait for redirect to home page (Clerk's afterSignOutUrl is "/")
          await page.waitForURL('/', { 
            timeout: testEnv.timeouts.navigation 
          });
        } else {
          // Fallback: navigate to home page directly
          await page.goto('/');
        }
      } else {
        // Alternative: navigate to home page directly
        await page.goto('/');
      }

      // Verify sign out was successful by checking authentication state
      const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
      if (isAuthenticated) {
        // If still authenticated, try navigating to home page again
        await page.goto('/');
        await page.waitForTimeout(1000);
      }
    } catch (error) {
      // If all else fails, just navigate to home page
      console.warn('Sign out failed, using fallback method:', error);
      await page.goto('/');
    }
  }

  /**
   * Check if user is authenticated on the current page without navigation
   */
  static async isAuthenticatedOnCurrentPage(page: Page): Promise<boolean> {
    try {
      const currentUrl = page.url();
      
      // If we're on sign-in page, definitely not authenticated
      if (currentUrl.includes('/sign-in')) {
        return false;
      }
      
      // Check for user button as primary indicator
      const userButtonSelectors = [
        '.cl-userButton',
        '[data-clerk-component="UserButton"]',
        '.cl-userButtonTrigger'
      ];
      
      for (const selector of userButtonSelectors) {
        try {
          const userButton = page.locator(selector).first();
          if (await userButton.isVisible({ timeout: 2000 })) {
            return true;
          }
        } catch {
          continue;
        }
      }
      
      // Fallback: check if we're on a protected route and not redirected to sign-in
      const protectedRoutes = ['/dashboard', '/subscriptions', '/providers', '/labels', '/family', '/profile'];
      const isOnProtectedRoute = protectedRoutes.some(route => currentUrl.includes(route));
      
      if (isOnProtectedRoute) {
        // Wait a bit to see if we get redirected to sign-in
        await page.waitForTimeout(1000);
        const updatedUrl = page.url();
        return !updatedUrl.includes('/sign-in');
      }
      
      return false;
    } catch {
      return false;
    }
  }

  /**
   * Check if user is authenticated (legacy method that navigates)
   */
  static async isAuthenticated(page: Page): Promise<boolean> {
    try {
      // First check: try to access a protected route and see if we get redirected
      const currentUrl = page.url();
      
      // If we're already on sign-in page, definitely not authenticated
      if (currentUrl.includes('/sign-in')) {
        return false;
      }
      
      // Try to navigate to dashboard and see if we get redirected
      await page.goto('/dashboard');
      await page.waitForTimeout(2000); // Give time for potential redirect
      
      const newUrl = page.url();
      
      // If we're redirected to sign-in, not authenticated
      if (newUrl.includes('/sign-in')) {
        return false;
      }
      
      // If we're on dashboard, check for user button as secondary confirmation
      if (newUrl.includes('/dashboard')) {
        const userButtonSelectors = [
          '.cl-userButton',
          '[data-clerk-component="UserButton"]',
          '.cl-userButtonTrigger'
        ];
        
        for (const selector of userButtonSelectors) {
          try {
            const userButton = page.locator(selector).first();
            if (await userButton.isVisible({ timeout: 1000 })) {
              return true;
            }
          } catch {
            continue;
          }
        }
        
        // If on dashboard but no user button found, still consider authenticated
        // (user button might not have loaded yet)
        return true;
      }
      
      return false;
    } catch {
      return false;
    }
  }

  /**
   * Wait for Clerk to be fully loaded
   */
  static async waitForClerkLoaded(page: Page): Promise<void> {
    await page.waitForSelector('[data-clerk-component="SignIn"]', {
      timeout: testEnv.timeouts.navigation
    });
    
    // Also wait for the form to be interactive
    await page.waitForLoadState('networkidle');
  }

  /**
   * Get current user session information
   */
  static async getSessionInfo(page: Page): Promise<any> {
    return await page.evaluate(() => {
      // Access Clerk session information from window object
      return (window as any).__clerk_session_info || null;
    });
  }
}

/**
 * Extended test with authentication fixtures
 */
export const test = base.extend<AuthFixtures>({
  /**
   * Authenticated page fixture - signs in the default test user
   */
  authenticatedPage: async ({ page }, use) => {
    const testUser: TestUser = {
      email: testEnv.testUser.email,
      password: testEnv.testUser.password,
      role: 'owner'
    };

    // Sign in the user
    await ClerkAuthHelper.signInUser(page, testUser);

    // Provide the authenticated page to the test
    await use(page);

    // Cleanup: sign out after test (simplified to avoid timeouts)
    try {
      // Check if network is available first
      const isOffline = await page.evaluate(() => !navigator.onLine);
      if (!isOffline) {
        // Just navigate to home page to trigger logout
        await page.goto('/');
        await page.waitForTimeout(1000);
      } else {
        console.log('Network is offline, skipping cleanup navigation');
      }
    } catch (error) {
      console.warn('Failed to sign out user during cleanup:', error.message || error);
    }
  },

  /**
   * Unauthenticated page fixture - ensures user is signed out
   */
  unauthenticatedPage: async ({ page }, use) => {
    try {
      // Navigate to home page first
      await page.goto('/');
      
      // Check if user is already signed out
      const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
      
      if (isAuthenticated) {
        // Try to sign out using Clerk's sign out functionality
        await ClerkAuthHelper.signOutUser(page);
      }
    } catch (error) {
      console.log('Sign out failed in fixture, continuing:', error);
    }
    
    // Navigate to sign-in page to ensure clean state
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);

    // Provide the unauthenticated page to the test
    await use(page);
  },

  /**
   * Family owner page fixture - signs in as family owner
   */
  familyOwnerPage: async ({ page }, use) => {
    const familyOwner: TestUser = testEnv.familyOwner ? {
      email: testEnv.familyOwner.email,
      password: testEnv.familyOwner.password,
      role: 'owner'
    } : {
      email: testEnv.testUser.email,
      password: testEnv.testUser.password,
      role: 'owner'
    };

    // Sign in the family owner
    await ClerkAuthHelper.signInUser(page, familyOwner);

    // Provide the authenticated page to the test
    await use(page);

    // Cleanup: sign out after test (simplified to avoid timeouts)
    try {
      // Check if network is available first
      const isOffline = await page.evaluate(() => !navigator.onLine);
      if (!isOffline) {
        await page.goto('/');
        await page.waitForTimeout(1000);
      } else {
        console.log('Network is offline, skipping cleanup navigation');
      }
    } catch (error) {
      console.warn('Failed to sign out family owner during cleanup:', error.message || error);
    }
  },

  /**
   * Family member page fixture - signs in as family member
   */
  familyMemberPage: async ({ page }, use) => {
    const familyMember: TestUser = testEnv.familyMember ? {
      email: testEnv.familyMember.email,
      password: testEnv.familyMember.password,
      role: 'member'
    } : {
      email: testEnv.testUser.email,
      password: testEnv.testUser.password,
      role: 'member'
    };

    // Sign in the family member
    await ClerkAuthHelper.signInUser(page, familyMember);

    // Provide the authenticated page to the test
    await use(page);

    // Cleanup: sign out after test (simplified to avoid timeouts)
    try {
      // Check if network is available first
      const isOffline = await page.evaluate(() => !navigator.onLine);
      if (!isOffline) {
        await page.goto('/');
        await page.waitForTimeout(1000);
      } else {
        console.log('Network is offline, skipping cleanup navigation');
      }
    } catch (error) {
      console.warn('Failed to sign out family member during cleanup:', error.message || error);
    }
  },

  /**
   * Test user fixture - provides current test user information
   */
  testUser: async ({}, use) => {
    const testUser: TestUser = {
      email: testEnv.testUser.email,
      password: testEnv.testUser.password,
      role: 'owner'
    };

    await use(testUser);
  },
});

/**
 * Export expect for convenience
 */
export { expect } from '@playwright/test';

/**
 * Export AuthFixtures interface for type checking
 */
export type { AuthFixtures };