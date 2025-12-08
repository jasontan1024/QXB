package contract

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// ERC20 标准接口的 ABI（Application Binary Interface）
// 这里只包含我们需要的函数：balanceOf(address) 和 decimals()
const erc20ABI = `[
	{
		"constant": true,
		"inputs": [{"name": "_owner", "type": "address"}],
		"name": "balanceOf",
		"outputs": [{"name": "balance", "type": "uint256"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "decimals",
		"outputs": [{"name": "", "type": "uint8"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "symbol",
		"outputs": [{"name": "", "type": "string"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [{"name": "", "type": "string"}],
		"type": "function"
	}
]`

// GetERC20Balance 读取 ERC20 代币合约中指定地址的余额
// 这是智能合约交互的典型示例：调用合约的只读函数
func GetERC20Balance(client *ethclient.Client, contractAddressHex, ownerAddressHex string) (string, error) {
	// 解析合约地址和所有者地址
	contractAddress := common.HexToAddress(contractAddressHex)
	ownerAddress := common.HexToAddress(ownerAddressHex)

	// 解析 ERC20 ABI
	contractABI, err := abi.JSON(strings.NewReader(erc20ABI))
	if err != nil {
		return "", fmt.Errorf("解析 ABI 失败: %w", err)
	}

	// 准备调用 balanceOf(address) 函数
	// 这是 ERC20 标准函数，用于查询指定地址的代币余额
	data, err := contractABI.Pack("balanceOf", ownerAddress)
	if err != nil {
		return "", fmt.Errorf("打包函数调用失败: %w", err)
	}

	// 创建合约调用（只读操作，不会改变链上状态）
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	msg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: data,
	}

	result, err := client.CallContract(ctx, msg, nil)
	if err != nil {
		return "", fmt.Errorf("调用合约失败: %w", err)
	}

	// 解析返回结果（uint256 类型的余额）
	var balance *big.Int
	err = contractABI.UnpackIntoInterface(&balance, "balanceOf", result)
	if err != nil {
		return "", fmt.Errorf("解析返回结果失败: %w", err)
	}

	// 尝试获取代币的小数位数（decimals）
	// 大多数 ERC20 代币使用 18 位小数（与 ETH 相同）
	decimalsData, _ := contractABI.Pack("decimals")
	decimalsMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: decimalsData,
	}
	decimalsResult, err := client.CallContract(ctx, decimalsMsg, nil)
	var decimals uint8 = 18 // 默认值
	if err == nil {
		contractABI.UnpackIntoInterface(&decimals, "decimals", decimalsResult)
	}

	// 格式化余额（考虑小数位数）
	// 例如：如果余额是 1000000000000000000，decimals 是 18，则显示为 1.0
	divisor := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(decimals)), nil))
	balanceFloat := new(big.Float).SetInt(balance)
	formattedBalance := new(big.Float).Quo(balanceFloat, divisor)

	return formattedBalance.Text('f', 6), nil
}

