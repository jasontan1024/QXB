package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"lbtc/internal/api"
	"lbtc/internal/config"
)

func main() {
	// 创建 API 服务器
	server := api.NewServer(config.EthereumRPCURL)

	// 设置路由
	server.SetupRoutes()

	// 启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// 包装路由器以处理 CORS
	corsHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

		// 处理预检请求（必须在路由之前）
		// 对于所有路径的 OPTIONS 请求都返回成功
		if r.Method == "OPTIONS" {
			log.Printf("处理 OPTIONS 请求: %s", r.URL.Path)
			w.WriteHeader(http.StatusOK)
			return
		}

		// 继续处理其他请求
		server.Router.ServeHTTP(w, r)
	})

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: corsHandler,
	}

	// 启动服务器（在 goroutine 中）
	go func() {
		log.Printf("🚀 API 服务器启动在端口 %s", port)
		log.Printf("📖 API 文档: http://localhost:%s/api/docs", port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("正在关闭服务器...")

	// 优雅关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("服务器关闭失败: %v", err)
	}

	log.Println("服务器已关闭")
}
