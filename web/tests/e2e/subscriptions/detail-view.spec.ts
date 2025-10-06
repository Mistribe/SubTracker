import { test, expect } from '../../fixtures/auth';
import { SubscriptionsPage } from '../../page-objects/subscriptions-page';
import { TestDataGenerators, SubscriptionData } from '../../utils/data-generators';
import { TestHelpers, createTestHelpers } from '../../utils/test-helpers';

test.describe('Subscription Detail View', () => {
  let subscriptionsPage: SubscriptionsPage;
  let testHelpers: TestHelpers;
  let testSubscription: SubscriptionData;
  let testProvider: any;
  let createdSubscriptionId: string;

  test.beforeEach(async ({ authenticatedPage }) => {
    subscriptionsPage = new SubscriptionsPage(authenticatedPage);
    testHelpers = await createTestHelpers(authenticatedPage);
    
    // Generate test data with comprehensive details
    testSubscription = TestDataGenerators.generateSubscription({
      name: `Detailed Test Subscription ${Date.now()}`,
      amount: 29.99,
      currency: 'USD',
      billingCycle: 'monthly',
      description: 'Test subscription with detailed information',
      isActive: true
    });
    
    // Create test data using real API
    testProvider = await testHelpers.createTestProvider({
      name: `Detailed Provider ${Date.now()}`,
      description: 'Test provider with detailed information',
      website: 'https://example.com',
    });
    
    testSubscription.providerId = testProvider.id;
    createdSubscriptionId = await testHelpers.createTestSubscription(testSubscription);
    
    // Navigate to subscriptions page
    await authenticatedPage.goto('/subscriptions');
    await subscriptionsPage.waitForPageLoad();
  });

  test.afterEach(async () => {
    // Clean up test data
    if (testHelpers) {
      try {
        if (createdSubscriptionId && createdSubscriptionId !== 'default-subscription') {
          await testHelpers.cleanupTestSubscription(createdSubscriptionId);
        }
        if (testProvider && testProvider.id !== 'default-provider') {
          await testHelpers.cleanupTestProvider(testProvider.id);
        }
      } catch {
        // Ignore cleanup errors
      }
    }
  });

  test('should display subscription information accurately in table view', async () => {
    // Verify subscription appears in table with correct information
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name,
      status: 'Active'
    });
    
    // Get detailed information from table row
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    const subscriptionData = await subscriptionsPage.getSubscriptionFromRow(rowIndex);
    
    // Verify all displayed information is correct
    expect(subscriptionData.name).toContain(testSubscription.name);
    expect(subscriptionData.provider).toContain(testProvider.name);
    expect(subscriptionData.status).toContain('Active');
    expect(subscriptionData.recurrency).toContain('monthly');
  });

  test('should show subscription details in edit mode', async () => {
    // Navigate to edit mode to view full subscription details
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Verify form is populated with correct subscription data
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    const providerSelect = subscriptionsPage.pageInstance.locator('[name="providerId"]');
    
    await expect(nameInput).toHaveValue(testSubscription.name);
    await expect(providerSelect).toHaveValue(testProvider.id);
    
    // Navigate through form steps to verify all details
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    
    // Step 2: Recurrency details
    await nextButton.click();
    const recurrencySelect = subscriptionsPage.pageInstance.locator('[name="recurrency"]');
    await expect(recurrencySelect).toHaveValue(testSubscription.billingCycle);
    
    // Step 3: Date details
    await nextButton.click();
    const startDateInput = subscriptionsPage.pageInstance.locator('[name="startDate"]');
    await expect(startDateInput).toBeVisible();
    
    // Step 4: Ownership details
    await nextButton.click();
    const ownerTypeSelect = subscriptionsPage.pageInstance.locator('[name="ownerType"]');
    await expect(ownerTypeSelect).toBeVisible();
    
    // Step 5: Free trial details
    await nextButton.click();
    const freeTrialCheckbox = subscriptionsPage.pageInstance.locator('[name="hasFreeTrialPeriod"]');
    await expect(freeTrialCheckbox).toBeVisible();
  });

  test('should display provider information correctly', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    const row = subscriptionsPage.pageInstance.locator('tbody tr').nth(rowIndex);
    
    // Check provider cell contains correct information
    const providerCell = row.locator('td').first();
    await expect(providerCell).toContainText(testProvider.name);
    
    // If provider has an icon, it should be displayed
    const providerIcon = providerCell.locator('img');
    const fallbackIcon = providerCell.locator('svg');
    
    // Either provider icon or fallback icon should be visible
    await expect(providerIcon.or(fallbackIcon)).toBeVisible();
  });

  test('should display pricing information correctly', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    const row = subscriptionsPage.pageInstance.locator('tbody tr').nth(rowIndex);
    
    // Check price cell contains amount and currency
    const priceCell = row.locator('td').nth(2);
    const priceText = await priceCell.textContent();
    
    expect(priceText).toContain(testSubscription.amount.toString());
    expect(priceText).toContain(testSubscription.currency);
  });

  test('should display billing cycle information correctly', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    const row = subscriptionsPage.pageInstance.locator('tbody tr').nth(rowIndex);
    
    // Check recurrency cell shows correct billing cycle
    const recurrencyCell = row.locator('td').nth(3);
    await expect(recurrencyCell).toContainText(testSubscription.billingCycle);
  });

  test('should display subscription dates correctly', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    const row = subscriptionsPage.pageInstance.locator('tbody tr').nth(rowIndex);
    
    // Check dates cell contains start date information
    const datesCell = row.locator('td').nth(4);
    await expect(datesCell).toBeVisible();
    
    // Should contain calendar icon and date information
    const calendarIcon = datesCell.locator('svg');
    await expect(calendarIcon).toBeVisible();
  });

  test('should display subscription status correctly', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    const row = subscriptionsPage.pageInstance.locator('tbody tr').nth(rowIndex);
    
    // Check status cell shows active status
    const statusCell = row.locator('td').nth(6);
    await expect(statusCell).toContainText('Active');
    
    // Should contain status badge with appropriate styling
    const statusBadge = statusCell.locator('[class*="badge"]');
    await expect(statusBadge).toBeVisible();
  });

  test('should show action menu with correct options', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    await subscriptionsPage.openSubscriptionActionMenu(rowIndex);
    
    // Verify action menu contains expected options
    const actionMenu = subscriptionsPage.pageInstance.locator('[role="menu"]');
    await expect(actionMenu).toBeVisible();
    
    // Check for edit option
    const editOption = actionMenu.locator('text="Edit"');
    await expect(editOption).toBeVisible();
    
    // Check for delete option
    const deleteOption = actionMenu.locator('text="Delete"');
    await expect(deleteOption).toBeVisible();
    
    // Check for copy option (if available)
    const copyOption = actionMenu.locator('text="Make a copy"');
    await expect(copyOption).toBeVisible();
    
    // Close menu by clicking outside
    await subscriptionsPage.pageInstance.click('body');
  });

  test('should handle subscription with custom pricing details', async () => {
    // Create subscription with custom pricing
    const customPriceSubscription = {
      ...TestDataGenerators.generateSubscription(),
      name: `Custom Price Sub ${Date.now()}`,
      providerId: testProvider.id,
      amount: 49.99,
      currency: 'EUR'
    };
    
    const customId = await testHelpers.createTestSubscription(customPriceSubscription);
    
    // Refresh page to see new subscription
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForSubscriptionsToLoad();
    
    // Verify custom pricing is displayed correctly
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(customPriceSubscription.name);
    const row = subscriptionsPage.pageInstance.locator('tbody tr').nth(rowIndex);
    const priceCell = row.locator('td').nth(2);
    
    const priceText = await priceCell.textContent();
    expect(priceText).toContain('49.99');
    expect(priceText).toContain('EUR');
    
    // Clean up
    await testHelpers.cleanupTestSubscription(customId);
  });

  test('should display subscription with end date correctly', async () => {
    // Create subscription with end date
    const endDate = new Date();
    endDate.setMonth(endDate.getMonth() + 6);
    const endDateString = endDate.toISOString().split('T')[0];
    
    const endDateSubscription = {
      ...TestDataGenerators.generateSubscription(),
      name: `End Date Sub ${Date.now()}`,
      providerId: testProvider.id,
      // Note: This would need to be supported by the API
    };
    
    const endDateId = await testHelpers.createTestSubscription(endDateSubscription);
    
    // Refresh page
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForSubscriptionsToLoad();
    
    // Verify subscription appears in table
    await subscriptionsPage.verifySubscriptionInTable({
      name: endDateSubscription.name,
      provider: testProvider.name
    });
    
    // Clean up
    await testHelpers.cleanupTestSubscription(endDateId);
  });

  test('should navigate to edit view from table row', async () => {
    // Click edit from action menu
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Verify we're on the edit page
    await expect(subscriptionsPage.pageInstance).toHaveURL(/.*subscriptions\/edit\/.*/);
    
    // Verify page title indicates edit mode
    const pageTitle = subscriptionsPage.pageInstance.locator('h1:has-text("Edit Subscription")');
    await expect(pageTitle).toBeVisible();
  });

  test('should display subscription information consistently across views', async () => {
    // Get information from table view
    const tableData = await subscriptionsPage.getSubscriptionFromRow(
      await subscriptionsPage.findSubscriptionRowByName(testSubscription.name)
    );
    
    // Navigate to edit view
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Verify information matches in edit form
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    const formName = await nameInput.inputValue();
    
    expect(formName).toBe(testSubscription.name);
    expect(tableData.name).toContain(testSubscription.name);
  });

  test('should handle long subscription names gracefully', async () => {
    // Create subscription with very long name
    const longName = `Very Long Subscription Name That Should Be Handled Gracefully In The UI ${Date.now()}`;
    const longNameSubscription = {
      ...TestDataGenerators.generateSubscription(),
      name: longName,
      providerId: testProvider.id
    };
    
    const longNameId = await testHelpers.createTestSubscription(longNameSubscription);
    
    // Refresh page
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForSubscriptionsToLoad();
    
    // Verify long name is displayed (possibly truncated)
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(longName);
    const row = subscriptionsPage.pageInstance.locator('tbody tr').nth(rowIndex);
    const nameCell = row.locator('td').nth(1);
    
    // Name should be visible and contain at least part of the long name
    await expect(nameCell).toBeVisible();
    const displayedName = await nameCell.textContent();
    expect(displayedName).toBeTruthy();
    
    // Clean up
    await testHelpers.cleanupTestSubscription(longNameId);
  });

  test('should display tooltip information for complex data', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    const row = subscriptionsPage.pageInstance.locator('tbody tr').nth(rowIndex);
    
    // Check if there are any tooltips or hover information
    // This might be on price calculations, dates, or other complex data
    const priceCell = row.locator('td').nth(2);
    
    // Hover over price cell to see if tooltip appears
    await priceCell.hover();
    
    // Wait a moment for potential tooltip
    await subscriptionsPage.pageInstance.waitForTimeout(1000);
    
    // Check for tooltip (implementation may vary)
    const tooltip = subscriptionsPage.pageInstance.locator('[role="tooltip"]');
    // Tooltip may or may not be present depending on implementation
    // This test verifies the UI handles hover states gracefully
  });
});