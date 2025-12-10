import { Page, expect } from '@playwright/test';

const API_BASE_URL = process.env.REACT_APP_API_URL || 'http://localhost:8080';

export async function registerUser(page: Page, email: string, password: string) {
  // 监听所有网络请求以便调试
  page.on('response', (response) => {
    const url = response.url();
    if (url.includes('/api/')) {
      console.log(`[API Response] ${url} - Status: ${response.status()}`);
      if (response.status() !== 200) {
        response.text().then(text => console.log(`[API Error] ${text}`)).catch(() => {});
      }
    }
  });
  
  page.on('requestfailed', (request) => {
    console.log(`[Request Failed] ${request.url()} - ${request.failure()?.errorText}`);
  });
  
  page.on('console', (msg) => {
    if (msg.type() === 'error') {
      console.log(`[Console Error] ${msg.text()}`);
    }
  });
  
  await page.goto('/register', { waitUntil: 'domcontentloaded' });
  await page.waitForSelector('input[type="email"]', { timeout: 5000 });
  
  await page.fill('input[type="email"]', email);
  
  // 填写第一个密码输入框
  const passwordInputs = page.locator('input[type="password"]');
  const count = await passwordInputs.count();
  console.log(`[Debug] Found ${count} password inputs`);
  
  if (count >= 1) {
    await passwordInputs.nth(0).fill(password);
  }
  if (count >= 2) {
    await passwordInputs.nth(1).fill(password);
  }
  
  // 等待按钮可点击
  const submitButton = page.locator('button[type="submit"]');
  await submitButton.waitFor({ state: 'visible', timeout: 5000 });
  
  // 点击提交按钮
  await submitButton.click();
  
  // 等待页面导航或错误消息
  try {
    await page.waitForURL('/dashboard', { timeout: 8000 });
  } catch (e) {
    // 如果导航失败，检查是否有错误消息
    const errorMsg = await page.locator('.error-message').textContent().catch(() => null);
    if (errorMsg) {
      console.log(`[Registration Error] ${errorMsg}`);
    }
    // 检查当前 URL
    const currentUrl = page.url();
    console.log(`[Current URL] ${currentUrl}`);
    throw e;
  }
}

export async function loginUser(page: Page, email: string, password: string) {
  await page.goto('/login', { waitUntil: 'domcontentloaded' });
  await page.waitForSelector('input[type="email"]', { timeout: 5000 });
  await page.fill('input[type="email"]', email);
  await page.fill('input[type="password"]', password);
  await page.click('button[type="submit"]');
  // 部分情况下网络慢，放宽等待时间
  await page.waitForURL('/dashboard', { timeout: 12000 });
}

export async function getBalance(page: Page): Promise<string> {
  const balanceElement = page.locator('.balance');
  await balanceElement.waitFor({ timeout: 5000 });
  const text = await balanceElement.textContent();
  return text?.trim() || '0';
}

export async function getRewardStatus(page: Page): Promise<{ canClaim: boolean; text: string }> {
  const rewardElement = page.locator('.can-claim, .cannot-claim').first();
  await rewardElement.waitFor({ timeout: 5000 });
  const text = await rewardElement.textContent() || '';
  return {
    canClaim: text.includes('可以领取'),
    text: text.trim(),
  };
}

export async function claimReward(page: Page, password: string) {
  const claimButton = page.locator('.claim-button');
  await claimButton.click();
  
  // 等待密码输入框出现
  await page.waitForSelector('.password-input', { timeout: 2000 });
  await page.fill('.password-input', password);
  
  // 点击确认领取
  await claimButton.click();
  
  // 等待响应（成功或错误消息）
  await Promise.race([
    page.waitForSelector('.error-message', { timeout: 3000 }).catch(() => null),
    page.waitForTimeout(2000),
  ]);
}

export async function transferTokens(
  page: Page,
  toAddress: string,
  amount: string,
  password: string
) {
  // 找到转账表单中的输入框
  const transferForm = page.locator('.info-card').filter({ hasText: '转账' });
  await transferForm.locator('input[placeholder="0x..."]').fill(toAddress);
  await transferForm.locator('input[placeholder="例如: 100"]').fill(amount);
  await transferForm.locator('input[type="password"]').fill(password);
  await transferForm.locator('button[type="submit"]').click();
  
  // 等待响应（成功或错误消息）
  await Promise.race([
    page.waitForSelector('.error-message', { timeout: 3000 }).catch(() => null),
    page.waitForTimeout(2000),
  ]);
}

export async function getAddress(page: Page): Promise<string> {
  const addressElement = page.locator('.address');
  await addressElement.waitFor({ timeout: 5000 });
  const text = await addressElement.textContent();
  return text?.trim() || '';
}

export async function logout(page: Page) {
  await page.click('.logout-button');
  await page.waitForURL('/login', { timeout: 5000 });
}

export async function clearDatabase() {
  // 注意：这需要后端提供清理接口，或者直接删除数据库文件
  // 为了测试，我们可以使用不同的邮箱来避免冲突
}

export function generateTestEmail(prefix: string = 'test'): string {
  const timestamp = Date.now();
  const random = Math.floor(Math.random() * 10000);
  return `${prefix}-${timestamp}-${random}@test.com`;
}

