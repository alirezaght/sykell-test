# Sykell E2E Tests

End-to-end tests for the Sykell web crawler application using Playwright. These tests verify the complete user journey from signup to crawling and result verification.

## Overview

The E2E tests cover the complete user workflow:
1. **User Authentication** - Login/Register flow
2. **URL Management** - Adding URLs for crawling
3. **Crawl Execution** - Starting crawl processes
4. **Real-time Updates** - Verifying live status updates
5. **Result Verification** - Checking crawl results in the dashboard

## Technology Stack

- **Playwright** - Modern E2E testing framework
- **TypeScript** - Type-safe test development
- **Node.js** - Runtime environment

## Prerequisites

- **Node.js 18+** (recommended: Node.js 20 LTS)
- **npm** or **yarn**
- **Running Sykell Application** (backend + frontend + database)

## Quick Start

### 1. Install Dependencies

```bash
cd e2e/

# Install dependencies
npm install

# Install Playwright browsers
npm run install-browsers
```

### 2. Environment Configuration

Configure the test environment in `.env`:

```bash
# Base URL for the frontend application
BASE_URL=http://localhost

# Test user credentials (must exist in the system)
TEST_EMAIL=testuser@example.com
TEST_PASSWORD=123456

# Test URLs for crawling
TEST_URL=https://google.com

# Timeouts (in milliseconds)
CRAWL_WAIT_TIMEOUT=10000
DEFAULT_TIMEOUT=5000

# CI flag
CI=false
```

### 3. Start the Application

**Before running tests, ensure the full application is running:**

```bash
# From the project root
docker compose up

# Or start individual services:
# - Frontend: http://localhost (port 80)
# - Backend: http://localhost:7070
# - Database: MySQL on port 3306
```

### 4. Run the Tests

```bash
# Run all tests (headless mode)
npm test

# Run with browser visible (headed mode)
npm run test:headed

# Run with Playwright UI (interactive mode)
npm run test:ui

# Run in debug mode (step-by-step debugging)
npm run test:debug
```

## Available Commands

| Command | Description |
|---------|-------------|
| `npm test` | Run all tests in headless mode |
| `npm run test:headed` | Run tests with browser visible |
| `npm run test:ui` | Open Playwright UI for interactive testing |
| `npm run test:debug` | Run tests in debug mode with step-by-step execution |
| `npm run show-report` | Open the HTML test report |
| `npm run install-browsers` | Install Playwright browser dependencies |

## Test Configuration

### Playwright Configuration (`playwright.config.ts`)

- **Test Timeout**: 60 seconds per test
- **Parallel Execution**: Tests run in parallel for faster execution
- **Browser Support**: Currently configured for Chromium (can enable Firefox/Safari)
- **Screenshots**: Captured on test failure
- **Videos**: Recorded for failed tests
- **Traces**: Generated on retry for debugging

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `BASE_URL` | Frontend application URL | `http://localhost` | Yes |
| `TEST_EMAIL` | Test user email (must exist) | - | Yes |
| `TEST_PASSWORD` | Test user password | - | Yes |
| `TEST_URL` | URL to test crawling | `https://google.com` | Yes |
| `CRAWL_WAIT_TIMEOUT` | Max wait time for crawl completion | `10000` | No |
| `DEFAULT_TIMEOUT` | Default action timeout | `5000` | No |
| `CI` | CI environment flag | `false` | No |

## Test Structure

### Current Test Suite

```
tests/
└── full-journey.spec.ts    # Complete user workflow test
```

### Full Journey Test

The main test covers the complete user workflow:

1. **Authentication**
   - Navigate to login page
   - Enter credentials
   - Verify successful login and redirect to dashboard

2. **URL Management**
   - Click "Add URL" button
   - Fill in URL form
   - Submit and verify URL appears in the list

3. **Crawl Execution**
   - Start crawl process
   - Monitor real-time status updates
   - Wait for completion

4. **Result Verification**
   - Check crawl results in the dashboard
   - Verify data accuracy and completeness

## Test Data Requirements

### User Account

The tests require a pre-existing user account in the system:

```bash
Email: testuser@example.com
Password: 123456
```

**To create a test user:**
1. Start the application
2. Navigate to the registration page
3. Create an account with the test credentials

### Test URLs

The tests use configurable URLs for crawling:
- Default: `https://google.com`
- Configurable via `TEST_URL` environment variable
- Should be publicly accessible URLs

## Running Tests in Different Modes

### 1. Headless Mode (Default)

```bash
npm test
```

- Fastest execution
- No browser window visible
- Suitable for CI/CD pipelines
- Screenshots and videos captured on failure

### 2. Headed Mode (Browser Visible)

```bash
npm run test:headed
```

- Browser window visible during test execution
- Useful for watching test execution
- Good for debugging test flow

### 3. UI Mode (Interactive)

```bash
npm run test:ui
```

- Opens Playwright Test UI
- Interactive test execution
- Step-by-step debugging
- Test result exploration
- Time travel debugging

### 4. Debug Mode

```bash
npm run test:debug
```

- Pauses execution at each step
- Allows inspector-style debugging
- Browser stays open for manual inspection
- Console output for debugging

## Test Reports and Artifacts

### HTML Report

After running tests, view the detailed report:

```bash
npm run show-report
```

The report includes:
- Test execution summary
- Screenshots of failures
- Video recordings
- Execution traces
- Performance metrics

### Artifacts Location

```
e2e/
├── test-results/        # Test execution artifacts
│   ├── screenshots/     # Failure screenshots
│   ├── videos/         # Test execution videos
│   └── traces/         # Execution traces
├── playwright-report/   # HTML test reports
└── *.png              # Debug screenshots
```

### Debugging Failed Tests

When tests fail, check these artifacts:

1. **Screenshots**: Visual state at failure point
2. **Videos**: Complete test execution recording
3. **Traces**: Interactive debugging with time travel
4. **Console Logs**: Application and test output

## CI/CD Integration

### GitHub Actions Example

```yaml
- name: Run E2E Tests
  run: |
    cd e2e
    npm ci
    npm run install-browsers
    npm test
  env:
    CI: true
    BASE_URL: http://localhost
    TEST_EMAIL: test@example.com
    TEST_PASSWORD: testpassword123
```

### CI Configuration

Set `CI=true` environment variable for CI environments:
- Enables retries on failure
- Optimizes for CI execution
- Reduces parallel workers