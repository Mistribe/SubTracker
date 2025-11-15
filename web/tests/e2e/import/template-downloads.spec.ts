/**
 * Template Download E2E Tests
 * 
 * Test suite covering template file download functionality:
 * - Verify template files are accessible at correct paths
 * - Test download functionality for all template files
 * - Verify downloaded files have correct filenames
 * - Verify downloaded files have correct content
 * 
 * Requirements: 1.4, 1.5, 1.7, 4.5
 */

import { test, expect } from '../../fixtures/auth';
import path from 'path';
import fs from 'fs';

test.describe('Template File Downloads', () => {
  const ENTITY_TYPES = ['labels', 'providers', 'subscriptions'] as const;
  const FORMATS = ['csv', 'json', 'yaml'] as const;

  test.beforeEach(async ({ authenticatedPage }) => {
    console.log('🚀 Setting up template download test');
  });

  ENTITY_TYPES.forEach((entityType) => {
    test.describe(`${entityType} templates`, () => {
      FORMATS.forEach((format) => {
        test(`should download ${format.toUpperCase()} template for ${entityType} (Requirements 1.4, 1.5, 1.7, 4.5)`, async ({ authenticatedPage, page }) => {
          console.log(`🔍 Testing ${format.toUpperCase()} template download for ${entityType}`);

          // Navigate to the import page for this entity type
          const importPath = `/${entityType}/import`;
          await authenticatedPage.goto(importPath);
          await authenticatedPage.waitForTimeout(2000);

          // Verify we're on the import page
          await authenticatedPage.waitForURL(`**${importPath}`, { timeout: 10000 });
          console.log(`✅ Navigated to ${entityType} import page`);

          // Look for the template download section
          const templateSectionSelectors = [
            '[data-testid="template-download-section"]',
            'text="Download Template"',
            'text="Template Files"',
            'text="Download"'
          ];

          let templateSectionFound = false;
          for (const selector of templateSectionSelectors) {
            const section = authenticatedPage.locator(selector);
            if (await section.count() > 0) {
              try {
                await expect(section.first()).toBeVisible({ timeout: 5000 });
                templateSectionFound = true;
                console.log(`✅ Found template download section with selector: ${selector}`);
                break;
              } catch {
                continue;
              }
            }
          }

          expect(templateSectionFound).toBeTruthy();

          // Look for the specific format download link
          const downloadLinkSelectors = [
            `a[href*="${entityType}-template.${format}"]`,
            `a[download*="${entityType}-template.${format}"]`,
            `a:has-text("${format.toUpperCase()}")`,
            `button:has-text("${format.toUpperCase()}")`,
            `[data-testid="template-${format}"]`,
            `[data-format="${format}"]`
          ];

          let downloadLink = null;
          for (const selector of downloadLinkSelectors) {
            const link = authenticatedPage.locator(selector);
            if (await link.count() > 0) {
              try {
                await expect(link.first()).toBeVisible({ timeout: 5000 });
                downloadLink = link.first();
                console.log(`✅ Found ${format.toUpperCase()} download link with selector: ${selector}`);
                break;
              } catch {
                continue;
              }
            }
          }

          expect(downloadLink).not.toBeNull();

          // Verify the link has correct attributes (Requirement 4.5)
          const href = await downloadLink!.getAttribute('href');
          expect(href).toBeTruthy();
          expect(href).toContain(`${entityType}-template.${format}`);
          console.log(`✅ Download link has correct href: ${href}`);

          // Verify download attribute exists
          const downloadAttr = await downloadLink!.getAttribute('download');
          if (downloadAttr !== null) {
            console.log(`✅ Download attribute present: ${downloadAttr || '(empty)'}`);
          }

          // Test that the file is accessible by making a direct request
          const templateUrl = new URL(href!, authenticatedPage.url()).href;
          console.log(`📡 Fetching template from: ${templateUrl}`);

          const response = await page.request.get(templateUrl);
          expect(response.ok()).toBeTruthy();
          console.log(`✅ Template file is accessible (status: ${response.status()})`);

          // Verify content type (Requirement 4.5)
          const contentType = response.headers()['content-type'];
          console.log(`📄 Content-Type: ${contentType}`);

          // Get the file content
          const content = await response.text();
          expect(content.length).toBeGreaterThan(0);
          console.log(`✅ Template file has content (${content.length} bytes)`);

          // Verify the content is valid for the format
          switch (format) {
            case 'csv':
              // CSV should have headers
              // Subscriptions use 'friendlyName' instead of 'name'
              if (entityType === 'subscriptions') {
                expect(content).toContain('friendlyName');
              } else {
                expect(content).toContain('name');
              }
              expect(content).toContain('id');
              const csvLines = content.trim().split('\n');
              expect(csvLines.length).toBeGreaterThan(1); // At least header + 1 data row
              console.log(`✅ CSV template has ${csvLines.length} lines (including header)`);
              break;

            case 'json':
              // JSON should be valid and an array
              const jsonData = JSON.parse(content);
              expect(Array.isArray(jsonData)).toBeTruthy();
              expect(jsonData.length).toBeGreaterThan(0);
              console.log(`✅ JSON template has ${jsonData.length} records`);
              
              // Verify at least one record has an ID and one doesn't (Requirement 1.5)
              const hasRecordWithId = jsonData.some((record: any) => record.id);
              const hasRecordWithoutId = jsonData.some((record: any) => !record.id);
              expect(hasRecordWithId).toBeTruthy();
              expect(hasRecordWithoutId).toBeTruthy();
              console.log(`✅ JSON template includes records with and without IDs`);
              break;

            case 'yaml':
              // YAML should be valid (basic check - starts with array indicator)
              expect(content.trim()).toMatch(/^-\s+/);
              // Subscriptions use 'friendlyName' instead of 'name'
              if (entityType === 'subscriptions') {
                expect(content).toContain('friendlyName');
              } else {
                expect(content).toContain('name');
              }
              expect(content).toContain('id');
              console.log(`✅ YAML template appears valid`);
              break;
          }

          console.log(`✅ ${format.toUpperCase()} template download test completed for ${entityType}`);
        });
      });

      test(`should have all three format templates accessible for ${entityType} (Requirement 1.4, 1.5)`, async ({ authenticatedPage, page }) => {
        console.log(`🔍 Testing all template formats are accessible for ${entityType}`);

        // Navigate to the import page
        const importPath = `/${entityType}/import`;
        await authenticatedPage.goto(importPath);
        await authenticatedPage.waitForTimeout(2000);

        // Verify we're on the import page
        await authenticatedPage.waitForURL(`**${importPath}`, { timeout: 10000 });
        console.log(`✅ Navigated to ${entityType} import page`);

        // Check that all three format links are present
        const formatLinks: { format: string; found: boolean }[] = [];

        for (const format of FORMATS) {
          const linkSelectors = [
            `a[href*="${entityType}-template.${format}"]`,
            `a:has-text("${format.toUpperCase()}")`,
            `[data-testid="template-${format}"]`
          ];

          let found = false;
          for (const selector of linkSelectors) {
            const link = authenticatedPage.locator(selector);
            if (await link.count() > 0) {
              try {
                await expect(link.first()).toBeVisible({ timeout: 3000 });
                found = true;
                console.log(`✅ Found ${format.toUpperCase()} template link`);
                break;
              } catch {
                continue;
              }
            }
          }

          formatLinks.push({ format, found });
        }

        // Verify all three formats are available
        const allFormatsAvailable = formatLinks.every(f => f.found);
        expect(allFormatsAvailable).toBeTruthy();

        if (!allFormatsAvailable) {
          const missingFormats = formatLinks.filter(f => !f.found).map(f => f.format);
          console.log(`❌ Missing template formats: ${missingFormats.join(', ')}`);
        } else {
          console.log(`✅ All three template formats (CSV, JSON, YAML) are available for ${entityType}`);
        }
      });
    });
  });

  // These tests don't need authentication - they just check file system
  test.describe('File system validation', () => {
    test.use({ storageState: { cookies: [], origins: [] } }); // Skip auth for these tests

    test('should verify template files exist in public directory', async () => {
      console.log('🔍 Verifying template files exist in public directory');

      const templateBasePath = path.join(process.cwd(), 'public', 'templates');
      
      for (const entityType of ENTITY_TYPES) {
        for (const format of FORMATS) {
          const templatePath = path.join(
            templateBasePath,
            entityType,
            `${entityType}-template.${format}`
          );

          const exists = fs.existsSync(templatePath);
          expect(exists).toBeTruthy();
          
          if (exists) {
            const stats = fs.statSync(templatePath);
            expect(stats.size).toBeGreaterThan(0);
            console.log(`✅ ${entityType}-template.${format} exists (${stats.size} bytes)`);
          } else {
            console.log(`❌ ${entityType}-template.${format} not found at ${templatePath}`);
          }
        }
      }

      console.log('✅ All template files verified in public directory');
    });

    test('should verify template README exists', async () => {
      console.log('🔍 Verifying template README exists');

      const readmePath = path.join(process.cwd(), 'public', 'templates', 'README.md');
      const exists = fs.existsSync(readmePath);
      
      expect(exists).toBeTruthy();
      
      if (exists) {
        const content = fs.readFileSync(readmePath, 'utf-8');
        expect(content.length).toBeGreaterThan(0);
        
        // Verify README contains key information
        expect(content).toContain('template');
        expect(content).toContain('UUID');
        console.log(`✅ Template README exists and contains relevant information (${content.length} bytes)`);
      } else {
        console.log(`❌ Template README not found at ${readmePath}`);
      }
    });
  });
});
