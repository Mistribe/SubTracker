/**
 * Authentication Happy Path Tests
 * 
 * Comprehensive test suite covering successful authentication workflows:
 * - Login with valid credentials
 * - Logout successfully
 * - Session persistence across navigation
 * - Access to protected routes after authentication
 * 
 * Requirements: 8.1, 8.2, 8.3, 8.4, 8.5
 * 
 * Testing Philosophy: Happy Path Only
 * ✅ Tests successful login and logout flows
 * ✅ Tests successful navigation to protected routes
 * ✅ Tests session persistence and authentication state
 * ❌ No authentication failures, invalid credentials, or error scenarios
 */

import { test, expect } from '../../fixtures';
import { ClerkAuthHelper } from '../../fixtures/auth';

test.describe('Authentication Happy Path', () => {
  test.beforeEach(async ({ page }) => {
    console.log('🚀 Setting up authentication test');
    
    // Ensure clean state by navigating to sign-in page
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);
    console.log('✅ Navigated to sign-in page and Clerk loaded');
  });

  test('should successfully login with valid credentials', async ({
    unauthenticatedPage: page,
    testUser
  }) => {
    console.log('🔐 Testing successful login flow');
    
    // Navigate to sign-in page
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);
    console.log('✅ Sign-in page loaded');

    // Fill in valid credentials
    await page.fill('input[name="identifier"]', testUser.email);
    await page.fill('input[name="password"]', testUser.password);
    console.log('✅ Filled in login credentials');

    // Click sign in button
    await page.click('button:has-text("Continue")');
    console.log('✅ Clicked sign in button');

    // Verify successful login and redirect to dashboard
    await page.waitForURL(/\/dashboard/, { timeout: 15000 });
    await expect(page).toHaveURL(/\/dashboard/);
    console.log('✅ Successfully redirected to dashboard');

    // Verify user is authenticated
    const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ User authentication confirmed');

    // Verify dashboard page loads correctly
    const dashboardTitle = page.locator('h1:has-text("Dashboard")');
    await expect(dashboardTitle).toBeVisible();
    console.log('✅ Dashboard page loaded successfully');
  });

  test('should successfully logout authenticated user', async ({ 
    authenticatedPage: page 
  }) => {
    console.log('🚪 Testing successful logout flow');
    
    // Verify user is authenticated and on dashboard
    await expect(page).toHaveURL(/\/dashboard/);
    const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Confirmed user is authenticated');

    // Sign out user using Clerk's user menu
    const userButton = page.locator('.cl-userButton, [data-clerk-component="UserButton"]').first();
    await userButton.waitFor({ state: 'visible', timeout: 5000 });
    await userButton.click();
    console.log('✅ Opened user menu');

    // Click sign out option
    const signOutButton = page.locator('button:has-text("Sign out"), .cl-userButtonPopoverActionButton:has-text("Sign out")').first();
    await signOutButton.waitFor({ state: 'visible', timeout: 5000 });
    await signOutButton.click();
    console.log('✅ Clicked sign out button');

    // Wait for redirect to home page
    await page.waitForURL('/', { timeout: 15000 });
    await expect(page).toHaveURL('/');
    console.log('✅ Successfully redirected to home page');

    // Verify user is no longer authenticated
    const isStillAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isStillAuthenticated).toBe(false);
    console.log('✅ User logout confirmed');
  });

  test('should maintain authentication across page navigation', async ({ 
    authenticatedPage: page 
  }) => {
    console.log('🧭 Testing authentication persistence across navigation');
    
    // Start on dashboard
    await expect(page).toHaveURL(/\/dashboard/);
    let isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Confirmed authentication on dashboard');

    // Navigate to subscriptions page
    await page.goto('/subscriptions');
    await expect(page).toHaveURL(/\/subscriptions/);
    isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Authentication maintained on subscriptions page');

    // Navigate to providers page
    await page.goto('/providers');
    await expect(page).toHaveURL(/\/providers/);
    isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Authentication maintained on providers page');

    // Navigate to profile page
    await page.goto('/profile');
    await expect(page).toHaveURL(/\/profile/);
    isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Authentication maintained on profile page');

    // Navigate back to dashboard
    await page.goto('/dashboard');
    await expect(page).toHaveURL(/\/dashboard/);
    isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Authentication maintained when returning to dashboard');
  });

  test('should access protected routes successfully when authenticated', async ({ 
    authenticatedPage: page 
  }) => {
    console.log('🔒 Testing access to protected routes');
    
    const protectedRoutes = [
      '/dashboard',
      '/subscriptions',
      '/providers',
      '/labels',
      '/family',
      '/profile'
    ];

    for (const route of protectedRoutes) {
      console.log(`🔍 Testing access to ${route}`);
      
      // Navigate to protected route
      await page.goto(route);
      await expect(page).toHaveURL(route);
      console.log(`✅ Successfully accessed ${route}`);
      
      // Verify user is still authenticated
      const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
      expect(isAuthenticated).toBe(true);
      console.log(`✅ Authentication confirmed on ${route}`);
      
      // Verify page loads without redirect to sign-in
      const currentUrl = page.url();
      expect(currentUrl).not.toContain('/sign-in');
      console.log(`✅ No redirect to sign-in from ${route}`);
    }
  });

  test('should maintain session after page refresh', async ({ 
    authenticatedPage: page 
  }) => {
    console.log('🔄 Testing session persistence after page refresh');
    
    // Verify initial authentication
    await expect(page).toHaveURL(/\/dashboard/);
    let isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Initial authentication confirmed');

    // Refresh the page
    await page.reload();
    await page.waitForLoadState('networkidle');
    console.log('✅ Page refreshed');

    // Verify still on dashboard (not redirected to sign-in)
    await expect(page).toHaveURL(/\/dashboard/);
    console.log('✅ Still on dashboard after refresh');

    // Verify authentication is maintained
    isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Authentication maintained after refresh');

    // Verify dashboard functionality is still available
    const dashboardTitle = page.locator('h1:has-text("Dashboard")');
    await expect(dashboardTitle).toBeVisible();
    console.log('✅ Dashboard functionality confirmed after refresh');
  });

  test('should complete full authentication cycle successfully', async ({
    unauthenticatedPage: page,
    testUser
  }) => {
    console.log('🔄 Testing complete authentication cycle (login → navigate → logout)');
    
    // Step 1: Login
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);
    await page.fill('input[name="identifier"]', testUser.email);
    await page.fill('input[name="password"]', testUser.password);
    await page.click('button:has-text("Continue")');
    await page.waitForURL(/\/dashboard/, { timeout: 15000 });
    console.log('✅ Step 1: Login successful');
    
    // Step 2: Navigate to different pages
    await page.goto('/subscriptions');
    await expect(page).toHaveURL(/\/subscriptions/);
    await page.goto('/providers');
    await expect(page).toHaveURL(/\/providers/);
    console.log('✅ Step 2: Navigation successful');
    
    // Step 3: Return to dashboard
    await page.goto('/dashboard');
    await expect(page).toHaveURL(/\/dashboard/);
    const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
    console.log('✅ Step 3: Return to dashboard successful');
    
    // Step 4: Logout
    const userButton = page.locator('.cl-userButton, [data-clerk-component="UserButton"]').first();
    await userButton.waitFor({ state: 'visible', timeout: 5000 });
    await userButton.click();
    const signOutButton = page.locator('button:has-text("Sign out"), .cl-userButtonPopoverActionButton:has-text("Sign out")').first();
    await signOutButton.waitFor({ state: 'visible', timeout: 5000 });
    await signOutButton.click();
    await page.waitForURL('/', { timeout: 15000 });
    console.log('✅ Step 4: Logout successful');
    
    // Step 5: Verify logout
    const isStillAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isStillAuthenticated).toBe(false);
    console.log('✅ Step 5: Logout verification successful');
    
    console.log('🎉 Complete authentication cycle test completed successfully');
  });
});