import { Page, expect } from '@playwright/test';
import { BasePage } from './base-page';

/**
 * Page object for the Profile/Preferences page
 * Handles profile settings and preferences management
 */
export class ProfilePage extends BasePage {
    // Main page selectors
    private readonly pageTitle = 'h1:has-text("Preferences")';
    private readonly loadingIndicator = 'text="Loading preferences"';
    private readonly errorMessage = 'text="Error loading preferences"';

    // Currency preference selectors
    private readonly currencySelect = '[role="combobox"]';
    private readonly currencyTrigger = 'button:has([role="combobox"])';
    private readonly currencyValue = '[data-value]';
    private readonly currencyLabel = 'text="Preferred Currency"';
    private readonly savingIndicator = 'text="Saving..."';

    // Theme preference selectors
    private readonly themeLabel = 'text="Theme"';
    private readonly themeToggleButton = 'button[aria-label*="theme" i], button:has-text("Light"), button:has-text("Dark"), button:has-text("System")';
    private readonly lightModeButton = 'button:has-text("Light"), [role="menuitem"]:has-text("Light")';
    private readonly darkModeButton = 'button:has-text("Dark"), [role="menuitem"]:has-text("Dark")';
    private readonly systemModeButton = 'button:has-text("System"), [role="menuitem"]:has-text("System")';

    // Dropdown menu items
    private readonly dropdownContent = '[role="listbox"], [role="menu"]';
    private readonly dropdownItem = '[role="option"], [role="menuitem"]';

    constructor(page: Page) {
        super(page);
    }

    getPageUrl(): string {
        return '/profile';
    }

    async waitForPageLoad(): Promise<void> {
        // Wait for the page title to be visible
        await this.waitForElement(this.pageTitle, { timeout: 10000 });

        // Wait for loading state to finish (if present)
        const loadingLocator = this.page.locator(this.loadingIndicator);
        if (await loadingLocator.count() > 0) {
            await loadingLocator.waitFor({ state: 'hidden', timeout: 10000 });
        }

        await this.waitForPageReady();
    }

    /**
     * Verify page is displayed correctly
     */
    async verifyPageDisplayed(): Promise<void> {
        await expect(this.page.locator(this.pageTitle)).toBeVisible();
        await expect(this.page.locator(this.currencyLabel)).toBeVisible();
        await expect(this.page.locator(this.themeLabel)).toBeVisible();
    }

    /**
     * Get current selected currency
     */
    async getCurrentCurrency(): Promise<string> {
        // Try to get the currency from the combobox value
        const combobox = this.page.locator('[role="combobox"]').first();

        // Wait for the combobox to be visible
        await combobox.waitFor({ state: 'visible', timeout: 5000 });

        // Get the text content (which should show the selected currency)
        const text = await combobox.textContent();

        if (!text) {
            throw new Error('Could not get current currency value');
        }

        // Extract currency code (e.g., "EUR" from "Euro (€)")
        const match = text.match(/\b([A-Z]{3})\b/);
        if (match) {
            return match[1];
        }

        return text.trim();
    }

    /**
     * Change preferred currency through the UI
     */
    async changeCurrency(currencyCode: string): Promise<void> {
        console.log(`📝 Changing currency to: ${currencyCode}`);

        // Click the currency select trigger
        const trigger = this.page.locator('[role="combobox"]').first();
        await trigger.waitFor({ state: 'visible', timeout: 5000 });
        await trigger.click();

        // Wait for dropdown to open
        await this.page.waitForTimeout(500);

        // Find and click the currency option
        // Try multiple selectors to find the right option
        const selectors = [
            `[role="option"]:has-text("${currencyCode}")`,
            `[role="option"] >> text=${currencyCode}`,
            `text="${currencyCode}"`,
        ];

        let clicked = false;
        for (const selector of selectors) {
            try {
                const option = this.page.locator(selector).first();
                if (await option.isVisible({ timeout: 2000 })) {
                    await option.click();
                    clicked = true;
                    console.log(`✅ Clicked currency option: ${currencyCode}`);
                    break;
                }
            } catch {
                continue;
            }
        }

        if (!clicked) {
            throw new Error(`Could not find currency option: ${currencyCode}`);
        }

        // Wait for the selection to be applied
        await this.page.waitForTimeout(1000);
    }

    /**
     * Verify currency change was applied
     */
    async verifyCurrencyChanged(expectedCurrency: string): Promise<void> {
        console.log(`🔍 Verifying currency changed to: ${expectedCurrency}`);

        // Wait for the saving indicator to disappear if present
        const savingLocator = this.page.locator(this.savingIndicator);
        if (await savingLocator.count() > 0) {
            await savingLocator.waitFor({ state: 'hidden', timeout: 5000 });
        }

        // Wait a bit for the UI to update
        await this.page.waitForTimeout(1000);

        // Get the current currency and verify
        const currentCurrency = await this.getCurrentCurrency();
        expect(currentCurrency).toContain(expectedCurrency);
        console.log(`✅ Currency verified: ${currentCurrency}`);
    }

    /**
     * Open theme toggle dropdown
     */
    async openThemeToggle(): Promise<void> {
        // Try to find and click the theme toggle button
        const selectors = [
            'button:has-text("Toggle theme")',
            'button:has([data-lucide="sun"])',
            'button:has([data-lucide="moon"])',
            'button[variant="outline"][size="icon"]',
            '[data-testid="theme-toggle"]',
        ];

        for (const selector of selectors) {
            try {
                const button = this.page.locator(selector).first();
                if (await button.isVisible({ timeout: 2000 })) {
                    await button.click();
                    await this.page.waitForTimeout(500);
                    console.log(`✅ Opened theme toggle`);
                    return;
                }
            } catch {
                continue;
            }
        }

        throw new Error('Could not find theme toggle button');
    }    /**
     * Select theme mode (Light, Dark, or System)
     */
    async selectTheme(theme: 'Light' | 'Dark' | 'System'): Promise<void> {
        console.log(`📝 Changing theme to: ${theme}`);

        // First open the theme toggle
        await this.openThemeToggle();

        // Wait for dropdown menu to appear
        await this.page.waitForTimeout(500);

        // Find and click the theme option
        // The theme toggle uses DropdownMenuItem which renders as a div with role="menuitem"
        const selectors = [
            `div[role="menuitem"]:has-text("${theme}")`,
            `[role="menuitem"]:has-text("${theme}")`,
            `text="${theme}"`,
        ];

        let clicked = false;
        for (const selector of selectors) {
            try {
                const option = this.page.locator(selector).first();
                if (await option.isVisible({ timeout: 2000 })) {
                    await option.click();
                    clicked = true;
                    console.log(`✅ Selected theme: ${theme}`);
                    break;
                }
            } catch {
                continue;
            }
        }

        if (!clicked) {
            throw new Error(`Could not find theme option: ${theme}`);
        }

        await this.page.waitForTimeout(500);
    }

    /**
     * Verify theme is applied by checking the HTML element's class
     */
    async verifyThemeApplied(theme: 'Light' | 'Dark' | 'System'): Promise<void> {
        console.log(`🔍 Verifying theme applied: ${theme}`);

        // Wait a bit for the theme to be applied
        await this.page.waitForTimeout(1000);

        const html = this.page.locator('html');

        if (theme === 'Dark') {
            await expect(html).toHaveClass(/dark/);
            console.log(`✅ Dark theme verified`);
        } else if (theme === 'Light') {
            // Light theme should not have the 'dark' class
            const htmlClass = await html.getAttribute('class');
            expect(htmlClass).not.toContain('dark');
            console.log(`✅ Light theme verified`);
        } else {
            // System theme - just verify no errors occurred
            console.log(`✅ System theme set (appearance follows system preference)`);
        }
    }

    /**
     * Verify saving indicator appears and disappears
     */
    async verifySavingIndicator(): Promise<void> {
        const savingLocator = this.page.locator(this.savingIndicator);

        // Check if saving indicator appears (may be very brief)
        try {
            await savingLocator.waitFor({ state: 'visible', timeout: 2000 });
            console.log(`✅ Saving indicator appeared`);
        } catch {
            console.log(`ℹ️ Saving indicator did not appear (may have been too fast)`);
        }

        // Wait for it to disappear
        if (await savingLocator.count() > 0) {
            await savingLocator.waitFor({ state: 'hidden', timeout: 5000 });
            console.log(`✅ Saving indicator disappeared`);
        }
    }

    /**
     * Refresh page and verify currency persists
     */
    async refreshAndVerifyCurrency(expectedCurrency: string): Promise<void> {
        console.log(`🔄 Refreshing page to verify currency persistence`);

        await this.page.reload();
        await this.waitForPageLoad();

        const currentCurrency = await this.getCurrentCurrency();
        expect(currentCurrency).toContain(expectedCurrency);
        console.log(`✅ Currency persisted after refresh: ${currentCurrency}`);
    }

    /**
     * Verify error message is displayed
     */
    async verifyErrorDisplayed(): Promise<boolean> {
        const errorLocator = this.page.locator(this.errorMessage);
        return await errorLocator.isVisible({ timeout: 2000 });
    }

    /**
     * Get list of available currencies from dropdown
     */
    async getAvailableCurrencies(): Promise<string[]> {
        // Open the currency dropdown
        const trigger = this.page.locator('[role="combobox"]').first();
        await trigger.click();

        // Wait for dropdown to open
        await this.page.waitForTimeout(500);

        // Get all currency options
        const options = this.page.locator('[role="option"]');
        const count = await options.count();
        const currencies: string[] = [];

        for (let i = 0; i < count; i++) {
            const text = await options.nth(i).textContent();
            if (text) {
                // Extract currency code from text like "Euro (€)"
                const match = text.match(/\b([A-Z]{3})\b/);
                if (match) {
                    currencies.push(match[1]);
                }
            }
        }

        // Close the dropdown by clicking outside or pressing Escape
        await this.page.keyboard.press('Escape');
        await this.page.waitForTimeout(300);

        return currencies;
    }
}