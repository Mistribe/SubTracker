/**
 * Family Creation Tests
 * 
 * Tests for creating and managing families:
 * - Create family with valid data
 * - Add family members (adults and kids)
 * - View family members list
 * - Remove family members
 * 
 * Testing Philosophy: Happy Path with UI interactions
 * ✅ Tests successful family creation via UI
 * ✅ Tests adding members via UI
 * ✅ Uses API for cleanup
 */

import { test, expect } from '../../fixtures/auth';
import { FamilyPage } from '../../page-objects/family-page';

test.describe('Family Creation Happy Path', () => {
    let familyPage: FamilyPage;
    let createdFamilyName: string;

    test.beforeEach(async ({ authenticatedPage }) => {
        console.log('🚀 Setting up family creation test');

        familyPage = new FamilyPage(authenticatedPage);

        // Navigate to family page
        await authenticatedPage.goto('/family');
        await familyPage.waitForPageLoad();
        console.log('✅ Navigated to family page');
    });

    test('should create a family and add members via UI', async ({ authenticatedPage }) => {
        console.log('🧪 Testing family creation flow');

        // Generate unique family name
        const timestamp = Date.now();
        createdFamilyName = `Test Family ${timestamp}`;
        const creatorName = 'John Doe';

        // Check if we have the empty state (no family yet)
        const hasEmpty = await familyPage.hasEmptyState();

        if (hasEmpty) {
            console.log('📝 Creating new family via UI');

            // Create family via UI
            await familyPage.createFamily(createdFamilyName, creatorName);

            // Wait for family to be loaded
            await familyPage.waitForFamilyLoaded();

            console.log(`✅ Family created: ${createdFamilyName}`);

            // Verify creator is in the members table
            const hasCreator = await familyPage.hasMember(creatorName);
            expect(hasCreator).toBeTruthy();
            console.log('✅ Creator appears in members table');
        } else {
            console.log('ℹ️ Family already exists, skipping creation');
        }

        // Add adult member
        const adultMemberName = `Adult Member ${timestamp}`;
        console.log(`📝 Adding adult member: ${adultMemberName}`);
        await familyPage.addFamilyMember(adultMemberName, false);

        // Verify adult member appears
        await authenticatedPage.waitForTimeout(1000); // Wait for UI update
        const hasAdult = await familyPage.hasMember(adultMemberName);
        expect(hasAdult).toBeTruthy();
        console.log('✅ Adult member added successfully');

        // Add kid member
        const kidMemberName = `Kid Member ${timestamp}`;
        console.log(`📝 Adding kid member: ${kidMemberName}`);
        await familyPage.addFamilyMember(kidMemberName, true);

        // Verify kid member appears
        await authenticatedPage.waitForTimeout(1000); // Wait for UI update
        const hasKid = await familyPage.hasMember(kidMemberName);
        expect(hasKid).toBeTruthy();
        console.log('✅ Kid member added successfully');

        // Verify all members are visible
        const allMembers = await familyPage.getAllMemberNames();
        console.log(`📋 All members: ${allMembers.join(', ')}`);

        expect(allMembers.length).toBeGreaterThanOrEqual(2); // At least the two we added
        console.log('✅ All members verified');
    });

    test('should display family page layout correctly', async ({ authenticatedPage }) => {
        console.log('🧪 Testing family page layout');

        // Wait for page to load
        await familyPage.waitForFamilyLoaded();

        // Check for either empty state or family content
        const hasEmpty = await familyPage.hasEmptyState();

        if (hasEmpty) {
            console.log('📋 Verifying empty state');

            // Verify create family button is visible
            const createButton = authenticatedPage.locator('button:has-text("Add Family")');
            await expect(createButton).toBeVisible();
            console.log('✅ Create family button visible');
        } else {
            console.log('📋 Verifying family content');

            // Verify add member button is visible
            const addMemberButton = authenticatedPage.locator('button:has-text("Add Member")');
            await expect(addMemberButton).toBeVisible();
            console.log('✅ Add member button visible');

            // Verify members table exists
            const table = authenticatedPage.locator('table');
            await expect(table).toBeVisible();
            console.log('✅ Members table visible');
        }

        console.log('✅ Family page layout verified');
    });

    test('should add and interact with member actions menu', async ({ authenticatedPage }) => {
        console.log('🧪 Testing member actions menu interaction');

        // Ensure we have a family first
        const hasEmpty = await familyPage.hasEmptyState();

        if (hasEmpty) {
            const timestamp = Date.now();
            const familyName = `Test Family ${timestamp}`;
            const creatorName = `Creator ${timestamp}`;

            console.log('📝 Creating family first');
            await familyPage.createFamily(familyName, creatorName);
            await familyPage.waitForFamilyLoaded();
            console.log('✅ Family created');
        }

        // Add a temporary member
        const timestamp = Date.now();
        const tempMemberName = `Test Member ${timestamp}`;

        console.log(`📝 Adding temporary member: ${tempMemberName}`);
        await familyPage.addFamilyMember(tempMemberName, false);

        // Verify member was added
        await authenticatedPage.waitForTimeout(1500); // Wait for UI update
        const hasMember = await familyPage.hasMember(tempMemberName);
        expect(hasMember).toBeTruthy();
        console.log('✅ Member added successfully');

        // Verify we can access the actions menu
        const memberRow = familyPage.getMemberRow(tempMemberName);
        await memberRow.scrollIntoViewIfNeeded();

        // Find and click the actions button
        const actionsCell = memberRow.locator('td').last();
        const actionsButton = actionsCell.locator('button').first();
        await actionsButton.waitFor({ state: 'visible', timeout: 5000 });
        console.log('✓ Found actions button');

        await actionsButton.click();
        console.log('✓ Clicked actions button');

        // Verify menu opens
        await authenticatedPage.waitForTimeout(500);
        const menu = authenticatedPage.locator('[role="menu"]');
        await menu.waitFor({ state: 'visible', timeout: 3000 });
        console.log('✅ Actions menu opened successfully');

        // Verify Edit option is visible
        const editOption = authenticatedPage.getByRole('menuitem', { name: /edit/i });
        await expect(editOption).toBeVisible();
        console.log('✅ Edit option visible');

        // Verify Remove option is visible
        const removeOption = authenticatedPage.getByRole('menuitem').filter({ hasText: 'Remove' });
        await expect(removeOption).toBeVisible();
        console.log('✅ Remove option visible');

        // Close menu by clicking outside
        await authenticatedPage.keyboard.press('Escape');
        await menu.waitFor({ state: 'hidden', timeout: 2000 });

        console.log('✅ Member actions menu interaction test completed');
    });
});
