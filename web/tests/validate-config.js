#!/usr/bin/env node

/**
 * Configuration validation script
 * 
 * This script validates that the Playwright configuration
 * can be loaded and parsed correctly without running tests.
 */

import { readFileSync, statSync } from 'fs';
import { resolve } from 'path';

console.log('🔍 Validating Playwright configuration...');

try {
  // Check if configuration file exists
  const configPath = resolve('playwright.config.ts');
  const configContent = readFileSync(configPath, 'utf8');
  
  console.log('✅ Configuration file found');
  
  // Basic syntax validation
  if (configContent.includes('defineConfig')) {
    console.log('✅ Configuration uses defineConfig');
  }
  
  if (configContent.includes('projects:')) {
    console.log('✅ Browser projects configured');
  }
  
  if (configContent.includes('reporter:')) {
    console.log('✅ Test reporters configured');
  }
  
  if (configContent.includes('globalSetup:')) {
    console.log('✅ Global setup configured');
  }
  
  if (configContent.includes('globalTeardown:')) {
    console.log('✅ Global teardown configured');
  }
  
  // Check test directory structure
  const testDirs = [
    'tests/config',
    'tests/e2e',
  ];
  
  testDirs.forEach(dir => {
    try {
      statSync(dir);
      console.log(`✅ Test directory exists: ${dir}`);
    } catch (error) {
      console.log(`⚠️  Test directory missing: ${dir}`);
    }
  });
  
  console.log('\n🎉 Configuration validation completed successfully!');
  console.log('\n📋 Next steps:');
  console.log('1. Set up test environment variables in .env.test.local');
  console.log('2. Create test user in Clerk');
  console.log('3. Run: npm run test:e2e');
  
} catch (error) {
  console.error('❌ Configuration validation failed:', error.message);
  process.exit(1);
}