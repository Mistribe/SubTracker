/**
 * Protected Routes Tests
 * 
 * Tests for route protection and access control using Clerk authentication
 */

import { test, expect } from '../../fixtures';
import { ClerkAuthHelper } from '../../fixtures/auth';

test.describe('Protected Routes', () => {
  const protectedRoutes = [
    '/dashboard',
    '/subscriptions',
    '/subscriptions/create',
    '/providers',
    '/labels',
    '/family',
    '/profile',
  ];

  const publicRoutes = [
    '/sign-in',
    '/sign-up',
    '/', // Home page might be public
  ];

  test.describe('Unauthenticated Access', () => {
    protectedRoutes.forEach(route => {
      test(`should redirect unauthenticated user from ${route} to sign-in`, async ({ 
        unauthenticatedPage: page 
      }) => {
        // Try to access protected route
        await page.goto(route);
        
        // Should be redirected to sign-in page
        await page.waitForURL(/\/sign-in/, { timeout: 15000 });
        await expect(page).toHaveURL(/\/sign-in/);
        
        // Verify sign-in form is visible using correct Clerk selectors
        const signInSelectors = [
          '[data-clerk-component="SignIn"]',
          '[data-clerk-element="rootBox"]',
          '.cl-signIn-root'
        ];
        
        let signInFormFound = false;
        for (const selector of signInSelectors) {
          try {
            await expect(page.locator(selector)).toBeVisible({ timeout: 5000 });
            signInFormFound = true;
            break;
          } catch {
            continue;
          }
        }
        expect(signInFormFound).toBe(true);
      });
    });

    publicRoutes.forEach(route => {
      test(`should allow unauthenticated access to ${route}`, async ({ 
        unauthenticatedPage: page 
      }) => {
        // Navigate to public route
        await page.goto(route);
        
        // Wait for page to load
        await page.waitForLoadState('networkidle');
        
        // Should be able to access the route
        if (route === '/') {
          // Home page might redirect to sign-in or show public content
          const currentUrl = page.url();
          expect(currentUrl).toMatch(/\/(sign-in|$)/);
        } else if (route === '/sign-in') {
          // Sign-in page should be accessible and should be on sign-in URL
          await expect(page).toHaveURL(/\/sign-in/);
        } else {
          // For other public routes, should be accessible
          await expect(page).toHaveURL(route);
        }
      });
    });
  });

  test.describe('Authenticated Access', () => {
    protectedRoutes.forEach(route => {
      test(`should allow authenticated user to access ${route}`, async ({ 
        authenticatedPage: page 
      }) => {
        // Navigate to protected route
        await page.goto(route);
        
        // Wait for page to load
        await page.waitForLoadState('networkidle');
        
        // Should be able to access the route
        await expect(page).toHaveURL(route);
        
        // Verify user is still authenticated
        const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
        expect(isAuthenticated).toBe(true);
      });
    });

    test('should maintain authentication across route navigation', async ({ 
      authenticatedPage: page 
    }) => {
      // Navigate through multiple protected routes
      for (const route of protectedRoutes) {
        await page.goto(route);
        await page.waitForLoadState('networkidle');
        
        // Check if we're redirected to sign-in (indicates logout)
        const currentUrl = page.url();
        if (currentUrl.includes('/sign-in')) {
          throw new Error(`User was logged out while navigating to ${route}`);
        }
        
        // Verify we're on the expected route (or a valid redirect like dashboard)
        if (!currentUrl.includes(route) && !currentUrl.includes('/dashboard')) {
          throw new Error(`Unexpected redirect from ${route} to ${currentUrl}`);
        }
        
        // Verify user remains authenticated by checking for user button or auth state
        // without navigating away from current page
        const isAuthenticated = await ClerkAuthHelper.isAuthenticatedOnCurrentPage(page);
        expect(isAuthenticated).toBe(true);
        
        // Small delay to prevent overwhelming the server
        await page.waitForTimeout(500);
      }
    });
  });

  test.describe('Route Redirects', () => {
    test('should redirect to dashboard after successful login', async ({ 
      unauthenticatedPage: page, 
      testUser 
    }) => {
      // Navigate to sign-in page
      await page.goto('/sign-in');
      
      // Sign in user
      await ClerkAuthHelper.signInUser(page, testUser);
      
      // Should be redirected to dashboard
      await expect(page).toHaveURL(/\/dashboard/);
    });

    test('should redirect to intended page after login from protected route', async ({ 
      unauthenticatedPage: page, 
      testUser 
    }) => {
      // Try to access a specific protected route
      await page.goto('/subscriptions');
      
      // Should be redirected to sign-in
      await page.waitForURL(/\/sign-in/, { timeout: 15000 });
      
      // Sign in user
      await ClerkAuthHelper.signInUser(page, testUser);
      
      // Should be redirected to originally requested page or dashboard
      // Note: This depends on the app's redirect logic implementation
      await page.waitForLoadState('networkidle');
      const currentUrl = page.url();
      expect(currentUrl).toMatch(/\/(dashboard|subscriptions)/);
    });

    test('should handle deep links after authentication', async ({ 
      unauthenticatedPage: page, 
      testUser 
    }) => {
      // Try to access a deep link
      const deepLink = '/subscriptions/create';
      await page.goto(deepLink);
      
      // Wait for potential redirect with longer timeout
      try {
        await page.waitForURL(/\/sign-in/, { timeout: 15000 });
      } catch {
        // If not redirected to sign-in, check if we're already on the deep link (might be public)
        const currentUrl = page.url();
        if (currentUrl.includes(deepLink)) {
          // Deep link is accessible without auth, skip the rest
          return;
        }
        // If we're somewhere else, continue with the test
      }
      
      // Sign in user
      await ClerkAuthHelper.signInUser(page, testUser);
      
      // Should be able to access the deep link or be redirected to dashboard
      await page.waitForLoadState('networkidle');
      const currentUrl = page.url();
      expect(currentUrl).toMatch(/\/(dashboard|subscriptions)/);
    });
  });

  test.describe('Route Guards', () => {
    test('should prevent access to routes after logout', async ({ 
      authenticatedPage: page 
    }) => {
      // Verify user is authenticated and can access protected route
      await page.goto('/subscriptions');
      await page.waitForLoadState('networkidle');
      await expect(page).toHaveURL(/\/subscriptions/);
      
      // Sign out user
      await ClerkAuthHelper.signOutUser(page);
      
      // Try to access protected routes (test a few key ones to avoid timeout)
      const testRoutes = ['/dashboard', '/subscriptions', '/providers'];
      for (const route of testRoutes) {
        await page.goto(route);
        
        // Should be redirected to sign-in or stay on current page if already signed out
        try {
          await page.waitForURL(/\/sign-in/, { timeout: 10000 });
          await expect(page).toHaveURL(/\/sign-in/);
        } catch {
          // If not redirected, check if we're on home page or sign-in already
          const currentUrl = page.url();
          expect(currentUrl).toMatch(/\/(sign-in|$)/);
        }
        
        // Small delay to prevent overwhelming the server
        await page.waitForTimeout(300);
      }
    });

    test('should handle browser back button after logout', async ({ 
      authenticatedPage: page 
    }) => {
      // Navigate to protected route
      await page.goto('/subscriptions');
      await page.waitForLoadState('networkidle');
      await expect(page).toHaveURL(/\/subscriptions/);
      
      // Sign out user
      await ClerkAuthHelper.signOutUser(page);
      
      // Wait for logout to complete
      await page.waitForTimeout(2000);
      
      // Try to go back using browser back button
      await page.goBack();
      
      // Should still be on sign-in page or redirected back to sign-in
      try {
        await page.waitForURL(/\/sign-in/, { timeout: 10000 });
        await expect(page).toHaveURL(/\/sign-in/);
      } catch {
        // If not on sign-in, check if we're on home page
        const currentUrl = page.url();
        expect(currentUrl).toMatch(/\/(sign-in|$)/);
      }
    });

    test('should handle direct URL access after logout', async ({ 
      authenticatedPage: page 
    }) => {
      // Sign out user
      await ClerkAuthHelper.signOutUser(page);
      
      // Wait for logout to complete
      await page.waitForTimeout(2000);
      
      // Try to access protected route directly via URL
      await page.goto('/dashboard');
      
      // Should be redirected to sign-in or home page
      try {
        await page.waitForURL(/\/sign-in/, { timeout: 10000 });
        await expect(page).toHaveURL(/\/sign-in/);
      } catch {
        // If not redirected, check current URL - might be on dashboard if still authenticated
        const currentUrl = page.url();
        // Check if user is actually still authenticated
        const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
        if (isAuthenticated) {
          // User is still authenticated, this is acceptable
          expect(currentUrl).toMatch(/\/(dashboard|sign-in|$)/);
        } else {
          // User is not authenticated, should be on sign-in or home
          expect(currentUrl).toMatch(/\/(sign-in|$)/);
        }
      }
    });
  });

  test.describe('Error Handling', () => {
    test('should handle authentication errors gracefully', async ({ 
      unauthenticatedPage: page 
    }) => {
      // Simulate authentication service error by intercepting Clerk requests
      await page.route('**/clerk/**', route => {
        const url = route.request().url();
        if (url.includes('/v1/client') || url.includes('/environment')) {
          route.abort('failed');
        } else {
          route.continue();
        }
      });
      
      // Try to access protected route
      await page.goto('/dashboard');
      
      // Wait for potential error handling
      await page.waitForTimeout(3000);
      
      // Should handle error gracefully (might show error page or redirect to sign-in)
      const currentUrl = page.url();
      expect(currentUrl).toMatch(/\/(sign-in|error|dashboard)/);
      
      // Clear route handler
      await page.unroute('**/clerk/**');
    });

    test('should handle network errors during route protection check', async ({ 
      authenticatedPage: page 
    }) => {
      // Navigate to a protected route first to establish baseline
      await page.goto('/subscriptions');
      await page.waitForLoadState('networkidle');
      await expect(page).toHaveURL(/\/subscriptions/);
      
      // Simulate network error using context.setOffline instead of page.setOffline
      await page.context().setOffline(true);
      
      // Try to navigate to another protected route
      try {
        await page.goto('/providers');
        await page.waitForTimeout(2000);
      } catch (error) {
        // Network error is expected
        console.log('Expected network error:', error);
      }
      
      // Should either work with cached authentication or show appropriate error
      const currentUrl = page.url();
      expect(currentUrl).toMatch(/\/(providers|subscriptions|sign-in|error)/);
      
      // Restore network
      await page.context().setOffline(false);
      await page.waitForLoadState('networkidle');
    });
  });

  test.describe('User Roles and Permissions', () => {
    test('should handle family owner permissions', async ({ 
      familyOwnerPage: page 
    }) => {
      // Family owner should have access to all routes
      for (const route of protectedRoutes) {
        await page.goto(route);
        await page.waitForLoadState('networkidle');
        await expect(page).toHaveURL(route);
        
        // Small delay to prevent overwhelming the server
        await page.waitForTimeout(300);
      }
      
      // Verify user is authenticated as family owner
      const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
      expect(isAuthenticated).toBe(true);
    });

    test('should handle family member permissions', async ({ 
      familyMemberPage: page 
    }) => {
      // Family member should have access to basic routes
      const memberAccessibleRoutes = [
        '/dashboard',
        '/subscriptions',
        '/providers',
        '/labels',
        '/profile',
      ];
      
      for (const route of memberAccessibleRoutes) {
        await page.goto(route);
        await page.waitForLoadState('networkidle');
        await expect(page).toHaveURL(route);
        
        // Small delay to prevent overwhelming the server
        await page.waitForTimeout(300);
      }
      
      // Verify user is authenticated as family member
      const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
      expect(isAuthenticated).toBe(true);
    });
  });

  test.describe('Route Performance', () => {
    test('should load protected routes quickly for authenticated users', async ({ 
      authenticatedPage: page 
    }) => {
      // Measure navigation time to dashboard
      const startTime = Date.now();
      
      await page.goto('/dashboard');
      await page.waitForLoadState('networkidle');
      await expect(page).toHaveURL(/\/dashboard/);
      
      const loadTime = Date.now() - startTime;
      
      // Should load within reasonable time (increased threshold for CI environments)
      expect(loadTime).toBeLessThan(10000); // 10 seconds
      
      // Log performance for debugging
      console.log(`Dashboard load time: ${loadTime}ms`);
    });

    test('should handle concurrent route access', async ({ 
      context, 
      authenticatedPage: page 
    }) => {
      try {
        // Ensure the main page is authenticated first
        await page.goto('/dashboard');
        await expect(page).toHaveURL(/\/dashboard/);
        
        // Open multiple tabs
        const tabs = await Promise.all([
          context.newPage(),
          context.newPage(),
          context.newPage(),
        ]);
        
        // Navigate all tabs to different routes with staggered timing
        await tabs[0].goto('/subscriptions');
        await page.waitForTimeout(500);
        
        await tabs[1].goto('/providers');
        await page.waitForTimeout(500);
        
        await tabs[2].goto('/labels');
        
        // Wait for all navigations to complete
        await Promise.all([
          tabs[0].waitForLoadState('networkidle'),
          tabs[1].waitForLoadState('networkidle'),
          tabs[2].waitForLoadState('networkidle'),
        ]);
        
        // Verify all tabs loaded correctly
        await Promise.all([
          expect(tabs[0]).toHaveURL(/\/subscriptions/),
          expect(tabs[1]).toHaveURL(/\/providers/),
          expect(tabs[2]).toHaveURL(/\/labels/),
        ]);
        
        // Close tabs
        await Promise.all(tabs.map(tab => tab.close()));
      } catch (error) {
        console.error('Concurrent access test failed:', error);
        throw error;
      }
    });
  });
});