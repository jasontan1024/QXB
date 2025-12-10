# QXB 代币项目

QXB（齐夏币）是一个基于以太坊 Sepolia 测试网的 ERC20 代币，包含每日奖励机制。

**合约地址**: `0x5068a014aC8e691Be53848FE5872cbA9f8C4dA17`

- [在 Etherscan 查看](https://sepolia.etherscan.io/address/0x5068a014aC8e691Be53848FE5872cbA9f8C4dA17)
## 功能特性

### 合约功能
- ✅ ERC20 标准代币（完全兼容 ERC20 标准）
- ✅ 标准转账功能（transfer, transferFrom）
- ✅ 授权管理（approve, increaseAllowance, decreaseAllowance）
- ✅ 代币铸造（mint，仅合约所有者可调用）
- ✅ 代币销毁（burn，任何持有者都可以销毁自己的代币）
- ✅ 每日奖励领取（每天 1 QXB，基于 UTC 日期）
- ✅ 奖励状态查询（canClaimDailyReward, getClaimDayInfo）

### 工具和 API
- ✅ Go API 服务器（RESTful API）
- ✅ 合约部署工具
- ✅ React 前端应用（用户注册/登录、代币管理、奖励领取、转账）

### 用户认证和钱包管理
- ✅ 用户注册和登录（JWT 认证）
- ✅ 私钥加密存储（使用 Argon2 + AES-GCM）
- ✅ 基于存储私钥的代币操作（无需客户端输入私钥）

## 快速开始

### 0. 编译合约（首次使用或更新合约后）

```bash
# 安装 Foundry（如果未安装）
curl -L https://foundry.paradigm.xyz | bash
foundryup

# 编译合约
forge build
```

### 1. 部署合约

```bash
# 设置私钥
export PRIVATE_KEY=你的私钥

# 部署（需要先编译合约）
go run ./cmd/deploy-direct
```

### 2. 启动 API 服务器

```bash
# 需要先编译合约（forge build）
go run ./cmd/api
```

API 服务运行在 `http://localhost:8080`

### 3. 启动前端应用

```bash
# 进入前端目录
cd web

# 安装依赖（首次运行）
npm install

# 启动开发服务器
npm start
```

前端应用运行在 `http://localhost:3000`

**💡 前端功能：**
- 用户注册：创建账户并自动生成以太坊地址
- 用户登录：使用邮箱和密码登录
- 代币管理：查看余额、奖励状态
- 奖励领取：使用存储的私钥领取每日奖励（需要输入密码）
- 代币转账：向其他地址转账（需要输入密码）

### 4. 运行 E2E 测试

```bash
# 进入前端目录
cd web

# 运行所有 E2E 测试（会自动启动前端服务器）
npm run test:e2e

# 以 UI 模式运行（推荐用于调试）
npm run test:e2e:ui

# 以调试模式运行
npm run test:e2e:debug

# 查看测试报告
npx playwright show-report
```

**💡 E2E 测试覆盖：**
- 用户注册和登录流程
- 余额和奖励状态查看
- 奖励领取功能
- 代币转账功能（用户A向用户B转账）
- 错误场景：错误密码、余额不足、未登录、无效地址等

**📋 手工验收清单：**
详细的验收测试步骤请参考 [MANUAL_TESTING_CHECKLIST.md](MANUAL_TESTING_CHECKLIST.md)

**💡 查看合约信息：**
- 在 Etherscan 上查看合约：https://sepolia.etherscan.io/address/0x5068a014aC8e691Be53848FE5872cbA9f8C4dA17
- 使用 API 查询代币信息：`GET /api/token/info`（合约地址已在配置中固定）

**🪙 添加到 MetaMask：**

1. 打开 MetaMask，点击"添加代币"
2. 切换到"自定义代币"标签页
3. 输入以下信息：
   - **合约地址**: `0x5068a014aC8e691Be53848FE5872cbA9f8C4dA17`
   - **代币符号**: `QXB`
   - **小数精度**: `18`
4. 点击"添加代币"

**⚠️ 注意事项：**

- 确保 MetaMask 已连接到 **Sepolia 测试网**（链 ID: 11155111）
- 如果小数精度字段显示为 0 且是灰色的，尝试：
  - 刷新 MetaMask 后重新添加
  - 确保网络选择为 Sepolia 测试网
  - 手动输入代币符号（QXB）和小数位数（18）
- 如果代币余额显示为 0，检查：
  - 是否在正确的网络（Sepolia 测试网）
  - 代币是否已正确添加（检查代币列表）
  - 尝试刷新 MetaMask

## 代币配置

- **名称**: 齐夏币
- **符号**: QXB
- **小数位数**: 18
- **总供应量**: 2,025 QXB（初始供应量，可通过 mint 增加）
- **合约版本**: 1.0.0
- **Solidity 版本**: 0.8.30

## API 端点

### 基础端点
- `GET /health` - 健康检查
- `GET /api/docs` - API 文档

### 代币相关

#### 查询代币信息
- **请求方法**: `GET`
- **请求路径**: `/api/token/info`
- **参数**: 无（合约地址已在配置中固定）

**响应示例：**
```json
{
  "success": true,
  "data": {
    "name": "齐夏币",
    "symbol": "QXB",
    "decimals": 18,
    "totalSupply": "2025000000000000000000",
    "version": "1.0.0"
  },
  "error": ""
}
```

#### 查询代币余额
- **请求方法**: `GET`
- **请求路径**: `/api/token/balance/<地址>`
- **路径参数**: 
  - `<地址>`: 以太坊地址（十六进制字符串，可以带或不带0x前缀）

**响应示例：**
```json
{
  "success": true,
  "data": {
    "address": "0x...",
    "balance": "1000000000000000000"
  },
  "error": ""
}
```

**使用示例：**
```bash
curl http://localhost:8080/api/token/balance/0x你的地址
```

### 每日奖励相关

#### 查询奖励状态
- **请求方法**: `GET`
- **请求路径**: `/api/reward/status/<地址>`
- **路径参数**: 
  - `<地址>`: 以太坊地址（十六进制字符串，可以带或不带0x前缀）

**响应示例：**
```json
{
  "success": true,
  "data": {
    "address": "0x...",
    "canClaim": true,
    "lastClaimDay": 19701,
    "nextClaimDay": 19702
  },
  "error": ""
}
```

**响应字段说明：**
- `canClaim` (bool): 是否可以领取奖励
- `lastClaimDay` (uint64): 上次领取的日期（UTC 天数，从 1970-01-01 开始计算）
- `nextClaimDay` (uint64): 下次可以领取的日期（UTC 天数）

**使用示例：**
```bash
curl http://localhost:8080/api/reward/status/0x你的地址
```
#### 领取每日奖励
- **请求方法**: `POST`
- **请求路径**: `/api/reward/claim`
- **Content-Type**: `application/json`

**请求体（JSON）：**

方式一：使用存储的私钥（推荐，需要登录）
```json
{
  "password": "你的密码"
}
```

方式二：直接提供私钥（向后兼容）
```json
{
  "privateKey": "你的私钥（十六进制字符串，可以带或不带0x前缀）"
}
```

**参数说明：**
- `password` (string, 可选): 用户密码，用于解密存储的私钥。需要先登录并在请求头中提供 JWT token
- `privateKey` (string, 可选): 用于签名交易的私钥，十六进制格式
  - 可以带 `0x` 前缀，也可以不带
  - 例如：`"0x1234567890abcdef..."` 或 `"1234567890abcdef..."`

**注意**：`password` 和 `privateKey` 至少需要提供一个。

**响应示例：**
```json
{
  "success": true,
  "data": {
    "txHash": "0xabc123...",
    "status": "pending"
  },
  "error": ""
}
```

**使用示例（curl）：**

使用存储的私钥（需要先登录获取 token）：
```bash
curl -X POST http://localhost:8080/api/reward/claim \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <你的JWT令牌>" \
  -d '{"password": "你的密码"}'
```

直接提供私钥：
```bash
curl -X POST http://localhost:8080/api/reward/claim \
  -H "Content-Type: application/json" \
  -d '{"privateKey": "你的私钥"}'
```

**注意**：合约地址已在配置文件中固定（`internal/config/config.go`），无需在 API 请求中传入。

### 认证相关

#### 用户注册
- **请求方法**: `POST`
- **请求路径**: `/api/auth/register`
- **Content-Type**: `application/json`

**请求体（JSON）：**
```json
{
  "email": "user@example.com",
  "password": "your_password"
}
```

**响应示例：**
```json
{
  "success": true,
  "data": {
    "user_id": 1,
    "email": "user@example.com",
    "address": "0x...",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "error": ""
}
```

**说明**：
- 注册时会自动生成以太坊密钥对
- 私钥使用用户密码加密后存储（Argon2 + AES-GCM）
- 返回 JWT token，可用于后续认证

#### 用户登录
- **请求方法**: `POST`
- **请求路径**: `/api/auth/login`
- **Content-Type**: `application/json`

**请求体（JSON）：**
```json
{
  "email": "user@example.com",
  "password": "your_password"
}
```

**响应示例：**
```json
{
  "success": true,
  "data": {
    "user_id": 1,
    "email": "user@example.com",
    "address": "0x...",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  },
  "error": ""
}
```

#### 获取当前用户信息
- **请求方法**: `GET`
- **请求路径**: `/api/auth/me`
- **需要认证**: 是（在请求头中提供 `Authorization: Bearer <token>`）

**响应示例：**
```json
{
  "success": true,
  "data": {
    "user_id": 1,
    "email": "user@example.com",
    "address": "0x..."
  },
  "error": ""
}
```

### 代币转账

#### 转账代币
- **请求方法**: `POST`
- **请求路径**: `/api/token/transfer`
- **Content-Type**: `application/json`
- **需要认证**: 是（在请求头中提供 `Authorization: Bearer <token>`）

**请求体（JSON）：**
```json
{
  "to": "0x接收地址",
  "amount": "1000000000000000000",
  "password": "你的密码"
}
```

**参数说明：**
- `to` (string, 必需): 接收代币的以太坊地址
- `amount` (string, 必需): 转账金额（以 wei 为单位，18 位小数）
  - 例如：`"1000000000000000000"` 表示 1 QXB
- `password` (string, 必需): 用户密码，用于解密存储的私钥

**响应示例：**
```json
{
  "success": true,
  "data": {
    "txHash": "0xabc123...",
    "status": "pending"
  },
  "error": ""
}
```

**错误响应：**
```json
{
  "success": false,
  "data": null,
  "error": "余额不足"
}
```

**使用示例（curl）：**
```bash
curl -X POST http://localhost:8080/api/token/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <你的JWT令牌>" \
  -d '{
    "to": "0x接收地址",
    "amount": "1000000000000000000",
    "password": "你的密码"
  }'
```

**📊 查看调用记录：**

1. **Etherscan 区块链浏览器**（推荐，最全面）：
   - 合约地址：https://sepolia.etherscan.io/address/0x5068a014aC8e691Be53848FE5872cbA9f8C4dA17
   - 可以查看：
     - **Transactions（交易）**：所有合约调用记录，包括转账、授权、领取奖励等
     - **Events（事件日志）**：Transfer、Approval、DailyRewardClaimed 等事件
     - **Token Holders（代币持有者）**：所有持有 QXB 的地址及余额
     - **Contract（合约）**：合约代码、ABI、验证状态
   - 点击任意交易哈希可以查看详细信息，包括：
     - Gas 费用
     - 交易状态（成功/失败）
     - 事件日志详情
     - 输入数据解码

2. **通过 MetaMask**：
   - 在 MetaMask 的"活动"标签页查看你的交易历史
   - 点击交易可以跳转到 Etherscan 查看详情
   - 可以看到你的所有转账、授权等操作

### API 响应格式
所有 API 响应均为 JSON 格式：
```json
{
  "success": true,
  "data": { ... },
  "error": ""
}
```

## 项目结构

```
QXB/
├── cmd/
│   ├── api/              # API 服务器
│   └── deploy-direct/    # 合约部署工具
├── contracts/
│   └── QXB.sol         # 代币合约
└── internal/
    ├── api/              # API 处理逻辑
    ├── blockchain/       # 区块链交互
    ├── contract/         # 合约交互
    └── config/           # 配置管理
```

## 配置

### 环境变量

创建 `.env` 文件：

```
PRIVATE_KEY=你的私钥
```

### 合约地址配置

合约地址在 `internal/config/config.go` 中配置：

```go
const QXBContractAddress = "0x5068a014aC8e691Be53848FE5872cbA9f8C4dA17"
```

如果需要使用不同的合约地址，修改此配置即可。

## 网络

- **测试网**: Sepolia
- **RPC URL**: https://ethereum-sepolia-rpc.publicnode.com
- **链 ID**: 11155111

## 安全说明

### 私钥加密存储

本项目实现了安全的私钥存储机制：

1. **密钥派生**：使用 Argon2 算法从用户密码派生加密密钥
   - 参数：time=1, memory=64KB, threads=1, keyLen=32
   - 每个用户使用独立的随机 salt

2. **私钥加密**：使用 AES-GCM 模式加密私钥
   - 加密密钥由用户密码通过 Argon2 派生
   - 每个私钥使用独立的随机 salt 和 nonce
   - 加密后的私钥以 base64 格式存储在数据库中

3. **密码验证**：使用 Argon2 派生密钥进行密码验证
   - 密码哈希和私钥加密使用不同的 salt
   - 密码错误会导致解密失败，无法访问私钥

4. **数据库存储**：
   - 数据库文件：`data/qxb.db`（SQLite）
   - 存储字段：`enc_priv_key`（加密私钥）、`enc_salt`（加密 salt）、`pass_salt`（密码 salt）
   - 私钥永远不会以明文形式存储或传输

5. **使用流程**：
   - 注册时：生成密钥对 → 使用密码加密私钥 → 存储加密后的私钥
   - 转账/领取时：用户输入密码 → 解密私钥 → 签名交易 → 立即清除内存中的私钥

**⚠️ 安全建议**：
- 使用强密码（至少 12 位，包含大小写字母、数字、特殊字符）
- 定期备份数据库文件
- 生产环境应使用更严格的 Argon2 参数
- 考虑使用硬件安全模块（HSM）或密钥管理服务（KMS）

## 注意事项

1. 部署需要支付 Gas 费用（约 0.001-0.01 ETH）
2. 私钥不要提交到 Git
3. 本项目仅用于学习和测试
4. 数据库文件（`data/qxb.db`）包含加密的私钥，请妥善保管
5. JWT token 默认有效期为 24 小时
