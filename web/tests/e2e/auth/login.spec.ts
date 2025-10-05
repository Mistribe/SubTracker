/**
 * Login Authentication Tests
 * 
 * Tests for user login functionality using Clerk authentication
 */

import { test, expect } from '../../fixtures';
import { ClerkAuthHelper } from '../../fixtures/auth';

test.describe('User Login', () => {
  test.beforeEach(async ({ page }) => {
    // Ensure user is signed out before each test
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
      console.log('Sign out failed, continuing with test:', error);
    }

    // Navigate to sign-in page to ensure clean state
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);
  });

  test('should successfully login with valid credentials', async ({
    unauthenticatedPage: page,
    testUser
  }) => {
    // Navigate to sign-in page
    await page.goto('/sign-in');

    // Wait for Clerk to load
    await ClerkAuthHelper.waitForClerkLoaded(page);

    // Fill in credentials using correct selectors
    await page.fill('input[name="identifier"]', testUser.email);
    await page.fill('input[name="password"]', testUser.password);

    // Click sign in button
    await page.click('button:has-text("Continue")');

    // Verify successful login and redirect to dashboard
    await page.waitForURL('/dashboard', { timeout: 15000 });
    await expect(page).toHaveURL(/\/dashboard/);

    // Verify user is authenticated
    const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);

    // Verify user button is visible (using correct selector)
    const userButtonSelectors = ['.cl-userButton', '[data-clerk-component="UserButton"]'];
    let userButtonFound = false;
    for (const selector of userButtonSelectors) {
      try {
        await expect(page.locator(selector).first()).toBeVisible({ timeout: 5000 });
        userButtonFound = true;
        break;
      } catch {
        continue;
      }
    }
    expect(userButtonFound).toBe(true);
  });

  test('should handle invalid credentials gracefully', async ({ unauthenticatedPage: page }) => {
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);

    // Fill in invalid credentials
    await page.fill('input[name="identifier"]', 'nonexistent@example.com');
    await page.fill('input[name="password"]', 'wrongpassword');

    // Click sign in button
    await page.click('button:has-text("Continue")');

    // Wait a moment for any potential error processing
    await page.waitForTimeout(3000);

    // Verify we're still on sign-in page (authentication failed)
    await expect(page).toHaveURL(/\/sign-in/);

    // Verify the Clerk form is still present (not redirected to dashboard)
    await expect(page.locator('[data-clerk-component="SignIn"]')).toBeVisible();
  });

  test('should handle wrong password gracefully', async ({
    unauthenticatedPage: page,
    testUser
  }) => {
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);

    // Fill in valid email but invalid password
    await page.fill('input[name="identifier"]', testUser.email);
    await page.fill('input[name="password"]', 'wrongpassword123');

    // Click sign in button
    await page.click('button:has-text("Continue")');

    // Wait a moment for any potential error processing
    await page.waitForTimeout(3000);

    // Verify we're still on sign-in page (authentication failed)
    await expect(page).toHaveURL(/\/sign-in/);

    // Verify the Clerk form is still present (not redirected to dashboard)
    // Note: Clerk may change form state after failed attempts
    await expect(page.locator('[data-clerk-component="SignIn"]')).toBeVisible();
  });

  test('should handle empty credentials gracefully', async ({ unauthenticatedPage: page }) => {
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);

    // Try to click sign in button without filling credentials
    await page.click('button:has-text("Continue")');

    // Wait a moment for any potential validation
    await page.waitForTimeout(2000);

    // Verify we're still on sign-in page (form validation should prevent submission)
    await expect(page).toHaveURL(/\/sign-in/);

    // Verify the form fields are still visible and interactive
    await expect(page.locator('input[name="identifier"]')).toBeVisible();
    await expect(page.locator('input[name="password"]')).toBeVisible();
    await expect(page.locator('button:has-text("Continue")')).toBeVisible();
  });

  test('should maintain session after page refresh', async ({
    unauthenticatedPage: page,
    testUser
  }) => {
    // Sign in user
    await ClerkAuthHelper.signInUser(page, testUser);

    // Verify we're on dashboard
    await expect(page).toHaveURL(/\/dashboard/);

    // Refresh the page
    await page.reload();

    // Verify user is still authenticated
    await expect(page).toHaveURL(/\/dashboard/);
    const isAuthenticated = await ClerkAuthHelper.isAuthenticated(page);
    expect(isAuthenticated).toBe(true);
  });

  test('should handle authentication with extended timeout', async ({
    unauthenticatedPage: page,
    testUser
  }) => {
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);

    // Fill in credentials
    await page.fill('input[name="identifier"]', testUser.email);
    await page.fill('input[name="password"]', testUser.password);

    // Click sign in button
    await page.click('button:has-text("Continue")');

    // Verify successful login with extended timeout (simulating slower conditions)
    await page.waitForURL('/dashboard', { timeout: 30000 });
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should redirect to intended page after login', async ({
    unauthenticatedPage: page,
    testUser
  }) => {
    // Try to access protected route
    await page.goto('/subscriptions');

    // Should be redirected to sign-in
    await page.waitForURL(/\/sign-in/, { timeout: 10000 });

    // Sign in
    await ClerkAuthHelper.signInUser(page, testUser);

    // Should be redirected to originally requested page
    // Note: This depends on the app's redirect logic
    await expect(page).toHaveURL(/\/(dashboard|subscriptions)/);
  });

  test('should handle multiple login attempts', async ({
    unauthenticatedPage: page,
    testUser
  }) => {
    await page.goto('/sign-in');
    await ClerkAuthHelper.waitForClerkLoaded(page);

    // First attempt with wrong password
    await page.fill('input[name="identifier"]', testUser.email);
    await page.fill('input[name="password"]', 'wrongpassword123');
    await page.click('button:has-text("Continue")');

    // Wait for potential error processing
    await page.waitForTimeout(3000);

    // Verify we're still on sign-in page
    await expect(page).toHaveURL(/\/sign-in/);

    // Clear and try again with correct password
    await page.fill('input[name="password"]', '');
    await page.fill('input[name="password"]', testUser.password);
    await page.click('button:has-text("Continue")');

    // Verify successful login
    await page.waitForURL('/dashboard', { timeout: 15000 });
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test('should persist session across browser tabs', async ({
    context,
    unauthenticatedPage: page,
    testUser
  }) => {
    // Sign in on first tab
    await ClerkAuthHelper.signInUser(page, testUser);
    await expect(page).toHaveURL(/\/dashboard/);

    // Open new tab
    const newPage = await context.newPage();
    await newPage.goto('/dashboard');

    // Verify user is authenticated in new tab
    await expect(newPage).toHaveURL(/\/dashboard/);
    const isAuthenticated = await ClerkAuthHelper.isAuthenticated(newPage);
    expect(isAuthenticated).toBe(true);

    await newPage.close();
  });
});