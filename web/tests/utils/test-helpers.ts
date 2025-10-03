import { Page, expect } from '@playwright/test';
import { TestApiClient, ApiTestHelpers } from './api-client';
import { TestDataGenerators, TestIdGenerator } from './data-generators';

/**
 * Common test utilities and helper functions
 */
export class TestHelpers {
  /**
   * Wait for element to be visible with custom timeout
   */
  static async waitForElement(page: Page, selector: string, timeout = 30000): Promise<void> {
    await page.locator(selector).waitFor({ state: 'visible', timeout });
  }

  /**
   * Wait for element to disappear
   */
  static async waitForElementToDisappear(page: Page, selector: string, timeout = 30000): Promise<void> {
    await page.locator(selector).waitFor({ state: 'detached', timeout });
  }

  /**
   * Fill form field with proper waiting
   */
  static async fillFormField(page: Page, selector: string, value: string): Promise<void> {
    await this.waitForElement(page, selector);
    await page.locator(selector).clear();
    await page.locator(selector).fill(value);
  }

  /**
   * Click button with proper waiting
   */
  static async clickButton(page: Page, selector: string): Promise<void> {
    await this.waitForElement(page, selector);
    await page.locator(selector).click();
  }

  /**
   * Select option from dropdown
   */
  static async selectOption(page: Page, selector: string, value: string): Promise<void> {
    await this.waitForElement(page, selector);
    await page.locator(selector).selectOption(value);
  }

  /**
   * Upload file to input
   */
  static async uploadFile(page: Page, selector: string, filePath: string): Promise<void> {
    await this.waitForElement(page, selector);
    await page.locator(selector).setInputFiles(filePath);
  }

  /**
   * Wait for API response
   */
  static async waitForApiResponse(page: Page, urlPattern: string | RegExp, timeout = 30000): Promise<void> {
    await page.waitForResponse(urlPattern, { timeout });
  }

  /**
   * Wait for page to load completely
   */
  static async waitForPageLoad(page: Page): Promise<void> {
    await page.waitForLoadState('networkidle');
  }

  /**
   * Take screenshot for debugging
   */
  static async takeScreenshot(page: Page, name: string): Promise<void> {
    await page.screenshot({ 
      path: `test-results/${name}-${Date.now()}.png`,
      fullPage: true 
    });
  }

  /**
   * Scroll element into view
   */
  static async scrollToElement(page: Page, selector: string): Promise<void> {
    await page.locator(selector).scrollIntoViewIfNeeded();
  }

  /**
   * Get element text content
   */
  static async getElementText(page: Page, selector: string): Promise<string> {
    await this.waitForElement(page, selector);
    const text = await page.locator(selector).textContent();
    return text?.trim() || '';
  }

  /**
   * Get element attribute value
   */
  static async getElementAttribute(page: Page, selector: string, attribute: string): Promise<string | null> {
    await this.waitForElement(page, selector);
    return await page.locator(selector).getAttribute(attribute);
  }

  /**
   * Check if element is visible
   */
  static async isElementVisible(page: Page, selector: string, timeout = 5000): Promise<boolean> {
    try {
      await page.locator(selector).waitFor({ state: 'visible', timeout });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Check if element exists in DOM
   */
  static async isElementPresent(page: Page, selector: string, timeout = 5000): Promise<boolean> {
    try {
      await page.locator(selector).waitFor({ state: 'attached', timeout });
      return true;
    } catch {
      return false;
    }
  }

  /**
   * Wait for condition to be true
   */
  static async waitForCondition(
    condition: () => Promise<boolean> | boolean,
    timeout = 30000,
    interval = 1000
  ): Promise<void> {
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
      
      await new Promise(resolve => setTimeout(resolve, interval));
    }

    throw new Error(`Condition not met within ${timeout}ms`);
  }

  /**
   * Retry operation with exponential backoff
   */
  static async retryOperation<T>(
    operation: () => Promise<T>,
    maxRetries = 3,
    baseDelay = 1000
  ): Promise<T> {
    let lastError: Error | undefined;

    for (let attempt = 1; attempt <= maxRetries; attempt++) {
      try {
        return await operation();
      } catch (error) {
        lastError = error as Error;
        
        if (attempt === maxRetries) {
          break;
        }
        
        const delay = baseDelay * Math.pow(2, attempt - 1);
        await new Promise(resolve => setTimeout(resolve, delay));
      }
    }

    throw new Error(`Operation failed after ${maxRetries} attempts. Last error: ${lastError?.message || 'Unknown error'}`);
  }

  /**
   * Generate unique test identifier
   */
  static generateTestId(prefix?: string): string {
    return TestIdGenerator.generate(prefix);
  }

  /**
   * Generate test email
   */
  static generateTestEmail(domain = 'test.example.com'): string {
    return TestIdGenerator.generateEmail(domain);
  }

  /**
   * Format currency for display
   */
  static formatCurrency(amount: number, currency = 'USD'): string {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency,
    }).format(amount);
  }

  /**
   * Format date for display
   */
  static formatDate(date: string | Date, format = 'MM/DD/YYYY'): string {
    const dateObj = typeof date === 'string' ? new Date(date) : date;
    
    switch (format) {
      case 'MM/DD/YYYY':
        return dateObj.toLocaleDateString('en-US');
      case 'YYYY-MM-DD':
        return dateObj.toISOString().split('T')[0];
      default:
        return dateObj.toLocaleDateString();
    }
  }

  /**
   * Generate random string
   */
  static generateRandomString(length = 10): string {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    let result = '';
    for (let i = 0; i < length; i++) {
      result += chars.charAt(Math.floor(Math.random() * chars.length));
    }
    return result;
  }

  /**
   * Generate random number within range
   */
  static generateRandomNumber(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min + 1)) + min;
  }

  /**
   * Sleep for specified duration
   */
  static async sleep(ms: number): Promise<void> {
    await new Promise(resolve => setTimeout(resolve, ms));
  }

  /**
   * Get current timestamp
   */
  static getCurrentTimestamp(): string {
    return new Date().toISOString();
  }

  /**
   * Compare two objects for equality (deep comparison)
   */
  static deepEqual(obj1: any, obj2: any): boolean {
    return JSON.stringify(obj1) === JSON.stringify(obj2);
  }

  /**
   * Clean up test data using API helpers
   */
  static async cleanupTestData(apiHelpers: ApiTestHelpers, testId: string): Promise<void> {
    try {
      await apiHelpers.cleanupTestData(testId);
    } catch (error) {
      console.warn(`Failed to cleanup test data for ${testId}:`, error);
    }
  }

  /**
   * Verify API health before running tests
   */
  static async verifyApiHealth(apiHelpers: ApiTestHelpers): Promise<void> {
    const isHealthy = await apiHelpers.verifyApiHealth();
    if (!isHealthy) {
      throw new Error('API is not healthy. Tests cannot proceed.');
    }
  }
}

/**
 * Assertion helpers for common test scenarios
 */
export class TestAssertions {
  /**
   * Assert element is visible
   */
  static async assertElementVisible(page: Page, selector: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toBeVisible();
  }

  /**
   * Assert element is hidden
   */
  static async assertElementHidden(page: Page, selector: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toBeHidden();
  }

  /**
   * Assert element contains text
   */
  static async assertElementContainsText(page: Page, selector: string, text: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toContainText(text);
  }

  /**
   * Assert element has exact text
   */
  static async assertElementHasText(page: Page, selector: string, text: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toHaveText(text);
  }

  /**
   * Assert element has attribute
   */
  static async assertElementHasAttribute(page: Page, selector: string, attribute: string, value: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toHaveAttribute(attribute, value);
  }

  /**
   * Assert element has class
   */
  static async assertElementHasClass(page: Page, selector: string, className: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toHaveClass(new RegExp(className));
  }

  /**
   * Assert element count
   */
  static async assertElementCount(page: Page, selector: string, count: number, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toHaveCount(count);
  }

  /**
   * Assert page URL
   */
  static async assertPageUrl(page: Page, url: string | RegExp, message?: string): Promise<void> {
    await expect(page, message).toHaveURL(url);
  }

  /**
   * Assert page title
   */
  static async assertPageTitle(page: Page, title: string | RegExp, message?: string): Promise<void> {
    await expect(page, message).toHaveTitle(title);
  }

  /**
   * Assert input value
   */
  static async assertInputValue(page: Page, selector: string, value: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toHaveValue(value);
  }

  /**
   * Assert checkbox is checked
   */
  static async assertCheckboxChecked(page: Page, selector: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toBeChecked();
  }

  /**
   * Assert checkbox is unchecked
   */
  static async assertCheckboxUnchecked(page: Page, selector: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).not.toBeChecked();
  }

  /**
   * Assert element is enabled
   */
  static async assertElementEnabled(page: Page, selector: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toBeEnabled();
  }

  /**
   * Assert element is disabled
   */
  static async assertElementDisabled(page: Page, selector: string, message?: string): Promise<void> {
    await expect(page.locator(selector), message).toBeDisabled();
  }

  /**
   * Assert API response success
   */
  static assertApiSuccess(response: { success: boolean; error?: string }, message?: string): void {
    expect(response.success, message || response.error).toBe(true);
  }

  /**
   * Assert API response failure
   */
  static assertApiFailure(response: { success: boolean; error?: string }, message?: string): void {
    expect(response.success, message).toBe(false);
  }

  /**
   * Assert array contains item
   */
  static assertArrayContains<T>(array: T[], item: T, message?: string): void {
    expect(array, message).toContain(item);
  }

  /**
   * Assert array length
   */
  static assertArrayLength<T>(array: T[], length: number, message?: string): void {
    expect(array, message).toHaveLength(length);
  }

  /**
   * Assert object has property
   */
  static assertObjectHasProperty(obj: any, property: string, message?: string): void {
    expect(obj, message).toHaveProperty(property);
  }

  /**
   * Assert values are equal
   */
  static assertEqual<T>(actual: T, expected: T, message?: string): void {
    if (message) {
      expect(actual, message).toBe(expected);
    } else {
      expect(actual).toBe(expected);
    }
  }

  /**
   * Assert values are not equal
   */
  static assertNotEqual<T>(actual: T, expected: T, message?: string): void {
    if (message) {
      expect(actual, message).not.toBe(expected);
    } else {
      expect(actual).not.toBe(expected);
    }
  }

  /**
   * Assert value is truthy
   */
  static assertTruthy(value: any, message?: string): void {
    expect(value, message).toBeTruthy();
  }

  /**
   * Assert value is falsy
   */
  static assertFalsy(value: any, message?: string): void {
    expect(value, message).toBeFalsy();
  }
}

/**
 * Environment utilities for test configuration
 */
export class TestEnvironment {
  /**
   * Get environment variable with default value
   */
  static getEnvVar(name: string, defaultValue?: string): string {
    const value = process.env[name];
    if (value === undefined && defaultValue === undefined) {
      throw new Error(`Environment variable ${name} is required but not set`);
    }
    return value || defaultValue!;
  }

  /**
   * Check if running in CI environment
   */
  static isCI(): boolean {
    return process.env.CI === 'true' || process.env.GITHUB_ACTIONS === 'true';
  }

  /**
   * Check if running in debug mode
   */
  static isDebugMode(): boolean {
    return process.env.DEBUG === 'true' || process.env.PLAYWRIGHT_DEBUG === 'true';
  }

  /**
   * Get base URL for tests
   */
  static getBaseUrl(): string {
    return this.getEnvVar('BASE_URL', 'http://localhost:3000');
  }

  /**
   * Get API base URL for tests
   */
  static getApiBaseUrl(): string {
    return this.getEnvVar('API_BASE_URL', 'http://localhost:3000');
  }

  /**
   * Get test timeout
   */
  static getTestTimeout(): number {
    return parseInt(this.getEnvVar('TEST_TIMEOUT', '30000'), 10);
  }

  /**
   * Get test user credentials
   */
  static getTestUserCredentials(): { email: string; password: string } {
    return {
      email: this.getEnvVar('TEST_USER_EMAIL'),
      password: this.getEnvVar('TEST_USER_PASSWORD'),
    };
  }

  /**
   * Get Clerk configuration
   */
  static getClerkConfig(): { publishableKey: string; secretKey: string } {
    return {
      publishableKey: this.getEnvVar('CLERK_PUBLISHABLE_KEY'),
      secretKey: this.getEnvVar('CLERK_SECRET_KEY'),
    };
  }
}