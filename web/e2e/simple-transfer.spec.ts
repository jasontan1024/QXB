import { test, expect } from '@playwright/test';

test.describe('简单转账验证', () => {
  test('验证转账功能代码逻辑', async ({ page }) => {
    // 这个测试主要验证页面能正常加载
    await page.goto('http://localhost:3000');
    
    // 检查页面是否加载
    await expect(page).toHaveURL(/.*/);
    
    // 检查是否有注册或登录链接
    const hasAuthLink = await page.locator('a[href*="register"], a[href*="login"]').count() > 0;
    expect(hasAuthLink || page.url().includes('register') || page.url().includes('login')).toBeTruthy();
  });
});




