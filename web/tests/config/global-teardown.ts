import { FullConfig } from '@playwright/test';
import { GlobalTestSetup } from '../fixtures';

/**
 * Global teardown for Playwright tests
 * 
 * This teardown runs once after all tests complete and handles:
 * - Test data cleanup
 * - Resource cleanup
 * - Test result summary
 * - Artifact organization
 */
async function globalTeardown(config: FullConfig) {
  console.log('\n🧹 Starting Playwright E2E Test Suite Teardown...');
  
  try {
    // Cleanup test data and resources
    console.log('🧹 Cleaning up test data and resources...');
    await GlobalTestSetup.cleanup();
    
    // Get test summary
    const testSummary = GlobalTestSetup.getTestSummary();
    console.log('📊 Test Data Summary:');
    console.log(`   Test data cleaned: ${JSON.stringify(testSummary.testData)}`);
    console.log(`   Active sessions: ${testSummary.activeSessions}`);
    
    // Log test execution summary
    console.log('\n📊 Test Execution Summary:');
    console.log(`   Configuration: ${config.configFile || 'default'}`);
    console.log(`   Projects: ${config.projects?.length || 0}`);
    console.log(`   Workers: ${config.workers || 'auto'}`);
    
    // Cleanup test artifacts if needed
    if (process.env.CLEANUP_ARTIFACTS === 'true') {
      console.log('🗑️  Cleaning up test artifacts...');
      // Add artifact cleanup logic here if needed
      // For now, we keep artifacts for debugging
    }
    
    // Log artifact locations
    console.log('\n📁 Test Artifacts:');
    console.log('   HTML Report: test-results/html-report/index.html');
    console.log('   JUnit Report: test-results/junit-report.xml');
    console.log('   Screenshots/Videos: test-results/artifacts/');
    
    // Environment-specific cleanup
    if (process.env.CI) {
      console.log('🔧 CI Environment - artifacts will be uploaded by CI system');
    } else {
      console.log('💻 Local Environment - artifacts saved locally for debugging');
    }
    
    console.log('\n✅ Teardown completed successfully!');
    
  } catch (error) {
    console.error('❌ Error during teardown:', error);
    // Don't fail the entire test suite due to teardown issues
  }
}

export default globalTeardown;