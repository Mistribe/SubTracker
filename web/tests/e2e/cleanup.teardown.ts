/**
 * Cleanup Teardown for E2E Tests
 * 
 * This teardown file handles cleanup after all tests complete
 */

import { test as teardown } from '@playwright/test';
import { GlobalTestSetup } from '../fixtures';
import { ClerkTestUserManager } from '../utils/clerk-user-manager';
import { SessionManager } from '../utils/session-manager';

/**
 * Cleanup all test data and sessions
 */
teardown('cleanup test data', async () => {
  console.log('🧹 Starting test data cleanup...');
  
  try {
    // Cleanup all test data
    await GlobalTestSetup.cleanup();
    
    // Get final summary
    const summary = GlobalTestSetup.getTestSummary();
    console.log('📊 Final cleanup summary:', summary);
    
    console.log('✅ Test data cleanup completed successfully');
    
  } catch (error) {
    console.error('❌ Error during test data cleanup:', error);
    // Don't fail the test suite due to cleanup issues
  }
});

/**
 * Reset all user sessions
 */
teardown('reset user sessions', async () => {
  console.log('🔄 Resetting user sessions...');
  
  try {
    // Reset all sessions
    await ClerkTestUserManager.resetAllSessions();
    SessionManager.clearAllSessions();
    
    console.log('✅ User sessions reset completed');
    
  } catch (error) {
    console.error('❌ Error during session reset:', error);
  }
});

/**
 * Cleanup expired sessions
 */
teardown('cleanup expired sessions', async () => {
  console.log('⏰ Cleaning up expired sessions...');
  
  try {
    // Cleanup expired sessions
    SessionManager.cleanupExpiredSessions();
    
    const activeSessions = SessionManager.getActiveSessions();
    console.log(`📊 Active sessions remaining: ${activeSessions.length}`);
    
    console.log('✅ Expired sessions cleanup completed');
    
  } catch (error) {
    console.error('❌ Error during expired sessions cleanup:', error);
  }
});

/**
 * Final cleanup verification
 */
teardown('verify cleanup completion', async () => {
  console.log('🔍 Verifying cleanup completion...');
  
  try {
    // Verify all test users are cleaned up (except default ones)
    const allUsers = await ClerkTestUserManager.getAllTestUsers();
    const defaultUserIds = ['default-test-user', 'family-owner', 'family-member'];
    const remainingTestUsers = allUsers.filter(user => !defaultUserIds.includes(user.id));
    
    if (remainingTestUsers.length > 0) {
      console.warn(`⚠️  ${remainingTestUsers.length} test users were not cleaned up:`, 
        remainingTestUsers.map(u => u.email));
    }
    
    // Verify no active sessions remain
    const activeSessions = SessionManager.getActiveSessions();
    if (activeSessions.length > 0) {
      console.warn(`⚠️  ${activeSessions.length} active sessions remain`);
    }
    
    console.log('✅ Cleanup verification completed');
    
  } catch (error) {
    console.error('❌ Error during cleanup verification:', error);
  }
});