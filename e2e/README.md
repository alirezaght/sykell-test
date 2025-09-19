# Sykell E2E Tests

This directory contains end-to-end tests for the Sykell application using Playwright.

## Setup

1. Install dependencies:
```bash
npm install
```

2. Install Playwright browsers:
```bash
npm run install-browsers
```

3. Configure environment variables:
Copy `.env` and update the values as needed for your test environment.

4. Ensure your application is running:
- Backend API: http://localhost:8080
- Frontend: http://localhost:3000
- Database and Temporal services should be running

## Running Tests

### Basic test execution:
```bash
# Run all tests
npm test

# Run tests with browser UI visible
npm run test:headed

# Run tests with Playwright UI mode
npm run test:ui

# Debug tests step by step
npm run test:debug
```

### View test reports:
```bash
npm run show-report
```

## Test Structure

### Main Test: `full-journey.spec.ts`
This test covers the complete user workflow:
1. **User Registration** - Sign up with test credentials
2. **User Login** - Log out and log back in to verify login flow
3. **Add URL** - Add a URL (https://google.com) to crawl
4. **Start Crawl** - Initiate the crawling process
5. **Wait & Verify** - Wait 10 seconds and verify results appear in the table

### Page Objects
- `AuthPage.ts` - Handles login and registration interactions
- `DashboardPage.ts` - Handles dashboard operations (add URL, start crawl, view results)

## Configuration

### Environment Variables
- `BASE_URL` - Frontend application URL (default: http://localhost:3000)
- `TEST_EMAIL` - Email for test user (default: test@example.com)
- `TEST_PASSWORD` - Password for test user (default: testpassword123)
- `TEST_URL` - URL to crawl in tests (default: https://google.com)
- `CRAWL_WAIT_TIMEOUT` - Time to wait for crawl completion (default: 10000ms)

### Playwright Configuration
- Tests run against Chrome, Firefox, and Safari
- Screenshots and videos captured on failure
- Automatic retry on CI environments
- HTML reporter for detailed results

## Docker Integration

If running the application with Docker Compose, ensure all services are up:
```bash
cd ..
docker-compose up -d
```

Then run the tests:
```bash
npx playwright install   
npm test
```

## Troubleshooting

1. **Application not responding**: Ensure frontend (port 3000) and backend (port 8080) are running
2. **Database connection issues**: Check that MySQL and Temporal services are running
3. **Test timeout**: Increase `CRAWL_WAIT_TIMEOUT` if your crawling process takes longer
4. **Element not found**: Check if the UI has changed and update the page object selectors

## CI/CD Integration

The tests are configured to run in CI environments with:
- Reduced parallelism for stability
- Automatic retries on failure
- Screenshot and video capture on failures
- HTML report generation