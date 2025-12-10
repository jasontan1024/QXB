package config

import (
	"os"

	"github.com/joho/godotenv"
)

// 区块链连接配置
const (
	// Sepolia 测试网络 RPC 节点
	// 使用公共节点，无需 API Key
	// 如果需要更高的速率限制，可以注册 Infura 或 Alchemy 获取 API Key
	// 其他可选节点：
	// - https://rpc.sepolia.org (公共节点，可能不稳定)
	// - https://ethereum-sepolia-rpc.publicnode.com (公共节点)
	// - https://sepolia.infura.io/v3/YOUR_API_KEY (需要注册)
	EthereumRPCURL = "https://ethereum-sepolia-rpc.publicnode.com" // Sepolia 测试网公共 RPC 节点

	// QXB 合约地址（Sepolia 测试网）
	QXBContractAddress = "0x5068a014aC8e691Be53848FE5872cbA9f8C4dA17"

	// SQLite 数据库路径
	DefaultDBPath = "data/app.db"
)

var (
	// 是否已初始化
	initialized bool
)

// init 初始化配置，加载 .env 文件
func init() {
	LoadEnv()
}

// LoadEnv 加载 .env 文件到环境变量
// 如果 .env 文件不存在，不会报错（允许通过系统环境变量设置）
func LoadEnv() {
	if initialized {
		return
	}
	// 尝试加载 .env 文件，如果文件不存在也不报错
	_ = godotenv.Load()
	initialized = true
}

// GetDBPath 获取数据库路径
func GetDBPath() string {
	LoadEnv()
	if v := os.Getenv("DB_PATH"); v != "" {
		return v
	}
	return DefaultDBPath
}

// GetJWTSecret 获取 JWT 密钥
func GetJWTSecret() string {
	LoadEnv()
	return os.Getenv("JWT_SECRET")
}

// GetPrivateKey 从环境变量获取私钥
// ⚠️⚠️⚠️ 安全警告 ⚠️⚠️⚠️
// 私钥是非常敏感的信息，必须通过环境变量设置！
// 请确保：
// 1. 在 .env 文件中设置 PRIVATE_KEY（不要提交 .env 到 Git）
// 2. 或者通过系统环境变量设置：export PRIVATE_KEY=你的私钥
// 3. 不要在主网使用测试私钥
// 4. 生产环境应该使用密钥管理服务
func GetPrivateKey() string {
	// 确保已加载 .env 文件
	LoadEnv()
	return os.Getenv("PRIVATE_KEY")
}

// NetworkInfo 网络信息
type NetworkInfo struct {
	Name    string
	ChainID uint64
}

// GetNetworkName 根据链 ID 获取网络名称
func GetNetworkName(chainID uint64) string {
	switch chainID {
	case 1:
		return "以太坊主网"
	case 11155111:
		return "Sepolia 测试网"
	default:
		return "未知网络"
	}
}
