import { Page, expect } from '@playwright/test';
import { BasePage } from './base-page';
import { ProviderData } from '../utils/data-generators';

/**
 * Page object for the Providers page
 * Handles provider CRUD operations, search functionality, and navigation
 */
export class ProvidersPage extends BasePage {
  // Main page selectors
  private readonly pageTitle = 'h1:has-text("Providers")';
  private readonly searchInput = 'input[placeholder="Search..."]';
  private readonly addProviderButton = 'button:has-text("Add Provider"), button:has-text("Create Provider")';

  // Provider form selectors (for creation/editing)
  private readonly providerForm = '[role="dialog"] form, [role="dialog"] #provider-form';
  private readonly nameInput = '[role="dialog"] input[id="name"]';
  private readonly descriptionInput = '[role="dialog"] textarea[id="description"]';
  private readonly websiteInput = '[role="dialog"] input[id="url"]';
  private readonly submitButton = '[role="dialog"] button:has-text("Create Provider"), [role="dialog"] button:has-text("Update Provider"), [role="dialog"] button:has-text("Save")';

  constructor(page: Page) {
    super(page);
  }

  getPageUrl(): string {
    return '/providers';
  }

  async waitForPageLoad(): Promise<void> {
    // Wait for either the page title or the add provider button to be visible
    await Promise.race([
      this.waitForElement(this.pageTitle),
      this.waitForElement(this.addProviderButton)
    ]);
    await this.waitForPageReady();
  }

  /**
   * Search for providers using the search input
   */
  async searchProviders(query: string): Promise<void> {
    const searchInput = this.page.locator(this.searchInput);
    if (await searchInput.count() > 0) {
      await this.fillInput(this.searchInput, query);
      // Wait for search results to load
      await this.page.waitForTimeout(1000);
      await this.waitForPageReady();
    }
  }

  /**
   * Clear the search input
   */
  async clearSearch(): Promise<void> {
    const searchInput = this.page.locator(this.searchInput);
    if (await searchInput.count() > 0) {
      await this.fillInput(this.searchInput, '');
      await this.waitForPageReady();
    }
  }

  /**
   * Click the Add Provider button to open creation modal
   */
  async clickAddProvider(): Promise<void> {
    const addButton = this.page.locator(this.addProviderButton);
    if (await addButton.count() > 0) {
      await this.clickElement(this.addProviderButton);
      // Wait for modal to open
      await this.waitForElement('[role="dialog"]');
    }
  }

  /**
   * Fill the provider form with the provided data
   */
  async fillProviderForm(data: ProviderData): Promise<void> {
    console.log(`Filling provider form with data: ${data.name}`);
    
    // Fill name field
    await this.fillInput(this.nameInput, data.name);
    
    // Fill description if provided
    if (data.description) {
      const descInput = this.page.locator(this.descriptionInput);
      if (await descInput.count() > 0) {
        await this.fillInput(this.descriptionInput, data.description);
      }
    }
    
    // Fill website if provided
    if (data.website) {
      const websiteInput = this.page.locator(this.websiteInput);
      if (await websiteInput.count() > 0) {
        await this.fillInput(this.websiteInput, data.website);
      }
    }
    
    console.log('Provider form filled successfully');
  }

  /**
   * Submit the provider form
   */
  async submitProviderForm(): Promise<void> {
    await this.clickElement(this.submitButton);
    console.log('Provider form submitted');
  }

  /**
   * Update provider form with partial data (for editing)
   */
  async updateProviderForm(updates: Partial<ProviderData>): Promise<void> {
    console.log('Updating provider form with changes');
    
    if (updates.name) {
      await this.page.locator(this.nameInput).fill('');
      await this.fillInput(this.nameInput, updates.name);
    }
    
    if (updates.description) {
      const descInput = this.page.locator(this.descriptionInput);
      if (await descInput.count() > 0) {
        await descInput.fill('');
        await this.fillInput(this.descriptionInput, updates.description);
      }
    }
    
    if (updates.website) {
      const websiteInput = this.page.locator(this.websiteInput);
      if (await websiteInput.count() > 0) {
        await websiteInput.fill('');
        await this.fillInput(this.websiteInput, updates.website);
      }
    }
    
    console.log('Provider form updated successfully');
  }

  /**
   * Edit a provider by name
   */
  async editProviderByName(providerName: string): Promise<void> {
    console.log(`Editing provider: ${providerName}`);
    
    // Wait for page to be ready and providers to load
    await this.waitForPageReady();
    await this.page.waitForTimeout(2000); // Give time for providers to load
    
    // Use search to find the provider more reliably
    await this.searchProviders(providerName);
    await this.page.waitForTimeout(1000);
    
    // Find the provider card containing the name - use the actual Card component structure
    const providerCard = this.page.locator('[data-slot="card"]').filter({ hasText: providerName });
    await expect(providerCard).toBeVisible();
    
    // Look for the dropdown menu button (MoreVertical icon button)
    const dropdownButton = providerCard.locator('button:has(svg)').last(); // The MoreVertical button is typically the last button
    
    // Check if dropdown button exists (provider might not be editable)
    const dropdownCount = await dropdownButton.count();
    if (dropdownCount === 0) {
      console.log('No dropdown button found - provider might not be editable');
      throw new Error(`Provider "${providerName}" does not have an edit dropdown menu - it might not be editable`);
    }
    
    await expect(dropdownButton).toBeVisible();
    
    // Click the dropdown button
    await dropdownButton.click();
    console.log('Clicked dropdown button');
    
    // Wait for dropdown menu to appear and click Edit option
    await this.page.waitForTimeout(1000);
    const editOption = this.page.locator('text="Edit Provider"');
    
    // Check if edit option is available
    const editOptionCount = await editOption.count();
    if (editOptionCount === 0) {
      console.log('Edit Provider option not found in dropdown menu');
      // Try to find any dropdown menu items for debugging
      const allMenuItems = this.page.locator('[role="menuitem"]');
      const menuItemCount = await allMenuItems.count();
      console.log(`Found ${menuItemCount} menu items in dropdown`);
      
      for (let i = 0; i < menuItemCount; i++) {
        const itemText = await allMenuItems.nth(i).textContent();
        console.log(`Menu item ${i}: "${itemText}"`);
      }
      
      throw new Error(`Edit Provider option not found in dropdown menu for provider "${providerName}"`);
    }
    
    await expect(editOption).toBeVisible();
    await editOption.click();
    console.log('Clicked Edit Provider option');
    
    // Wait for edit dialog to open - this is where the backend issue occurs
    await this.page.waitForTimeout(2000);
    const editDialog = this.page.locator('[role="dialog"]:has-text("Edit Provider")');
    
    try {
      await expect(editDialog).toBeVisible({ timeout: 10000 });
      console.log('Provider edit dialog opened');
      
      // Wait for the form inside the dialog
      const formInDialog = editDialog.locator('form, #provider-form');
      await expect(formInDialog).toBeVisible({ timeout: 5000 });
      console.log('Provider edit form loaded');
    } catch (error) {
      console.log('⚠️ Edit dialog failed to open - this is a known backend issue');
      throw new Error(`Edit dialog failed to open for provider "${providerName}" - backend issue`);
    }
  }

  /**
   * Delete a provider by name
   */
  async deleteProviderByName(providerName: string): Promise<void> {
    console.log(`Deleting provider: ${providerName}`);
    
    // Wait for page to be ready and providers to load
    await this.waitForPageReady();
    await this.page.waitForTimeout(2000); // Give time for providers to load
    
    // Use search to find the provider more reliably
    await this.searchProviders(providerName);
    await this.page.waitForTimeout(1000);
    
    // Find the provider card containing the name - use the actual Card component structure
    const providerCard = this.page.locator('[data-slot="card"]').filter({ hasText: providerName });
    await expect(providerCard).toBeVisible();
    
    // Look for the dropdown menu button (MoreVertical icon button)
    const dropdownButton = providerCard.locator('button:has(svg)').last(); // The MoreVertical button is typically the last button
    
    // Check if dropdown button exists (provider might not be deletable)
    const dropdownCount = await dropdownButton.count();
    if (dropdownCount === 0) {
      console.log('No dropdown button found - provider might not be deletable');
      throw new Error(`Provider "${providerName}" does not have a delete dropdown menu - it might not be deletable`);
    }
    
    await expect(dropdownButton).toBeVisible();
    
    // Click the dropdown button
    await dropdownButton.click();
    
    // Wait for dropdown menu to appear and click Remove option
    await this.page.waitForTimeout(1000);
    const deleteOption = this.page.locator('text="Remove Provider"');
    await expect(deleteOption).toBeVisible();
    await deleteOption.click();
    
    // Wait for confirmation dialog to appear
    await this.page.waitForTimeout(1000);
    const confirmDialog = this.page.locator('[role="alertdialog"]');
    await expect(confirmDialog).toBeVisible();
    
    // Click the Delete button to confirm
    const deleteButton = confirmDialog.locator('button:has-text("Delete")');
    await expect(deleteButton).toBeVisible();
    await deleteButton.click();
    
    // Wait for deletion to complete
    await this.page.waitForTimeout(2000);
    console.log('Provider deleted successfully');
  }

  /**
   * View provider details by name
   */
  async viewProviderDetails(providerName: string): Promise<void> {
    console.log(`Viewing details for provider: ${providerName}`);
    
    // Use search to find the provider more reliably
    await this.searchProviders(providerName);
    await this.page.waitForTimeout(1000);
    
    // Find the provider card and click on it (but not on the dropdown button)
    const providerCard = this.page.locator('[data-slot="card"]').filter({ hasText: providerName });
    await expect(providerCard).toBeVisible();
    
    // Click on the card title area to navigate to details
    const cardTitle = providerCard.locator('[data-slot="card-title"]').first();
    await expect(cardTitle).toBeVisible();
    await cardTitle.click();
    
    console.log('Provider details opened');
  }

  /**
   * Verify provider appears in the list
   */
  async verifyProviderInList(provider: { name: string; description?: string }): Promise<void> {
    console.log(`Verifying provider in list: ${provider.name}`);
    
    // Wait for providers to load
    await this.waitForPageReady();
    await this.page.waitForTimeout(2000);
    
    // Try using search to find the provider more reliably
    console.log('Using search to find the provider...');
    await this.searchProviders(provider.name);
    
    // Now try to find the provider by name in the card structure
    const providerCard = this.page.locator('[data-slot="card"]').filter({ hasText: provider.name });
    
    // If not found with search, try refreshing and searching again
    if (await providerCard.count() === 0) {
      console.log('Provider not found with search, refreshing page...');
      await this.page.reload();
      await this.waitForPageLoad();
      await this.page.waitForTimeout(3000);
      
      // Search again after refresh
      await this.searchProviders(provider.name);
      await this.page.waitForTimeout(2000);
    }
    
    // Now verify the provider is visible
    await expect(providerCard).toBeVisible({ timeout: 15000 });
    
    if (provider.description) {
      const descriptionElement = providerCard.locator(`text="${provider.description}"`);
      await expect(descriptionElement).toBeVisible();
    }
    
    // Clear search after verification
    await this.clearSearch();
    
    console.log('Provider verified in list');
  }

  /**
   * Verify provider is not in the list
   */
  async verifyProviderNotInList(providerName: string): Promise<void> {
    console.log(`Verifying provider not in list: ${providerName}`);
    
    const providerElement = this.page.locator(`text="${providerName}"`);
    await expect(providerElement).not.toBeVisible();
    
    console.log('Provider confirmed not in list');
  }

  /**
   * Verify provider details on detail page
   */
  async verifyProviderDetails(provider: { name: string; description?: string; website?: string }): Promise<void> {
    console.log(`Verifying provider details: ${provider.name}`);
    
    // Verify name
    const nameElement = this.page.locator(`h2:has-text("${provider.name}")`);
    await expect(nameElement).toBeVisible();
    
    // Verify description if provided
    if (provider.description) {
      const descElement = this.page.locator(`p:has-text("${provider.description}")`);
      await expect(descElement).toBeVisible();
    }
    
    // Verify website if provided (it's shown as a "Website" button)
    if (provider.website) {
      const websiteButton = this.page.locator('button:has-text("Website")');
      await expect(websiteButton).toBeVisible();
    }
    
    console.log('Provider details verified');
  }

  /**
   * Get all provider names currently visible
   */
  async getAllProviderNames(): Promise<string[]> {
    // Wait for providers to load
    await this.waitForPageReady();
    await this.page.waitForTimeout(1000); // Give time for providers to render
    
    // Get provider names from card titles using the actual Card component structure
    const cardTitles = this.page.locator('[data-slot="card-title"]');
    const count = await cardTitles.count();
    
    const names: string[] = [];
    for (let i = 0; i < count; i++) {
      const text = await cardTitles.nth(i).textContent();
      if (text && text.trim().length > 0) {
        names.push(text.trim());
      }
    }
    
    return names;
  }
}