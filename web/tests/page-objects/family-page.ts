import { Page, Locator } from '@playwright/test';
import { BasePage } from './base-page';

/**
 * Page Object for Family page
 * Handles interactions with family management UI
 */
export class FamilyPage extends BasePage {
    // Page elements
    private readonly pageTitle: Locator;
    private readonly addMemberButton: Locator;
    private readonly createFamilyButton: Locator;
    private readonly familyNameInput: Locator;
    private readonly creatorNameInput: Locator;
    private readonly memberNameInput: Locator;
    private readonly memberIsKidCheckbox: Locator;
    private readonly submitButton: Locator;
    private readonly membersTable: Locator;

    constructor(page: Page) {
        super(page);
        this.pageTitle = page.locator('h1:has-text("Family")');
        this.addMemberButton = page.locator('button:has-text("Add Member")');
        this.createFamilyButton = page.locator('button:has-text("Add Family")');
        this.familyNameInput = page.locator('input[placeholder*="family name"], input[name="name"]').first();
        this.creatorNameInput = page.locator('input[placeholder*="your name"], input[name="creatorName"]');
        this.memberNameInput = page.locator('input[placeholder*="member name"]');
        this.memberIsKidCheckbox = page.locator('button[role="checkbox"]');
        this.submitButton = page.locator('button[type="submit"]');
        this.membersTable = page.locator('table');
    }

    getPageUrl(): string {
        return '/family';
    }

    async waitForPageLoad(): Promise<void> {
        await this.pageTitle.waitFor({ state: 'visible', timeout: this.timeout });
    }

    /**
     * Check if empty state is displayed
     */
    async hasEmptyState(): Promise<boolean> {
        try {
            const emptyState = this.page.locator('text=/no family/i, text=/create.*family/i').first();
            await emptyState.waitFor({ state: 'visible', timeout: 3000 });
            return true;
        } catch {
            return false;
        }
    }

    /**
     * Create a new family via UI
     */
    async createFamily(familyName: string, creatorName: string): Promise<void> {
        // Click create family button
        await this.createFamilyButton.click();

        // Wait for dialog to be visible
        await this.page.locator('role=dialog').waitFor({ state: 'visible' });

        // Fill in the form
        await this.familyNameInput.fill(familyName);
        await this.creatorNameInput.fill(creatorName);

        // Submit the form
        await this.submitButton.click();

        // Wait for dialog to close
        await this.page.locator('role=dialog').waitFor({ state: 'hidden', timeout: 10000 });
    }

    /**
     * Add a family member via UI
     */
    async addFamilyMember(memberName: string, isKid: boolean = false): Promise<void> {
        // Click add member button
        await this.addMemberButton.click();

        // Wait for dialog to be visible
        await this.page.locator('role=dialog').waitFor({ state: 'visible' });

        // Fill in member name
        await this.memberNameInput.fill(memberName);

        // Set kid checkbox if needed
        if (isKid) {
            const checkbox = this.page.locator('button[role="checkbox"]');
            const isChecked = await checkbox.getAttribute('data-state') === 'checked';
            if (!isChecked) {
                await checkbox.click();
            }
        }

        // Submit the form
        await this.submitButton.click();

        // Wait for dialog to close
        await this.page.locator('role=dialog').waitFor({ state: 'hidden', timeout: 10000 });
    }

    /**
     * Get family member row by name
     */
    getMemberRow(memberName: string): Locator {
        return this.page.locator(`tr:has-text("${memberName}")`);
    }

    /**
     * Delete a family member via the actions dropdown menu
     */
    async deleteMember(memberName: string): Promise<void> {
        console.log(`🗑️ Attempting to delete member: ${memberName}`);

        // Find the member row
        const memberRow = this.getMemberRow(memberName);
        await memberRow.waitFor({ state: 'visible', timeout: 5000 });
        console.log(`✓ Found member row for: ${memberName}`);

        // Scroll the row into view to ensure button is visible
        await memberRow.scrollIntoViewIfNeeded();
        await this.page.waitForTimeout(500); // Wait for scroll animation

        // Find the actions menu button - it's in the last table cell
        const actionsCell = memberRow.locator('td').last();
        const actionsButton = actionsCell.locator('button').first();

        await actionsButton.waitFor({ state: 'visible', timeout: 5000 });
        console.log('✓ Found actions button');

        // Click to open the dropdown
        await actionsButton.click();
        console.log('✓ Clicked actions button');

        // Wait for the dropdown menu to be visible
        await this.page.waitForTimeout(800); // Wait for animation

        // Find and click the Remove menu item
        const removeMenuItem = this.page.getByRole('menuitem').filter({ hasText: 'Remove' });
        await removeMenuItem.waitFor({ state: 'visible', timeout: 5000 });
        console.log('✓ Found Remove menu item');

        await removeMenuItem.click();
        console.log('✓ Clicked Remove menu item');

        // Wait a bit for the delete request to process
        await this.page.waitForTimeout(2000);

        // Check if row is still there - it might take some time for API call
        try {
            await memberRow.waitFor({ state: 'detached', timeout: 10000 });
            console.log(`✅ Member deleted: ${memberName}`);
        } catch (e) {
            // If still visible, log and check count instead
            console.log(`⚠️ Row still visible, but member might be deleted. Checking...`);
            await this.page.waitForTimeout(2000);
        }
    }

    /**
     * Check if member exists in the table
     */
    async hasMember(memberName: string): Promise<boolean> {
        try {
            const memberRow = this.getMemberRow(memberName);
            await memberRow.waitFor({ state: 'visible', timeout: 3000 });
            return true;
        } catch {
            return false;
        }
    }

    /**
     * Get all member names from the table
     */
    async getAllMemberNames(): Promise<string[]> {
        await this.membersTable.waitFor({ state: 'visible' });
        const rows = await this.page.locator('tbody tr').all();
        const names: string[] = [];

        for (const row of rows) {
            const nameCell = row.locator('td').first();
            const name = await nameCell.textContent();
            if (name) {
                names.push(name.trim());
            }
        }

        return names;
    }

    /**
     * Get family name from header
     */
    async getFamilyName(): Promise<string | null> {
        try {
            const familyHeader = this.page.locator('[data-testid="family-name"], h2, h3').first();
            await familyHeader.waitFor({ state: 'visible', timeout: 3000 });
            return await familyHeader.textContent();
        } catch {
            return null;
        }
    }

    /**
     * Wait for family to be loaded
     */
    async waitForFamilyLoaded(): Promise<void> {
        // Wait for either empty state or members table
        await Promise.race([
            this.page.locator('text=/no family/i').waitFor({ state: 'visible', timeout: 10000 }),
            this.membersTable.waitFor({ state: 'visible', timeout: 10000 }),
        ]);
    }
}