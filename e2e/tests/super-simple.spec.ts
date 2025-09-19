import { test, expect } from '@playwright/test';

test.describe('Super Simple Login Test', () => {
  test('just try to fill form', async ({ page }) => {
    console.log('Starting super simple test...');
    
    // Go to login page
    await page.goto('/login');
    console.log('Navigated to login page');
    
    // Try to fill email immediately without any waits
    console.log('Attempting to fill email...');
    await page.fill('input[type="email"]', 'alirezaght@gmail.com');
    console.log('Email filled successfully');
    
    // Try to fill password
    console.log('Attempting to fill password...');
    await page.fill('input[type="password"]', '123456');
    console.log('Password filled successfully');
    
    // Try to click submit
    console.log('Attempting to click submit...');
    await page.click('button[type="submit"]');
    console.log('Submit clicked successfully');
    
    // Take a screenshot to see what happened
    await page.screenshot({ path: 'after-login-attempt.png' });
    console.log('Screenshot taken');
    
    console.log('Test completed');
  });
});