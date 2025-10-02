/**
 * Authentication Setup for E2E Tests
 * 
 * This setup file handles authentication state preparation
 * for tests that require authenticated users
 */

import { test as setup, expect } from '@playwright/test';
import { ClerkAuthHelper } from '../fixtures/auth';
import { ClerkTestUserManager } from '../utils/clerk-user-manager';
import { SessionManager } from '../utils/session-manager';
import { getTestEnvironment } from '../config/test-env';

const testEnv = getTestEnvironment();

/**
 * Setup authenticated state for default test user
 */
setup('authenticate default user', async ({ page }) => {
  console.log('🔐 Setting up authentication for default test user...');
  
  try {
    // Get default test user
    const testUser = await ClerkTestUserManager.getDefaultTestUser();
    
    // Sign in the user
    await ClerkAuthHelper.signInUser(page, testUser);
    
    // Verify authentication was successful
    await expect(page).toHaveURL(/\/dashboard/);
    
    // Setup session persistence
    await SessionManager.setupSessionPersistence(page, testUser);
    
    console.log('✅ Default user authentication setup completed');
    
  } catch (error) {
    console.error('❌ Failed to setup default user authentication:', error);
    throw error;
  }
});

/**
 * Setup authenticated state for family owner (if configured)
 */
setup('authenticate family owner', async ({ page }) => {
  if (!testEnv.familyOwner) {
    console.log('⏭️  Skipping family owner setup - not configured');
    return;
  }
  
  console.log('🔐 Setting up authentication for family owner...');
  
  try {
    // Get family owner test user
    const familyOwner = await ClerkTestUserManager.getFamilyOwnerUser();
    
    // Sign in the family owner
    await ClerkAuthHelper.signInUser(page, familyOwner);
    
    // Verify authentication was successful
    await expect(page).toHaveURL(/\/dashboard/);
    
    // Setup session persistence
    await SessionManager.setupSessionPersistence(page, familyOwner);
    
    console.log('✅ Family owner authentication setup completed');
    
  } catch (error) {
    console.error('❌ Failed to setup family owner authentication:', error);
    throw error;
  }
});

/**
 * Setup authenticated state for family member (if configured)
 */
setup('authenticate family member', async ({ page }) => {
  if (!testEnv.familyMember) {
    console.log('⏭️  Skipping family member setup - not configured');
    return;
  }
  
  console.log('🔐 Setting up authentication for family member...');
  
  try {
    // Get family member test user
    const familyMember = await ClerkTestUserManager.getFamilyMemberUser();
    
    // Sign in the family member
    await ClerkAuthHelper.signInUser(page, familyMember);
    
    // Verify authentication was successful
    await expect(page).toHaveURL(/\/dashboard/);
    
    // Setup session persistence
    await SessionManager.setupSessionPersistence(page, familyMember);
    
    console.log('✅ Family member authentication setup completed');
    
  } catch (error) {
    console.error('❌ Failed to setup family member authentication:', error);
    throw error;
  }
});

/**
 * Verify Clerk authentication system is working
 */
setup('verify clerk system', async ({ page }) => {
  console.log('🔍 Verifying Clerk authentication system...');
  
  try {
    // Navigate to sign-in page
    await page.goto('/sign-in');
    
    // Wait for Clerk to load
    await ClerkAuthHelper.waitForClerkLoaded(page);
    
    // Verify Clerk elements are present
    await expect(page.locator('[data-clerk-component="SignIn"]')).toBeVisible();
    await expect(page.locator('input[name="identifier"], input[type="email"]').first()).toBeVisible();
    await expect(page.locator('input[name="password"], input[type="password"]').first()).toBeVisible();
    
    // Check for submit button (may be hidden initially)
    const submitButton = page.locator('button[type="submit"], .cl-formButtonPrimary').first();
    await expect(submitButton).toBeAttached(); // Just check it exists, not necessarily visible
    
    console.log('✅ Clerk authentication system verification completed');
    
  } catch (error) {
    console.error('❌ Clerk authentication system verification failed:', error);
    throw error;
  }
});