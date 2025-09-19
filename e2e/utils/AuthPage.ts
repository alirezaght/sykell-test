import { Page, Locator } from '@playwright/test';

export class AuthPage {
  readonly page: Page;
  readonly emailInput: Locator;
  readonly passwordInput: Locator;
  readonly confirmPasswordInput: Locator;
  readonly submitButton: Locator;
  readonly loginLink: Locator;
  readonly registerLink: Locator;
  readonly errorMessage: Locator;

  constructor(page: Page) {
    this.page = page;
    this.emailInput = page.locator('input[type="email"]');
    this.passwordInput = page.locator('input[type="password"]').first();
    this.confirmPasswordInput = page.locator('input[type="password"]').nth(1);
    this.submitButton = page.locator('button[type="submit"]');
    this.loginLink = page.getByText('Login');
    this.registerLink = page.getByText('Sign up');
    this.errorMessage = page.locator('[role="alert"], .error-message, .text-red-500');
  }

  async goto(path: 'login' | 'register') {
    await this.page.goto(`/${path}`);
  }

  async login(email: string, password: string) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.submitButton.click();
  }

  async register(email: string, password: string) {
    await this.emailInput.fill(email);
    await this.passwordInput.fill(password);
    await this.confirmPasswordInput.fill(password);
    await this.submitButton.click();
  }

  async waitForRedirect() {
    // Wait for navigation away from login/register page
    await this.page.waitForURL('**/dashboard');
  }
}