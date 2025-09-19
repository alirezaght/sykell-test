import { test, expect } from '@playwright/test';

test.describe('Minimal Login Test', () => {
  test('minimal navigation test', async ({ page }) => {
    console.log('Starting minimal test...');
    
    // Just try to navigate and see what happens
    console.log('Navigating to /login...');
    await page.goto('/login');
    
    console.log('Taking screenshot...');
    await page.screenshot({ path: 'minimal-test.png' });
    
    console.log('Getting URL...');
    const url = page.url();
    console.log('Current URL:', url);
    
    console.log('Getting page title...');
    const title = await page.title();
    console.log('Page title:', title);
    
    console.log('Test completed successfully');
    expect(url).toContain('localhost');
  });
});