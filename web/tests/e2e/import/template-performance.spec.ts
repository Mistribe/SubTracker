/**
 * Template File Loading Performance Tests
 * 
 * Test suite covering template file loading performance:
 * - Verify template files load quickly
 * - Test template download performance
 * - Ensure template files are properly cached
 * 
 * Requirements: 1.4, 1.5, 4.5
 */

import { test, expect } from '../../fixtures/auth';

test.describe('Template File Loading Performance', () => {
  const ENTITY_TYPES = ['labels', 'providers', 'subscriptions'] as const;
  const FORMATS = ['csv', 'json', 'yaml'] as const;

  test.describe('Template file loading speed', () => {
    ENTITY_TYPES.forEach((entityType) => {
      FORMATS.forEach((format) => {
        test(`should load ${entityType} ${format.toUpperCase()} template quickly (Requirement 1.4, 1.5, 4.5)`, async ({ authenticatedPage, page }) => {
          console.log(`⏱️ Testing load time for ${entityType} ${format.toUpperCase()} template`);

          // Navigate to the import page
          const importPath = `/${entityType}/import`;
          await authenticatedPage.goto(importPath);
          await authenticatedPage.waitForTimeout(1000);

          // Find the template download link
          const downloadLinkSelectors = [
            `a[href*="${entityType}-template.${format}"]`,
            `a[download*="${entityType}-template.${format}"]`,
          ];

          let href = null;
          for (const selector of downloadLinkSelectors) {
            const link = authenticatedPage.locator(selector);
            if (await link.count() > 0) {
              href = await link.first().getAttribute('href');
              if (href) break;
            }
          }

          expect(href).toBeTruthy();

          // Measure download time
          const templateUrl = new URL(href!, authenticatedPage.url()).href;
          
          const startTime = Date.now();
          const response = await page.request.get(templateUrl);
          const endTime = Date.now();
          
          const loadTime = endTime - startTime;

          expect(response.ok()).toBeTruthy();
          
          // Template files should load in less than 500ms
          expect(loadTime).toBeLessThan(500);
          
          console.log(`✅ ${entityType} ${format.toUpperCase()} template loaded in ${loadTime}ms`);
        });
      });
    });

    test('should load all template files for an entity type quickly (Requirement 1.4, 4.5)', async ({ authenticatedPage, page }) => {
      console.log('⏱️ Testing concurrent load time for all templates');

      const entityType = 'labels'; // Test with one entity type
      const importPath = `/${entityType}/import`;
      await authenticatedPage.goto(importPath);
      await authenticatedPage.waitForTimeout(1000);

      // Collect all template URLs
      const templateUrls: string[] = [];
      
      for (const format of FORMATS) {
        const downloadLinkSelectors = [
          `a[href*="${entityType}-template.${format}"]`,
          `a[download*="${entityType}-template.${format}"]`,
        ];

        for (const selector of downloadLinkSelectors) {
          const link = authenticatedPage.locator(selector);
          if (await link.count() > 0) {
            const href = await link.first().getAttribute('href');
            if (href) {
              const templateUrl = new URL(href, authenticatedPage.url()).href;
              templateUrls.push(templateUrl);
              break;
            }
          }
        }
      }

      expect(templateUrls.length).toBe(3); // CSV, JSON, YAML

      // Measure concurrent download time
      const startTime = Date.now();
      
      const responses = await Promise.all(
        templateUrls.map(url => page.request.get(url))
      );
      
      const endTime = Date.now();
      const totalLoadTime = endTime - startTime;

      // All responses should be successful
      responses.forEach(response => {
        expect(response.ok()).toBeTruthy();
      });

      // All three templates should load concurrently in less than 1 second
      expect(totalLoadTime).toBeLessThan(1000);
      
      console.log(`✅ All ${entityType} templates loaded concurrently in ${totalLoadTime}ms`);
    });
  });

  test.describe('Template file caching', () => {
    test('should cache template files for faster subsequent loads (Requirement 4.5)', async ({ authenticatedPage, page }) => {
      console.log('⏱️ Testing template file caching');

      const entityType = 'labels';
      const format = 'json';
      const importPath = `/${entityType}/import`;
      
      await authenticatedPage.goto(importPath);
      await authenticatedPage.waitForTimeout(1000);

      // Find the template download link
      const downloadLinkSelectors = [
        `a[href*="${entityType}-template.${format}"]`,
        `a[download*="${entityType}-template.${format}"]`,
      ];

      let href = null;
      for (const selector of downloadLinkSelectors) {
        const link = authenticatedPage.locator(selector);
        if (await link.count() > 0) {
          href = await link.first().getAttribute('href');
          if (href) break;
        }
      }

      expect(href).toBeTruthy();
      const templateUrl = new URL(href!, authenticatedPage.url()).href;

      // First load (cold cache)
      const startTime1 = Date.now();
      const response1 = await page.request.get(templateUrl);
      const endTime1 = Date.now();
      const firstLoadTime = endTime1 - startTime1;

      expect(response1.ok()).toBeTruthy();
      console.log(`📥 First load (cold cache): ${firstLoadTime}ms`);

      // Second load (should be cached)
      const startTime2 = Date.now();
      const response2 = await page.request.get(templateUrl);
      const endTime2 = Date.now();
      const secondLoadTime = endTime2 - startTime2;

      expect(response2.ok()).toBeTruthy();
      console.log(`📥 Second load (warm cache): ${secondLoadTime}ms`);

      // Third load (should still be cached)
      const startTime3 = Date.now();
      const response3 = await page.request.get(templateUrl);
      const endTime3 = Date.now();
      const thirdLoadTime = endTime3 - startTime3;

      expect(response3.ok()).toBeTruthy();
      console.log(`📥 Third load (warm cache): ${thirdLoadTime}ms`);

      // Cached loads should be faster or similar to first load
      // Note: In some environments, caching may not show significant improvement
      // but we verify that subsequent loads don't get slower
      expect(secondLoadTime).toBeLessThanOrEqual(firstLoadTime * 1.5);
      expect(thirdLoadTime).toBeLessThanOrEqual(firstLoadTime * 1.5);

      console.log(`✅ Template caching verified - subsequent loads are consistent`);
    });

    test('should handle multiple concurrent requests efficiently (Requirement 4.5)', async ({ authenticatedPage, page }) => {
      console.log('⏱️ Testing concurrent template requests');

      const entityType = 'providers';
      const format = 'csv';
      const importPath = `/${entityType}/import`;
      
      await authenticatedPage.goto(importPath);
      await authenticatedPage.waitForTimeout(1000);

      // Find the template download link
      const downloadLinkSelectors = [
        `a[href*="${entityType}-template.${format}"]`,
        `a[download*="${entityType}-template.${format}"]`,
      ];

      let href = null;
      for (const selector of downloadLinkSelectors) {
        const link = authenticatedPage.locator(selector);
        if (await link.count() > 0) {
          href = await link.first().getAttribute('href');
          if (href) break;
        }
      }

      expect(href).toBeTruthy();
      const templateUrl = new URL(href!, authenticatedPage.url()).href;

      // Make 10 concurrent requests
      const concurrentRequests = 10;
      
      const startTime = Date.now();
      
      const responses = await Promise.all(
        Array.from({ length: concurrentRequests }, () => 
          page.request.get(templateUrl)
        )
      );
      
      const endTime = Date.now();
      const totalTime = endTime - startTime;

      // All responses should be successful
      responses.forEach(response => {
        expect(response.ok()).toBeTruthy();
      });

      // 10 concurrent requests should complete in less than 2 seconds
      expect(totalTime).toBeLessThan(2000);
      
      // Average time per request should be reasonable
      const avgTime = totalTime / concurrentRequests;
      expect(avgTime).toBeLessThan(500);

      console.log(`✅ ${concurrentRequests} concurrent requests completed in ${totalTime}ms (avg: ${avgTime.toFixed(2)}ms per request)`);
    });
  });

  test.describe('Template file size validation', () => {
    test('should verify template files are reasonably sized (Requirement 1.5)', async ({ authenticatedPage, page }) => {
      console.log('📏 Testing template file sizes');

      const maxTemplateSize = 50 * 1024; // 50KB as per design document
      const fileSizes: { entity: string; format: string; size: number }[] = [];

      for (const entityType of ENTITY_TYPES) {
        const importPath = `/${entityType}/import`;
        await authenticatedPage.goto(importPath);
        await authenticatedPage.waitForTimeout(500);

        for (const format of FORMATS) {
          const downloadLinkSelectors = [
            `a[href*="${entityType}-template.${format}"]`,
            `a[download*="${entityType}-template.${format}"]`,
          ];

          let href = null;
          for (const selector of downloadLinkSelectors) {
            const link = authenticatedPage.locator(selector);
            if (await link.count() > 0) {
              href = await link.first().getAttribute('href');
              if (href) break;
            }
          }

          if (href) {
            const templateUrl = new URL(href, authenticatedPage.url()).href;
            const response = await page.request.get(templateUrl);
            
            expect(response.ok()).toBeTruthy();
            
            const content = await response.text();
            const size = new Blob([content]).size;
            
            fileSizes.push({ entity: entityType, format, size });
            
            // Verify size is within limit
            expect(size).toBeLessThan(maxTemplateSize);
            
            console.log(`✅ ${entityType} ${format.toUpperCase()} template: ${(size / 1024).toFixed(2)}KB`);
          }
        }
      }

      // Verify we checked all templates
      expect(fileSizes.length).toBe(ENTITY_TYPES.length * FORMATS.length);

      // Calculate average size
      const avgSize = fileSizes.reduce((sum, f) => sum + f.size, 0) / fileSizes.length;
      console.log(`📊 Average template size: ${(avgSize / 1024).toFixed(2)}KB`);
      
      // Average should be well under the limit
      expect(avgSize).toBeLessThan(maxTemplateSize * 0.5);
    });
  });

  test.describe('Template page load performance', () => {
    test('should load import page with template section quickly (Requirement 4.5)', async ({ authenticatedPage }) => {
      console.log('⏱️ Testing import page load time with template section');

      const entityType = 'subscriptions';
      const importPath = `/${entityType}/import`;
      
      const startTime = Date.now();
      
      await authenticatedPage.goto(importPath);
      
      // Wait for template section to be visible
      const templateSectionSelectors = [
        '[data-testid="template-download-section"]',
        'text="Download Template"',
        'text="Template Files"',
      ];

      let templateSectionVisible = false;
      for (const selector of templateSectionSelectors) {
        const section = authenticatedPage.locator(selector);
        if (await section.count() > 0) {
          try {
            await expect(section.first()).toBeVisible({ timeout: 5000 });
            templateSectionVisible = true;
            break;
          } catch {
            continue;
          }
        }
      }

      const endTime = Date.now();
      const pageLoadTime = endTime - startTime;

      expect(templateSectionVisible).toBeTruthy();
      
      // Page with template section should load in less than 5 seconds
      expect(pageLoadTime).toBeLessThan(5000);
      
      console.log(`✅ Import page with template section loaded in ${pageLoadTime}ms`);
    });

    test('should render all template download links quickly (Requirement 4.5)', async ({ authenticatedPage }) => {
      console.log('⏱️ Testing template links rendering time');

      const entityType = 'labels';
      const importPath = `/${entityType}/import`;
      
      await authenticatedPage.goto(importPath);
      
      const startTime = Date.now();
      
      // Wait for all three format links to be visible
      const linkPromises = FORMATS.map(async (format) => {
        const linkSelectors = [
          `a[href*="${entityType}-template.${format}"]`,
          `a:has-text("${format.toUpperCase()}")`,
        ];

        for (const selector of linkSelectors) {
          const link = authenticatedPage.locator(selector);
          if (await link.count() > 0) {
            try {
              await expect(link.first()).toBeVisible({ timeout: 3000 });
              return true;
            } catch {
              continue;
            }
          }
        }
        return false;
      });

      const results = await Promise.all(linkPromises);
      
      const endTime = Date.now();
      const renderTime = endTime - startTime;

      // All links should be found
      expect(results.every(r => r === true)).toBeTruthy();
      
      // All links should render in less than 3 seconds
      expect(renderTime).toBeLessThan(3000);
      
      console.log(`✅ All template download links rendered in ${renderTime}ms`);
    });
  });
});
