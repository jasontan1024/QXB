package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// BlockInfo 区块信息结构体
type BlockInfo struct {
	Number     *big.Int
	Hash       common.Hash
	ParentHash common.Hash
	Time       uint64
	TxCount    uint
	Difficulty *big.Int
	Coinbase   common.Address
	Size       uint64
}

// GetLatestBlockNumber 获取区块链上最新的区块号
// 区块号是一个递增的数字，表示区块在链上的位置
func GetLatestBlockNumber(client *ethclient.Client) (*big.Int, error) {
	// BlockNumber 返回最新的区块号
	blockNumber, err := client.BlockNumber(context.Background())
	if err != nil {
		return nil, fmt.Errorf("获取区块号失败: %w", err)
	}

	// 将 uint64 转换为 *big.Int（大整数类型，用于处理大数字）
	return big.NewInt(int64(blockNumber)), nil
}

// GetBlockInfo 获取指定区块的详细信息
// 使用 HeaderByNumber 避免解析不支持的交易类型
func GetBlockInfo(client *ethclient.Client, blockNumber *big.Int) (*BlockInfo, error) {
	// 添加超时和重试机制
	var header *types.Header
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		// 创建带超时的 context（15秒超时）
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)

		// HeaderByNumber 获取区块头信息（不包含交易数据）
		// 这样可以避免解析不支持的交易类型
		header, err = client.HeaderByNumber(ctx, blockNumber)
		cancel()

		if err == nil {
			break
		}

		// 如果不是最后一次重试，等待后重试
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
		}
	}

	if err != nil {
		return nil, fmt.Errorf("获取区块信息失败: %w", err)
	}

	// 使用 RPC 调用获取交易数量（只获取交易哈希，不解析完整交易）
	// 这样可以避免交易类型不支持的问题
	txCount, err := getTransactionCountViaRPC(client, blockNumber)
	if err != nil {
		// 如果获取交易数量失败，设置为 0
		txCount = 0
	}

	return &BlockInfo{
		Number:     header.Number,
		Hash:       header.Hash(),
		ParentHash: header.ParentHash,
		Time:       header.Time,
		TxCount:    txCount,
		Difficulty: header.Difficulty,
		Coinbase:   header.Coinbase,
		Size:       0, // 区块大小需要完整区块数据，这里设为 0
	}, nil
}

// PrintBlockInfo 打印区块的详细信息
// 帮助理解区块的数据结构
func PrintBlockInfo(blockInfo *BlockInfo) {
	fmt.Printf("  区块号: %s\n", blockInfo.Number.String())
	fmt.Printf("  区块哈希: %s\n", blockInfo.Hash.Hex())
	fmt.Printf("  父区块哈希: %s\n", blockInfo.ParentHash.Hex())
	fmt.Printf("  时间戳: %s\n", time.Unix(int64(blockInfo.Time), 0).Format("2006-01-02 15:04:05"))
	fmt.Printf("  交易数量: %d\n", blockInfo.TxCount)
	fmt.Printf("  区块难度: %s\n", blockInfo.Difficulty.String())
	fmt.Printf("  矿工地址: %s\n", blockInfo.Coinbase.Hex())
	if blockInfo.Size > 0 {
		fmt.Printf("  区块大小: %d bytes\n", blockInfo.Size)
	}
}

// getTransactionCountViaRPC 通过 RPC 调用获取区块中的交易数量
// 使用 eth_getBlockByNumber RPC 方法，只请求交易哈希列表（不解析完整交易）
func getTransactionCountViaRPC(client *ethclient.Client, blockNumber *big.Int) (uint, error) {
	// 获取底层 RPC 客户端
	rpcClient := client.Client()

	// 构建 RPC 请求参数
	blockNumHex := fmt.Sprintf("0x%x", blockNumber)

	// 定义返回结构体（只包含交易哈希列表）
	var result struct {
		Transactions []common.Hash `json:"transactions"`
	}

	// 调用 eth_getBlockByNumber，第二个参数 false 表示只返回交易哈希
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := rpcClient.CallContext(ctx, &result, "eth_getBlockByNumber", blockNumHex, false)
	if err != nil {
		return 0, err
	}

	return uint(len(result.Transactions)), nil
}

