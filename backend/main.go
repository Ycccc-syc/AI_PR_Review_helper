package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qiniudemo/config"
	"qiniudemo/logger"
	"qiniudemo/repository"
	"qiniudemo/router"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	gin.SetMode(cfg.GinMode)
	logger.Init(cfg.LogLevel)
	defer logger.Sync()

	db, err := repository.InitDB(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	if err != nil {
		logger.Log.Fatal("数据库初始化失败", zap.Error(err))
	}
	defer db.Close()

	if _, err := repository.InitRedis(cfg.RedisAddr, cfg.RedisPassword); err != nil {
		logger.Log.Warn("Redis 初始化失败，缓存功能不可用", zap.Error(err))
	}

	r := router.Setup(cfg, db)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Log.Info("server starting", zap.String("port", cfg.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Info("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Fatal("forced shutdown", zap.Error(err))
	}
	logger.Log.Info("server exited")
}
