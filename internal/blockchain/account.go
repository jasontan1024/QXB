package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// GetAccountBalance 查询指定以太坊地址的余额
// 余额以 Wei 为单位（1 ETH = 10^18 Wei）
func GetAccountBalance(client *ethclient.Client, addressHex string) (string, error) {
	// common.HexToAddress 将十六进制字符串转换为 Address 类型
	address := common.HexToAddress(addressHex)

	// 添加超时和重试机制，提高查询的稳定性
	var balance *big.Int
	var err error
	maxRetries := 3

	for i := 0; i < maxRetries; i++ {
		// 创建带超时的 context（10秒超时）
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// BalanceAt 获取指定地址在最新区块的余额
		// nil 参数表示查询最新区块的余额
		balance, err = client.BalanceAt(ctx, address, nil)
		cancel()

		if err == nil {
			break
		}

		// 如果不是最后一次重试，等待后重试
		if i < maxRetries-1 {
			time.Sleep(time.Duration(i+1) * time.Second)
			fmt.Printf("  查询余额失败，正在重试 (%d/%d)...\n", i+1, maxRetries-1)
		}
	}

	if err != nil {
		return "", fmt.Errorf("查询余额失败: %w", err)
	}

	// 将 Wei 转换为 ETH（除以 10^18）
	// 使用浮点数格式化，保留 4 位小数
	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(balance), big.NewFloat(1e18))
	return ethBalance.Text('f', 4), nil
}

// WatchNewBlocks 监听新区块的产生
// 这是一个实时监听示例，展示如何响应链上事件
func WatchNewBlocks(client *ethclient.Client, duration time.Duration) {
	// 获取当前区块号作为起始点
	startBlock, err := client.BlockNumber(context.Background())
	if err != nil {
		fmt.Printf("获取起始区块失败: %v\n", err)
		return
	}

	fmt.Printf("  开始监听，当前区块: %d\n", startBlock)

	// 创建一个定时器，在指定时间后停止监听
	timer := time.NewTimer(duration)
	defer timer.Stop()

	// 创建一个 ticker，每 2 秒检查一次新区块
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	lastBlock := startBlock

	for {
		select {
		case <-timer.C:
			// 时间到了，停止监听
			fmt.Println("  监听时间结束")
			return

		case <-ticker.C:
			// 每 2 秒检查一次是否有新区块
			currentBlock, err := client.BlockNumber(context.Background())
			if err != nil {
				fmt.Printf("获取区块号失败: %v\n", err)
				continue
			}

			// 如果发现新区块，打印信息
			if currentBlock > lastBlock {
				newBlocks := currentBlock - lastBlock
				fmt.Printf("  ✓ 发现 %d 个新区块！最新区块号: %d\n", newBlocks, currentBlock)
				lastBlock = currentBlock
			}
		}
	}
}

