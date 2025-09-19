import { Page, Locator } from '@playwright/test';

export class DashboardPage {
  readonly page: Page;
  readonly addUrlButton: Locator;
  readonly urlInput: Locator;
  readonly addButton: Locator;
  readonly urlTable: Locator;
  readonly startCrawlButton: Locator;
  readonly stopCrawlButton: Locator;
  readonly urlTableRows: Locator;
  readonly statusColumn: Locator;
  readonly logoutButton: Locator;
  readonly profileButton: Locator;

  constructor(page: Page) {
    this.page = page;
    this.addUrlButton = page.getByText('Add URL');
    this.urlInput = page.locator('input[placeholder*="URL"], input[name="url"]');
    this.addButton = page.getByRole('button', { name: 'Add' });
    this.urlTable = page.locator('table, [data-testid="url-table"]');
    this.urlTableRows = page.locator('table tbody tr, [data-testid="url-table"] tbody tr');
    this.startCrawlButton = page.getByRole('button', { name: /start/i });
    this.stopCrawlButton = page.getByRole('button', { name: /stop/i });
    this.statusColumn = page.locator('td:has-text("Status"), [data-testid="status"]');
    this.logoutButton = page.getByText('Logout');
    this.profileButton = page.getByText('Profile');
  }

  async goto() {
    await this.page.goto('/dashboard');
  }

  async addUrl(url: string) {
    // Open add URL modal/form by clicking the first "Add URL" button
    await this.addUrlButton.click();
    
    // Wait for the modal/form to appear
    await this.page.waitForSelector('input[placeholder*="URL"], input[name="url"]', { timeout: 5000 });
    
    // Fill in the URL
    await this.urlInput.fill(url);
    
    // Submit the form using the form submit button (more specific selector)
    await this.page.click('form button[type="submit"]');
    
    // Wait for the URL to appear in the table
    await this.page.waitForSelector(`text=${url}`, { timeout: 5000 });
  }

  async startCrawl(url: string) {
    // Find the row containing the URL and click start button
    const row = this.page.locator(`tr:has-text("${url}")`);
    const startButton = row.locator('button:has-text("Start"), button:has-text("start")').first();
    await startButton.click();
  }

  async waitForCrawlStatus(url: string, status: string, timeout: number = 10000) {
    // Wait for the crawl status to change to the expected status
    const row = this.page.locator(`tr:has-text("${url}")`);
    await row.locator(`text=${status}`).waitFor({ timeout });
  }

  async getUrlTableData() {
    // Wait for table to be visible
    await this.urlTable.waitFor();
    
    // Extract table data
    const rows = await this.urlTableRows.count();
    const data = [];
    
    for (let i = 0; i < rows; i++) {
      const row = this.urlTableRows.nth(i);
      const cells = row.locator('td');
      const cellCount = await cells.count();
      
      const rowData: string[] = [];
      for (let j = 0; j < cellCount; j++) {
        const cellText = await cells.nth(j).textContent();
        rowData.push(cellText?.trim() || '');
      }
      data.push(rowData);
    }
    
    return data;
  }

  async logout() {
    await this.logoutButton.click();
    await this.page.waitForURL('**/login');
  }
}