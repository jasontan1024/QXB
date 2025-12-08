package contract

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ContractInfo 合约信息结构体
type ContractInfo struct {
	Address common.Address
	Code    []byte
	Balance *big.Int
	HasCode bool
}

// GetContractInfo 获取智能合约的基本信息
// 包括：合约地址、是否有代码、合约余额等
func GetContractInfo(client *ethclient.Client, contractAddressHex string) (*ContractInfo, error) {
	contractAddress := common.HexToAddress(contractAddressHex)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 获取合约代码（如果地址是合约，会有代码；如果是普通账户，代码为空）
	code, err := client.CodeAt(ctx, contractAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("获取合约代码失败: %w", err)
	}

	// 获取合约的 ETH 余额
	balance, err := client.BalanceAt(ctx, contractAddress, nil)
	if err != nil {
		return nil, fmt.Errorf("获取合约余额失败: %w", err)
	}

	return &ContractInfo{
		Address: contractAddress,
		Code:    code,
		Balance: balance,
		HasCode: len(code) > 0,
	}, nil
}

// PrintContractInfo 打印合约信息
func PrintContractInfo(info *ContractInfo) {
	fmt.Printf("  合约地址: %s\n", info.Address.Hex())
	fmt.Printf("  是否为合约: %v\n", info.HasCode)
	if info.HasCode {
		fmt.Printf("  合约代码长度: %d bytes\n", len(info.Code))
	}

	// 显示合约的 ETH 余额
	ethBalance := new(big.Float).Quo(new(big.Float).SetInt(info.Balance), big.NewFloat(1e18))
	fmt.Printf("  合约 ETH 余额: %s ETH\n", ethBalance.Text('f', 6))
}

