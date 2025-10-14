/**
 * Profile Preferences Happy Path Tests
 * 
 * Comprehensive test suite covering successful preferences management workflows:
 * - Display preferences page correctly
 * - View current preferences
 * - Change preferred currency
 * - Change theme preferences
 * - Verify preference persistence
 * 
 * Requirements: 5.1, 5.2
 * 
 * Testing Philosophy: Happy Path Only
 * ✅ Tests successful preference updates with valid data
 * ✅ Tests successful UI interactions
 * ✅ Tests persistence across page refreshes
 * ❌ No error testing, validation failures, or edge cases
 */

import { test, expect } from '../../fixtures/auth';
import { ProfilePage } from '../../page-objects/profile-page';

test.describe('Profile Preferences Happy Path', () => {
    let profilePage: ProfilePage;

    test.beforeEach(async ({ authenticatedPage }) => {
        console.log('🚀 Setting up profile preferences test');

        profilePage = new ProfilePage(authenticatedPage);

        // Navigate to profile/preferences page
        await authenticatedPage.goto('/profile');
        await profilePage.waitForPageLoad();
        console.log('✅ Navigated to preferences page');
    });

    test('Display preferences page layout correctly', async () => {
        console.log('🧪 Test: Display preferences page layout correctly');

        // Verify page is displayed with all essential elements
        await profilePage.verifyPageDisplayed();

        console.log('✅ Test passed: Preferences page layout displayed correctly');
    });

    test('View current preferences successfully', async () => {
        console.log('🧪 Test: View current preferences successfully');

        // Get the current currency preference
        const currentCurrency = await profilePage.getCurrentCurrency();
        console.log(`📋 Current currency: ${currentCurrency}`);

        // Verify we got a valid currency code
        expect(currentCurrency).toBeTruthy();
        expect(currentCurrency.length).toBeGreaterThan(0);

        console.log('✅ Test passed: Current preferences retrieved successfully');
    });

    test('Change preferred currency successfully', async () => {
        console.log('🧪 Test: Change preferred currency successfully');

        // Get the initial currency
        const initialCurrency = await profilePage.getCurrentCurrency();
        console.log(`📋 Initial currency: ${initialCurrency}`);

        // Get available currencies
        const availableCurrencies = await profilePage.getAvailableCurrencies();
        console.log(`📋 Available currencies: ${availableCurrencies.join(', ')}`);

        // Find a different currency to switch to
        const newCurrency = availableCurrencies.find(c => c !== initialCurrency);
        if (!newCurrency) {
            console.log('⚠️ Only one currency available, skipping test');
            test.skip();
            return;
        }

        console.log(`📝 Changing currency from ${initialCurrency} to ${newCurrency}`);

        // Change the currency through UI
        await profilePage.changeCurrency(newCurrency);

        // Verify the change was applied
        await profilePage.verifyCurrencyChanged(newCurrency);

        console.log('✅ Test passed: Currency changed successfully');
    });

    test('Verify currency preference persists after page refresh', async ({ authenticatedPage }) => {
        console.log('🧪 Test: Verify currency preference persists after page refresh');

        // Get the initial currency
        const initialCurrency = await profilePage.getCurrentCurrency();
        console.log(`📋 Initial currency: ${initialCurrency}`);

        // Get available currencies
        const availableCurrencies = await profilePage.getAvailableCurrencies();

        // Find a different currency to switch to
        const newCurrency = availableCurrencies.find(c => c !== initialCurrency);
        if (!newCurrency) {
            console.log('⚠️ Only one currency available, skipping test');
            test.skip();
            return;
        }

        // Change the currency
        await profilePage.changeCurrency(newCurrency);
        await profilePage.verifyCurrencyChanged(newCurrency);

        // Refresh the page and verify persistence
        await profilePage.refreshAndVerifyCurrency(newCurrency);

        console.log('✅ Test passed: Currency preference persisted after refresh');
    });

    test('Change theme to dark mode successfully', async ({ authenticatedPage }) => {
        console.log('🧪 Test: Change theme to dark mode successfully');

        // Select dark theme
        await profilePage.selectTheme('Dark');

        // Verify dark theme is applied
        await profilePage.verifyThemeApplied('Dark');

        console.log('✅ Test passed: Dark theme applied successfully');
    });

    test('Change theme to light mode successfully', async ({ authenticatedPage }) => {
        console.log('🧪 Test: Change theme to light mode successfully');

        // First switch to dark to ensure we're changing from a known state
        await profilePage.selectTheme('Dark');
        await profilePage.verifyThemeApplied('Dark');

        // Now switch to light
        await profilePage.selectTheme('Light');

        // Verify light theme is applied
        await profilePage.verifyThemeApplied('Light');

        console.log('✅ Test passed: Light theme applied successfully');
    });

    test('Verify theme preference persists after page refresh', async ({ authenticatedPage }) => {
        console.log('🧪 Test: Verify theme preference persists after page refresh');

        // Select dark theme
        await profilePage.selectTheme('Dark');
        await profilePage.verifyThemeApplied('Dark');

        // Refresh the page
        console.log('🔄 Refreshing page to verify theme persistence');
        await authenticatedPage.reload();
        await profilePage.waitForPageLoad();

        // Verify dark theme persists
        await profilePage.verifyThemeApplied('Dark');

        console.log('✅ Test passed: Theme preference persisted after refresh');
    });

    test('Change multiple preferences in sequence', async () => {
        console.log('🧪 Test: Change multiple preferences in sequence');

        // Get initial state
        const initialCurrency = await profilePage.getCurrentCurrency();
        console.log(`📋 Initial currency: ${initialCurrency}`);

        // Get available currencies
        const availableCurrencies = await profilePage.getAvailableCurrencies();
        const newCurrency = availableCurrencies.find(c => c !== initialCurrency);

        if (newCurrency) {
            // Change currency
            console.log(`📝 Step 1: Changing currency to ${newCurrency}`);
            await profilePage.changeCurrency(newCurrency);
            await profilePage.verifyCurrencyChanged(newCurrency);
        }

        // Change theme to dark
        console.log(`📝 Step 2: Changing theme to Dark`);
        await profilePage.selectTheme('Dark');
        await profilePage.verifyThemeApplied('Dark');

        // Change theme back to light
        console.log(`📝 Step 3: Changing theme to Light`);
        await profilePage.selectTheme('Light');
        await profilePage.verifyThemeApplied('Light');

        console.log('✅ Test passed: Multiple preferences changed successfully');
    });
});
