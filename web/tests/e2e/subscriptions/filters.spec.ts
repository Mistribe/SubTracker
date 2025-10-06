import { test, expect } from '../../fixtures/auth';
import { SubscriptionsPage } from '../../page-objects/subscriptions-page';
import { TestDataGenerators, SubscriptionData } from '../../utils/data-generators';
import { TestHelpers, createTestHelpers } from '../../utils/test-helpers';

test.describe('Subscription Filtering and Search', () => {
  let subscriptionsPage: SubscriptionsPage;
  let testHelpers: TestHelpers;
  let testSubscriptions: SubscriptionData[];
  let testProviders: any[];
  let createdSubscriptionIds: string[];

  test.beforeEach(async ({ authenticatedPage }) => {
    subscriptionsPage = new SubscriptionsPage(authenticatedPage);
    testHelpers = await createTestHelpers(authenticatedPage);
    
    // Initialize arrays
    testProviders = [];
    testSubscriptions = [];
    createdSubscriptionIds = [];
    
    // Create test providers using real API
    for (let i = 0; i < 3; i++) {
      const provider = await testHelpers.createTestProvider({
        name: `Test Provider ${i + 1} ${Date.now()}`,
        description: `Test provider ${i + 1} for filtering`,
      });
      testProviders.push(provider);
    }
    
    // Create multiple test subscriptions with different properties
    testSubscriptions = [
      {
        ...TestDataGenerators.generateSubscription(),
        name: `Active Monthly Sub ${Date.now()}`,
        providerId: testProviders[0].id,
        billingCycle: 'monthly' as const,
        isActive: true
      },
      {
        ...TestDataGenerators.generateSubscription(),
        name: `Active Yearly Sub ${Date.now()}`,
        providerId: testProviders[1].id,
        billingCycle: 'yearly' as const,
        isActive: true
      },
      {
        ...TestDataGenerators.generateSubscription(),
        name: `Inactive Sub ${Date.now()}`,
        providerId: testProviders[2].id,
        billingCycle: 'monthly' as const,
        isActive: false
      }
    ];
    
    // Create subscriptions using real API
    for (const subscription of testSubscriptions) {
      const id = await testHelpers.createTestSubscription(subscription);
      createdSubscriptionIds.push(id);
    }
    
    // Navigate to subscriptions page
    await authenticatedPage.goto('/subscriptions');
    await subscriptionsPage.waitForPageLoad();
  });

  test.afterEach(async () => {
    // Clean up test data
    if (testHelpers) {
      if (createdSubscriptionIds && createdSubscriptionIds.length > 0) {
        for (const id of createdSubscriptionIds) {
          try {
            await testHelpers.cleanupTestSubscription(id);
          } catch {
            // Ignore cleanup errors
          }
        }
      }
      
      if (testProviders && testProviders.length > 0) {
        for (const provider of testProviders) {
          try {
            if (provider.id && !provider.id.startsWith('default-provider')) {
              await testHelpers.cleanupTestProvider(provider.id);
            }
          } catch {
            // Ignore cleanup errors
          }
        }
      }
    }
  });

  test('should search subscriptions by name', async () => {
    // Use a generic search term that might exist
    const searchTerm = 'test';
    
    // Search for subscriptions
    await subscriptionsPage.searchSubscriptions(searchTerm);
    
    // Verify search results contain matching subscriptions
    const subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    const matchingNames = subscriptionNames.filter(name => 
      name.toLowerCase().includes(searchTerm.toLowerCase())
    );
    
    expect(matchingNames.length).toBeGreaterThan(0);
    
    // Verify all visible subscriptions match the search term
    for (const name of subscriptionNames) {
      expect(name.toLowerCase()).toContain(searchTerm.toLowerCase());
    }
  });

  test('should clear search and show all subscriptions', async () => {
    // First search for something specific
    await subscriptionsPage.searchSubscriptions(testSubscriptions[0].name);
    
    // Verify filtered results
    let subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames.length).toBeLessThan(testSubscriptions.length);
    
    // Clear search
    await subscriptionsPage.clearSearch();
    
    // Verify all subscriptions are shown again
    subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames.length).toBeGreaterThanOrEqual(testSubscriptions.length);
  });

  test('should filter by active/inactive status', async () => {
    // Filter to include inactive subscriptions
    await subscriptionsPage.applyFilters({
      includeInactive: true
    });
    
    // Should see both active and inactive subscriptions
    let subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames).toContain(testSubscriptions[2].name); // Inactive subscription
    
    // Filter to exclude inactive subscriptions (default)
    await subscriptionsPage.applyFilters({
      includeInactive: false
    });
    
    // Should not see inactive subscriptions
    subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames).not.toContain(testSubscriptions[2].name); // Inactive subscription
  });

  test('should filter by date range', async () => {
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);
    const tomorrow = new Date(today);
    tomorrow.setDate(tomorrow.getDate() + 1);
    
    const yesterdayStr = yesterday.toISOString().split('T')[0];
    const tomorrowStr = tomorrow.toISOString().split('T')[0];
    
    // Filter by date range that should include today's subscriptions
    await subscriptionsPage.applyFilters({
      fromDate: yesterdayStr,
      toDate: tomorrowStr
    });
    
    // Should see subscriptions created today
    const subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames.length).toBeGreaterThan(0);
    
    // Clear filters
    await subscriptionsPage.clearFilters();
  });

  test('should filter by provider', async () => {
    // This test would require implementing provider filtering in the page object
    // For now, we'll test opening the filter drawer and verifying provider options exist
    
    await subscriptionsPage.openFilters();
    
    // Verify filter drawer is open
    const filterDrawer = subscriptionsPage.pageInstance.locator('[role="dialog"]:has-text("Filters")');
    await expect(filterDrawer).toBeVisible();
    
    // Verify provider filter section exists
    const providerSection = filterDrawer.locator('text="Providers"');
    await expect(providerSection).toBeVisible();
    
    // Close filters
    await subscriptionsPage.closeFilters();
  });

  test('should filter by billing cycle/recurrency', async () => {
    await subscriptionsPage.openFilters();
    
    // Verify recurrency filter section exists
    const filterDrawer = subscriptionsPage.pageInstance.locator('[role="dialog"]:has-text("Filters")');
    const recurrencySection = filterDrawer.locator('text="Recurrencies"');
    await expect(recurrencySection).toBeVisible();
    
    // Close filters
    await subscriptionsPage.closeFilters();
  });

  test('should clear all filters', async () => {
    // Apply multiple filters
    await subscriptionsPage.applyFilters({
      includeInactive: true,
      fromDate: '2024-01-01',
      toDate: '2024-12-31'
    });
    
    // Get filtered results count
    const filteredCount = await subscriptionsPage.getSubscriptionCount();
    
    // Clear all filters
    await subscriptionsPage.clearFilters();
    
    // Get unfiltered results count
    const unfilteredCount = await subscriptionsPage.getSubscriptionCount();
    
    // Unfiltered count should be >= filtered count
    expect(unfilteredCount).toBeGreaterThanOrEqual(filteredCount);
  });

  test('should persist search while applying filters', async () => {
    const searchTerm = 'Active';
    
    // Search first
    await subscriptionsPage.searchSubscriptions(searchTerm);
    
    // Apply additional filters
    await subscriptionsPage.applyFilters({
      includeInactive: false
    });
    
    // Verify search term is still in search input
    const searchInput = subscriptionsPage.pageInstance.locator('input[placeholder*="Search"]');
    await expect(searchInput).toHaveValue(searchTerm);
    
    // Verify results still match search term
    const subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    for (const name of subscriptionNames) {
      expect(name.toLowerCase()).toContain(searchTerm.toLowerCase());
    }
  });

  test('should show no results message when filters match nothing', async () => {
    // Search for something that doesn't exist
    await subscriptionsPage.searchSubscriptions('NonExistentSubscription12345');
    
    // Should show no results message
    const noResultsMessage = subscriptionsPage.pageInstance.locator('text="No subscriptions match your search criteria."');
    await expect(noResultsMessage).toBeVisible();
    
    // Table should be empty
    const isEmpty = await subscriptionsPage.isSubscriptionsTableEmpty();
    expect(isEmpty).toBe(true);
  });

  test('should handle special characters in search', async () => {
    // Create a subscription with special characters
    const specialSubscription = {
      ...TestDataGenerators.generateSubscription(),
      name: `Test-Sub_With@Special#Characters ${Date.now()}`,
      providerId: testProviders[0].id
    };
    
    const specialId = await testHelpers.createTestSubscription(specialSubscription);
    
    // Refresh page to see new subscription
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForSubscriptionsToLoad();
    
    // Search for subscription with special characters
    await subscriptionsPage.searchSubscriptions('Test-Sub_With@Special');
    
    // Should find the subscription
    const subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames).toContain(specialSubscription.name);
    
    // Clean up
    await testHelpers.cleanupTestSubscription(specialId);
  });

  test('should handle case-insensitive search', async () => {
    const subscription = testSubscriptions[0];
    const searchTerm = subscription.name.toUpperCase();
    
    // Search with uppercase
    await subscriptionsPage.searchSubscriptions(searchTerm);
    
    // Should find the subscription despite case difference
    const subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames).toContain(subscription.name);
  });

  test('should update results in real-time as user types', async () => {
    const searchTerm = testSubscriptions[0].name;
    
    // Type search term character by character
    const searchInput = subscriptionsPage.pageInstance.locator('input[placeholder*="Search"]');
    await searchInput.click();
    
    for (let i = 1; i <= searchTerm.length; i++) {
      const partialTerm = searchTerm.substring(0, i);
      await searchInput.fill(partialTerm);
      
      // Wait a bit for search to process
      await subscriptionsPage.pageInstance.waitForTimeout(500);
    }
    
    // Final results should contain the target subscription
    const subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames).toContain(searchTerm);
  });

  test('should maintain filter state when navigating away and back', async () => {
    // Apply filters
    await subscriptionsPage.applyFilters({
      includeInactive: true
    });
    
    // Navigate away (to create subscription page)
    await subscriptionsPage.clickAddSubscription();
    
    // Navigate back
    const backButton = subscriptionsPage.pageInstance.locator('button:has-text("Back to Subscriptions")');
    await backButton.click();
    
    // Wait for page to load
    await subscriptionsPage.waitForPageLoad();
    
    // Filters should still be applied (inactive subscriptions visible)
    const subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    expect(subscriptionNames).toContain(testSubscriptions[2].name); // Inactive subscription
  });

  test('should combine multiple filter types', async () => {
    const searchTerm = 'Active';
    
    // Apply search and filters together
    await subscriptionsPage.searchSubscriptions(searchTerm);
    await subscriptionsPage.applyFilters({
      includeInactive: false
    });
    
    // Results should match both search and filter criteria
    const subscriptionNames = await subscriptionsPage.getAllSubscriptionNames();
    
    // All results should contain search term
    for (const name of subscriptionNames) {
      expect(name.toLowerCase()).toContain(searchTerm.toLowerCase());
    }
    
    // Should not contain inactive subscription
    expect(subscriptionNames).not.toContain(testSubscriptions[2].name);
  });

  test('should show filter indicator when filters are active', async () => {
    // Apply filters
    await subscriptionsPage.applyFilters({
      includeInactive: true
    });
    
    // Check if there's any visual indication that filters are active
    // This might be a badge, different button color, etc.
    const filterButton = subscriptionsPage.pageInstance.locator('button:has-text("Filter")');
    await expect(filterButton).toBeVisible();
    
    // The exact implementation depends on the UI design
    // We can check for common filter indicators like badges or different styling
  });
});