#!/usr/bin/env ts-node

/**
 * Simple script to verify API connectivity for e2e tests
 * Run with: npx ts-node tests/utils/verify-api.ts
 */

import { createTestApiClient } from './api-client';

async function verifyApiConnection() {
  console.log('🔍 Verifying API connection...');
  
  const apiUrl = process.env.PLAYWRIGHT_API_URL || process.env.API_BASE_URL || 'http://localhost:8080';
  console.log(`📡 API URL: ${apiUrl}`);
  
  const apiClient = createTestApiClient({ baseUrl: apiUrl });
  
  try {
    await apiClient.initialize();
    console.log('✅ API client initialized successfully');
    
    const healthResponse = await apiClient.healthCheck();
    console.log('🏥 Health check response:', {
      success: healthResponse.success,
      status: healthResponse.status,
      data: healthResponse.data,
      error: healthResponse.error
    });
    
    if (healthResponse.success) {
      console.log('✅ API is healthy and ready for e2e tests');
    } else {
      console.log('❌ API health check failed');
      process.exit(1);
    }
    
  } catch (error) {
    console.error('❌ Failed to connect to API:', error);
    process.exit(1);
  } finally {
    await apiClient.dispose();
  }
}

// Run the verification
verifyApiConnection().catch(error => {
  console.error('❌ Verification failed:', error);
  process.exit(1);
});