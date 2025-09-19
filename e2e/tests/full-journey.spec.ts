import { test, expect } from '@playwright/test';

// Test configuration
const TEST_EMAIL = process.env.TEST_EMAIL || 'test@example.com';
const TEST_PASSWORD = process.env.TEST_PASSWORD || 'testpassword123';
const TEST_URL = process.env.TEST_URL || 'https://google.com';
const CRAWL_WAIT_TIMEOUT = parseInt(process.env.CRAWL_WAIT_TIMEOUT || '10000');

test.describe('Full User Journey E2E Test', () => {
  // Ensure clean state before each test
  test.beforeEach(async ({ page }) => {
    // Clear all browser data for a fresh start
    await page.context().clearCookies();
    await page.context().clearPermissions();

    console.log('Cleared all browser state - starting fresh test...');
  });

  test('login -> add url -> start crawl -> verify results', async ({ page }) => {

    // Login process
    await page.goto('/login');

    const currentUrl = page.url();

    if (currentUrl.includes('/dashboard')) {
      console.log('Already logged in - redirected to dashboard');
    } else if (currentUrl.includes('/login')) {
      console.log('On login page - performing login...');
      await page.fill('input[type="email"]', TEST_EMAIL);
      await page.fill('input[type="password"]', TEST_PASSWORD);
      await page.click('button[type="submit"]');

      await page.waitForURL('**/dashboard', { timeout: 10000 });
      console.log('Login successful - now on dashboard');
    }

    await expect(page).toHaveURL(/.*dashboard/);

    // Add URL process
    console.log('Adding URL to crawl...');

    // Click the "Add URL" button to open modal/form
    await page.click('button:has-text("Add URL")');

    // Wait a moment for any UI changes
    await page.waitForTimeout(1000);


    // Based on the console output, we can see Input 5 is the URL input with type="url" and placeholder="https://example.com"
    // Let's target it directly
    const urlInputSelector = 'input[type="url"][placeholder="https://example.com"]';

    // Verify it exists and is visible
    const urlInput = page.locator(urlInputSelector);
    await expect(urlInput).toBeVisible();

    // Fill in the URL - use simple approach
    console.log('Filling URL input...');

    // Clear and fill the input
    await urlInput.clear();
    await urlInput.fill(TEST_URL);

    // Submit the form - look for submit button near the URL input
    console.log('Looking for submit button...');


    // Try to find the submit button with more specific selectors
    let submitButton;
    const submitSelectors = [
      'button[type="submit"]:has-text("Add URL")', // Most specific - the submit button
      'button[type="submit"]', // Any submit button
      'button:has-text("Add URL"):not(:first-of-type)', // Second "Add URL" button
      'button:has-text("Add"):last'  // Last button with "Add" text
    ];

    for (const selector of submitSelectors) {
      const count = await page.locator(selector).count();
      if (count > 0) {
        const elements = await page.locator(selector).all();
        for (const element of elements) {
          const isVisible = await element.isVisible();
          const buttonType = await element.getAttribute('type');
          if (isVisible && buttonType === 'submit') {
            submitButton = element;
            console.log(`Found submit button with selector: ${selector}, type: ${buttonType}`);
            break;
          }
        }
        if (submitButton) break;
      }
    }

    if (!submitButton) {
      throw new Error('Could not find visible submit button');
    }

    console.log('Clicking submit button...');
    await submitButton.click();

    // Wait a moment for the form submission to process
    await page.waitForTimeout(2000);


    // Wait for the URL to appear in the table
    console.log('Waiting for URL to appear in table...');
    try {
      await page.waitForSelector(`text=${TEST_URL}`, { timeout: 10000 });
      console.log('URL found in page content');
    } catch (error) {
      console.log('URL not found in page content, checking table data directly...');
    }

    // Check if URL appears in the table using simpler approach
    const tableElements = await page.locator('table, tbody, .table').all();

    let urlFoundInTable = false;

    // Try to find the URL text directly in the page
    const urlInPage = await page.locator(`text=${TEST_URL}`).count();
    if (urlInPage > 0) {
      urlFoundInTable = true;
      console.log('URL found in page content');
    }

    expect(urlFoundInTable).toBe(true);
    console.log('URL added successfully');

    // Start crawl process
    console.log('Starting crawl process...');

    // First, let's see what the table row looks like
    const targetUrlRow = page.locator(`tr:has-text("${TEST_URL}")`);
    await expect(targetUrlRow).toBeVisible();

    // Debug: Get all text content in the row
    const rowText = await targetUrlRow.allTextContents();
    console.log('URL row contents:', rowText);

    // First, select the checkbox for the URL row
    console.log('Selecting checkbox for the URL...');
    const checkbox = targetUrlRow.locator('input[type="checkbox"]');
    await expect(checkbox).toBeVisible();
    await checkbox.check();
    console.log('Checkbox selected');

    // Wait a moment for the UI to update and show the start button
    await page.waitForTimeout(1000);


    // The start button might be outside the row, in a toolbar or action area
    const startSelector = 'button:has-text("Start Crawl")';

    let startButton;

    const count = await page.locator(startSelector).count();
    if (count > 0) {
      const elements = await page.locator(startSelector).all();
      for (const element of elements) {
        const isVisible = await element.isVisible();
        if (isVisible) {
          startButton = element;
        }
      }
    }


    if (!startButton) {     
      throw new Error('Could not find start button after selecting checkbox');
    }

    await startButton.click();
    console.log('Start button clicked - waiting for status change...');

    // Verify crawl started (status should change)
    console.log('Waiting for crawl to start...');
    await expect(async () => {
      // Look for status changes in the row containing our URL
      const urlRow = page.locator(`tr:has-text("${TEST_URL}")`);
      await expect(urlRow).toBeVisible();

      // Check if there's any indication that crawling started
      const statusTexts = await urlRow.locator('td').allTextContents();
      console.log('Row contents:', statusTexts);

      const hasStartedStatus = statusTexts.some(text => {
        const cellText = text.toLowerCase();
        return cellText.includes('running');
      });
      expect(hasStartedStatus).toBe(true);
    }).toPass({ timeout: 10000 });
    console.log('Crawl started successfully');

    // Wait for crawl completion
    console.log('Waiting for crawl to complete...');
    await page.waitForTimeout(CRAWL_WAIT_TIMEOUT);

    // Check that some crawl data exists in the table
    console.log('Checking for crawl completion...');
    await expect(async () => {
      const urlRow = page.locator(`tr:has-text("${TEST_URL}")`);
      await expect(urlRow).toBeVisible();

      const statusTexts = await urlRow.locator('td').allTextContents();

      // Check for completion indicators or data (more flexible)
      const hasResults = statusTexts.some(text => {
        const cellText = text.toLowerCase();
        return cellText.includes('done')
      });
      expect(hasResults).toBe(true);
    }).toPass({ timeout: 20000 });

    console.log('Crawl completed successfully');

    // Verify crawl results
    console.log('Verifying crawl results...');
    const urlRow = page.locator(`tr:has-text("${TEST_URL}")`);
    await expect(urlRow).toBeVisible();

    const finalStatusTexts = await urlRow.locator('td').allTextContents();

    // Verify the row contains meaningful data
    expect(finalStatusTexts.length).toBeGreaterThan(1); // Should have multiple columns

    // Check that at least one cell contains non-empty, meaningful data
    const hasNonEmptyData = finalStatusTexts.some(text =>
      text.trim().length > 0 &&
      text !== '-' &&
      text !== 'N/A'
    );
    expect(hasNonEmptyData).toBe(true);

    console.log('E2E test completed successfully! ✅');
  });

  // Cleanup test - remove the test user (optional)
  test.afterEach(async ({ page }) => {
    // Optionally clean up test data
    // This could involve API calls to remove the test user
    // or database cleanup
  });
});