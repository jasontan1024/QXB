# QXB API 文档

本文档详细说明了 QXB 代币项目的所有 API 端点。

**Base URL**: `http://localhost:8080`

所有 API 响应均为 JSON 格式，统一使用以下结构：

```json
{
  "success": true,
  "data": { ... },
  "error": ""
}
```

## 基础端点

### 健康检查

- **请求方法**: `GET`
- **请求路径**: `/health`
- **需要认证**: 否

**响应示例：**
```json
{
  "success": true,
  "data": {
    "status": "ok",
    "service": "QXB API"
  }
}
```

### API 文档

- **请求方法**: `GET`
- **请求路径**: `/api/docs`
- **需要认证**: 否

返回 API 文档的文本格式说明。

## 代币相关

### 查询代币信息

- **请求方法**: `GET`
- **请求路径**: `/api/token/info`
- **需要认证**: 否
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

**使用示例：**
```bash
curl http://localhost:8080/api/token/info
```

### 查询代币余额

- **请求方法**: `GET`
- **请求路径**: `/api/token/balance/<地址>`
- **需要认证**: 否
- **路径参数**: 
  - `<地址>`: 以太坊地址（十六进制字符串，可以带或不带0x前缀）

**响应示例：**
```json
{
  "success": true,
  "data": {
    "address": "0x...",
    "balance": "1000000000000000000",
    "symbol": "QXB"
  },
  "error": ""
}
```

**使用示例：**
```bash
curl http://localhost:8080/api/token/balance/0x你的地址
```

### 转账代币

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

**注意事项：**
- 转账前会自动检查并补充 ETH 余额（如果余额不足）
- 不能转账给自己
- 需要确保账户有足够的代币余额和 ETH（用于支付 Gas）

## 每日奖励相关

### 查询奖励状态

- **请求方法**: `GET`
- **请求路径**: `/api/reward/status/<地址>`
- **需要认证**: 否
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

### 领取每日奖励

- **请求方法**: `POST`
- **请求路径**: `/api/reward/claim`
- **Content-Type**: `application/json`
- **需要认证**: 可选（如果使用 password 方式，需要 JWT token）

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

**注意事项：**
- 每个地址每天只能领取一次奖励（1 QXB）
- 领取前会自动检查并补充 ETH 余额（如果余额不足）
- 合约地址已在配置文件中固定（`internal/config/config.go`），无需在 API 请求中传入

## 认证相关

### 用户注册

- **请求方法**: `POST`
- **请求路径**: `/api/auth/register`
- **Content-Type**: `application/json`
- **需要认证**: 否

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
- 如果邮箱已被注册，返回 409 错误

**使用示例：**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your_password"
  }'
```

### 用户登录

- **请求方法**: `POST`
- **请求路径**: `/api/auth/login`
- **Content-Type**: `application/json`
- **需要认证**: 否

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

**使用示例：**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "your_password"
  }'
```

**错误响应：**
```json
{
  "success": false,
  "error": "邮箱或密码错误"
}
```

### 获取当前用户信息

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

**使用示例：**
```bash
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer <你的JWT令牌>"
```

**错误响应：**
```json
{
  "success": false,
  "error": "无效或过期的令牌"
}
```

## 其他端点

### 获取作者简历

- **请求方法**: `GET`
- **请求路径**: `/api/resume`
- **需要认证**: 否

**响应示例：**
```json
{
  "success": true,
  "data": {
    "content": "# 作者简历\n\n..."
  },
  "error": ""
}
```

**说明**：返回合约中存储的作者简历（Markdown 格式）

## 认证说明

### JWT Token 使用

大部分需要认证的 API 需要在请求头中提供 JWT token：

```
Authorization: Bearer <你的JWT令牌>
```

**Token 获取方式：**
- 注册用户时自动返回 token
- 登录时返回 token
- Token 有效期为 24 小时

**Token 格式：**
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJleHAiOjE2ODk5OTk5OTl9.signature
```

## 错误处理

### 常见错误码

- `400 Bad Request`: 请求参数错误
- `401 Unauthorized`: 未认证或 token 无效
- `403 Forbidden`: 权限不足
- `404 Not Found`: 资源不存在
- `409 Conflict`: 资源冲突（如邮箱已注册）
- `500 Internal Server Error`: 服务器内部错误

### 错误响应格式

```json
{
  "success": false,
  "error": "错误描述信息"
}
```

## 查看交易记录

### Etherscan 区块链浏览器

推荐使用 Etherscan 查看所有链上交易：

- **合约地址**: https://sepolia.etherscan.io/address/0x5068a014aC8e691Be53848FE5872cbA9f8C4dA17

可以查看：
- **Transactions（交易）**：所有合约调用记录，包括转账、授权、领取奖励等
- **Events（事件日志）**：Transfer、Approval、DailyRewardClaimed 等事件
- **Token Holders（代币持有者）**：所有持有 QXB 的地址及余额
- **Contract（合约）**：合约代码、ABI、验证状态

点击任意交易哈希可以查看详细信息，包括：
- Gas 费用
- 交易状态（成功/失败）
- 事件日志详情
- 输入数据解码

### MetaMask

在 MetaMask 的"活动"标签页查看你的交易历史，点击交易可以跳转到 Etherscan 查看详情。

## 注意事项

1. **合约地址**：合约地址已在配置文件中固定（`internal/config/config.go`），无需在 API 请求中传入
2. **ETH 余额**：转账和领取奖励需要 ETH 支付 Gas 费用，系统会自动检查并补充 ETH（如果配置了 PRIVATE_KEY）
3. **交易确认**：链上交易需要等待确认，可能需要几秒到几分钟
4. **测试网限制**：Sepolia 测试网可能有速率限制
5. **私钥安全**：使用存储私钥方式时，密码不会发送到服务器，仅在服务器端用于解密私钥
