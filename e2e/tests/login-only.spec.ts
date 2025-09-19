import { test, expect } from '@playwright/test';
import { AuthPage } from '../utils/AuthPage';

// Test configuration
const TEST_EMAIL = process.env.TEST_EMAIL || 'test@example.com';
const TEST_PASSWORD = process.env.TEST_PASSWORD || 'testpassword123';

test.describe('Login Test', () => {
  test.beforeEach(async ({ page }) => {
    // Clear all browser data for a fresh start
    await page.context().clearCookies();
    await page.context().clearPermissions();
    
    console.log('Cleared all browser state - starting fresh test...');
  });

  test('login flow test', async ({ page }) => {
    const authPage = new AuthPage(page);

    console.log('Going directly to login page...');
    
    // Go directly to login page
    await page.goto('/login');
    
    // Check where we landed
    const currentUrl = page.url();
    console.log('Current URL after going to /login:', currentUrl);
    
    // If we're redirected to dashboard, we're already logged in
    if (currentUrl.includes('/dashboard')) {
      console.log('Already logged in - redirected to dashboard');
      return;
    }
    
    // If we're on login page, perform login
    if (currentUrl.includes('/login')) {
      console.log('On login page - performing login...');
      
      // Try direct form filling (like in super-simple test)
      console.log('Filling email field...');
      await page.fill('input[type="email"]', TEST_EMAIL);
      
      console.log('Filling password field...');
      await page.fill('input[type="password"]', TEST_PASSWORD);
      
      console.log('Clicking submit button...');
      await page.click('button[type="submit"]');
      
      // Wait for redirect to dashboard
      console.log('Waiting for redirect to dashboard...');
      await page.waitForURL('**/dashboard', { timeout: 10000 });
      await expect(page).toHaveURL(/.*dashboard/);
      
      console.log('Login successful - now on dashboard');
      return;
    }
    
    throw new Error(`Unexpected location after going to /login: ${currentUrl}`);
  });
});