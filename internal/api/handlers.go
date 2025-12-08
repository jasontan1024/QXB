package api

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gorilla/mux"
)

// TokenInfo 代币信息
type TokenInfo struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Decimals    uint8  `json:"decimals"`
	TotalSupply string `json:"totalSupply"`
	Version     string `json:"version,omitempty"`
}

// BalanceInfo 余额信息
type BalanceInfo struct {
	Address string `json:"address"`
	Balance string `json:"balance"`
	Symbol  string `json:"symbol"`
}

// RewardStatus 奖励状态
type RewardStatus struct {
	Address      string `json:"address"`
	CanClaim     bool   `json:"canClaim"`
	LastClaimDay uint64 `json:"lastClaimDay"`
	NextClaimDay uint64 `json:"nextClaimDay"`
}

// ClaimRequest 领取奖励请求
type ClaimRequest struct {
	PrivateKey string `json:"privateKey"`
}

// ClaimResponse 领取奖励响应
type ClaimResponse struct {
	TxHash string `json:"txHash"`
	Status string `json:"status"`
}

// 获取代币信息
func (s *Server) handleTokenInfo(w http.ResponseWriter, r *http.Request) {
	contract := s.ContractAddress
	ctx := context.Background()

	// 查询代币名称
	name, err := s.callString(ctx, contract, "name")
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("查询名称失败: %v", err))
		return
	}

	// 查询代币符号
	symbol, err := s.callString(ctx, contract, "symbol")
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("查询符号失败: %v", err))
		return
	}

	// 查询小数位数
	decimals, err := s.callUint8(ctx, contract, "decimals")
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("查询小数位数失败: %v", err))
		return
	}

	// 查询总供应量
	totalSupply, err := s.callUint256(ctx, contract, "totalSupply")
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("查询总供应量失败: %v", err))
		return
	}

	// 格式化总供应量
	decimalsInt := big.NewInt(int64(decimals))
	divisor := new(big.Int).Exp(big.NewInt(10), decimalsInt, nil)
	totalSupplyFloat := new(big.Float).Quo(new(big.Float).SetInt(totalSupply), new(big.Float).SetInt(divisor))

	// 查询版本（可选）
	version, _ := s.callString(ctx, contract, "version")

	info := TokenInfo{
		Name:        name,
		Symbol:      symbol,
		Decimals:    decimals,
		TotalSupply: totalSupplyFloat.Text('f', 6),
		Version:     version,
	}

	respondSuccess(w, info)
}

// 查询代币余额
func (s *Server) handleTokenBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if !common.IsHexAddress(address) {
		respondError(w, http.StatusBadRequest, "无效的地址")
		return
	}

	contract := s.ContractAddress
	userAddr := common.HexToAddress(address)
	ctx := context.Background()

	// 查询余额
	balance, err := s.callUint256WithParam(ctx, contract, "balanceOf", userAddr)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("查询余额失败: %v", err))
		return
	}

	// 查询小数位数
	decimals, err := s.callUint8(ctx, contract, "decimals")
	if err != nil {
		decimals = 18 // 默认值
	}

	// 格式化余额
	decimalsInt := big.NewInt(int64(decimals))
	divisor := new(big.Int).Exp(big.NewInt(10), decimalsInt, nil)
	balanceFloat := new(big.Float).Quo(new(big.Float).SetInt(balance), new(big.Float).SetInt(divisor))

	// 查询符号
	symbol, _ := s.callString(ctx, contract, "symbol")

	info := BalanceInfo{
		Address: address,
		Balance: balanceFloat.Text('f', 6),
		Symbol:  symbol,
	}

	respondSuccess(w, info)
}

// 查询每日奖励状态
func (s *Server) handleRewardStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	address := vars["address"]

	if !common.IsHexAddress(address) {
		respondError(w, http.StatusBadRequest, "无效的地址")
		return
	}

	contract := s.ContractAddress
	userAddr := common.HexToAddress(address)
	ctx := context.Background()

	// 调用 canClaimDailyReward 函数，返回 (bool canClaim, uint256 nextClaimDay)
	data, err := s.Contract.ABI.Pack("canClaimDailyReward", userAddr)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("打包调用失败: %v", err))
		return
	}

	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}

	result, err := s.Client.CallContract(ctx, msg, nil)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("调用合约失败: %v", err))
		return
	}

	// 解析返回值：canClaimDailyReward 返回 (bool, uint256)
	// 第一个 32 字节是 bool (canClaim)，第二个 32 字节是 uint256 (nextClaimDay)
	var canClaim bool
	var nextClaimDay *big.Int

	if len(result) >= 64 {
		// 解析 bool (第一个 32 字节的最后一位)
		canClaimBytes := result[31:32]
		canClaim = canClaimBytes[0] != 0

		// 解析 uint256 (第二个 32 字节)
		nextClaimDay = new(big.Int).SetBytes(result[32:64])
	} else {
		// 如果结果长度不足，使用默认值
		canClaim = true
		nextClaimDay = big.NewInt(0)
	}

	// 获取 lastClaimDay 用于显示
	data2, err := s.Contract.ABI.Pack("lastClaimDay", userAddr)
	if err == nil {
		msg2 := ethereum.CallMsg{
			To:   &contract,
			Data: data2,
		}
		result2, err2 := s.Client.CallContract(ctx, msg2, nil)
		if err2 == nil && len(result2) >= 32 {
			lastClaimDay := new(big.Int).SetBytes(result2[len(result2)-32:])
			status := RewardStatus{
				Address:      address,
				CanClaim:     canClaim,
				LastClaimDay: lastClaimDay.Uint64(),
				NextClaimDay: nextClaimDay.Uint64(),
			}
			respondSuccess(w, status)
			return
		}
	}

	// 如果获取 lastClaimDay 失败，仍然返回 canClaim 和 nextClaimDay
	status := RewardStatus{
		Address:      address,
		CanClaim:     canClaim,
		LastClaimDay: 0,
		NextClaimDay: nextClaimDay.Uint64(),
	}

	respondSuccess(w, status)
}

// 领取每日奖励
func (s *Server) handleClaimReward(w http.ResponseWriter, r *http.Request) {
	var req ClaimRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "无效的请求体")
		return
	}

	if req.PrivateKey == "" {
		respondError(w, http.StatusBadRequest, "私钥不能为空")
		return
	}

	// 解析私钥
	privateKeyHex := req.PrivateKey
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("无效的私钥: %v", err))
		return
	}

	// 获取地址
	fromAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	contract := s.ContractAddress
	ctx := context.Background()

	// 获取 nonce
	nonce, err := s.Client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("获取 nonce 失败: %v", err))
		return
	}

	// 获取链 ID
	chainID, err := s.Client.NetworkID(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("获取链 ID 失败: %v", err))
		return
	}

	// 获取 Gas 价格
	gasPrice, err := s.Client.SuggestGasPrice(ctx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("获取 Gas 价格失败: %v", err))
		return
	}

	// 打包 claimDailyReward 调用
	data, err := s.Contract.ABI.Pack("claimDailyReward")
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("打包调用失败: %v", err))
		return
	}

	// 估算 Gas
	msg := ethereum.CallMsg{
		From: fromAddress,
		To:   &contract,
		Data: data,
	}
	gasLimit, err := s.Client.EstimateGas(ctx, msg)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("估算 Gas 失败: %v", err))
		return
	}

	// 创建交易
	tx := types.NewTransaction(nonce, contract, big.NewInt(0), gasLimit, gasPrice, data)

	// 签名交易
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("签名交易失败: %v", err))
		return
	}

	// 发送交易
	err = s.Client.SendTransaction(ctx, signedTx)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("发送交易失败: %v", err))
		return
	}

	respondSuccess(w, ClaimResponse{
		TxHash: signedTx.Hash().Hex(),
		Status: "pending",
	})
}

// 辅助函数：调用返回 string 的函数
func (s *Server) callString(ctx context.Context, contract common.Address, method string) (string, error) {
	data, err := s.Contract.ABI.Pack(method)
	if err != nil {
		return "", err
	}

	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}

	result, err := s.Client.CallContract(ctx, msg, nil)
	if err != nil {
		return "", err
	}

	var value string
	err = s.Contract.ABI.UnpackIntoInterface(&value, method, result)
	return value, err
}

// 辅助函数：调用返回 uint8 的函数
func (s *Server) callUint8(ctx context.Context, contract common.Address, method string) (uint8, error) {
	data, err := s.Contract.ABI.Pack(method)
	if err != nil {
		return 0, err
	}

	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}

	result, err := s.Client.CallContract(ctx, msg, nil)
	if err != nil {
		return 0, err
	}

	var value uint8
	err = s.Contract.ABI.UnpackIntoInterface(&value, method, result)
	return value, err
}

// 辅助函数：调用返回 uint256 的函数
func (s *Server) callUint256(ctx context.Context, contract common.Address, method string) (*big.Int, error) {
	data, err := s.Contract.ABI.Pack(method)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}

	result, err := s.Client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	var value *big.Int
	err = s.Contract.ABI.UnpackIntoInterface(&value, method, result)
	return value, err
}

// 辅助函数：调用带参数的返回 uint256 的函数
func (s *Server) callUint256WithParam(ctx context.Context, contract common.Address, method string, param interface{}) (*big.Int, error) {
	data, err := s.Contract.ABI.Pack(method, param)
	if err != nil {
		return nil, err
	}

	msg := ethereum.CallMsg{
		To:   &contract,
		Data: data,
	}

	result, err := s.Client.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, err
	}

	// 直接解析字节，绕过 ABI 解析问题
	// 对于返回单个 uint256 的方法，直接解析 32 字节
	if len(result) == 0 {
		return big.NewInt(0), nil
	}
	
	// 取最后 32 字节（标准 uint256 长度）
	startIdx := 0
	if len(result) > 32 {
		startIdx = len(result) - 32
	}
	
	value := new(big.Int).SetBytes(result[startIdx:])
	return value, nil
}

