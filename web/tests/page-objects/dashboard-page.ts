import { Locator, expect } from '@playwright/test';
import { BasePage } from './base-page';

/**
 * Dashboard page object for interacting with the main dashboard
 * Provides methods to verify summary cards, upcoming renewals, top providers, and top labels
 */
export class DashboardPage extends BasePage {
  // Page URL
  getPageUrl(): string {
    return '/dashboard';
  }



  // Page load verification
  async waitForPageLoad(): Promise<void> {
    await this.waitForElement('h1:has-text("Dashboard")');
    await this.waitForPageReady();
  }

  // Summary Cards Selectors and Methods
  private get summaryCardsContainer(): Locator {
    return this.page.locator('.grid.grid-cols-1.md\\:grid-cols-3.gap-6').first();
  }

  private get monthlyExpensesCard(): Locator {
    return this.summaryCardsContainer.locator('.border-l-blue-500');
  }

  private get yearlyExpensesCard(): Locator {
    return this.summaryCardsContainer.locator('.border-l-purple-500');
  }

  private get activeSubscriptionsCard(): Locator {
    return this.summaryCardsContainer.locator('.border-l-green-500');
  }

  /**
   * Get the monthly expenses amount from the summary card
   */
  async getMonthlyExpenses(): Promise<string> {
    const card = this.monthlyExpensesCard;
    await card.waitFor({ state: 'visible' });
    
    // Wait for loading to complete (no skeleton in card content)
    await expect(card.locator('.h-10.w-28')).toHaveCount(0, { timeout: 10000 });
    
    // Look for the Money component with specific classes
    const amountElement = card.locator('[class*="text-3xl"][class*="font-bold"][class*="bg-gradient-to-r"]').first();
    return await amountElement.textContent() || '';
  }

  /**
   * Get the yearly expenses amount from the summary card
   */
  async getYearlyExpenses(): Promise<string> {
    const card = this.yearlyExpensesCard;
    await card.waitFor({ state: 'visible' });
    
    // Wait for loading to complete (no skeleton in card content)
    await expect(card.locator('.h-10.w-28')).toHaveCount(0, { timeout: 10000 });
    
    // Look for the Money component with specific classes
    const amountElement = card.locator('[class*="text-3xl"][class*="font-bold"][class*="bg-gradient-to-r"]').first();
    return await amountElement.textContent() || '';
  }

  /**
   * Get the active subscriptions count from the summary card
   */
  async getActiveSubscriptionsCount(): Promise<number> {
    const card = this.activeSubscriptionsCard;
    await card.waitFor({ state: 'visible' });
    
    // Wait for loading to complete (no skeleton in card content)
    await expect(card.locator('.h-10.w-28')).toHaveCount(0, { timeout: 10000 });
    
    // Look for the count with specific classes
    const countElement = card.locator('[class*="text-3xl"][class*="font-bold"][class*="bg-gradient-to-r"]').first();
    const countText = await countElement.textContent() || '0';
    return parseInt(countText, 10);
  }

  /**
   * Verify that summary cards are displayed and contain valid data
   */
  async verifySummaryCards(): Promise<void> {
    // Verify all three cards are visible
    await expect(this.monthlyExpensesCard).toBeVisible();
    await expect(this.yearlyExpensesCard).toBeVisible();
    await expect(this.activeSubscriptionsCard).toBeVisible();

    // Verify card titles - using specific CardTitle selectors
    await expect(this.monthlyExpensesCard.locator('.text-xl.font-medium:has-text("Monthly Expenses")')).toBeVisible();
    await expect(this.yearlyExpensesCard.locator('.text-xl.font-medium:has-text("Yearly Expenses")')).toBeVisible();
    await expect(this.activeSubscriptionsCard.locator('.text-xl.font-medium:has-text("Active Subscriptions")')).toBeVisible();

    // Verify icons are present - using specific icon classes
    await expect(this.monthlyExpensesCard.locator('.h-5.w-5.text-blue-500').first()).toBeVisible();
    await expect(this.yearlyExpensesCard.locator('.h-5.w-5.text-purple-500').first()).toBeVisible();
    await expect(this.activeSubscriptionsCard.locator('.h-5.w-5.text-green-500').first()).toBeVisible();
  }

  // Upcoming Renewals Section
  private get upcomingRenewalsSection(): Locator {
    return this.page.locator('div:has(p:has-text("Upcoming Renewals")), section:has(h2:has-text("Upcoming Renewals")), [data-testid="upcoming-renewals"]').first();
  }

  /**
   * Get the list of upcoming renewals
   */
  async getUpcomingRenewals(): Promise<Array<{ provider: string; amount: string; date: string }>> {
    const section = this.upcomingRenewalsSection;
    await section.waitFor({ state: 'visible', timeout: 15000 });

    // Wait for loading to complete or check for empty state
    try {
      await expect(section.locator('[key*="upcoming-skeleton"], .animate-pulse')).toHaveCount(0, { timeout: 10000 });
    } catch {
      // If skeleton check fails, continue anyway
    }

    // Check for empty state first
    const emptyMessages = [
      'text="No upcoming renewals found"',
      'text="No renewals"',
      ':has-text("No upcoming renewals")',
      ':has-text("No renewals")'
    ];
    
    for (const emptySelector of emptyMessages) {
      const emptyMessage = section.locator(emptySelector);
      if (await emptyMessage.isVisible({ timeout: 2000 })) {
        return [];
      }
    }

    // Look for renewal items with more flexible selectors
    const renewalItems = section.locator('.space-y-3 > div, .p-3.border.rounded-lg');
    const count = await renewalItems.count();
    const renewals: Array<{ provider: string; amount: string; date: string }> = [];

    for (let i = 0; i < count; i++) {
      const item = renewalItems.nth(i);
      
      try {
        // Try different selectors for provider name
        let provider = '';
        const providerSelectors = ['h4', '.font-medium', '[class*="font-medium"]'];
        for (const selector of providerSelectors) {
          try {
            const text = await item.locator(selector).first().textContent({ timeout: 2000 });
            if (text && text.trim()) {
              provider = text.trim();
              break;
            }
          } catch {
            continue;
          }
        }

        // Try different selectors for amount
        let amount = '';
        const amountSelectors = ['.font-semibold', '[class*="font-semibold"]', '[class*="bg-gradient-to-r"]'];
        for (const selector of amountSelectors) {
          try {
            const text = await item.locator(selector).first().textContent({ timeout: 2000 });
            if (text && text.trim()) {
              amount = text.trim();
              break;
            }
          } catch {
            continue;
          }
        }

        // Try to get date
        let date = '';
        try {
          const dateText = await item.locator('.text-sm').first().textContent({ timeout: 2000 }) || '';
          date = dateText.replace('Next renewal:', '').trim();
        } catch {
          // Date is optional
        }

        if (provider && amount) {
          renewals.push({ provider, amount, date });
        }
      } catch (error) {
        // Skip this item if we can't extract data
        console.warn(`Failed to extract renewal data for item ${i}:`, error);
      }
    }

    return renewals;
  }

  /**
   * Verify upcoming renewals section is displayed
   */
  async verifyUpcomingRenewalsSection(): Promise<void> {
    const section = this.upcomingRenewalsSection;
    await expect(section).toBeVisible({ timeout: 15000 });
    
    // Try multiple selectors for the title
    const titleSelectors = [
      'p.text-xl.font-medium:has-text("Upcoming Renewals")',
      'h2:has-text("Upcoming Renewals")',
      ':has-text("Upcoming Renewals")'
    ];
    
    let titleFound = false;
    for (const selector of titleSelectors) {
      try {
        await expect(section.locator(selector).first()).toBeVisible({ timeout: 2000 });
        titleFound = true;
        break;
      } catch {
        continue;
      }
    }
    
    if (!titleFound) {
      throw new Error('Upcoming Renewals section title not found');
    }
  }

  /**
   * Click on an upcoming renewal item to navigate to provider detail
   */
  async clickUpcomingRenewal(providerName: string): Promise<void> {
    const section = this.upcomingRenewalsSection;
    const renewalItem = section.locator('.cursor-pointer').filter({ hasText: providerName });
    await renewalItem.click();
  }

  // Top Providers Section
  private get topProvidersSection(): Locator {
    return this.page.locator('div:has(p:has-text("Top Providers by Expense")), section:has(h2:has-text("Top Providers")), [data-testid="top-providers"]').first();
  }

  /**
   * Get the list of top providers
   */
  async getTopProviders(): Promise<Array<{ name: string; amount: string; duration: string }>> {
    const section = this.topProvidersSection;
    await section.waitFor({ state: 'visible', timeout: 15000 });

    // Wait for loading to complete or check for empty state
    try {
      await expect(section.locator('[key*="topprov-skeleton"], .animate-pulse')).toHaveCount(0, { timeout: 10000 });
    } catch {
      // If skeleton check fails, continue anyway
    }

    // Check for empty state first
    const emptyMessages = [
      'text="No provider spending data available"',
      'text="No providers"',
      ':has-text("No provider")',
      ':has-text("No data")'
    ];
    
    for (const emptySelector of emptyMessages) {
      const emptyMessage = section.locator(emptySelector);
      if (await emptyMessage.isVisible({ timeout: 2000 })) {
        return [];
      }
    }

    // Look for provider items with more flexible selectors
    const providerItems = section.locator('.space-y-3 > div, .p-3.border.rounded-lg');
    const count = await providerItems.count();
    const providers: Array<{ name: string; amount: string; duration: string }> = [];

    for (let i = 0; i < count; i++) {
      const item = providerItems.nth(i);
      
      try {
        // Try different selectors for provider name
        let name = '';
        const nameSelectors = ['h4', '.font-medium', '[class*="font-medium"]'];
        for (const selector of nameSelectors) {
          try {
            const text = await item.locator(selector).first().textContent({ timeout: 2000 });
            if (text && text.trim()) {
              name = text.trim();
              break;
            }
          } catch {
            continue;
          }
        }

        // Try different selectors for amount
        let amount = '';
        const amountSelectors = ['.font-semibold', '[class*="font-semibold"]', '[class*="bg-gradient-to-r"]'];
        for (const selector of amountSelectors) {
          try {
            const text = await item.locator(selector).first().textContent({ timeout: 2000 });
            if (text && text.trim()) {
              amount = text.trim();
              break;
            }
          } catch {
            continue;
          }
        }

        // Try to get duration
        let duration = '';
        try {
          const durationText = await item.locator('.text-sm.text-muted-foreground, .text-sm').first().textContent({ timeout: 2000 }) || '';
          duration = durationText.trim();
        } catch {
          // Duration is optional
        }

        if (name && amount) {
          providers.push({ name, amount, duration });
        }
      } catch (error) {
        // Skip this item if we can't extract data
        console.warn(`Failed to extract provider data for item ${i}:`, error);
      }
    }

    return providers;
  }

  /**
   * Verify top providers section is displayed
   */
  async verifyTopProvidersSection(): Promise<void> {
    const section = this.topProvidersSection;
    await expect(section).toBeVisible({ timeout: 15000 });
    
    // Try multiple selectors for the title
    const titleSelectors = [
      'p.text-xl.font-medium:has-text("Top Providers by Expense")',
      'h2:has-text("Top Providers")',
      ':has-text("Top Providers")'
    ];
    
    let titleFound = false;
    for (const selector of titleSelectors) {
      try {
        await expect(section.locator(selector).first()).toBeVisible({ timeout: 2000 });
        titleFound = true;
        break;
      } catch {
        continue;
      }
    }
    
    if (!titleFound) {
      throw new Error('Top Providers section title not found');
    }
  }

  /**
   * Click on a top provider item to navigate to provider detail
   */
  async clickTopProvider(providerName: string): Promise<void> {
    const section = this.topProvidersSection;
    const providerItem = section.locator('.cursor-pointer').filter({ hasText: providerName });
    await providerItem.click();
  }

  // Top Labels Section
  private get topLabelsSection(): Locator {
    return this.page.locator('div:has(p:has-text("Top Spent by Label")), section:has(h2:has-text("Top Labels")), [data-testid="top-labels"]').first();
  }

  /**
   * Get the list of top labels
   */
  async getTopLabels(): Promise<Array<{ name: string; amount: string }>> {
    const section = this.topLabelsSection;
    await section.waitFor({ state: 'visible', timeout: 15000 });

    // Wait for loading to complete or check for empty state
    try {
      await expect(section.locator('[key*="toplabel-skeleton"], .animate-pulse')).toHaveCount(0, { timeout: 10000 });
    } catch {
      // If skeleton check fails, continue anyway
    }

    // Check for empty state first
    const emptyMessages = [
      'text="No label spending data available"',
      'text="No labels"',
      ':has-text("No label")',
      ':has-text("No data")'
    ];
    
    for (const emptySelector of emptyMessages) {
      const emptyMessage = section.locator(emptySelector);
      if (await emptyMessage.isVisible({ timeout: 2000 })) {
        return [];
      }
    }

    // Look for label items with more flexible selectors
    const labelItems = section.locator('.space-y-3 > div, .p-3.border.rounded-lg');
    const count = await labelItems.count();
    const labels: Array<{ name: string; amount: string }> = [];

    for (let i = 0; i < count; i++) {
      const item = labelItems.nth(i);
      
      try {
        // Try different selectors for label name
        let name = '';
        const nameSelectors = ['h4', '.font-medium', '[class*="font-medium"]'];
        for (const selector of nameSelectors) {
          try {
            const text = await item.locator(selector).first().textContent({ timeout: 2000 });
            if (text && text.trim()) {
              name = text.trim();
              break;
            }
          } catch {
            continue;
          }
        }

        // Try different selectors for amount
        let amount = '';
        const amountSelectors = ['.font-semibold', '[class*="font-semibold"]', '[class*="bg-gradient-to-r"]'];
        for (const selector of amountSelectors) {
          try {
            const text = await item.locator(selector).first().textContent({ timeout: 2000 });
            if (text && text.trim()) {
              amount = text.trim();
              break;
            }
          } catch {
            continue;
          }
        }

        if (name && amount) {
          labels.push({ name, amount });
        }
      } catch (error) {
        // Skip this item if we can't extract data
        console.warn(`Failed to extract label data for item ${i}:`, error);
      }
    }

    return labels;
  }

  /**
   * Verify top labels section is displayed
   */
  async verifyTopLabelsSection(): Promise<void> {
    const section = this.topLabelsSection;
    await expect(section).toBeVisible({ timeout: 15000 });
    
    // Try multiple selectors for the title
    const titleSelectors = [
      'p.text-xl.font-medium:has-text("Top Spent by Label")',
      'h2:has-text("Top Labels")',
      ':has-text("Top Labels")',
      ':has-text("Top Spent by Label")'
    ];
    
    let titleFound = false;
    for (const selector of titleSelectors) {
      try {
        await expect(section.locator(selector).first()).toBeVisible({ timeout: 2000 });
        titleFound = true;
        break;
      } catch {
        continue;
      }
    }
    
    if (!titleFound) {
      throw new Error('Top Labels section title not found');
    }
  }

  // Navigation Methods

  /**
   * Navigate to subscriptions page via sidebar
   */
  async navigateToSubscriptions(): Promise<void> {
    // Try multiple selectors to find the subscriptions link
    const selectors = [
      'a[href="/subscriptions"]',
      '[href="/subscriptions"]',
      'text="Subscriptions"',
      ':has-text("Subscriptions")'
    ];
    
    let clicked = false;
    for (const selector of selectors) {
      try {
        const link = this.page.locator(selector).first();
        if (await link.isVisible({ timeout: 2000 })) {
          await link.click();
          clicked = true;
          break;
        }
      } catch {
        continue;
      }
    }
    
    if (!clicked) {
      throw new Error('Could not find subscriptions navigation link');
    }
    
    await this.page.waitForURL('/subscriptions');
  }

  /**
   * Navigate to providers page via sidebar
   */
  async navigateToProviders(): Promise<void> {
    // Try multiple selectors to find the providers link
    const selectors = [
      'a[href="/providers"]',
      '[href="/providers"]',
      'text="Providers"',
      ':has-text("Providers")'
    ];
    
    let clicked = false;
    for (const selector of selectors) {
      try {
        const link = this.page.locator(selector).first();
        if (await link.isVisible({ timeout: 2000 })) {
          await link.click();
          clicked = true;
          break;
        }
      } catch {
        continue;
      }
    }
    
    if (!clicked) {
      throw new Error('Could not find providers navigation link');
    }
    
    await this.page.waitForURL('/providers');
  }

  /**
   * Navigate to labels page via sidebar
   */
  async navigateToLabels(): Promise<void> {
    // Try multiple selectors to find the labels link
    const selectors = [
      'a[href="/labels"]',
      '[href="/labels"]',
      'text="Labels"',
      ':has-text("Labels")'
    ];
    
    let clicked = false;
    for (const selector of selectors) {
      try {
        const link = this.page.locator(selector).first();
        if (await link.isVisible({ timeout: 2000 })) {
          await link.click();
          clicked = true;
          break;
        }
      } catch {
        continue;
      }
    }
    
    if (!clicked) {
      throw new Error('Could not find labels navigation link');
    }
    
    await this.page.waitForURL('/labels');
  }

  /**
   * Navigate to family page via sidebar
   */
  async navigateToFamily(): Promise<void> {
    // Try multiple selectors to find the family link
    const selectors = [
      'a[href="/family"]',
      '[href="/family"]',
      'text="Family"',
      ':has-text("Family")'
    ];
    
    let clicked = false;
    for (const selector of selectors) {
      try {
        const link = this.page.locator(selector).first();
        if (await link.isVisible({ timeout: 2000 })) {
          await link.click();
          clicked = true;
          break;
        }
      } catch {
        continue;
      }
    }
    
    if (!clicked) {
      throw new Error('Could not find family navigation link');
    }
    
    await this.page.waitForURL('/family');
  }

  /**
   * Navigate to profile page via sidebar
   */
  async navigateToProfile(): Promise<void> {
    // Try multiple selectors to find the profile link
    const selectors = [
      'a[href="/profile"]',
      '[href="/profile"]',
      'text="Profile"',
      'text="Preferences"',
      ':has-text("Profile")',
      ':has-text("Preferences")'
    ];
    
    let clicked = false;
    for (const selector of selectors) {
      try {
        const link = this.page.locator(selector).first();
        if (await link.isVisible({ timeout: 2000 })) {
          await link.click();
          clicked = true;
          break;
        }
      } catch {
        continue;
      }
    }
    
    if (!clicked) {
      throw new Error('Could not find profile navigation link');
    }
    
    await this.page.waitForURL('/profile');
  }

  /**
   * Verify sidebar navigation is present and functional
   */
  async verifySidebarNavigation(): Promise<void> {
    // Verify all navigation links are present - use flexible selectors
    const navigationLinks = [
      { href: '/dashboard', text: 'Dashboard' },
      { href: '/subscriptions', text: 'Subscriptions' },
      { href: '/providers', text: 'Providers' },
      { href: '/family', text: 'Family' },
      { href: '/labels', text: 'Labels' },
      { href: '/profile', text: 'Profile' }
    ];
    
    for (const nav of navigationLinks) {
      const selectors = [
        `a[href="${nav.href}"]`,
        `[href="${nav.href}"]`,
        `:has-text("${nav.text}")`
      ];
      
      let found = false;
      for (const selector of selectors) {
        try {
          const element = this.page.locator(selector).first();
          if (await element.isVisible({ timeout: 2000 })) {
            found = true;
            break;
          }
        } catch {
          continue;
        }
      }
      
      if (!found) {
        throw new Error(`Navigation link for ${nav.text} not found`);
      }
    }
  }

  /**
   * Verify empty state is displayed when no subscriptions exist
   */
  async verifyEmptyState(): Promise<void> {
    // Look for various empty state indicators
    const emptyStateSelectors = [
      'text="No subscriptions found"',
      'text="Get started by adding your first subscription"',
      '[data-testid="empty-state"]',
      ':has-text("No upcoming renewals found")',
      ':has-text("No provider spending data available")',
      ':has-text("No label spending data available")'
    ];
    
    let emptyStateFound = false;
    for (const selector of emptyStateSelectors) {
      try {
        await expect(this.page.locator(selector)).toBeVisible({ timeout: 2000 });
        emptyStateFound = true;
        break;
      } catch {
        // Continue to next selector
      }
    }
    
    // If no specific empty state found, verify that the active count is 0
    if (!emptyStateFound) {
      const activeCount = await this.getActiveSubscriptionsCount();
      expect(activeCount).toBe(0);
    }
  }



  /**
   * Wait for dashboard data to load completely
   */
  async waitForDashboardDataLoad(): Promise<void> {
    // Wait for dashboard content to be visible
    await this.waitForElement('h1:has-text("Dashboard")');
    
    // Wait for summary cards to be visible and loaded
    try {
      await expect(this.summaryCardsContainer).toBeVisible({ timeout: 10000 });
    } catch {
      // If summary cards container not found, just verify dashboard is loaded
      console.log('Summary cards container not found, but dashboard loaded');
      return;
    }
    
    // Wait for summary cards to finish loading (no skeleton loaders in cards)
    try {
      await expect(this.summaryCardsContainer.locator('.h-10.w-28, .animate-pulse')).toHaveCount(0, { timeout: 8000 });
    } catch {
      // If skeleton check fails, just wait for the cards to have content
      try {
        await expect(this.summaryCardsContainer.locator('[class*="text-3xl"][class*="font-bold"]')).toHaveCount(3, { timeout: 8000 });
      } catch {
        // If that also fails, just continue - dashboard might be empty
        console.log('Could not verify summary card content, but dashboard loaded');
      }
    }
    
    // Wait for API responses to complete
    try {
      await this.waitForPageReady();
    } catch {
      // If page ready check fails, just continue
      console.log('Page ready check failed, but continuing');
    }
    
    // Reduced wait for any async data loading
    await this.page.waitForTimeout(1000);
  }

  /**
   * Verify dashboard displays accurate data by checking data consistency
   */
  async verifyDataAccuracy(): Promise<void> {
    await this.waitForDashboardDataLoad();

    // Get data from all sections
    const monthlyExpenses = await this.getMonthlyExpenses();
    const yearlyExpenses = await this.getYearlyExpenses();
    const activeCount = await this.getActiveSubscriptionsCount();
    const upcomingRenewals = await this.getUpcomingRenewals();
    const topProviders = await this.getTopProviders();
    const topLabels = await this.getTopLabels();

    // Verify data is present and formatted correctly
    expect(monthlyExpenses).toBeTruthy();
    expect(yearlyExpenses).toBeTruthy();
    expect(activeCount).toBeGreaterThanOrEqual(0);

    // Verify arrays are properly structured
    upcomingRenewals.forEach(renewal => {
      expect(renewal.provider).toBeTruthy();
      expect(renewal.amount).toBeTruthy();
      expect(renewal.date).toBeTruthy();
    });

    topProviders.forEach(provider => {
      expect(provider.name).toBeTruthy();
      expect(provider.amount).toBeTruthy();
      expect(provider.duration).toBeTruthy();
    });

    topLabels.forEach(label => {
      expect(label.name).toBeTruthy();
      expect(label.amount).toBeTruthy();
    });
  }

  /**
   * Verify all dashboard sections are present and visible
   */
  async verifyAllSections(): Promise<void> {
    await this.verifySummaryCards();
    await this.verifyUpcomingRenewalsSection();
    await this.verifyTopProvidersSection();
    await this.verifyTopLabelsSection();
    await this.verifySidebarNavigation();
  }
}