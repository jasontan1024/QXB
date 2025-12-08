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

**⚠️ 重要：必须使用方法 1（自动添加工具）**

MetaMask 在手动添加时可能无法正确读取 decimals（小数精度），导致：
- 小数精度字段显示为 0 且变为灰色（不可编辑）
- 代币余额显示为 0（即使实际有余额）

**方法 1：使用自动添加工具（强烈推荐）⭐⭐⭐**

**⚠️ 重要：不能直接用 `file://` 协议打开，需要使用 HTTP 服务器**

1. **启动本地 HTTP 服务器**（在项目根目录执行）：
   ```bash
   # 使用 Python（推荐）
   python3 -m http.server 8000
   
   # 或者使用 Go（如果已安装）
   go run -m http.server 8000
   ```

2. **在浏览器中访问**：
   ```
   http://localhost:8000/add-token.html
   ```
   （不要使用 `file://` 协议打开，MetaMask 在 file:// 协议下可能无法正常工作）

3. 确保 MetaMask 已连接到 **Sepolia 测试网**

4. 点击"添加到 MetaMask"按钮

5. 在 MetaMask 弹窗中确认添加

6. ✅ 这会自动设置正确的 decimals（18），余额会正确显示

**方法 2：手动添加（不推荐，可能有问题）**
如果必须手动添加：
1. 打开 MetaMask，点击"添加代币"
2. 切换到"自定义代币"标签页
3. 输入合约地址：`0xFF96cF72Cc4FCb67C61e0E43924723fA88765A06`
4. **如果小数精度字段显示为 0 且是灰色的**：
   - 删除代币，改用方法 1（自动添加工具）
   - 或者尝试刷新 MetaMask 后重新添加

**⚠️ 如果代币显示余额为 0 或小数精度为灰色：**
- 这通常是因为 MetaMask 无法从链上读取 decimals，导致显示为 0 且字段变为灰色（不可编辑）
- **解决方法**：
  1. 在 MetaMask 中删除该代币（点击代币右侧的菜单，选择"隐藏代币"）
  2. **使用 `add-token.html` 文件重新添加**（这是最可靠的方法）
  3. `add-token.html` 会通过 MetaMask API 直接设置正确的 decimals（18），绕过 MetaMask 的自动读取问题

**⚠️ 常见问题：**

1. **代币已添加但列表中不显示？**
   - 检查 MetaMask 是否在 **Sepolia 测试网**（不是主网）
   - 检查代币余额：如果余额为 0，MetaMask 可能默认隐藏，需要：
     - 点击代币列表右上角的"导入代币"或"管理代币"
     - 或者滚动到底部查看隐藏的代币
   - 尝试刷新 MetaMask 或重新打开扩展
   - 确认代币网络：代币必须添加到 Sepolia 测试网，不能添加到主网

2. **"下一步"按钮是灰色的？**
   - 必须使用"自定义代币"标签页（不要用"搜索"标签页）
   - 确保输入了完整的合约地址、代币符号（QXB）和小数位数（18）

3. **确保网络正确：**
   - MetaMask 必须连接到 **Sepolia 测试网**（链 ID: 11155111）
   - 如果连接到主网或其他网络，代币不会显示

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
