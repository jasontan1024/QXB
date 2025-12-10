package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"lbtc/internal/config"
)

// QXB 合约 ABI（仅包含 setResume 函数）
const qxbABI = `[{"constant":false,"inputs":[{"name":"_resume","type":"string"}],"name":"setResume","outputs":[],"type":"function"}]`

func trim0x(s string) string {
	return strings.TrimPrefix(s, "0x")
}

func main() {
	resumeFile := flag.String("file", "", "简历 Markdown 文件路径（必填）")
	flag.Parse()

	if *resumeFile == "" {
		log.Fatal("必须指定 --file 参数（简历 Markdown 文件路径）")
	}

	// 读取简历文件
	resumeContent, err := ioutil.ReadFile(*resumeFile)
	if err != nil {
		log.Fatalf("读取简历文件失败: %v", err)
	}

	resumeText := strings.TrimSpace(string(resumeContent))
	if resumeText == "" {
		log.Fatal("简历内容为空")
	}

	fmt.Printf("📄 简历内容长度: %d 字符\n", len(resumeText))
	fmt.Println()

	privHex := config.GetPrivateKey()
	if privHex == "" {
		log.Fatal("缺少 PRIVATE_KEY 环境变量（合约拥有者私钥）")
	}

	client, err := ethclient.Dial(config.EthereumRPCURL)
	if err != nil {
		log.Fatalf("连接 RPC 失败: %v", err)
	}
	defer client.Close()

	privateKey, err := crypto.HexToECDSA(trim0x(privHex))
	if err != nil {
		log.Fatalf("解析私钥失败: %v", err)
	}

	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)
	contractAddr := common.HexToAddress(config.QXBContractAddress)

	fmt.Printf("合约地址: %s\n", contractAddr.Hex())
	fmt.Printf("发送地址: %s\n", fromAddr.Hex())
	fmt.Println()

	// 解析 ABI
	parsedABI, err := abi.JSON(strings.NewReader(qxbABI))
	if err != nil {
		log.Fatalf("解析 ABI 失败: %v", err)
	}

	// 编码 setResume 调用
	data, err := parsedABI.Pack("setResume", resumeText)
	if err != nil {
		log.Fatalf("编码调用失败: %v", err)
	}

	// 获取链 ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatalf("获取链 ID 失败: %v", err)
	}

	// 获取 nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	if err != nil {
		log.Fatalf("获取 nonce 失败: %v", err)
	}

	// 估算 Gas
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatalf("获取 Gas 价格失败: %v", err)
	}

	// 估算 Gas Limit
	msg := ethereum.CallMsg{
		From: fromAddr,
		To:   &contractAddr,
		Data: data,
	}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	if err != nil {
		log.Fatalf("估算 Gas Limit 失败: %v", err)
	}

	// 创建交易
	tx := types.NewTransaction(nonce, contractAddr, big.NewInt(0), gasLimit, gasPrice, data)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatalf("签名交易失败: %v", err)
	}

	fmt.Println("📝 交易信息：")
	fmt.Printf("  Nonce: %d\n", nonce)
	fmt.Printf("  Gas Price: %s Gwei\n", new(big.Float).Quo(new(big.Float).SetInt(gasPrice), big.NewFloat(1e9)).Text('f', 2))
	fmt.Printf("  Gas Limit: %d\n", gasLimit)
	fmt.Println()

	// 发送交易
	fmt.Println("🚀 发送交易到区块链...")
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("发送交易失败: %v", err)
	}

	fmt.Printf("✅ 交易已发送！\n")
	fmt.Printf("交易哈希: %s\n", signedTx.Hash().Hex())
	fmt.Println()

	// 等待交易确认
	fmt.Println("⏳ 等待交易确认...")
	ctx := context.Background()
	receipt, err := waitForTransaction(ctx, client, signedTx.Hash())
	if err != nil {
		log.Fatalf("等待交易确认失败: %v", err)
	}

	if receipt.Status == 0 {
		log.Fatal("❌ 交易失败！")
	}

	fmt.Println("✅ 简历写入成功！")
	fmt.Printf("区块号: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas 使用: %d\n", receipt.GasUsed)
	fmt.Println()
	fmt.Printf("📝 在 Etherscan 查看: https://sepolia.etherscan.io/tx/%s\n", signedTx.Hash().Hex())
}

// waitForTransaction 等待交易确认
func waitForTransaction(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(ctx, txHash)
		if err == nil {
			return receipt, nil
		}
		if err.Error() == "not found" {
			// 交易还未确认，继续等待
			continue
		}
		return nil, err
	}
}
