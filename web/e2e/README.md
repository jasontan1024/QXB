# E2E 测试说明

本目录包含使用 Playwright 编写的端到端测试。

## 运行测试

### 前置条件

1. 确保后端 API 服务器正在运行（`http://localhost:8080`）
2. 确保前端开发服务器可以启动（测试会自动启动）

### 运行所有测试

```bash
npm run test:e2e
```

### 运行特定测试文件

```bash
npx playwright test e2e/user-flow.spec.ts
npx playwright test e2e/error-scenarios.spec.ts
```

### 以 UI 模式运行（推荐用于调试）

```bash
npx playwright test --ui
```

### 查看测试报告

```bash
npx playwright show-report
```

## 测试场景

### user-flow.spec.ts

1. **场景1**: 注册用户A → 登录 → 查看余额/奖励 → 领取奖励
2. **场景2**: 注册用户B → 登录 → 查看余额
3. **场景3**: 用户A 向用户B 转账 → 验证双方余额变化

### error-scenarios.spec.ts

1. **错误密码登录**: 验证使用错误密码无法登录
2. **余额不足转账**: 验证余额不足时转账失败
3. **未登录访问**: 验证未登录时访问 Dashboard 会被重定向
4. **转账给自己**: 验证转账给自己时的处理
5. **无效地址转账**: 验证使用无效地址时转账失败
6. **错误密码领取奖励**: 验证使用错误密码无法领取奖励

## 注意事项

1. 测试使用随机生成的邮箱地址，避免测试之间的冲突
2. 链上交易需要等待确认，测试中使用了适当的等待时间
3. 某些测试场景（如余额变化验证）可能需要等待链上交易确认
4. 测试数据库会在每次测试后保留数据，实际使用中可能需要清理

## 调试

如果测试失败，可以：

1. 使用 `--debug` 模式运行：
   ```bash
   npx playwright test --debug
   ```

2. 查看测试截图和视频（在 `test-results/` 目录）

3. 使用 Playwright Inspector：
   ```bash
   PWDEBUG=1 npx playwright test
   ```




