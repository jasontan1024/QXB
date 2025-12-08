# QXB 代币项目

QXB（齐夏币）是一个基于以太坊 Sepolia 测试网的 ERC20 代币，包含每日奖励机制。

**合约地址**: `0xFF96cF72Cc4FCb67C61e0E43924723fA88765A06`

- [在 Etherscan 查看](https://sepolia.etherscan.io/address/0xFF96cF72Cc4FCb67C61e0E43924723fA88765A06)
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

**💡 查看合约信息：**
- 在 Etherscan 上查看合约：https://sepolia.etherscan.io/address/0xFF96cF72Cc4FCb67C61e0E43924723fA88765A06
- 使用 API 查询代币信息：`GET /api/token/info`（合约地址已在配置中固定）

**🪙 添加到 MetaMask：**

1. 打开 MetaMask，点击"添加代币"
2. 切换到"自定义代币"标签页
3. 输入以下信息：
   - **合约地址**: `0xFF96cF72Cc4FCb67C61e0E43924723fA88765A06`
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
- `POST /api/reward/claim` - 领取每日奖励

**领取奖励请求参数：**

请求方法：`POST`  
请求路径：`/api/reward/claim`  
Content-Type: `application/json`

请求体（JSON）：
```json
{
  "privateKey": "你的私钥（十六进制字符串，可以带或不带0x前缀）"
}
```

**参数说明：**
- `privateKey` (string, 必需): 用于签名交易的私钥，十六进制格式
  - 可以带 `0x` 前缀，也可以不带
  - 例如：`"0x1234567890abcdef..."` 或 `"1234567890abcdef..."`

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
  "error": "私钥不能为空"
}
```

**使用示例（curl）：**
```bash
curl -X POST http://localhost:8080/api/reward/claim \
  -H "Content-Type: application/json" \
  -d '{"privateKey": "你的私钥"}'
```

**注意**：合约地址已在配置文件中固定（`internal/config/config.go`），无需在 API 请求中传入。

**📊 查看调用记录：**

1. **Etherscan 区块链浏览器**（推荐，最全面）：
   - 合约地址：https://sepolia.etherscan.io/address/0xFF96cF72Cc4FCb67C61e0E43924723fA88765A06
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
const QXBContractAddress = "0xFF96cF72Cc4FCb67C61e0E43924723fA88765A06"
```

如果需要使用不同的合约地址，修改此配置即可。

## 网络

- **测试网**: Sepolia
- **RPC URL**: https://ethereum-sepolia-rpc.publicnode.com
- **链 ID**: 11155111

## 注意事项

1. 部署需要支付 Gas 费用（约 0.001-0.01 ETH）
2. 私钥不要提交到 Git
3. 本项目仅用于学习和测试
