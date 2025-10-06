/**
 * Provider Management Happy Path Tests
 * 
 * Comprehensive test suite covering successful provider workflows:
 * - Create provider with valid data
 * - Edit provider successfully
 * - Delete provider with confirmation
 * - View provider details
 * - Search for providers
 * - Navigate between provider pages
 * 
 * Requirements: 2.1, 2.2, 2.3, 2.4, 2.5
 * 
 * Testing Philosophy: Happy Path Only
 * ✅ Tests successful CRUD operations with valid data
 * ✅ Tests successful navigation and UI interactions
 * ✅ Tests working search functionality
 * ❌ No error testing, validation failures, or edge cases
 */

import { test, expect } from '../../fixtures/auth';
import { ProvidersPage } from '../../page-objects/providers-page';
import { TestDataGenerators, ProviderData } from '../../utils/data-generators';
import { TestHelpers, createTestHelpers } from '../../utils/test-helpers';

test.describe('Provider Management Happy Path', () => {
  let providersPage: ProvidersPage;
  let testHelpers: TestHelpers;
  let createdProviderIds: string[] = [];

  test.beforeEach(async ({ authenticatedPage }) => {
    console.log('🚀 Setting up provider management test');
    
    providersPage = new ProvidersPage(authenticatedPage);
    testHelpers = await createTestHelpers(authenticatedPage);
    
    // Navigate to providers page
    await authenticatedPage.goto('/providers');
    await providersPage.waitForPageLoad();
    console.log('✅ Navigated to providers page');
  });

  test.afterEach(async () => {
    console.log('🧹 Starting cleanup process');
    
    // Clean up created providers
    for (const providerId of createdProviderIds) {
      try {
        await testHelpers.cleanupTestProvider(providerId);
        console.log(`✅ Cleaned up provider: ${providerId}`);
      } catch {
        console.log(`ℹ️ Provider already cleaned up: ${providerId}`);
      }
    }
    
    // Reset arrays
    createdProviderIds = [];
    console.log('✅ Cleanup completed');
  });

  test('should display providers page layout correctly', async () => {
    console.log('🔍 Testing providers page layout');
    
    // Verify page title is visible
    const pageTitle = providersPage.pageInstance.locator('h1');
    await expect(pageTitle).toBeVisible();
    console.log('✅ Page title is visible');
    
    // Verify add provider button is visible and enabled
    const addButton = providersPage.pageInstance.locator('button:has-text("Add Provider"), button:has-text("Create Provider")');
    if (await addButton.count() > 0) {
      await expect(addButton.first()).toBeVisible();
      await expect(addButton.first()).toBeEnabled();
      console.log('✅ Add provider button is visible and enabled');
    }
    
    // Verify search input is visible and interactive
    const searchInput = providersPage.pageInstance.locator('input[placeholder*="Search"], input[type="search"]');
    if (await searchInput.count() > 0) {
      await expect(searchInput.first()).toBeVisible();
      await expect(searchInput.first()).toBeEnabled();
      console.log('✅ Search input is visible and enabled');
    }
    
    // Verify main content area is present
    const mainContent = providersPage.pageInstance.locator('main, [role="main"], .main-content, .container');
    if (await mainContent.count() > 0) {
      await expect(mainContent.first()).toBeVisible();
      console.log('✅ Main content area is visible');
    }
  });

  test('should create a new provider successfully', async () => {
    console.log('🆕 Testing provider creation');
    
    // Generate test provider data
    const providerData = TestDataGenerators.generateProvider({
      name: `Test Provider ${Date.now()}`,
      description: 'Test provider for E2E testing',
      website: 'https://example.com'
    });
    console.log(`📝 Generated provider data: ${providerData.name}`);
    
    // Click add provider button
    await providersPage.clickAddProvider();
    console.log('✅ Clicked add provider button');
    
    // Verify navigation to create page
    await expect(providersPage.pageInstance).toHaveURL(/.*providers\/create/);
    console.log('✅ Navigated to create provider page');
    
    // Fill out the provider form
    await providersPage.fillProviderForm(providerData);
    console.log('✅ Filled provider form');
    
    // Submit the form
    await providersPage.submitProviderForm();
    console.log('✅ Submitted provider form');
    
    // Wait for successful navigation back to providers list
    await providersPage.pageInstance.waitForURL('**/providers', { timeout: 15000 });
    console.log('✅ Successfully navigated back to providers list');
    
    // Verify provider appears in the list
    await providersPage.verifyProviderInList({
      name: providerData.name,
      description: providerData.description
    });
    console.log(`✅ Verified provider appears in list: ${providerData.name}`);
    
    // Get the created provider ID for cleanup
    const providerIdResult = await testHelpers.getProviderIdByName(providerData.name);
    if (providerIdResult) {
      createdProviderIds.push(providerIdResult);
      console.log(`📝 Added provider to cleanup list: ${providerIdResult}`);
    }
  });

  test('should edit an existing provider successfully', async () => {
    console.log('✏️ Testing provider editing');
    
    // First create a provider to edit
    const originalData = TestDataGenerators.generateProvider({
      name: `Original Provider ${Date.now()}`,
      description: 'Original description',
      website: 'https://original.com'
    });
    
    const providerResult = await testHelpers.createTestProvider(originalData);
    createdProviderIds.push(providerResult.id);
    console.log(`✅ Created provider to edit: ${originalData.name}`);
    
    // Refresh page to see the new provider
    await providersPage.pageInstance.reload();
    await providersPage.waitForPageLoad();
    
    // Edit the provider
    await providersPage.editProviderByName(originalData.name);
    console.log('✅ Opened provider for editing');
    
    // Update the provider data
    const updatedData = {
      name: `${originalData.name} - Edited`,
      description: 'Updated description',
      website: 'https://updated.com'
    };
    
    await providersPage.updateProviderForm(updatedData);
    console.log('✅ Updated provider form');
    
    // Submit the changes
    await providersPage.submitProviderForm();
    console.log('✅ Submitted provider changes');
    
    // Wait for navigation back to list
    await providersPage.pageInstance.waitForURL('**/providers', { timeout: 15000 });
    console.log('✅ Navigated back to providers list');
    
    // Verify the updated provider appears in the list
    await providersPage.verifyProviderInList({
      name: updatedData.name,
      description: updatedData.description
    });
    console.log(`✅ Verified updated provider in list: ${updatedData.name}`);
  });

  test('should delete a provider successfully', async () => {
    console.log('🗑️ Testing provider deletion');
    
    // First create a provider to delete
    const providerData = TestDataGenerators.generateProvider({
      name: `Provider to Delete ${Date.now()}`,
      description: 'This provider will be deleted',
      website: 'https://delete.com'
    });
    
    const providerResult = await testHelpers.createTestProvider(providerData);
    console.log(`✅ Created provider to delete: ${providerData.name}`);
    
    // Refresh page to see the new provider
    await providersPage.pageInstance.reload();
    await providersPage.waitForPageLoad();
    
    // Delete the provider
    await providersPage.deleteProviderByName(providerData.name);
    console.log('✅ Initiated provider deletion');
    
    // Verify provider is no longer in the list
    await providersPage.verifyProviderNotInList(providerData.name);
    console.log(`✅ Verified provider removed from list: ${providerData.name}`);
    
    // Note: No need to add to cleanup list since it's already deleted
  });

  test('should view provider details successfully', async () => {
    console.log('👁️ Testing provider detail view');
    
    // Create a provider with detailed information
    const detailedData = TestDataGenerators.generateProvider({
      name: `Detailed Provider ${Date.now()}`,
      description: 'Provider with comprehensive details',
      website: 'https://detailed-provider.com'
    });
    
    const providerResult = await testHelpers.createTestProvider(detailedData);
    createdProviderIds.push(providerResult.id);
    console.log(`✅ Created detailed provider: ${detailedData.name}`);
    
    // Refresh page to see the new provider
    await providersPage.pageInstance.reload();
    await providersPage.waitForPageLoad();
    
    // Click on provider to view details
    await providersPage.viewProviderDetails(detailedData.name);
    console.log('✅ Opened provider details');
    
    // Verify we're on the detail page
    await expect(providersPage.pageInstance).toHaveURL(/.*providers\/[^\/]+/);
    console.log('✅ Navigated to provider detail page');
    
    // Verify provider details are displayed
    await providersPage.verifyProviderDetails({
      name: detailedData.name,
      description: detailedData.description,
      website: detailedData.website
    });
    console.log('✅ Verified provider details are displayed correctly');
  });

  test('should search for providers successfully', async () => {
    console.log('🔍 Testing provider search functionality');
    
    // Create multiple providers with different names
    const providers = [
      TestDataGenerators.generateProvider({
        name: `Netflix Provider ${Date.now()}`,
        description: 'Streaming service provider'
      }),
      TestDataGenerators.generateProvider({
        name: `Spotify Provider ${Date.now()}`,
        description: 'Music streaming provider'
      }),
      TestDataGenerators.generateProvider({
        name: `Adobe Provider ${Date.now()}`,
        description: 'Creative software provider'
      })
    ];
    
    // Create all providers
    for (const provider of providers) {
      const providerResult = await testHelpers.createTestProvider(provider);
      createdProviderIds.push(providerResult.id);
      console.log(`✅ Created provider: ${provider.name}`);
    }
    
    // Refresh page to see all providers
    await providersPage.pageInstance.reload();
    await providersPage.waitForPageLoad();
    
    // Search for "Netflix" provider
    await providersPage.searchProviders('Netflix');
    console.log('✅ Performed search for "Netflix"');
    
    // Verify search results contain Netflix provider
    const searchResults = await providersPage.getAllProviderNames();
    const netflixProvider = providers.find(p => p.name.includes('Netflix'));
    expect(searchResults).toContain(netflixProvider!.name);
    console.log('✅ Verified Netflix provider appears in search results');
    
    // Clear search to show all providers again
    await providersPage.clearSearch();
    console.log('✅ Cleared search');
    
    // Verify all providers are visible again
    const allResults = await providersPage.getAllProviderNames();
    for (const provider of providers) {
      expect(allResults).toContain(provider.name);
    }
    console.log('✅ Verified all providers visible after clearing search');
  });

  test('should navigate between provider pages successfully', async () => {
    console.log('🧭 Testing navigation between provider pages');
    
    // Start on providers list page
    await expect(providersPage.pageInstance).toHaveURL(/.*providers(?!\/)/);
    console.log('✅ Confirmed on providers list page');
    
    // Navigate to create provider page
    await providersPage.clickAddProvider();
    await expect(providersPage.pageInstance).toHaveURL(/.*providers\/create/);
    console.log('✅ Successfully navigated to create provider page');
    
    // Navigate back to providers list
    const backButton = providersPage.pageInstance.locator('button:has-text("Back"), a:has-text("Back")');
    if (await backButton.count() > 0) {
      await backButton.first().click();
    } else {
      // Alternative: use browser back button
      await providersPage.pageInstance.goBack();
    }
    
    await expect(providersPage.pageInstance).toHaveURL(/.*providers(?!\/create)/);
    console.log('✅ Successfully navigated back to providers list');
    
    // Verify page is functional after navigation
    const addButton = providersPage.pageInstance.locator('button:has-text("Add Provider"), button:has-text("Create Provider")');
    if (await addButton.count() > 0) {
      await expect(addButton.first()).toBeVisible();
      await expect(addButton.first()).toBeEnabled();
      console.log('✅ Verified page is functional after navigation');
    }
  });

  test('should complete full provider lifecycle successfully', async () => {
    console.log('🔄 Testing complete provider lifecycle (create → edit → view → delete)');
    
    // Step 1: Create provider
    const providerData = TestDataGenerators.generateProvider({
      name: `Lifecycle Provider ${Date.now()}`,
      description: 'Provider for lifecycle testing',
      website: 'https://lifecycle.com'
    });
    
    await providersPage.clickAddProvider();
    await providersPage.fillProviderForm(providerData);
    await providersPage.submitProviderForm();
    await providersPage.pageInstance.waitForURL('**/providers', { timeout: 15000 });
    console.log('✅ Step 1: Created provider successfully');
    
    // Step 2: Edit provider
    await providersPage.editProviderByName(providerData.name);
    const updatedName = `${providerData.name} - Updated`;
    await providersPage.updateProviderForm({ 
      name: updatedName,
      description: 'Updated description'
    });
    await providersPage.submitProviderForm();
    await providersPage.pageInstance.waitForURL('**/providers', { timeout: 15000 });
    console.log('✅ Step 2: Edited provider successfully');
    
    // Step 3: View provider details
    await providersPage.viewProviderDetails(updatedName);
    await expect(providersPage.pageInstance).toHaveURL(/.*providers\/[^\/]+/);
    await providersPage.verifyProviderDetails({
      name: updatedName,
      description: 'Updated description'
    });
    console.log('✅ Step 3: Viewed provider details successfully');
    
    // Navigate back to list
    await providersPage.pageInstance.goBack();
    await providersPage.waitForPageLoad();
    
    // Step 4: Delete provider
    await providersPage.deleteProviderByName(updatedName);
    await providersPage.verifyProviderNotInList(updatedName);
    console.log('✅ Step 4: Deleted provider successfully');
    
    console.log('🎉 Complete provider lifecycle test completed successfully');
  });
});