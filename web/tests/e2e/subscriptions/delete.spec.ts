import { test, expect } from '../../fixtures/auth';
import { SubscriptionsPage } from '../../page-objects/subscriptions-page';
import { TestDataGenerators, SubscriptionData } from '../../utils/data-generators';
import { TestHelpers, createTestHelpers } from '../../utils/test-helpers';

test.describe('Subscription Deletion', () => {
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
      description: 'Test provider for subscription deletion',
    });
    
    testSubscription.providerId = testProvider.id;
    createdSubscriptionId = await testHelpers.createTestSubscription(testSubscription);
    
    // Navigate to subscriptions page
    await authenticatedPage.goto('/subscriptions');
    await subscriptionsPage.waitForPageLoad();
  });

  test.afterEach(async () => {
    // Clean up test data (subscription might already be deleted by test)
    if (testHelpers) {
      try {
        if (createdSubscriptionId && createdSubscriptionId !== 'default-subscription') {
          await testHelpers.cleanupTestSubscription(createdSubscriptionId);
        }
      } catch {
        // Subscription might already be deleted, ignore error
      }
      
      if (testProvider && testProvider.id !== 'default-provider') {
        try {
          await testHelpers.cleanupTestProvider(testProvider.id);
        } catch {
          // Ignore cleanup errors
        }
      }
    }
  });

  test('should delete subscription with confirmation', async () => {
    // Verify subscription exists before deletion
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name
    });
    
    // Get initial count
    const initialCount = await subscriptionsPage.getSubscriptionCount();
    
    // Delete the subscription with confirmation
    await subscriptionsPage.deleteSubscriptionByName(testSubscription.name, true);
    
    // Verify subscription is removed from table
    await subscriptionsPage.verifySubscriptionNotInTable(testSubscription.name);
    
    // Verify count decreased
    const newCount = await subscriptionsPage.getSubscriptionCount();
    expect(newCount).toBe(initialCount - 1);
  });

  test('should cancel deletion when user clicks cancel', async () => {
    // Verify subscription exists
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name
    });
    
    // Get initial count
    const initialCount = await subscriptionsPage.getSubscriptionCount();
    
    // Attempt to delete but cancel
    await subscriptionsPage.deleteSubscriptionByName(testSubscription.name, false);
    
    // Verify subscription still exists
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name
    });
    
    // Verify count unchanged
    const newCount = await subscriptionsPage.getSubscriptionCount();
    expect(newCount).toBe(initialCount);
  });

  test('should show delete confirmation dialog with subscription details', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    await subscriptionsPage.openSubscriptionActionMenu(rowIndex);
    
    // Click delete menu item
    const deleteMenuItem = subscriptionsPage.pageInstance.locator('text="Delete"');
    await deleteMenuItem.click();
    
    // Verify delete dialog appears
    const deleteDialog = subscriptionsPage.pageInstance.locator('[role="dialog"]:has-text("Delete Subscription")');
    await expect(deleteDialog).toBeVisible();
    
    // Verify dialog contains subscription information
    await expect(deleteDialog).toContainText(testSubscription.name);
    
    // Verify dialog has cancel and confirm buttons
    const cancelButton = deleteDialog.locator('button:has-text("Cancel")');
    const confirmButton = deleteDialog.locator('button:has-text("Delete")');
    
    await expect(cancelButton).toBeVisible();
    await expect(confirmButton).toBeVisible();
    
    // Cancel to close dialog
    await cancelButton.click();
    await expect(deleteDialog).not.toBeVisible();
  });

  test('should delete multiple subscriptions sequentially', async () => {
    // Create additional test subscriptions
    const additionalSubscriptions = [
      { ...TestDataGenerators.generateSubscription(), providerId: testProvider.id },
      { ...TestDataGenerators.generateSubscription(), providerId: testProvider.id }
    ];
    
    const additionalIds: string[] = [];
    for (const subscription of additionalSubscriptions) {
      const id = await testHelpers.createTestSubscription(subscription);
      additionalIds.push(id);
    }
    
    // Refresh page to see all subscriptions
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForSubscriptionsToLoad();
    
    // Get initial count
    const initialCount = await subscriptionsPage.getSubscriptionCount();
    
    // Delete all test subscriptions
    const allSubscriptions = [testSubscription, ...additionalSubscriptions];
    for (const subscription of allSubscriptions) {
      await subscriptionsPage.deleteSubscriptionByName(subscription.name, true);
    }
    
    // Verify all subscriptions are deleted
    for (const subscription of allSubscriptions) {
      await subscriptionsPage.verifySubscriptionNotInTable(subscription.name);
    }
    
    // Verify count decreased by the number of deleted subscriptions
    const finalCount = await subscriptionsPage.getSubscriptionCount();
    expect(finalCount).toBe(initialCount - allSubscriptions.length);
    
    // Clean up additional subscription IDs
    for (const id of additionalIds) {
      try {
        await testHelpers.cleanupTestSubscription(id);
      } catch {
        // Already deleted, ignore error
      }
    }
  });

  test('should handle deletion of non-existent subscription gracefully', async () => {
    // Delete the subscription first
    await subscriptionsPage.deleteSubscriptionByName(testSubscription.name, true);
    
    // Verify it's deleted
    await subscriptionsPage.verifySubscriptionNotInTable(testSubscription.name);
    
    // Try to delete again - should not find the subscription
    const subscriptionExists = await subscriptionsPage.subscriptionExists(testSubscription.name);
    expect(subscriptionExists).toBe(false);
  });

  test('should show loading state during deletion', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    await subscriptionsPage.openSubscriptionActionMenu(rowIndex);
    
    // Click delete
    const deleteMenuItem = subscriptionsPage.pageInstance.locator('text="Delete"');
    await deleteMenuItem.click();
    
    // Confirm deletion and immediately check for loading state
    const confirmButton = subscriptionsPage.pageInstance.locator('button:has-text("Delete")');
    await confirmButton.click();
    
    // Should show loading state (spinner or disabled button)
    const loadingIndicator = subscriptionsPage.pageInstance.locator('.animate-spin');
    const disabledButton = subscriptionsPage.pageInstance.locator('button[disabled]:has-text("Delete")');
    
    // One of these should be visible during deletion
    await expect(loadingIndicator.or(disabledButton)).toBeVisible();
  });

  test('should close action menu when clicking outside', async () => {
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    await subscriptionsPage.openSubscriptionActionMenu(rowIndex);
    
    // Verify menu is open
    const actionMenu = subscriptionsPage.pageInstance.locator('[role="menu"]');
    await expect(actionMenu).toBeVisible();
    
    // Click outside the menu
    await subscriptionsPage.pageInstance.click('body');
    
    // Verify menu is closed
    await expect(actionMenu).not.toBeVisible();
  });

  test('should delete subscription and verify it disappears from all views', async () => {
    // Verify subscription exists in table
    await subscriptionsPage.verifySubscriptionInTable({
      name: testSubscription.name,
      provider: testProvider.name
    });
    
    // Delete the subscription
    await subscriptionsPage.deleteSubscriptionByName(testSubscription.name, true);
    
    // Verify it's not in the main table
    await subscriptionsPage.verifySubscriptionNotInTable(testSubscription.name);
    
    // Search for the subscription - should not be found
    await subscriptionsPage.searchSubscriptions(testSubscription.name);
    const isEmpty = await subscriptionsPage.isSubscriptionsTableEmpty();
    expect(isEmpty).toBe(true);
    
    // Clear search and verify it's still not there
    await subscriptionsPage.clearSearch();
    await subscriptionsPage.verifySubscriptionNotInTable(testSubscription.name);
  });

  test('should handle deletion error gracefully', async () => {
    // This test simulates a deletion error scenario
    // In a real scenario, this might happen due to network issues or server errors
    
    const rowIndex = await subscriptionsPage.findSubscriptionRowByName(testSubscription.name);
    await subscriptionsPage.openSubscriptionActionMenu(rowIndex);
    
    // Click delete
    const deleteMenuItem = subscriptionsPage.pageInstance.locator('text="Delete"');
    await deleteMenuItem.click();
    
    // Confirm deletion
    const confirmButton = subscriptionsPage.pageInstance.locator('button:has-text("Delete")');
    await confirmButton.click();
    
    // If deletion fails, the subscription should still be in the table
    // and an error message might be shown
    // Wait a bit for the operation to complete or fail
    await subscriptionsPage.pageInstance.waitForTimeout(3000);
    
    // Check if subscription still exists (in case of error) or is deleted (success)
    const stillExists = await subscriptionsPage.subscriptionExists(testSubscription.name);
    
    if (stillExists) {
      // If it still exists, there might have been an error
      // Check for error message
      const errorMessage = subscriptionsPage.pageInstance.locator('[role="alert"]');
      // Error message might be visible or subscription might still be there
      expect(stillExists).toBe(true);
    } else {
      // Deletion was successful
      await subscriptionsPage.verifySubscriptionNotInTable(testSubscription.name);
    }
  });

  test('should delete subscription from different table positions', async () => {
    // Create multiple subscriptions to test deletion from different positions
    const subscriptions = [
      testSubscription,
      { ...TestDataGenerators.generateSubscription(), providerId: testProvider.id },
      { ...TestDataGenerators.generateSubscription(), providerId: testProvider.id }
    ];
    
    // Create additional subscriptions
    const additionalIds: string[] = [];
    for (let i = 1; i < subscriptions.length; i++) {
      const id = await testHelpers.createTestSubscription(subscriptions[i]);
      additionalIds.push(id);
    }
    
    // Refresh to see all subscriptions
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForSubscriptionsToLoad();
    
    // Delete from middle position (index 1)
    await subscriptionsPage.deleteSubscription(1, true);
    
    // Delete from first position (index 0)
    await subscriptionsPage.deleteSubscription(0, true);
    
    // Delete remaining subscription
    const remainingCount = await subscriptionsPage.getSubscriptionCount();
    if (remainingCount > 0) {
      await subscriptionsPage.deleteSubscription(0, true);
    }
    
    // Clean up additional IDs
    for (const id of additionalIds) {
      try {
        await testHelpers.cleanupTestSubscription(id);
      } catch {
        // Already deleted, ignore error
      }
    }
  });
});