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

## 部署方式

### 合约部署

#### 前置准备

1. **安装 Foundry 工具链（仅用于编译合约）**
   - Foundry 用于编译 Solidity 合约，生成字节码和 ABI
   - 部署程序会读取 Foundry 编译后的 JSON 文件（`out/QXB.sol/QXB.json`）
   - 安装完成后使用 `forge build` 编译合约

2. **准备部署账户**
   - 需要准备一个 Sepolia 测试网账户
   - 账户需要有足够的 Sepolia ETH 用于支付 Gas 费用（建议至少 0.01 ETH）
   - 设置环境变量 `PRIVATE_KEY` 为部署账户的私钥

3. **配置网络**
   - 默认部署到 Sepolia 测试网
   - RPC 节点：https://ethereum-sepolia-rpc.publicnode.com
   - 链 ID：11155111

#### 部署步骤

1. **编译合约**
   - 使用 Foundry 编译 QXB 合约：`forge build`
   - 编译后会生成 `out/QXB.sol/QXB.json` 文件
   - 该文件包含合约字节码和 ABI，供部署程序使用

2. **执行部署**
   - 使用项目自带的 Go 部署程序进行部署：`go run ./cmd/deploy-direct`
   - 部署程序会读取 Foundry 编译后的 JSON 文件
   - 自动处理交易签名、Gas 估算和交易发送
   - 等待交易确认后输出合约地址

3. **配置合约地址**
   - 部署成功后，将合约地址配置到项目中
   - 更新 `internal/config/config.go` 中的合约地址
   - 重启 API 服务器使配置生效

**说明**：虽然编译合约需要 Foundry，但实际的部署过程完全由项目自带的 Go 程序完成，不依赖 Foundry 的部署功能。

### 项目部署

项目支持两种部署方式：本地部署和 Docker 部署。

#### 方式一：本地部署

**前置要求**：
- Go 1.21+ 环境
- Node.js 20+ 和 npm
- 已部署合约并配置合约地址

**部署步骤**：

1. **配置环境变量**
   - 创建 `.env` 文件（可参考 `.env.example`）
   - 设置 `PRIVATE_KEY`（用于自动转账 ETH 功能）
   - 设置 `JWT_SECRET`（JWT 密钥，可选）

2. **启动后端 API 服务器**
   - 进入项目根目录
   - 运行：`go run ./cmd/api`
   - API 服务将在 `http://localhost:8080` 启动

3. **启动前端应用**
   - 进入 `web` 目录
   - 安装依赖：`npm install`（首次运行）
   - 启动开发服务器：`npm start`
   - 前端应用将在 `http://localhost:3000` 启动

4. **验证部署**
   - 访问 `http://localhost:3000` 查看前端界面
   - 访问 `http://localhost:8080/health` 检查后端健康状态
   - 访问 `http://localhost:8080/api/docs` 查看 API 文档

**数据存储**：
- 数据库文件默认存储在 `data/app.db`
- 可通过环境变量 `DB_PATH` 自定义数据库路径

#### 方式二：Docker 部署

**前置要求**：
- Docker 和 Docker Compose
- 已部署合约并配置合约地址

**部署步骤**：

1. **配置环境变量**
   - 创建 `.env` 文件（可参考 `.env.example`）
   - 设置 `PRIVATE_KEY`（用于自动转账 ETH 功能）
   - 设置 `JWT_SECRET`（JWT 密钥，可选）

2. **构建和启动服务**
   - 在项目根目录运行：`docker-compose up -d`
   - Docker Compose 会自动构建并启动 API 和 Web 服务

3. **查看服务状态**
   - 查看日志：`docker-compose logs -f`
   - 查看服务状态：`docker-compose ps`
   - 停止服务：`docker-compose down`

4. **验证部署**
   - 访问 `http://localhost:3000` 查看前端界面
   - 访问 `http://localhost:8080/health` 检查后端健康状态

**Docker 配置说明**：
- **API 服务**（`Dockerfile.api`）：
  - 基于 Go 1.25 Alpine 镜像
  - 数据目录挂载到 `./data`
  - 端口映射：`8080:8080`
  
- **Web 服务**（`Dockerfile.web`）：
  - 基于 Node.js 20 Alpine 镜像
  - 构建时注入 API 地址
  - 使用 `serve` 提供静态文件服务
  - 端口映射：`3000:3000`

- **数据持久化**：
  - 数据库文件存储在 `./data` 目录
  - 通过 Docker volume 挂载，确保数据持久化

**注意事项**：
- 首次构建可能需要较长时间（下载依赖和编译）
- 确保 `data` 目录有写入权限
- 生产环境建议使用 `.env` 文件管理敏感信息，不要将私钥提交到代码仓库

### 合约功能说明

**QXB 代币合约**实现了以下功能：

- **ERC20 标准代币**
  - 完全兼容 ERC20 标准
  - 支持标准转账和授权功能

- **代币管理**
  - 代币名称：齐夏币
  - 代币符号：QXB
  - 小数位数：18
  - 初始总供应量：2,025 QXB（分配给部署者）

- **每日奖励机制**
  - 每个地址每天可以领取 1 QXB
  - 基于 UTC 日期，每天只能领取一次
  - 提供奖励状态查询功能

- **代币操作**
  - 标准转账功能（transfer）
  - 授权转账功能（transferFrom）
  - 授权管理（approve, increaseAllowance, decreaseAllowance）
  - 代币铸造（mint，仅合约所有者）
  - 代币销毁（burn，任何持有者）

- **作者简历功能**
  - 合约所有者可以设置 Markdown 格式的简历
  - 任何人都可以读取简历内容

### 部署后配置

部署完成后，需要：

1. **更新配置**
   - 将部署得到的合约地址写入配置文件
   - 确保 API 服务器使用正确的合约地址

2. **验证部署**
   - 在 Etherscan 上查看合约详情
   - 验证合约代码和状态
   - 确认初始代币分配正确

3. **测试功能**
   - 测试代币转账功能
   - 测试每日奖励领取功能
   - 验证所有 API 接口正常工作

## 快速开始

### 1. 部署合约

按照 [合约部署](#合约部署) 章节的说明，完成合约的编译和部署。

### 2. 部署项目

选择以下两种方式之一部署项目：

**方式一：本地部署**
- 按照 [本地部署](#方式一本地部署) 章节的说明进行部署
- 适合开发和测试环境

**方式二：Docker 部署**
- 按照 [Docker 部署](#方式二docker-部署) 章节的说明进行部署
- 适合生产环境和容器化部署

### 3. 访问应用

- **前端应用**：http://localhost:3000
- **API 服务**：http://localhost:8080
- **API 文档**：http://localhost:8080/api/docs
- **健康检查**：http://localhost:8080/health

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

## API 文档

详细的 API 文档请参考 [API.md](API.md)

**主要 API 端点概览：**

- **基础端点**
  - `GET /health` - 健康检查
  - `GET /api/docs` - API 文档

- **代币相关**
  - `GET /api/token/info` - 查询代币信息
  - `GET /api/token/balance/<地址>` - 查询代币余额
  - `POST /api/token/transfer` - 转账代币（需要认证）

- **每日奖励相关**
  - `GET /api/reward/status/<地址>` - 查询奖励状态
  - `POST /api/reward/claim` - 领取每日奖励

- **认证相关**
  - `POST /api/auth/register` - 用户注册
  - `POST /api/auth/login` - 用户登录
  - `GET /api/auth/me` - 获取当前用户信息（需要认证）

- **其他**
  - `GET /api/resume` - 获取作者简历

所有 API 响应均为 JSON 格式，统一使用以下结构：

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
