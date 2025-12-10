package main

import (
	"context"
	"flag"
	"fmt"
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

// 标准 ERC20 transfer ABI
const erc20ABI = `[{"constant":false,"inputs":[{"name":"to","type":"address"},{"name":"value","type":"uint256"}],"name":"transfer","outputs":[{"name":"","type":"bool"}],"type":"function"}]`

func must(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}

func main() {
	toFlag := flag.String("to", "", "接收地址")
	amountFlag := flag.String("amount", "10000000000000000000", "转账金额(wei)，默认10 QXB (10 * 1e18)")
	flag.Parse()

	if *toFlag == "" {
		log.Fatal("必须指定 --to 地址")
	}

	privHex := config.GetPrivateKey()
	if privHex == "" {
		log.Fatal("缺少 PRIVATE_KEY 环境变量（合约拥有者私钥）")
	}

	client, err := ethclient.Dial(config.EthereumRPCURL)
	must(err, "连接 RPC 失败")

	privateKey, err := crypto.HexToECDSA(trim0x(privHex))
	must(err, "解析私钥失败")
	fromAddr := crypto.PubkeyToAddress(privateKey.PublicKey)

	to := common.HexToAddress(*toFlag)
	amount, ok := new(big.Int).SetString(*amountFlag, 10)
	if !ok {
		log.Fatal("amount 解析失败")
	}

	parsedABI, err := abi.JSON(strings.NewReader(erc20ABI))
	must(err, "解析 ABI 失败")
	data, err := parsedABI.Pack("transfer", to, amount)
	must(err, "打包数据失败")

	nonce, err := client.PendingNonceAt(context.Background(), fromAddr)
	must(err, "获取 nonce 失败")
	gasPrice, err := client.SuggestGasPrice(context.Background())
	must(err, "获取 gasPrice 失败")

	contract := common.HexToAddress(config.QXBContractAddress)
	msg := ethereum.CallMsg{From: fromAddr, To: &contract, Data: data}
	gasLimit, err := client.EstimateGas(context.Background(), msg)
	must(err, "估算 Gas 失败")

	tx := types.NewTransaction(nonce, contract, big.NewInt(0), gasLimit, gasPrice, data)

	chainID, err := client.NetworkID(context.Background())
	must(err, "获取 chainID 失败")

	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	must(err, "签名交易失败")

	err = client.SendTransaction(context.Background(), signedTx)
	must(err, "发送交易失败")

	fmt.Printf("转账提交成功: txHash=%s\n", signedTx.Hash().Hex())
	fmt.Printf("from=%s -> to=%s amount=%s wei\n", fromAddr.Hex(), to.Hex(), amount.String())
}

func trim0x(s string) string {
	if len(s) > 1 && (s[0:2] == "0x" || s[0:2] == "0X") {
		return s[2:]
	}
	return s
}
