import { chromium, FullConfig } from '@playwright/test';
import * as dotenv from 'dotenv';
import { GlobalTestSetup } from '../fixtures';

/**
 * Global setup for Playwright tests
 * 
 * This setup runs once before all tests and handles:
 * - Environment validation
 * - Test data preparation
 * - Authentication setup
 * - Application health checks
 */
async function globalSetup(config: FullConfig) {
  console.log('🚀 Starting Playwright E2E Test Suite Setup...');
  
  // Load environment variables
  const envFile = process.env.NODE_ENV === 'test' ? '.env.test' : '.env';
  dotenv.config({ path: envFile });
  dotenv.config({ path: `${envFile}.local`, override: true });
  
  // Validate required environment variables
  const requiredEnvVars = [
    'VITE_CLERK_PUBLISHABLE_KEY',
    'TEST_USER_EMAIL',
    'TEST_USER_PASSWORD',
  ];
  
  const missingEnvVars = requiredEnvVars.filter(envVar => !process.env[envVar]);
  
  if (missingEnvVars.length > 0) {
    console.error('❌ Missing required environment variables:');
    missingEnvVars.forEach(envVar => {
      console.error(`   - ${envVar}`);
    });
    console.error('\nPlease check your .env.test file or environment configuration.');
    process.exit(1);
  }
  
  // Get base URL from config or environment
  const baseURL = process.env.PLAYWRIGHT_BASE_URL || config.projects[0]?.use?.baseURL || 'http://localhost:5173';
  
  console.log(`📍 Base URL: ${baseURL}`);
  console.log(`🌍 Environment: ${process.env.NODE_ENV || 'development'}`);
  
  // Perform application health check
  try {
    console.log('🔍 Performing application health check...');
    
    const browser = await chromium.launch();
    const context = await browser.newContext();
    const page = await context.newPage();
    
    // Check if the application is accessible
    const response = await page.goto(baseURL, { 
      waitUntil: 'networkidle',
      timeout: 30000 
    });
    
    if (!response || !response.ok()) {
      throw new Error(`Application not accessible at ${baseURL}. Status: ${response?.status()}`);
    }
    
    // Check if Clerk is loaded (basic check)
    try {
      await page.waitForFunction(
        () => (window as any).Clerk !== undefined,
        { timeout: 10000 }
      );
      console.log('✅ Clerk authentication system detected');
    } catch (error) {
      console.warn('⚠️  Clerk authentication system not detected - some tests may fail');
    }
    
    await browser.close();
    console.log('✅ Application health check passed');
    
  } catch (error) {
    console.error('❌ Application health check failed:', error);
    
    if (!process.env.CI) {
      console.log('\n💡 Make sure your development server is running:');
      console.log('   npm run dev');
      console.log('\n   Or check if the PLAYWRIGHT_BASE_URL is correct in your .env.test file');
    }
    
    process.exit(1);
  }
  
  // Log test configuration
  console.log('\n📋 Test Configuration:');
  console.log(`   Workers: ${config.workers || 'auto'}`);
  console.log(`   Projects: ${config.projects?.length || 0}`);
  
  // Initialize test utilities
  try {
    console.log('🔧 Initializing test utilities...');
    await GlobalTestSetup.initialize();
    console.log('✅ Test utilities initialized successfully');
  } catch (error) {
    console.error('❌ Failed to initialize test utilities:', error);
    process.exit(1);
  }

  console.log('\n🎯 Setup completed successfully!\n');
}

export default globalSetup;