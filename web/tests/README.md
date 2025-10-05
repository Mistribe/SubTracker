# Playwright E2E Tests

This directory contains the end-to-end test suite for the SubTracker application using Playwright.

## Setup

### Prerequisites

1. **Node.js**: Ensure you have Node.js installed (version 18 or higher)
2. **Dependencies**: Install project dependencies:
   ```bash
   npm install
   ```
3. **Playwright Browsers**: Install Playwright browsers:
   ```bash
   npx playwright install
   ```

### Environment Configuration

1. **Copy environment file**:
   ```bash
   cp .env.test.example .env.test.local
   ```

2. **Configure test environment variables** in `.env.test.local`:
   ```bash
   # Clerk Configuration
   VITE_CLERK_PUBLISHABLE_KEY=your_clerk_publishable_key
   CLERK_SECRET_KEY=your_clerk_secret_key
   
   # Application URLs
   PLAYWRIGHT_BASE_URL=http://localhost:5173
   PLAYWRIGHT_API_URL=http://localhost:8080
   
   # Test User Credentials
   TEST_USER_EMAIL=test@example.com
   TEST_USER_PASSWORD=testpassword123
   ```

3. **Create test user**: Create a test user in your Clerk application with the credentials specified above.

## Running Tests

### Development

```bash
# Run all tests
npm run test:e2e

# Run tests in headed mode (see browser)
npm run test:e2e:headed

# Run tests with UI mode
npm run test:e2e:ui

# Run specific test file
npx playwright test tests/e2e/auth/login.spec.ts

# Run tests in specific browser
npx playwright test --project=chromium
```

### Debugging

```bash
# Debug mode (step through tests)
npm run test:e2e:debug

# Run with trace viewer
npx playwright test --trace on

# Show test report
npm run test:e2e:report
```

### CI/CD

```bash
# Run tests in CI mode
CI=true npm run test:e2e
```

## Test Structure

```
tests/
├── config/                 # Test configuration files
│   ├── global-setup.ts     # Global setup (runs once before all tests)
│   ├── global-teardown.ts  # Global teardown (runs once after all tests)
│   └── test-env.ts         # Environment configuration utilities
├── e2e/                    # End-to-end test files
│   ├── auth.setup.ts       # Authentication setup
│   ├── cleanup.teardown.ts # Cleanup teardown
│   ├── auth/               # Authentication tests
│   ├── dashboard/          # Dashboard tests
│   ├── subscriptions/      # Subscription management tests
│   ├── providers/          # Provider management tests
│   ├── labels/             # Label management tests
│   ├── family/             # Family management tests
│   └── profile/            # Profile and preferences tests
├── fixtures/               # Test fixtures and utilities
├── page-objects/           # Page Object Model classes
└── utils/                  # Test utilities and helpers
```

## Configuration

The Playwright configuration (`playwright.config.ts`) supports:

- **Multiple environments**: Development, test, staging, production
- **Multiple browsers**: Chromium, Firefox, WebKit
- **Device testing**: Desktop, tablet, mobile viewports
- **Parallel execution**: Configurable workers and test isolation
- **Retry logic**: Automatic retries on failure
- **Rich reporting**: HTML, JUnit, JSON reports
- **Artifacts**: Screenshots, videos, traces on failure

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `VITE_CLERK_PUBLISHABLE_KEY` | Clerk publishable key | Yes |
| `CLERK_SECRET_KEY` | Clerk secret key | No |
| `PLAYWRIGHT_BASE_URL` | Application base URL | No (defaults to localhost:5173) |
| `PLAYWRIGHT_API_URL` | API base URL | No (defaults to localhost:8080) |
| `TEST_USER_EMAIL` | Test user email | Yes |
| `TEST_USER_PASSWORD` | Test user password | Yes |
| `FAMILY_OWNER_EMAIL` | Family owner test user email | No |
| `FAMILY_MEMBER_EMAIL` | Family member test user email | No |
| `CI` | CI environment flag | No |
| `NODE_ENV` | Node environment | No |

## Best Practices

1. **Test Isolation**: Each test should be independent and not rely on other tests
2. **Page Objects**: Use page object pattern for maintainable test code
3. **Reliable Selectors**: Use data-testid attributes for stable element selection
4. **Proper Waits**: Use Playwright's built-in waiting mechanisms
5. **Test Data**: Generate unique test data to avoid conflicts
6. **Cleanup**: Clean up test data after each test
7. **Error Handling**: Handle expected errors and edge cases

## Troubleshooting

### Common Issues

1. **Tests failing locally**:
   - Ensure development server is running (`npm run dev`)
   - Check environment variables in `.env.test.local`
   - Verify test user exists in Clerk

2. **Authentication issues**:
   - Verify Clerk configuration
   - Check test user credentials
   - Ensure Clerk is properly loaded in the application

3. **Timeout errors**:
   - Increase timeouts in `playwright.config.ts`
   - Check network conditions
   - Verify application performance

4. **Flaky tests**:
   - Use proper wait conditions
   - Avoid fixed delays (`page.waitForTimeout`)
   - Ensure test data isolation

### Getting Help

- Check the [Playwright documentation](https://playwright.dev/)
- Review test logs and artifacts in `test-results/`
- Use Playwright's trace viewer for debugging
- Check the HTML report for detailed test results