import { test } from '../../fixtures/auth';

test.describe('Debug Subscriptions Page', () => {
  test('should debug authenticated subscription pages with error checking', async ({ authenticatedPage }) => {
    // Listen for console errors and network failures
    const consoleErrors: string[] = [];
    const networkErrors: string[] = [];
    
    authenticatedPage.on('console', msg => {
      if (msg.type() === 'error') {
        consoleErrors.push(msg.text());
      }
    });
    
    authenticatedPage.on('response', response => {
      if (!response.ok()) {
        networkErrors.push(`${response.status()} ${response.url()}`);
      }
    });
    
    // Navigate to subscriptions page
    await authenticatedPage.goto('/subscriptions');
    
    // Wait a bit for any async operations
    await authenticatedPage.waitForTimeout(3000);
    
    // Take a screenshot
    await authenticatedPage.screenshot({ path: 'debug-subscriptions-auth.png', fullPage: true });
    
    console.log('Subscriptions page title:', await authenticatedPage.title());
    console.log('Subscriptions page URL:', authenticatedPage.url());
    console.log('Console errors:', consoleErrors);
    console.log('Network errors:', networkErrors);
    
    // Check for any text content
    const bodyText = await authenticatedPage.locator('body').textContent();
    console.log('Body text (first 500 chars):', bodyText?.substring(0, 500));
    
    // Look for loading states
    const loadingElements = await authenticatedPage.locator('.animate-spin').all();
    console.log('Found loading elements:', loadingElements.length);
    
    // Look for error states
    const errorElements = await authenticatedPage.locator('[role="alert"], .text-red-500').all();
    console.log('Found error elements:', errorElements.length);
    
    // Look for any h1 elements
    const h1Elements = await authenticatedPage.locator('h1').all();
    console.log('Found h1 elements:', h1Elements.length);
    
    for (let i = 0; i < h1Elements.length; i++) {
      const text = await h1Elements[i].textContent();
      console.log(`H1 ${i + 1}:`, text);
    }
    
    // Check for React components by looking for common patterns
    const reactElements = await authenticatedPage.locator('[data-testid], [class*="react"], [class*="component"]').all();
    console.log('Found potential React elements:', reactElements.length);
    
    // Try to navigate to create subscription
    try {
      await authenticatedPage.goto('/subscriptions/create');
      
      // Wait for page to load
      await authenticatedPage.waitForTimeout(3000);
      
      await authenticatedPage.screenshot({ path: 'debug-create-subscription.png', fullPage: true });
      
      console.log('Create subscription page title:', await authenticatedPage.title());
      console.log('Create subscription page URL:', authenticatedPage.url());
      console.log('Console errors after navigation:', consoleErrors);
      console.log('Network errors after navigation:', networkErrors);
      
      // Check for form elements
      const formElements = await authenticatedPage.locator('form').all();
      console.log('Found forms:', formElements.length);
      
      const inputElements = await authenticatedPage.locator('input').all();
      console.log('Found inputs:', inputElements.length);
      
      const selectElements = await authenticatedPage.locator('select').all();
      console.log('Found selects:', selectElements.length);
      
      // Check for any divs that might contain form content
      const divElements = await authenticatedPage.locator('div').all();
      console.log('Found divs:', divElements.length);
      
      // Check for step indicators (the form uses a wizard)
      const stepElements = await authenticatedPage.locator('text="Basic Information"').all();
      console.log('Found step elements:', stepElements.length);
      
      // Check for loading states on create page
      const createLoadingElements = await authenticatedPage.locator('.animate-spin').all();
      console.log('Found loading elements on create page:', createLoadingElements.length);
      
      // Check for error states on create page
      const createErrorElements = await authenticatedPage.locator('[role="alert"], .text-red-500').all();
      console.log('Found error elements on create page:', createErrorElements.length);
      
    } catch (error) {
      console.log('Error navigating to create subscription:', error);
    }
  });
});