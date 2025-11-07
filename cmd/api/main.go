package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyeokjun/eodini/config"
	"github.com/hyeokjun/eodini/internal/handler"
	"github.com/hyeokjun/eodini/pkg/logger"
)

// 📝 설명: 애플리케이션 진입점
// 🎯 실무 포인트: Graceful Shutdown, 설정 로드, 로거 초기화
// ⚠️ 주의사항: 서버 시작 전 설정 검증 필수

// @title						Eodini API
// @version					1.0
// @description				통학/통원 차량 관리 시스템 API
// @termsOfService				http://swagger.io/terms/
// @contact.name				API Support
// @contact.email				support@eodini.com
// @license.name				Apache 2.0
// @license.url				http://www.apache.org/licenses/LICENSE-2.0.html
// @host						localhost:8080
// @BasePath					/api/v1
// @schemes					http https
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
// @description				JWT Bearer token (추후 구현 예정)
func main() {
	// 1. 설정 로드
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. 로거 초기화
	initLogger(cfg)
	logger.Info("Starting Eodini API Server", map[string]interface{}{
		"version":     "0.1.0",
		"environment": cfg.Server.Environment,
		"port":        cfg.Server.Port,
	})

	// 3. Gin 모드 설정
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// 4. 라우터 설정
	router := handler.SetupRouter()

	// 5. HTTP 서버 설정
	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	// 6. 서버 시작 (고루틴)
	go func() {
		logger.Infof("Server listening on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// 7. Graceful Shutdown 대기
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...", nil)

	// 8. Graceful Shutdown (최대 30초 대기)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("Server exited gracefully", nil)
}

// initLogger - 로거 초기화
func initLogger(cfg *config.Config) {
	// 로그 레벨 설정
	switch cfg.Log.Level {
	case "debug":
		logger.SetLevel(logger.DebugLevel)
	case "info":
		logger.SetLevel(logger.InfoLevel)
	case "warn":
		logger.SetLevel(logger.WarnLevel)
	case "error":
		logger.SetLevel(logger.ErrorLevel)
	default:
		logger.SetLevel(logger.InfoLevel)
	}

	// TODO: 로그 포맷 설정 (json/text)
	// 현재는 기본 text 포맷 사용
}
