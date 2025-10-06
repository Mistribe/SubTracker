/**
 * Subscription Management Happy Path Tests
 * 
 * Comprehensive test suite covering all successful subscription workflows:
 * - Create subscription with valid data
 * - Edit subscription successfully  
 * - Delete subscription with confirmation
 * - View subscription details
 * - Navigate between subscription pages
 * - Basic search functionality
 * 
 * Requirements: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6
 * 
 * Testing Philosophy: Happy Path Only
 * ✅ Tests successful CRUD operations with valid data
 * ✅ Tests successful navigation and UI interactions
 * ✅ Tests working search and basic filtering
 * ❌ No error testing, validation failures, or edge cases
 */

import { test, expect } from '../../fixtures/auth';
import { SubscriptionsPage } from '../../page-objects/subscriptions-page';
import { TestDataGenerators, SubscriptionData } from '../../utils/data-generators';
import { TestHelpers, createTestHelpers } from '../../utils/test-helpers';

test.describe('Subscription Management Happy Path', () => {
  let subscriptionsPage: SubscriptionsPage;
  let testHelpers: TestHelpers;
  let testProvider: { id: string; name: string };
  let createdSubscriptionIds: string[] = [];

  test.beforeEach(async ({ authenticatedPage }) => {
    console.log('🚀 Setting up subscription management test');
    
    subscriptionsPage = new SubscriptionsPage(authenticatedPage);
    testHelpers = await createTestHelpers(authenticatedPage);
    
    // Create a test provider for our subscriptions
    testProvider = await testHelpers.createTestProvider({
      name: `Test Provider ${Date.now()}`,
      description: 'Test provider for subscription management',
    });
    console.log(`✅ Created test provider: ${testProvider.name}`);
    
    // Navigate to subscriptions page
    await authenticatedPage.goto('/subscriptions');
    await subscriptionsPage.waitForPageLoad();
    console.log('✅ Navigated to subscriptions page');
  });

  test.afterEach(async () => {
    console.log('🧹 Starting cleanup process');
    
    // Clean up created subscriptions
    for (const subscriptionId of createdSubscriptionIds) {
      try {
        await testHelpers.cleanupTestSubscription(subscriptionId);
        console.log(`✅ Cleaned up subscription: ${subscriptionId}`);
      } catch {
        console.log(`ℹ️ Subscription already cleaned up: ${subscriptionId}`);
      }
    }
    
    // Clean up test provider
    try {
      await testHelpers.cleanupTestProvider(testProvider.id);
      console.log(`✅ Cleaned up test provider: ${testProvider.name}`);
    } catch {
      console.log(`ℹ️ Provider already cleaned up: ${testProvider.name}`);
    }
    
    // Reset arrays
    createdSubscriptionIds = [];
    console.log('✅ Cleanup completed');
  });

  test('should display subscriptions page layout correctly', async () => {
    console.log('🔍 Testing subscriptions page layout');
    
    // Verify page title is visible
    const pageTitle = subscriptionsPage.pageInstance.locator('h1');
    await expect(pageTitle).toBeVisible();
    console.log('✅ Page title is visible');
    
    // Verify add subscription button is visible and enabled
    const addButton = subscriptionsPage.pageInstance.locator('button:has-text("Add Subscription")');
    await expect(addButton).toBeVisible();
    await expect(addButton).toBeEnabled();
    console.log('✅ Add subscription button is visible and enabled');
    
    // Verify search input is visible and interactive
    const searchInput = subscriptionsPage.pageInstance.locator('input[placeholder*="Search"], input[type="search"]');
    await expect(searchInput).toBeVisible();
    await expect(searchInput).toBeEnabled();
    console.log('✅ Search input is visible and enabled');
    
    // Verify table structure is present
    const table = subscriptionsPage.pageInstance.locator('table');
    await expect(table).toBeVisible();
    console.log('✅ Subscriptions table is visible');
  });

  test('should create a new subscription successfully', async () => {
    console.log('🆕 Testing subscription creation');
    
    // Generate test subscription data
    const subscriptionData = TestDataGenerators.generateSubscription({
      name: `Test Subscription ${Date.now()}`,
      providerId: testProvider.id,
      amount: 9.99,
      currency: 'USD',
      billingCycle: 'monthly'
    });
    console.log(`📝 Generated subscription data: ${subscriptionData.name}`);
    
    // Click add subscription button
    await subscriptionsPage.clickAddSubscription();
    console.log('✅ Clicked add subscription button');
    
    // Verify navigation to create page
    await expect(subscriptionsPage.pageInstance).toHaveURL(/.*subscriptions\/create/);
    console.log('✅ Navigated to create subscription page');
    
    // Fill out the subscription form using helper methods
    await subscriptionsPage.fillSubscriptionForm(subscriptionData);
    console.log('✅ Filled subscription form');
    
    // Submit the form
    await subscriptionsPage.submitSubscriptionForm();
    console.log('✅ Submitted subscription form');
    
    // Wait for successful navigation back to subscriptions list
    await subscriptionsPage.pageInstance.waitForURL('**/subscriptions', { timeout: 15000 });
    console.log('✅ Successfully navigated back to subscriptions list');
    
    // Verify subscription appears in the table
    await subscriptionsPage.verifySubscriptionInTable({
      name: subscriptionData.name,
      provider: testProvider.name
    });
    console.log(`✅ Verified subscription appears in table: ${subscriptionData.name}`);
    
    // Get the created subscription ID for cleanup
    const subscriptionId = await testHelpers.getSubscriptionIdByName(subscriptionData.name);
    if (subscriptionId) {
      createdSubscriptionIds.push(subscriptionId);
      console.log(`📝 Added subscription to cleanup list: ${subscriptionId}`);
    }
  });

  test('should edit an existing subscription successfully', async () => {
    console.log('✏️ Testing subscription editing');
    
    // First create a subscription to edit
    const originalData = TestDataGenerators.generateSubscription({
      name: `Original Subscription ${Date.now()}`,
      providerId: testProvider.id,
      amount: 15.99,
      currency: 'USD',
      billingCycle: 'monthly'
    });
    
    const subscriptionId = await testHelpers.createTestSubscription(originalData);
    createdSubscriptionIds.push(subscriptionId);
    console.log(`✅ Created subscription to edit: ${originalData.name}`);
    
    // Refresh page to see the new subscription
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForPageLoad();
    
    // Edit the subscription
    await subscriptionsPage.editSubscriptionByName(originalData.name);
    console.log('✅ Opened subscription for editing');
    
    // Update the subscription data
    const updatedData = {
      name: `${originalData.name} - Edited`,
      amount: 19.99
    };
    
    await subscriptionsPage.updateSubscriptionForm(updatedData);
    console.log('✅ Updated subscription form');
    
    // Submit the changes
    await subscriptionsPage.submitSubscriptionForm();
    console.log('✅ Submitted subscription changes');
    
    // Wait for navigation back to list
    await subscriptionsPage.pageInstance.waitForURL('**/subscriptions', { timeout: 15000 });
    console.log('✅ Navigated back to subscriptions list');
    
    // Verify the updated subscription appears in the table
    await subscriptionsPage.verifySubscriptionInTable({
      name: updatedData.name,
      provider: testProvider.name
    });
    console.log(`✅ Verified updated subscription in table: ${updatedData.name}`);
  });

  test('should delete a subscription successfully', async () => {
    console.log('🗑️ Testing subscription deletion');
    
    // First create a subscription to delete
    const subscriptionData = TestDataGenerators.generateSubscription({
      name: `Subscription to Delete ${Date.now()}`,
      providerId: testProvider.id,
      amount: 12.99,
      currency: 'USD',
      billingCycle: 'monthly'
    });
    
    const subscriptionId = await testHelpers.createTestSubscription(subscriptionData);
    console.log(`✅ Created subscription to delete: ${subscriptionData.name}`);
    
    // Refresh page to see the new subscription
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForPageLoad();
    
    // Delete the subscription
    await subscriptionsPage.deleteSubscriptionByName(subscriptionData.name);
    console.log('✅ Initiated subscription deletion');
    
    // Verify subscription is no longer in the table
    await subscriptionsPage.verifySubscriptionNotInTable(subscriptionData.name);
    console.log(`✅ Verified subscription removed from table: ${subscriptionData.name}`);
    
    // Note: No need to add to cleanup list since it's already deleted
  });

  test('should view subscription details successfully', async () => {
    console.log('👁️ Testing subscription detail view');
    
    // Create a subscription with detailed information
    const detailedData = TestDataGenerators.generateSubscription({
      name: `Detailed Subscription ${Date.now()}`,
      providerId: testProvider.id,
      amount: 25.99,
      currency: 'USD',
      billingCycle: 'yearly',
      description: 'Test subscription with detailed information'
    });
    
    const subscriptionId = await testHelpers.createTestSubscription(detailedData);
    createdSubscriptionIds.push(subscriptionId);
    console.log(`✅ Created detailed subscription: ${detailedData.name}`);
    
    // Refresh page to see the new subscription
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForPageLoad();
    
    // Click on subscription to view details
    await subscriptionsPage.viewSubscriptionDetails(detailedData.name);
    console.log('✅ Opened subscription details');
    
    // Verify we're on the detail page
    await expect(subscriptionsPage.pageInstance).toHaveURL(/.*subscriptions\/[^\/]+/);
    console.log('✅ Navigated to subscription detail page');
    
    // Verify subscription details are displayed
    await subscriptionsPage.verifySubscriptionDetails({
      name: detailedData.name,
      provider: testProvider.name,
      amount: detailedData.amount.toString(),
      billingCycle: detailedData.billingCycle
    });
    console.log('✅ Verified subscription details are displayed correctly');
  });

  test('should search for subscriptions successfully', async () => {
    console.log('🔍 Testing subscription search functionality');
    
    // Create multiple subscriptions with different names
    const subscriptions = [
      TestDataGenerators.generateSubscription({
        name: `Netflix Test ${Date.now()}`,
        providerId: testProvider.id,
        amount: 15.99
      }),
      TestDataGenerators.generateSubscription({
        name: `Spotify Test ${Date.now()}`,
        providerId: testProvider.id,
        amount: 9.99
      }),
      TestDataGenerators.generateSubscription({
        name: `Adobe Test ${Date.now()}`,
        providerId: testProvider.id,
        amount: 29.99
      })
    ];
    
    // Create all subscriptions
    for (const subscription of subscriptions) {
      const id = await testHelpers.createTestSubscription(subscription);
      createdSubscriptionIds.push(id);
      console.log(`✅ Created subscription: ${subscription.name}`);
    }
    
    // Refresh page to see all subscriptions
    await subscriptionsPage.pageInstance.reload();
    await subscriptionsPage.waitForPageLoad();
    
    // Search for "Netflix" subscription
    await subscriptionsPage.searchSubscriptions('Netflix');
    console.log('✅ Performed search for "Netflix"');
    
    // Verify search results contain Netflix subscription
    const searchResults = await subscriptionsPage.getAllSubscriptionNames();
    const netflixSubscription = subscriptions.find(s => s.name.includes('Netflix'));
    expect(searchResults).toContain(netflixSubscription!.name);
    console.log('✅ Verified Netflix subscription appears in search results');
    
    // Clear search to show all subscriptions again
    await subscriptionsPage.clearSearch();
    console.log('✅ Cleared search');
    
    // Verify all subscriptions are visible again
    const allResults = await subscriptionsPage.getAllSubscriptionNames();
    for (const subscription of subscriptions) {
      expect(allResults).toContain(subscription.name);
    }
    console.log('✅ Verified all subscriptions visible after clearing search');
  });

  test('should navigate between subscription pages successfully', async () => {
    console.log('🧭 Testing navigation between subscription pages');
    
    // Start on subscriptions list page
    await expect(subscriptionsPage.pageInstance).toHaveURL(/.*subscriptions(?!\/)/);
    console.log('✅ Confirmed on subscriptions list page');
    
    // Navigate to create subscription page
    await subscriptionsPage.clickAddSubscription();
    await expect(subscriptionsPage.pageInstance).toHaveURL(/.*subscriptions\/create/);
    console.log('✅ Successfully navigated to create subscription page');
    
    // Navigate back to subscriptions list
    const backButton = subscriptionsPage.pageInstance.locator('button:has-text("Back"), a:has-text("Back")');
    if (await backButton.count() > 0) {
      await backButton.first().click();
    } else {
      // Alternative: use browser back button
      await subscriptionsPage.pageInstance.goBack();
    }
    
    await expect(subscriptionsPage.pageInstance).toHaveURL(/.*subscriptions(?!\/create)/);
    console.log('✅ Successfully navigated back to subscriptions list');
    
    // Verify page is functional after navigation
    const addButton = subscriptionsPage.pageInstance.locator('button:has-text("Add Subscription")');
    await expect(addButton).toBeVisible();
    await expect(addButton).toBeEnabled();
    console.log('✅ Verified page is functional after navigation');
  });

  test('should handle empty search results gracefully', async () => {
    console.log('🔍 Testing empty search results handling');
    
    // Search for something that doesn't exist
    const nonExistentSearch = `NonExistent${Date.now()}`;
    await subscriptionsPage.searchSubscriptions(nonExistentSearch);
    console.log(`✅ Searched for non-existent term: ${nonExistentSearch}`);
    
    // Verify empty state is handled gracefully
    const isEmpty = await subscriptionsPage.isSubscriptionsTableEmpty();
    expect(isEmpty).toBe(true);
    console.log('✅ Verified empty search results are handled gracefully');
    
    // Clear search to return to normal state
    await subscriptionsPage.clearSearch();
    console.log('✅ Cleared search to return to normal state');
  });

  test('should complete full subscription lifecycle successfully', async () => {
    console.log('🔄 Testing complete subscription lifecycle (create → edit → view → delete)');
    
    // Step 1: Create subscription
    const subscriptionData = TestDataGenerators.generateSubscription({
      name: `Lifecycle Test ${Date.now()}`,
      providerId: testProvider.id,
      amount: 19.99,
      currency: 'USD',
      billingCycle: 'monthly'
    });
    
    await subscriptionsPage.clickAddSubscription();
    await subscriptionsPage.fillSubscriptionForm(subscriptionData);
    await subscriptionsPage.submitSubscriptionForm();
    await subscriptionsPage.pageInstance.waitForURL('**/subscriptions', { timeout: 15000 });
    console.log('✅ Step 1: Created subscription successfully');
    
    // Step 2: Edit subscription
    await subscriptionsPage.editSubscriptionByName(subscriptionData.name);
    const updatedName = `${subscriptionData.name} - Updated`;
    await subscriptionsPage.updateSubscriptionForm({ name: updatedName });
    await subscriptionsPage.submitSubscriptionForm();
    await subscriptionsPage.pageInstance.waitForURL('**/subscriptions', { timeout: 15000 });
    console.log('✅ Step 2: Edited subscription successfully');
    
    // Step 3: View subscription details
    await subscriptionsPage.viewSubscriptionDetails(updatedName);
    await expect(subscriptionsPage.pageInstance).toHaveURL(/.*subscriptions\/[^\/]+/);
    await subscriptionsPage.verifySubscriptionDetails({
      name: updatedName,
      provider: testProvider.name
    });
    console.log('✅ Step 3: Viewed subscription details successfully');
    
    // Navigate back to list
    await subscriptionsPage.pageInstance.goBack();
    await subscriptionsPage.waitForPageLoad();
    
    // Step 4: Delete subscription
    await subscriptionsPage.deleteSubscriptionByName(updatedName);
    await subscriptionsPage.verifySubscriptionNotInTable(updatedName);
    console.log('✅ Step 4: Deleted subscription successfully');
    
    console.log('🎉 Complete subscription lifecycle test completed successfully');
  });
});