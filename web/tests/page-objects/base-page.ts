import { Page, Locator, expect } from '@playwright/test';

/**
 * Base page object class providing common functionality for all page objects
 * Implements reliable element waiting, interaction utilities, and error handling
 */
export abstract class BasePage {
  protected readonly page: Page;
  protected readonly timeout = 30000; // 30 seconds default timeout

  constructor(page: Page) {
    this.page = page;
  }

  /**
   * Get the page instance for direct access when needed
   */
  get pageInstance(): Page {
    return this.page;
  }

  /**
   * Abstract method to get the page URL - must be implemented by subclasses
   */
  abstract getPageUrl(): string;

  /**
   * Abstract method to wait for page load - must be implemented by subclasses
   */
  abstract waitForPageLoad(): Promise<void>;

  /**
   * Navigate to the page and wait for it to load
   */
  async navigateToPage(): Promise<void> {
    try {
      await this.page.goto(this.getPageUrl());
      await this.waitForPageLoad();
    } catch (error) {
      throw new Error(`Failed to navigate to ${this.getPageUrl()}: ${error}`);
    }
  }

  /**
   * Wait for an element to be visible and ready for interaction
   */
  async waitForElement(selector: string, options?: { timeout?: number }): Promise<Locator> {
    try {
      const element = this.page.locator(selector);
      await element.waitFor({ 
        state: 'visible', 
        timeout: options?.timeout || this.timeout 
      });
      return element;
    } catch (error) {
      throw new Error(`Element with selector "${selector}" not found or not visible: ${error}`);
    }
  }

  /**
   * Wait for an element to be attached to DOM (may not be visible)
   */
  async waitForElementAttached(selector: string, options?: { timeout?: number }): Promise<Locator> {
    try {
      const element = this.page.locator(selector);
      await element.waitFor({ 
        state: 'attached', 
        timeout: options?.timeout || this.timeout 
      });
      return element;
    } catch (error) {
      throw new Error(`Element with selector "${selector}" not found in DOM: ${error}`);
    }
  }

  /**
   * Click an element with retry logic and proper waiting
   */
  async clickElement(selector: string, options?: { timeout?: number; force?: boolean }): Promise<void> {
    try {
      const element = await this.waitForElement(selector, { timeout: options?.timeout });
      await element.click({ force: options?.force });
    } catch (error) {
      throw new Error(`Failed to click element "${selector}": ${error}`);
    }
  }

  /**
   * Fill an input field with proper waiting and clearing
   */
  async fillInput(selector: string, value: string, options?: { timeout?: number; clear?: boolean }): Promise<void> {
    try {
      const element = await this.waitForElement(selector, { timeout: options?.timeout });
      
      if (options?.clear !== false) {
        await element.clear();
      }
      
      await element.fill(value);
    } catch (error) {
      throw new Error(`Failed to fill input "${selector}" with value "${value}": ${error}`);
    }
  }

  /**
   * Get text content from an element
   */
  async getElementText(selector: string, options?: { timeout?: number }): Promise<string> {
    try {
      const element = await this.waitForElement(selector, { timeout: options?.timeout });
      const text = await element.textContent();
      return text?.trim() || '';
    } catch (error) {
      throw new Error(`Failed to get text from element "${selector}": ${error}`);
    }
  }

  /**
   * Get attribute value from an element
   */
  async getElementAttribute(selector: string, attribute: string, options?: { timeout?: number }): Promise<string | null> {
    try {
      const element = await this.waitForElement(selector, { timeout: options?.timeout });
      return await element.getAttribute(attribute);
    } catch (error) {
      throw new Error(`Failed to get attribute "${attribute}" from element "${selector}": ${error}`);
    }
  }

  /**
   * Check if an element is visible
   */
  async isElementVisible(selector: string, options?: { timeout?: number }): Promise<boolean> {
    try {
      const element = this.page.locator(selector);
      await element.waitFor({ 
        state: 'visible', 
        timeout: options?.timeout || 5000 
      });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Check if an element exists in DOM (may not be visible)
   */
  async isElementPresent(selector: string, options?: { timeout?: number }): Promise<boolean> {
    try {
      const element = this.page.locator(selector);
      await element.waitFor({ 
        state: 'attached', 
        timeout: options?.timeout || 5000 
      });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Wait for an element to disappear
   */
  async waitForElementToDisappear(selector: string, options?: { timeout?: number }): Promise<void> {
    try {
      const element = this.page.locator(selector);
      await element.waitFor({ 
        state: 'detached', 
        timeout: options?.timeout || this.timeout 
      });
    } catch (error) {
      throw new Error(`Element "${selector}" did not disappear: ${error}`);
    }
  }

  /**
   * Select option from dropdown/select element
   */
  async selectOption(selector: string, value: string | { label?: string; value?: string; index?: number }, options?: { timeout?: number }): Promise<void> {
    try {
      const element = await this.waitForElement(selector, { timeout: options?.timeout });
      
      if (typeof value === 'string') {
        await element.selectOption(value);
      } else if (value.value) {
        await element.selectOption({ value: value.value });
      } else if (value.label) {
        await element.selectOption({ label: value.label });
      } else if (value.index !== undefined) {
        await element.selectOption({ index: value.index });
      }
    } catch (error) {
      throw new Error(`Failed to select option in "${selector}": ${error}`);
    }
  }

  /**
   * Upload file to input element
   */
  async uploadFile(selector: string, filePath: string, options?: { timeout?: number }): Promise<void> {
    try {
      const element = await this.waitForElement(selector, { timeout: options?.timeout });
      await element.setInputFiles(filePath);
    } catch (error) {
      throw new Error(`Failed to upload file to "${selector}": ${error}`);
    }
  }

  /**
   * Scroll element into view
   */
  async scrollToElement(selector: string, options?: { timeout?: number }): Promise<void> {
    try {
      const element = await this.waitForElement(selector, { timeout: options?.timeout });
      await element.scrollIntoViewIfNeeded();
    } catch (error) {
      throw new Error(`Failed to scroll to element "${selector}": ${error}`);
    }
  }

  /**
   * Wait for page to load completely
   */
  async waitForPageReady(): Promise<void> {
    try {
      await this.page.waitForLoadState('networkidle');
    } catch (error) {
      throw new Error(`Page did not reach ready state: ${error}`);
    }
  }

  /**
   * Wait for API response
   */
  async waitForApiResponse(urlPattern: string | RegExp, options?: { timeout?: number }): Promise<void> {
    try {
      await this.page.waitForResponse(urlPattern, { 
        timeout: options?.timeout || this.timeout 
      });
    } catch (error) {
      throw new Error(`API response for "${urlPattern}" not received: ${error}`);
    }
  }

  /**
   * Get current page URL
   */
  getCurrentUrl(): string {
    return this.page.url();
  }

  /**
   * Get page title
   */
  async getPageTitle(): Promise<string> {
    return await this.page.title();
  }

  /**
   * Take screenshot for debugging
   */
  async takeScreenshot(name?: string): Promise<Buffer> {
    const screenshotName = name || `screenshot-${Date.now()}`;
    return await this.page.screenshot({ 
      path: `test-results/${screenshotName}.png`,
      fullPage: true 
    });
  }

  /**
   * Execute JavaScript in browser context
   */
  async executeScript<T>(script: string | Function, ...args: any[]): Promise<T> {
    try {
      return await this.page.evaluate(script, ...args);
    } catch (error) {
      throw new Error(`Failed to execute script: ${error}`);
    }
  }

  /**
   * Wait for specific condition to be true
   */
  async waitForCondition(
    condition: () => Promise<boolean> | boolean,
    options?: { timeout?: number; interval?: number }
  ): Promise<void> {
    const timeout = options?.timeout || this.timeout;
    const interval = options?.interval || 1000;
    const startTime = Date.now();

    while (Date.now() - startTime < timeout) {
      try {
        const result = await condition();
        if (result) {
          return;
        }
      } catch {
        // Continue waiting if condition throws
      }
      
      await this.page.waitForTimeout(interval);
    }

    throw new Error(`Condition not met within ${timeout}ms`);
  }

  /**
   * Handle common error scenarios with retry logic
   */
  async withRetry<T>(
    operation: () => Promise<T>,
    options?: { maxRetries?: number; delay?: number }
  ): Promise<T> {
    const maxRetries = options?.maxRetries || 3;
    const delay = options?.delay || 1000;
    let lastError: Error;

    for (let attempt = 1; attempt <= maxRetries; attempt++) {
      try {
        return await operation();
      } catch (error) {
        lastError = error as Error;
        
        if (attempt === maxRetries) {
          break;
        }
        
        await this.page.waitForTimeout(delay);
      }
    }

    throw new Error(`Operation failed after ${maxRetries} attempts. Last error: ${lastError.message}`);
  }
}