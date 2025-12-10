package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gorilla/mux"

	"lbtc/internal/auth"
	"lbtc/internal/config"
	"lbtc/internal/storage"
)

// Server API 服务器结构
type Server struct {
	Router          *mux.Router
	Client          *ethclient.Client
	Contract        *ContractService
	ContractAddress common.Address // 固定的合约地址
	AuthService     *auth.Service  // 认证服务
}

// ContractService 合约服务
type ContractService struct {
	Client *ethclient.Client
	ABI    abi.ABI
}

// QXB 合约 ABI
const qxbABI = `[
	{
		"constant": true,
		"inputs": [],
		"name": "name",
		"outputs": [{"name": "", "type": "string"}],
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
		"name": "decimals",
		"outputs": [{"name": "", "type": "uint8"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "totalSupply",
		"outputs": [{"name": "", "type": "uint256"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [{"name": "", "type": "address"}],
		"name": "balanceOf",
		"outputs": [{"name": "", "type": "uint256"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [{"name": "_user", "type": "address"}],
		"name": "canClaimDailyReward",
		"outputs": [
			{"name": "canClaim", "type": "bool"},
			{"name": "nextClaimDay", "type": "uint256"}
		],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [{"name": "_user", "type": "address"}],
		"name": "getClaimDayInfo",
		"outputs": [
			{"name": "lastDay", "type": "uint256"},
			{"name": "currentDay", "type": "uint256"}
		],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [{"name": "", "type": "address"}],
		"name": "lastClaimDay",
		"outputs": [{"name": "", "type": "uint256"}],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "version",
		"outputs": [{"name": "", "type": "string"}],
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [{"name": "_resume", "type": "string"}],
		"name": "setResume",
		"outputs": [],
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "getResume",
		"outputs": [{"name": "", "type": "string"}],
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [],
		"name": "claimDailyReward",
		"outputs": [{"name": "success", "type": "bool"}],
		"type": "function"
	}
]`

// NewServer 创建新的 API 服务器
func NewServer(rpcURL string) *Server {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("连接区块链失败: %v", err)
	}

	// 初始化数据库（使用 GORM）
	db, err := storage.OpenGORM(config.GetDBPath())
	if err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化认证服务
	authService, err := auth.NewService(db)
	if err != nil {
		log.Fatalf("初始化认证服务失败: %v", err)
	}

	// 使用内置 ABI（包含最新接口）
	contractABI, err := abi.JSON(strings.NewReader(qxbABI))
	if err != nil {
		log.Fatalf("解析 ABI 失败: %v", err)
	}

	return &Server{
		Router:          mux.NewRouter(),
		Client:          client,
		ContractAddress: common.HexToAddress(config.QXBContractAddress),
		Contract: &ContractService{
			Client: client,
			ABI:    contractABI,
		},
		AuthService: authService,
	}
}

// SetupRoutes 设置路由
func (s *Server) SetupRoutes() {
	// 启用 CORS（必须在所有路由之前）
	s.Router.Use(corsMiddleware)

	// 为所有路径添加 OPTIONS 处理（通配符匹配）
	s.Router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.WriteHeader(http.StatusOK)
	})

	// API 文档
	s.Router.HandleFunc("/api/docs", s.handleDocs).Methods("GET")

	// 健康检查
	s.Router.HandleFunc("/health", s.handleHealth).Methods("GET")

	// 代币信息 API
	api := s.Router.PathPrefix("/api").Subrouter()
	api.Use(corsMiddleware) // 子路由也需要 CORS

	// 代币相关
	api.HandleFunc("/token/info", s.handleTokenInfo).Methods("GET")
	api.HandleFunc("/token/balance/{address}", s.handleTokenBalance).Methods("GET")
	api.HandleFunc("/resume", s.handleResume).Methods("GET")

	// 每日奖励相关
	api.HandleFunc("/reward/status/{address}", s.handleRewardStatus).Methods("GET")
	api.HandleFunc("/reward/claim", s.optionalAuthMiddleware(s.handleClaimReward)).Methods("POST")

	// 认证相关
	api.HandleFunc("/auth/register", s.handleRegister).Methods("POST")
	api.HandleFunc("/auth/login", s.handleLogin).Methods("POST")
	api.HandleFunc("/auth/me", s.authMiddleware(s.handleMe)).Methods("GET")

	// 代币转账（需要认证）
	api.HandleFunc("/token/transfer", s.authMiddleware(s.handleTransfer)).Methods("POST")
}

// Response 通用响应结构
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// 响应辅助函数
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

func respondSuccess(w http.ResponseWriter, data interface{}) {
	respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}

// CORS 中间件
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 设置 CORS 头
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// 处理预检请求
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// 健康检查
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	respondSuccess(w, map[string]string{
		"status":  "ok",
		"service": "QXB API",
	})
}

// API 文档
func (s *Server) handleDocs(w http.ResponseWriter, r *http.Request) {
	docs := `
# QXB API 文档

## 基础信息
- Base URL: http://localhost:8080/api
- 所有响应格式: JSON

## 端点列表

### 1. 获取代币信息
GET /api/token/info

说明: 获取 QXB 代币的基本信息（合约地址已在配置中固定）

示例:
GET /api/token/info

响应:
{
  "success": true,
  "data": {
    "name": "齐夏币",
    "symbol": "QXB",
    "decimals": 18,
    "totalSupply": "2025.0",
    "version": "1.0.0"
  }
}

### 2. 查询代币余额
GET /api/token/balance/{address}

参数:
- address: 查询的地址

示例:
GET /api/token/balance/0x405e2ea956ea490bf3d4bd734dc334a1d42b35b9

响应:
{
  "success": true,
  "data": {
    "address": "0x405e2ea956ea490bf3d4bd734dc334a1d42b35b9",
    "balance": "1000.5",
    "symbol": "QXB"
  }
}

### 3. 查询每日奖励状态
GET /api/reward/status/{address}

参数:
- address: 查询的地址

示例:
GET /api/reward/status/0xe628ce9c1def02fa8958d081bbda75b4a9907955

响应:
{
  "success": true,
  "data": {
    "address": "0x405e2ea956ea490bf3d4bd734dc334a1d42b35b9",
    "canClaim": true,
    "lastClaimDay": 19723,
    "nextClaimDay": 0
  }
}

### 4. 领取每日奖励
POST /api/reward/claim

请求体:
{
  "privateKey": "你的私钥"
}

响应:
{
  "success": true,
  "data": {
    "txHash": "0x...",
    "status": "pending"
  }
}
`
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, docs)
}
