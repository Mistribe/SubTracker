# E2E Tests - Real API Integration

This directory contains end-to-end tests that use **real API calls** instead of mocked data. The tests are designed to work with a live backend server and provide comprehensive testing of the subscription management functionality.

## 🚀 Quick Start

### Prerequisites

1. **Backend Server**: Ensure your backend API server is running on `http://localhost:8080` (or set `PLAYWRIGHT_API_URL`)
2. **Frontend Server**: Ensure your frontend is running on `http://localhost:5173` (or set `PLAYWRIGHT_BASE_URL`)
3. **Authentication**: Valid test user credentials in `.env.test`

### Environment Setup

Create or update your `.env.test` file:

```bash
# Clerk Configuration
VITE_CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
CLERK_SECRET_KEY=your_clerk_secret_key

# Test Application URLs
PLAYWRIGHT_BASE_URL=http://localhost:5173
PLAYWRIGHT_API_URL=http://localhost:8080

# Test User Credentials
TEST_USER_EMAIL=your_test_user@example.com
TEST_USER_PASSWORD=your_test_password

# Environment
NODE_ENV=test
```

### Running Tests

```bash
# Run all e2e tests
npm run test:e2e

# Run specific test suite
npx playwright test tests/e2e/subscriptions/

# Run specific test file
npx playwright test tests/e2e/subscriptions/create.spec.ts

# Run tests in headed mode (see browser)
npx playwright test --headed

# Run tests with debug mode
npx playwright test --debug
```

## 🧪 Test Structure

### Test Files

- **`create.spec.ts`** - Tests subscription creation through the UI wizard
- **`delete.spec.ts`** - Tests subscription deletion with confirmation dialogs
- **`edit.spec.ts`** - Tests subscription editing and updates
- **`detail-view.spec.ts`** - Tests subscription display and information accuracy
- **`filters.spec.ts`** - Tests search and filtering functionality
- **`ui-only.spec.ts`** - Tests UI components without API dependencies
- **`debug.spec.ts`** - Debug utilities for troubleshooting

### Key Features

✅ **Real API Integration** - All tests use actual API calls, no mocking
✅ **Automatic Cleanup** - Test data is automatically created and cleaned up
✅ **Authentication** - Tests run with authenticated users using Clerk
✅ **Error Handling** - Robust error handling and meaningful error messages
✅ **Data Isolation** - Each test creates its own test data to avoid conflicts

## 🔧 API Integration

### How It Works

1. **Authentication**: Tests extract Clerk session tokens from authenticated pages
2. **Test Data Creation**: Each test creates its own providers and subscriptions via API
3. **UI Interaction**: Tests interact with the actual UI components
4. **Data Verification**: Tests verify data through both UI and API responses
5. **Cleanup**: All created test data is automatically cleaned up after each test

### API Client Features

- **Multiple Health Check Endpoints**: Tries various health check URLs
- **Flexible Status Validation**: Accepts multiple health status formats
- **Automatic Token Management**: Handles authentication token extraction and usage
- **Comprehensive Error Handling**: Provides clear error messages for debugging

## 🛠️ Troubleshooting

### Common Issues

#### 1. API Connection Failed
```
❌ API connection failed: API health check failed
```

**Solutions:**
- Verify backend server is running on the correct port
- Check `PLAYWRIGHT_API_URL` environment variable
- Run the API verification script: `npx ts-node tests/utils/verify-api.ts`

#### 2. Authentication Token Issues
```
❌ Could not extract authentication token from page
```

**Solutions:**
- Verify Clerk configuration in `.env.test`
- Check test user credentials are valid
- Ensure user can successfully log in through the UI

#### 3. Test Data Cleanup Issues
```
⚠️ Failed to cleanup subscription: [name]
```

**Solutions:**
- This is usually not critical - cleanup failures are logged but don't fail tests
- Check API permissions for delete operations
- Verify test data was created successfully

### Debug Mode

Run tests in debug mode to step through issues:

```bash
npx playwright test --debug tests/e2e/subscriptions/debug.spec.ts
```

### API Verification

Test API connectivity independently:

```bash
npx ts-node tests/utils/verify-api.ts
```

## 📊 Test Data Management

### Data Creation Strategy

Each test follows this pattern:

1. **Setup**: Create test providers and subscriptions via API
2. **Execute**: Perform UI interactions and verifications
3. **Cleanup**: Delete all created test data

### Data Isolation

- Each test creates unique test data with timestamps
- Tests don't depend on existing data in the system
- Cleanup ensures no test data pollution

### Naming Convention

Test data uses descriptive names with timestamps:
- Providers: `Test Provider 1 1699123456789`
- Subscriptions: `Basic E2E Subscription 1699123456789`

## 🎯 Best Practices

### Writing New Tests

1. **Use Real API**: Always use `testHelpers` for data creation
2. **Unique Names**: Include timestamps in test data names
3. **Proper Cleanup**: Track created data for cleanup
4. **Error Handling**: Handle API failures gracefully
5. **Meaningful Assertions**: Verify both UI state and data accuracy

### Example Test Structure

```typescript
test('should perform subscription operation', async () => {
  // 1. Create test data via API
  const testProvider = await testHelpers.createTestProvider({
    name: `Test Provider ${Date.now()}`,
    description: 'Test provider description',
  });
  
  const testSubscription = TestDataGenerators.generateSubscription({
    providerId: testProvider.id,
    name: `Test Subscription ${Date.now()}`,
  });
  
  // 2. Perform UI interactions
  await subscriptionsPage.clickAddSubscription();
  await subscriptionsPage.fillSubscriptionForm(testSubscription);
  await subscriptionsPage.submitSubscriptionForm();
  
  // 3. Verify results
  await subscriptionsPage.verifySubscriptionInTable({
    name: testSubscription.name,
    provider: testProvider.name,
  });
  
  // 4. Track for cleanup (handled automatically in afterEach)
  createdSubscriptionsForCleanup.push(testSubscription.name);
});
```

## 🔍 Monitoring and Reporting

### Test Reports

Tests generate multiple report formats:
- **HTML Report**: `test-results/html-report/index.html`
- **JUnit Report**: `test-results/junit-report.xml`
- **JSON Report**: `test-results/test-results.json`

### Screenshots and Videos

Failed tests automatically capture:
- Screenshots: `test-results/artifacts/`
- Videos: `test-results/artifacts/`
- Traces: `test-results/artifacts/`

### Console Output

Tests provide detailed console output:
- ✅ Successful operations
- ⚠️ Warnings (non-critical issues)
- ❌ Errors (test failures)

## 🚀 CI/CD Integration

### Environment Variables

Set these in your CI/CD environment:

```bash
PLAYWRIGHT_BASE_URL=https://your-staging-app.com
PLAYWRIGHT_API_URL=https://your-staging-api.com
TEST_USER_EMAIL=ci-test-user@example.com
TEST_USER_PASSWORD=secure_password
VITE_CLERK_PUBLISHABLE_KEY=pk_test_...
CLERK_SECRET_KEY=sk_test_...
```

### CI Configuration

The tests are configured for CI environments:
- Headless mode by default
- Retry failed tests
- Generate comprehensive reports
- Capture artifacts on failure

## 📝 Contributing

When adding new tests:

1. Follow the existing patterns for API integration
2. Ensure proper test data cleanup
3. Add meaningful test descriptions
4. Include error handling
5. Update this README if adding new test categories

## 🆘 Support

If you encounter issues:

1. Check the troubleshooting section above
2. Run the API verification script
3. Check test logs and screenshots
4. Verify environment configuration
5. Run tests in debug mode for detailed inspection