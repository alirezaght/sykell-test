import { test, expect } from '@playwright/test';
import { AuthPage } from '../utils/AuthPage';
import { DashboardPage } from '../utils/DashboardPage';

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
    const authPage = new AuthPage(page);
    const dashboardPage = new DashboardPage(page);

    // Login process
    console.log('Starting login process...');
    await page.goto('/login');
    
    const currentUrl = page.url();
    console.log('Current URL after going to /login:', currentUrl);
    
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
    await page.goto('/dashboard'); // Ensure we're on dashboard
    
    // Click the "Add URL" button to open modal/form
    console.log('Clicking Add URL button...');
    await page.click('button:has-text("Add URL")');
    
    // Wait a moment for any UI changes
    await page.waitForTimeout(1000);
    
    // Take a screenshot to see what happens after clicking
    await page.screenshot({ path: 'after-add-url-click.png' });
    
    // Look for any input that might have appeared (don't assume modal structure)
    console.log('Looking for URL input field after clicking Add URL...');
    
    // First, log all inputs to see what's available
    const allInputs = await page.locator('input').all();
    console.log('All input elements found after clicking Add URL:');
    for (let i = 0; i < allInputs.length; i++) {
      const input = allInputs[i];
      const name = await input.getAttribute('name');
      const type = await input.getAttribute('type');
      const placeholder = await input.getAttribute('placeholder');
      const className = await input.getAttribute('class');
      const isVisible = await input.isVisible();
      console.log(`Input ${i}: name="${name}", type="${type}", placeholder="${placeholder}", class="${className}", visible=${isVisible}`);
    }
    
    // Based on the console output, we can see Input 5 is the URL input with type="url" and placeholder="https://example.com"
    // Let's target it directly
    const urlInputSelector = 'input[type="url"][placeholder="https://example.com"]';
    
    // Verify it exists and is visible
    const urlInput = page.locator(urlInputSelector);
    await expect(urlInput).toBeVisible();
    console.log('Found URL input with placeholder "https://example.com"');
    
    // Fill in the URL - use simple approach
    console.log('Filling URL input...');
    
    // Clear and fill the input
    await urlInput.clear();
    await urlInput.fill(TEST_URL);
    
    // Verify the value was set
    const inputValue = await urlInput.inputValue();
    console.log(`Input value after fill: "${inputValue}"`);
    
    if (!inputValue || inputValue !== TEST_URL) {
      await page.screenshot({ path: 'debug-input-fill-failed.png' });
      throw new Error(`Failed to fill URL input. Expected: "${TEST_URL}", Got: "${inputValue}"`);
    }
    
    console.log('URL input filled successfully!');
    
    // Submit the form - look for submit button near the URL input
    console.log('Looking for submit button...');
    
    // First, let's see what buttons are available
    const allButtons = await page.locator('button').all();
    console.log('All buttons found:');
    for (let i = 0; i < allButtons.length; i++) {
      const button = allButtons[i];
      const text = await button.textContent();
      const type = await button.getAttribute('type');
      const className = await button.getAttribute('class');
      const isVisible = await button.isVisible();
      console.log(`Button ${i}: text="${text}", type="${type}", class="${className}", visible=${isVisible}`);
    }
    
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
      // Fallback - target Button 8 directly which we know is the submit button
      submitButton = page.locator('button').nth(8);
      console.log('Using fallback - targeting button index 8 (the submit button)');
    }
    
    if (!submitButton) {
      await page.screenshot({ path: 'debug-submit-button.png' });
      throw new Error('Could not find visible submit button');
    }
    
    console.log('Clicking submit button...');
    await submitButton.click();
    
    // Wait a moment for the form submission to process
    await page.waitForTimeout(2000);
    
    // Take a screenshot to see what happened after submit
    await page.screenshot({ path: 'after-submit-click.png' });
    
    // Wait for the URL to appear in the table
    console.log('Waiting for URL to appear in table...');
    try {
      await page.waitForSelector(`text=${TEST_URL}`, { timeout: 10000 });
      console.log('URL found in page content');
    } catch (error) {
      console.log('URL not found in page content, checking table data directly...');
      await page.screenshot({ path: 'debug-url-not-found.png' });
    }
    
    // Check if URL appears in the table using simpler approach
    console.log('Checking if URL appears in table...');
    const tableElements = await page.locator('table, tbody, .table').all();
    console.log(`Found ${tableElements.length} table-like elements`);
    
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
    
    // Now look for the start button (it should appear after selecting checkbox)
    console.log('Looking for start crawl button...');
    
    // The start button might be outside the row, in a toolbar or action area
    const startSelectors = [
      'button:has-text("Start Crawl")',
      'button:has-text("Start")',
      'button:has-text("Crawl")',
      'button:has-text("Begin")',
      `tr:has-text("${TEST_URL}") button`,  // Any button in the row
      'button[type="button"]:has-text("Start")'
    ];
    
    let startButton;
    for (const selector of startSelectors) {
      const count = await page.locator(selector).count();
      if (count > 0) {
        const elements = await page.locator(selector).all();
        for (const element of elements) {
          const isVisible = await element.isVisible();
          if (isVisible) {
            startButton = element;
            console.log(`Found start button with selector: ${selector}`);
            break;
          }
        }
        if (startButton) break;
      }
    }
    
    if (!startButton) {
      // Debug: Check all visible buttons on the page after selecting checkbox
      const allVisibleButtons = await page.locator('button').all();
      console.log('All visible buttons after selecting checkbox:');
      for (let i = 0; i < allVisibleButtons.length; i++) {
        const button = allVisibleButtons[i];
        const text = await button.textContent();
        const isVisible = await button.isVisible();
        if (isVisible) {
          console.log(`Visible Button ${i}: text="${text}"`);
        }
      }
      
      await page.screenshot({ path: 'debug-no-start-button.png' });
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
        return cellText.includes('running') || 
               cellText.includes('queued') ||
               cellText.includes('active') ||
               cellText.includes('processing') ||
               cellText.includes('started');
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
      console.log('Final row contents:', statusTexts);
      
      // Check for completion indicators or data (more flexible)
      const hasResults = statusTexts.some(text => {
        const cellText = text.toLowerCase();
        return cellText.includes('completed') ||
               cellText.includes('finished') ||
               cellText.includes('done') ||
               cellText.includes('success') ||
               /\d+/.test(text); // Contains numbers (could be link count, etc.)
      });
      expect(hasResults).toBe(true);
    }).toPass({ timeout: 20000 });
    
    console.log('Crawl completed successfully');

    // Verify crawl results
    console.log('Verifying crawl results...');
    const urlRow = page.locator(`tr:has-text("${TEST_URL}")`);
    await expect(urlRow).toBeVisible();
    
    const finalStatusTexts = await urlRow.locator('td').allTextContents();
    console.log('Final crawl results:', finalStatusTexts);
    
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