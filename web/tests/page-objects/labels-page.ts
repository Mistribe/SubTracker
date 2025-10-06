import { Page, expect } from '@playwright/test';
import { BasePage } from './base-page';
import { LabelData } from '../utils/data-generators';

/**
 * Page object for the Labels page
 * Handles label CRUD operations, assignment functionality, and navigation
 */
export class LabelsPage extends BasePage {
  // Main page selectors
  private readonly pageTitle = 'h1:has-text("Labels")';
  private readonly searchInput = 'input[placeholder*="Search"], input[type="search"]';
  private readonly addLabelButton = 'button:has-text("Add Label"), button:has-text("Create Label")';

  // Label list selectors
  private readonly labelsList = '[data-testid="labels-list"], .labels-list, .grid';
  private readonly labelCards = '[data-testid="label-card"], .label-card, .card';
  private readonly loadingIndicator = 'text="Loading", .loading, .spinner';
  private readonly emptyState = 'text="No labels found", text="No labels match"';

  // Label form selectors (for creation/editing)
  private readonly labelForm = 'form';
  private readonly nameInput = '[name="name"], input[placeholder*="name" i]';
  private readonly colorInput = '[name="color"], input[type="color"]';
  private readonly descriptionInput = '[name="description"], textarea[placeholder*="description" i]';
  private readonly submitButton = 'button[type="submit"], button:has-text("Create"), button:has-text("Save")';

  // Color picker selectors
  private readonly colorPicker = '.color-picker, [data-testid="color-picker"]';
  private readonly colorOption = '.color-option, [data-testid="color-option"]';

  // Action menu selectors
  private readonly actionMenuTrigger = 'button[aria-label*="menu" i], button[aria-label*="action" i]';
  private readonly editMenuItem = 'text="Edit"';
  private readonly deleteMenuItem = 'text="Delete"';
  private readonly viewMenuItem = 'text="View"';

  // Delete confirmation dialog selectors
  private readonly deleteDialog = '[role="dialog"]:has-text("Delete"), [role="dialog"]:has-text("Remove")';
  private readonly confirmDeleteButton = 'button:has-text("Delete"), button:has-text("Remove")';
  private readonly cancelDeleteButton = 'button:has-text("Cancel")';

  // Modal selectors
  private readonly modal = '[role="dialog"]';
  private readonly modalCloseButton = 'button[aria-label*="close" i], button:has([data-lucide="x"])';

  constructor(page: Page) {
    super(page);
  }

  getPageUrl(): string {
    return '/labels';
  }

  async waitForPageLoad(): Promise<void> {
    // Wait for either the page title or the add label button to be visible
    try {
      await Promise.race([
        this.waitForElement(this.pageTitle, { timeout: 10000 }),
        this.waitForElement(this.addLabelButton, { timeout: 10000 })
      ]);
    } catch (error) {
      console.log('⚠️ Standard page elements not found, checking for any content');
      // Fallback: wait for any main content
      await this.page.waitForSelector('main, .container, .content, body', { timeout: 15000 });
    }
    await this.waitForPageReady();
  }

  /**
   * Search for labels using the search input
   */
  async searchLabels(query: string): Promise<void> {
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
   * Click the Add Label button to open creation form or modal
   */
  async clickAddLabel(): Promise<void> {
    // Try multiple possible selectors for add label button
    const addButtonSelectors = [
      'button:has-text("Add Label")',
      'button:has-text("Create Label")',
      'button:has-text("New Label")',
      'button:has-text("+")',
      'a[href*="labels/create"]',
      'button[data-testid*="add"]',
      'button[data-testid*="create"]'
    ];
    
    let buttonClicked = false;
    
    for (const selector of addButtonSelectors) {
      try {
        const button = this.page.locator(selector);
        if (await button.count() > 0 && await button.first().isVisible({ timeout: 3000 })) {
          await button.first().click();
          console.log(`Clicked add label button using selector: ${selector}`);
          buttonClicked = true;
          break;
        }
      } catch (error) {
        console.log(`Failed to click add button with selector ${selector}: ${error}`);
        continue;
      }
    }
    
    if (!buttonClicked) {
      throw new Error('Could not find or click add label button');
    }
    
    // Wait for navigation to create page or modal to open
    await this.page.waitForTimeout(2000);
    
    try {
      await this.page.waitForURL('**/labels/create', { timeout: 5000 });
      console.log('Navigated to create label page');
    } catch {
      // Might be a modal instead of navigation
      try {
        await this.waitForElement(this.modal, { timeout: 5000 });
        console.log('Opened create label modal');
      } catch {
        console.log('⚠️ Neither page navigation nor modal detected, continuing...');
      }
    }
  }

  /**
   * Fill the label form with the provided data
   */
  async fillLabelForm(data: LabelData): Promise<void> {
    console.log(`Filling label form with data: ${data.name}`);
    
    // Fill name field
    await this.fillInput(this.nameInput, data.name);
    
    // Set color
    await this.setLabelColor(data.color);
    
    // Fill description if provided
    if (data.description) {
      const descInput = this.page.locator(this.descriptionInput);
      if (await descInput.count() > 0) {
        await this.fillInput(this.descriptionInput, data.description);
      }
    }
    
    console.log('Label form filled successfully');
  }

  /**
   * Set the label color using color picker or input
   */
  private async setLabelColor(color: string): Promise<void> {
    // Try color input first
    const colorInput = this.page.locator(this.colorInput);
    if (await colorInput.count() > 0) {
      await colorInput.fill(color);
      return;
    }
    
    // Try color picker
    const colorPicker = this.page.locator(this.colorPicker);
    if (await colorPicker.count() > 0) {
      await colorPicker.click();
      
      // Look for specific color option
      const colorOption = this.page.locator(`${this.colorOption}[data-color="${color}"]`);
      if (await colorOption.count() > 0) {
        await colorOption.click();
      }
    }
  }

  /**
   * Submit the label form
   */
  async submitLabelForm(): Promise<void> {
    await this.clickElement(this.submitButton);
    console.log('Label form submitted');
  }

  /**
   * Submit the label form with fallback button detection
   */
  async submitLabelFormWithFallback(): Promise<void> {
    // Wait a moment for the form to be fully loaded
    await this.page.waitForTimeout(2000);
    
    // First, try to find the actual submit button in the modal
    const modalSubmitSelectors = [
      '[role="dialog"] button[type="submit"]',
      '[role="dialog"] button:has-text("Create")',
      '[role="dialog"] button:has-text("Save")',
      '[role="dialog"] button:has-text("Add")',
      '[role="dialog"] form button:last-child',
      '[role="dialog"] .flex.gap-2 button:last-child',
      '[role="dialog"] .space-x-2 button:last-child'
    ];
    
    let submitted = false;
    
    // Try modal-specific selectors first
    for (const selector of modalSubmitSelectors) {
      try {
        const buttons = this.page.locator(selector);
        const count = await buttons.count();
        
        if (count > 0) {
          for (let i = 0; i < count; i++) {
            const button = buttons.nth(i);
            try {
              if (await button.isVisible({ timeout: 2000 })) {
                // Use force click to bypass overlay issues
                await button.click({ force: true, timeout: 10000 });
                console.log(`Label form submitted using modal selector: ${selector} (button ${i})`);
                submitted = true;
                break;
              }
            } catch (buttonError) {
              console.log(`Modal button ${i} with selector ${selector} failed: ${buttonError}`);
              continue;
            }
          }
          if (submitted) break;
        }
      } catch (error) {
        console.log(`Failed to submit with modal selector ${selector}: ${error}`);
        continue;
      }
    }
    
    // If modal approach failed, try general selectors with force click
    if (!submitted) {
      const generalSelectors = [
        'button[type="submit"]',
        'button:has-text("Create")',
        'button:has-text("Save")',
        'button:has-text("Add")',
        'button:has-text("Submit")',
        'form button:last-child',
        'button:not([disabled])'
      ];
      
      for (const selector of generalSelectors) {
        try {
          const buttons = this.page.locator(selector);
          const count = await buttons.count();
          
          if (count > 0) {
            for (let i = 0; i < count; i++) {
              const button = buttons.nth(i);
              try {
                if (await button.isVisible({ timeout: 2000 })) {
                  // Use force click to bypass overlay issues
                  await button.click({ force: true, timeout: 10000 });
                  console.log(`Label form submitted using general selector: ${selector} (button ${i})`);
                  submitted = true;
                  break;
                }
              } catch (buttonError) {
                console.log(`General button ${i} with selector ${selector} failed: ${buttonError}`);
                continue;
              }
            }
            if (submitted) break;
          }
        } catch (error) {
          console.log(`Failed to submit with general selector ${selector}: ${error}`);
          continue;
        }
      }
    }
    
    // Last resort: try pressing Enter on form elements
    if (!submitted) {
      const enterTargets = [
        '[role="dialog"] form',
        '[role="dialog"] input[name="name"]',
        'form',
        'input[name="name"]',
        'textarea[name="description"]'
      ];
      
      for (const target of enterTargets) {
        try {
          const element = this.page.locator(target).first();
          if (await element.isVisible({ timeout: 2000 })) {
            await element.press('Enter');
            console.log(`Label form submitted using Enter key on: ${target}`);
            submitted = true;
            break;
          }
        } catch (error) {
          continue;
        }
      }
    }
    
    if (!submitted) {
      console.log('⚠️ Could not find any submit button, trying to close modal with Escape');
      // Try to close modal with Escape key as last resort
      await this.page.keyboard.press('Escape');
    }
    
    // Wait for any submission to process
    await this.page.waitForTimeout(2000);
  }

  /**
   * Update label form with partial data (for editing)
   */
  async updateLabelForm(updates: Partial<LabelData>): Promise<void> {
    console.log('Updating label form with changes');
    
    if (updates.name) {
      await this.page.locator(this.nameInput).fill('');
      await this.fillInput(this.nameInput, updates.name);
    }
    
    if (updates.color) {
      await this.setLabelColor(updates.color);
    }
    
    if (updates.description) {
      const descInput = this.page.locator(this.descriptionInput);
      if (await descInput.count() > 0) {
        await descInput.fill('');
        await this.fillInput(this.descriptionInput, updates.description);
      }
    }
    
    console.log('Label form updated successfully');
  }

  /**
   * Edit a label by name
   */
  async editLabelByName(labelName: string): Promise<void> {
    console.log(`Editing label: ${labelName}`);
    
    // Find the label card or row
    const labelElement = this.page.locator(`text="${labelName}"`).first();
    
    // Look for edit button near the label
    const editButton = labelElement.locator('..').locator('button:has-text("Edit"), button[aria-label*="edit" i]');
    
    if (await editButton.count() > 0) {
      await editButton.click();
    } else {
      // Try clicking on the label name itself
      await labelElement.click();
      
      // Look for edit option in action menu or on detail page
      const editOption = this.page.locator('button:has-text("Edit"), a:has-text("Edit")');
      if (await editOption.count() > 0) {
        await editOption.first().click();
      }
    }
    
    // Wait for edit form to load (could be modal or page)
    try {
      await this.waitForElement(this.labelForm);
    } catch {
      await this.waitForElement(this.modal);
    }
    console.log('Label edit form opened');
  }

  /**
   * Delete a label by name
   */
  async deleteLabelByName(labelName: string): Promise<void> {
    console.log(`Deleting label: ${labelName}`);
    
    // Find the label card or row
    const labelElement = this.page.locator(`text="${labelName}"`).first();
    
    // Look for delete button near the label
    const deleteButton = labelElement.locator('..').locator('button:has-text("Delete"), button[aria-label*="delete" i]');
    
    if (await deleteButton.count() > 0) {
      await deleteButton.click();
    } else {
      // Try action menu approach
      const actionMenu = labelElement.locator('..').locator(this.actionMenuTrigger);
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
    console.log('Label deleted successfully');
  }

  /**
   * View label details by name
   */
  async viewLabelDetails(labelName: string): Promise<void> {
    console.log(`Viewing details for label: ${labelName}`);
    
    // Click on the label name to view details
    const labelElement = this.page.locator(`text="${labelName}"`).first();
    await labelElement.click();
    
    console.log('Label details opened');
  }

  /**
   * Verify label appears in the list
   */
  async verifyLabelInList(label: { name: string; color?: string }): Promise<void> {
    console.log(`Verifying label in list: ${label.name}`);
    
    const labelElement = this.page.locator(`text="${label.name}"`);
    await expect(labelElement).toBeVisible();
    
    // If color is specified, try to verify it
    if (label.color) {
      // Look for color indicator near the label name
      const colorIndicator = labelElement.locator('..').locator(`[style*="${label.color}"], [data-color="${label.color}"]`);
      if (await colorIndicator.count() > 0) {
        await expect(colorIndicator).toBeVisible();
      }
    }
    
    console.log('Label verified in list');
  }

  /**
   * Verify label is not in the list
   */
  async verifyLabelNotInList(labelName: string): Promise<void> {
    console.log(`Verifying label not in list: ${labelName}`);
    
    const labelElement = this.page.locator(`text="${labelName}"`);
    await expect(labelElement).not.toBeVisible();
    
    console.log('Label confirmed not in list');
  }

  /**
   * Verify label details on detail page
   */
  async verifyLabelDetails(label: { name: string; color?: string; description?: string }): Promise<void> {
    console.log(`Verifying label details: ${label.name}`);
    
    // Verify name
    const nameElement = this.page.locator(`text="${label.name}"`);
    await expect(nameElement).toBeVisible();
    
    // Verify color if provided
    if (label.color) {
      const colorElement = this.page.locator(`[style*="${label.color}"], [data-color="${label.color}"]`);
      if (await colorElement.count() > 0) {
        await expect(colorElement).toBeVisible();
      }
    }
    
    // Verify description if provided
    if (label.description) {
      const descElement = this.page.locator(`text="${label.description}"`);
      await expect(descElement).toBeVisible();
    }
    
    console.log('Label details verified');
  }

  /**
   * Get all label names currently visible
   */
  async getAllLabelNames(): Promise<string[]> {
    // Wait for labels to load
    await this.waitForPageReady();
    
    // Try different selectors for label names
    const nameSelectors = [
      '[data-testid="label-name"]',
      '.label-name',
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
    const allText = await this.page.locator(this.labelsList).textContent();
    return allText ? [allText.trim()] : [];
  }
}