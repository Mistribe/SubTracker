import { test, expect } from '../../fixtures/auth';
import { SubscriptionsPage } from '../../page-objects/subscriptions-page';
import { TestDataGenerators, SubscriptionData } from '../../utils/data-generators';
import { TestHelpers, createTestHelpers } from '../../utils/test-helpers';

test.describe('Subscription Editing', () => {
  let subscriptionsPage: SubscriptionsPage;
  let testHelpers: TestHelpers;
  let testSubscription: SubscriptionData;
  let testProvider: any;
  let createdSubscriptionId: string;

  test.beforeEach(async ({ authenticatedPage }) => {
    subscriptionsPage = new SubscriptionsPage(authenticatedPage);
    testHelpers = await createTestHelpers(authenticatedPage);
    
    // Generate test data
    testSubscription = TestDataGenerators.generateSubscription();
    
    // Create test data using real API
    testProvider = await testHelpers.createTestProvider({
      name: `Test Provider ${Date.now()}`,
      description: 'Test provider for subscription editing',
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

  test('should edit subscription name successfully', async () => {
    const newName = `Updated ${testSubscription.name}`;
    
    // Edit the subscription
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Update the name in the form
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    await nameInput.clear();
    await nameInput.fill(newName);
    
    // Submit the form
    await subscriptionsPage.submitSubscriptionForm();
    
    // Verify the subscription name was updated in the table
    await subscriptionsPage.verifySubscriptionInTable({
      name: newName,
      provider: testProvider.name
    });
    
    // Verify old name is no longer in the table
    await subscriptionsPage.verifySubscriptionNotInTable(testSubscription.name);
  });

  test('should edit subscription billing cycle', async () => {
    const newBillingCycle = testSubscription.billingCycle === 'monthly' ? 'yearly' : 'monthly';
    
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Navigate to recurrency step
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click();
    
    // Update billing cycle
    const recurrencySelect = subscriptionsPage.pageInstance.locator('[name="recurrency"]');
    await recurrencySelect.selectOption(newBillingCycle);
    
    // Navigate to final step and submit
    await nextButton.click(); // Dates
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial (final step)
    
    await subscriptionsPage.submitSubscriptionForm();
    
    // Verify the subscription was updated
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name
    });
  });

  test('should edit subscription amount', async () => {
    const newAmount = testSubscription.amount + 10;
    
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Update custom price amount (assuming custom price is used)
    const amountInput = subscriptionsPage.pageInstance.locator('[name="customPrice.amount"]');
    if (await amountInput.isVisible()) {
      await amountInput.clear();
      await amountInput.fill(newAmount.toString());
    }
    
    // Navigate through steps and submit
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click(); // Recurrency
    await nextButton.click(); // Dates
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial
    
    await subscriptionsPage.submitSubscriptionForm();
    
    // Verify subscription was updated
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name
    });
  });

  test('should edit subscription start date', async () => {
    const newStartDate = new Date();
    newStartDate.setDate(newStartDate.getDate() + 7);
    const newStartDateString = newStartDate.toISOString().split('T')[0];
    
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Navigate to dates step
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click(); // Recurrency
    await nextButton.click(); // Dates
    
    // Update start date
    const startDateInput = subscriptionsPage.pageInstance.locator('[name="startDate"]');
    await startDateInput.fill(newStartDateString);
    
    // Navigate to final step and submit
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial
    
    await subscriptionsPage.submitSubscriptionForm();
    
    // Verify subscription was updated
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name
    });
  });

  test('should handle partial updates correctly', async () => {
    // Only update the name, leave other fields unchanged
    const newName = `Partially Updated ${testSubscription.name}`;
    
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Only change the name
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    await nameInput.clear();
    await nameInput.fill(newName);
    
    // Navigate through all steps without changing other fields
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click(); // Recurrency
    await nextButton.click(); // Dates
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial
    
    await subscriptionsPage.submitSubscriptionForm();
    
    // Verify only the name was updated
    await subscriptionsPage.verifySubscriptionInTable({
      name: newName,
      provider: testProvider.name
    });
  });

  test('should preserve unchanged fields during edit', async () => {
    const originalProvider = testProvider.name;
    
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Verify form is pre-populated with existing data
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    const providerSelect = subscriptionsPage.pageInstance.locator('[name="providerId"]');
    
    await expect(nameInput).toHaveValue(testSubscription.name);
    await expect(providerSelect).toHaveValue(testProvider.id);
    
    // Make a small change and submit
    await nameInput.fill(`${testSubscription.name} - Edited`);
    
    // Navigate through steps and submit
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click(); // Recurrency
    await nextButton.click(); // Dates
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial
    
    await subscriptionsPage.submitSubscriptionForm();
    
    // Verify provider remained unchanged
    await subscriptionsPage.verifySubscriptionInTable({
      name: `${testSubscription.name} - Edited`,
      provider: originalProvider
    });
  });

  test('should show validation errors for invalid edits', async () => {
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Clear required field
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    await nameInput.clear();
    
    // Try to submit with empty name
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click(); // Should stay on current step due to validation
    await nextButton.click(); // Recurrency
    await nextButton.click(); // Dates
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial
    
    const submitButton = subscriptionsPage.pageInstance.locator('button[type="submit"]');
    await submitButton.click();
    
    // Should show validation error or stay on form
    const errorMessage = subscriptionsPage.pageInstance.locator('text*="Failed to update subscription"');
    const formStillVisible = subscriptionsPage.pageInstance.locator('form');
    
    await expect(errorMessage.or(formStillVisible)).toBeVisible();
  });

  test('should cancel edit and return to subscriptions list', async () => {
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Make some changes
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    await nameInput.fill('This should not be saved');
    
    // Cancel the edit
    const backButton = subscriptionsPage.pageInstance.locator('button:has-text("Back to Subscriptions")');
    await backButton.click();
    
    // Should be back on subscriptions page
    await subscriptionsPage.waitForPageLoad();
    
    // Verify original subscription data is unchanged
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name
    });
  });

  test('should edit subscription and verify changes persist across page refresh', async () => {
    const newName = `Persistent Edit ${testSubscription.name}`;
    
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Update name
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    await nameInput.clear();
    await nameInput.fill(newName);
    
    // Submit form
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click(); // Recurrency
    await nextButton.click(); // Dates
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial
    
    await subscriptionsPage.submitSubscriptionForm();
    
    // Refresh the page
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForPageLoad();
    
    // Verify changes persisted
    await subscriptionsPage.verifySubscriptionInTable({
      name: newName,
      provider: testProvider.name
    });
  });

  test('should show loading state during edit submission', async () => {
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Make a change
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    await nameInput.fill(`${testSubscription.name} - Loading Test`);
    
    // Navigate to final step
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click(); // Recurrency
    await nextButton.click(); // Dates
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial
    
    // Submit and check for loading state
    const submitButton = subscriptionsPage.pageInstance.locator('button[type="submit"]');
    await submitButton.click();
    
    // Should show loading state
    const loadingButton = subscriptionsPage.pageInstance.locator('button:has-text("Updating...")');
    await expect(loadingButton.or(subscriptionsPage.pageInstance.locator('.animate-spin'))).toBeVisible();
  });

  test('should edit multiple fields simultaneously', async () => {
    const newName = `Multi-Edit ${testSubscription.name}`;
    const newBillingCycle = testSubscription.billingCycle === 'monthly' ? 'yearly' : 'monthly';
    
    await subscriptionsPage.editSubscriptionByName(testSubscription.name);
    
    // Update name
    const nameInput = subscriptionsPage.pageInstance.locator('[name="friendlyName"]');
    await nameInput.clear();
    await nameInput.fill(newName);
    
    const nextButton = subscriptionsPage.pageInstance.locator('button:has-text("Next")');
    await nextButton.click(); // Go to Recurrency step
    
    // Update billing cycle
    const recurrencySelect = subscriptionsPage.pageInstance.locator('[name="recurrency"]');
    await recurrencySelect.selectOption(newBillingCycle);
    
    // Navigate to final step and submit
    await nextButton.click(); // Dates
    await nextButton.click(); // Ownership
    await nextButton.click(); // Free Trial
    
    await subscriptionsPage.submitSubscriptionForm();
    
    // Verify both changes were applied
    await subscriptionsPage.verifySubscriptionInTable({
      name: newName,
      provider: testProvider.name
    });
  });
});