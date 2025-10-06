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
  private readonly searchInput = 'input[placeholder*="Search"], input[type="search"]';
  private readonly addProviderButton = 'button:has-text("Add Provider"), button:has-text("Create Provider")';

  // Provider list selectors
  private readonly providersList = '[data-testid="providers-list"], .providers-list, .grid';
  private readonly providerCards = '[data-testid="provider-card"], .provider-card, .card';
  private readonly loadingIndicator = 'text="Loading", .loading, .spinner';
  private readonly emptyState = 'text="No providers found", text="No providers match"';

  // Provider form selectors (for creation/editing)
  private readonly providerForm = 'form';
  private readonly nameInput = '[name="name"], input[placeholder*="name" i]';
  private readonly descriptionInput = '[name="description"], textarea[placeholder*="description" i]';
  private readonly websiteInput = '[name="website"], input[placeholder*="website" i], input[type="url"]';
  private readonly submitButton = 'button[type="submit"], button:has-text("Create"), button:has-text("Save")';

  // Action menu selectors
  private readonly actionMenuTrigger = 'button[aria-label*="menu" i], button[aria-label*="action" i]';
  private readonly editMenuItem = 'text="Edit"';
  private readonly deleteMenuItem = 'text="Delete"';
  private readonly viewMenuItem = 'text="View"';

  // Delete confirmation dialog selectors
  private readonly deleteDialog = '[role="dialog"]:has-text("Delete"), [role="dialog"]:has-text("Remove")';
  private readonly confirmDeleteButton = 'button:has-text("Delete"), button:has-text("Remove")';
  private readonly cancelDeleteButton = 'button:has-text("Cancel")';

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
   * Click the Add Provider button to navigate to creation form
   */
  async clickAddProvider(): Promise<void> {
    const addButton = this.page.locator(this.addProviderButton);
    if (await addButton.count() > 0) {
      await this.clickElement(this.addProviderButton);
      // Wait for navigation to create page or modal to open
      try {
        await this.page.waitForURL('**/providers/create', { timeout: 5000 });
      } catch {
        // Might be a modal instead of navigation
        await this.waitForElement('[role="dialog"]');
      }
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
    
    // Find the provider card or row
    const providerElement = this.page.locator(`text="${providerName}"`).first();
    
    // Look for edit button near the provider
    const editButton = providerElement.locator('..').locator('button:has-text("Edit"), button[aria-label*="edit" i]');
    
    if (await editButton.count() > 0) {
      await editButton.click();
    } else {
      // Try clicking on the provider name itself
      await providerElement.click();
      
      // Look for edit option in action menu or on detail page
      const editOption = this.page.locator('button:has-text("Edit"), a:has-text("Edit")');
      if (await editOption.count() > 0) {
        await editOption.first().click();
      }
    }
    
    // Wait for edit form to load
    await this.waitForElement(this.providerForm);
    console.log('Provider edit form opened');
  }

  /**
   * Delete a provider by name
   */
  async deleteProviderByName(providerName: string): Promise<void> {
    console.log(`Deleting provider: ${providerName}`);
    
    // Find the provider card or row
    const providerElement = this.page.locator(`text="${providerName}"`).first();
    
    // Look for delete button near the provider
    const deleteButton = providerElement.locator('..').locator('button:has-text("Delete"), button[aria-label*="delete" i]');
    
    if (await deleteButton.count() > 0) {
      await deleteButton.click();
    } else {
      // Try action menu approach
      const actionMenu = providerElement.locator('..').locator(this.actionMenuTrigger);
      if (await actionMenu.count() > 0) {
        await actionMenu.click();
        await this.clickElement(this.deleteMenuItem);
      }
    }
    
    // Confirm deletion
    await this.waitForElement(this.deleteDialog);
    await this.clickElement(this.confirmDeleteButton);
    
    // Wait for deletion to complete
    await this.page.waitForTimeout(2000);
    console.log('Provider deleted successfully');
  }

  /**
   * View provider details by name
   */
  async viewProviderDetails(providerName: string): Promise<void> {
    console.log(`Viewing details for provider: ${providerName}`);
    
    // Click on the provider name to view details
    const providerElement = this.page.locator(`text="${providerName}"`).first();
    await providerElement.click();
    
    console.log('Provider details opened');
  }

  /**
   * Verify provider appears in the list
   */
  async verifyProviderInList(provider: { name: string; description?: string }): Promise<void> {
    console.log(`Verifying provider in list: ${provider.name}`);
    
    const providerElement = this.page.locator(`text="${provider.name}"`);
    await expect(providerElement).toBeVisible();
    
    if (provider.description) {
      const descriptionElement = this.page.locator(`text="${provider.description}"`);
      await expect(descriptionElement).toBeVisible();
    }
    
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
    const nameElement = this.page.locator(`text="${provider.name}"`);
    await expect(nameElement).toBeVisible();
    
    // Verify description if provided
    if (provider.description) {
      const descElement = this.page.locator(`text="${provider.description}"`);
      await expect(descElement).toBeVisible();
    }
    
    // Verify website if provided
    if (provider.website) {
      const websiteElement = this.page.locator(`text="${provider.website}"`);
      await expect(websiteElement).toBeVisible();
    }
    
    console.log('Provider details verified');
  }

  /**
   * Get all provider names currently visible
   */
  async getAllProviderNames(): Promise<string[]> {
    // Wait for providers to load
    await this.waitForPageReady();
    
    // Try different selectors for provider names
    const nameSelectors = [
      '[data-testid="provider-name"]',
      '.provider-name',
      'h3',
      'h2',
      '.card h3',
      '.card h2'
    ];
    
    for (const selector of nameSelectors) {
      const elements = this.page.locator(selector);
      if (await elements.count() > 0) {
        const names = await elements.allTextContents();
        return names.filter(name => name.trim().length > 0);
      }
    }
    
    // Fallback: get all text content and filter
    const allText = await this.page.locator(this.providersList).textContent();
    return allText ? [allText.trim()] : [];
  }
}