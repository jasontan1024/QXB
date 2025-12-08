package blockchain

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"lbtc/internal/config"
)

// Connect 连接到以太坊区块链网络
// 返回一个 ethclient.Client 实例，用于后续的区块链交互
func Connect() (*ethclient.Client, error) {
	// ethclient.Dial 创建一个新的客户端连接到以太坊节点
	// 参数是 RPC 节点的 URL（HTTP 或 WebSocket）
	client, err := ethclient.Dial(config.EthereumRPCURL)
	if err != nil {
		return nil, fmt.Errorf("无法连接到 RPC 节点: %w", err)
	}

	// 验证连接是否正常（可选但推荐）
	// 通过获取链 ID 来测试连接
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("无法获取链 ID: %w", err)
	}

	// 显示网络信息
	networkName := config.GetNetworkName(chainID.Uint64())
	fmt.Printf("  网络: %s (链 ID: %s)\n", networkName, chainID.String())

	return client, nil
}

