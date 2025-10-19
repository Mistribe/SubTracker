/**
 * Profile Settings Happy Path Tests
 * 
 * Comprehensive test suite covering successful profile settings workflows:
 * - Navigate to profile settings
 * - View profile information
 * - Verify settings page accessibility
 * 
 * Requirements: 5.1
 * 
 * Testing Philosophy: Happy Path Only
 * ✅ Tests successful navigation and display
 * ✅ Tests page structure and accessibility
 * ❌ No error testing, validation failures, or edge cases
 * 
 * Note: The profile page in this application is primarily focused on preferences
 * (currency and theme). User profile information (name, email, avatar) is managed
 * through Clerk authentication service via the UserButton component in the sidebar.
 */

import { test, expect } from '../../fixtures/auth';
import { ProfilePage } from '../../page-objects/profile-page';

test.describe('Profile Settings Happy Path', () => {
    let profilePage: ProfilePage;

    test.beforeEach(async ({ authenticatedPage }) => {
        console.log('🚀 Setting up profile settings test');

        profilePage = new ProfilePage(authenticatedPage);

        // Navigate to profile page
        await authenticatedPage.goto('/profile');
        await profilePage.waitForPageLoad();
        console.log('✅ Navigated to profile settings page');
    });

    test('Navigate to profile settings page successfully', async () => {
        console.log('🧪 Test: Navigate to profile settings page successfully');

        // Verify we're on the profile page
        await expect(profilePage.pageInstance).toHaveURL(/\/profile/);

        // Verify page title is visible
        await profilePage.verifyPageDisplayed();

        console.log('✅ Test passed: Profile settings page navigation successful');
    });

    test('Display profile settings page structure correctly', async () => {
        console.log('🧪 Test: Display profile settings page structure correctly');

        // Verify page displays correctly
        await profilePage.verifyPageDisplayed();

        // Verify the page has the essential sections
        const currencyLabel = profilePage.pageInstance.locator('text="Preferred Currency"');
        const themeLabel = profilePage.pageInstance.locator('text="Theme"');

        await expect(currencyLabel).toBeVisible();
        await expect(themeLabel).toBeVisible();

        console.log('✅ Test passed: Profile settings page structure displayed correctly');
    });

    test('Access profile settings from sidebar navigation', async ({ authenticatedPage }) => {
        console.log('🧪 Test: Access profile settings from sidebar navigation');

        // Navigate away to dashboard first
        await authenticatedPage.goto('/dashboard');
        await authenticatedPage.waitForURL('/dashboard');

        // Find and click profile/preferences link in sidebar
        const profileLink = authenticatedPage.locator('a[href="/profile"]').first();
        await expect(profileLink).toBeVisible();
        await profileLink.click();

        // Verify navigation to profile page
        await authenticatedPage.waitForURL('/profile');
        await profilePage.waitForPageLoad();

        // Verify page loaded correctly
        await profilePage.verifyPageDisplayed();

        console.log('✅ Test passed: Profile settings accessible from sidebar');
    });

    test('Verify profile settings page after authentication', async ({ authenticatedPage }) => {
        console.log('🧪 Test: Verify profile settings page after authentication');

        // Verify the page is accessible
        await profilePage.verifyPageDisplayed();

        // Verify no error messages are shown
        const hasError = await profilePage.verifyErrorDisplayed();
        expect(hasError).toBe(false);

        // Verify currency preference is loaded
        const currentCurrency = await profilePage.getCurrentCurrency();
        expect(currentCurrency).toBeTruthy();
        console.log(`📋 Current currency loaded: ${currentCurrency}`);

        console.log('✅ Test passed: Profile settings page verified after authentication');
    });
});
