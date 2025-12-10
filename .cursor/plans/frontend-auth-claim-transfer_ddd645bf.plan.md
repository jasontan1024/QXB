---
name: frontend-auth-claim-transfer
overview: Add React frontend for register/login, reward claim, and token transfer, backed by new auth APIs.
todos: []
---

# Plan

1) Backend auth + wallet存储

- 新增 SQLite 持久化：users 表保存 email、密码哈希、链上地址、加密后的私钥
- 注册流程：生成密钥对，使用用户密码派生密钥（如 scrypt/argon2 + AES-GCM）加密私钥并存储
- 登录：验证密码，签发 JWT；GET /api/auth/me 返回当前用户基础信息
- 路由：在 [internal/api/server.go](internal/api/server.go) 添加 /api/auth/register、/api/auth/login、/api/auth/me

2) 代币操作使用后端存储的密钥

- 新增 POST /api/token/transfer：从 JWT 获取用户，使用用户密码（或会话内密钥）解密私钥，签名转账交易，使用固定合约地址
- 复用 /api/reward/claim：支持使用存储的私钥完成领取，无需客户端输入私钥
- 定义请求/响应和错误处理（密码错误、余额不足、链上失败等）

3) React 前端

- 创建 web/ React 应用：页面 Register、Login、Dashboard
- JWT 存储 localStorage；需要解密时提示输入密码
- Dashboard：展示余额、奖励状态，按钮调用后端完成领取；转账表单调用后端完成签名和发送

4) UI 流程与鉴权

- 注册：email+password → 返回地址（不返回私钥）
- 登录：获取 JWT，前端路由守卫
- Dashboard：加载 token info/balance/reward status；转账与领取调用后端
- 错误/加载态处理

5) 文档与测试

- 更新 README：前端启动方式，新增 auth/transfer 接口说明，私钥加密存储说明
- 自动化前端 E2E（Playwright/Cypress）：
- 注册用户A → 登录 → 查看余额/奖励 → 领取奖励
- 注册用户B → 登录 → 查看余额
- 用户A 向用户B 转账 → 验证双方余额变化
- 错误场景：错误密码、余额不足、未登录、转账给自己
- 保留关键手工验收清单