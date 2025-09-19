import { test, expect } from '@playwright/test';

test.describe('Connectivity Test', () => {
  test('basic connectivity check', async ({ page }) => {
    console.log('Testing basic connectivity...');
    
    // Try to connect to the base URL
    const response = await page.goto('/');
    console.log('Response status:', response?.status());
    console.log('Current URL:', page.url());
    
    // Wait a moment
    await page.waitForTimeout(2000);
    
    // Check if we can see any content
    const title = await page.title();
    console.log('Page title:', title);
    
    // Take a screenshot for debugging
    await page.screenshot({ path: 'connectivity-test.png' });
    
    // Basic assertion
    expect(response?.status()).toBeLessThan(400);
  });
});