package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/hyeokjun/eodini/internal/middleware"
)

// 📝 설명: API 라우터 설정
// 🎯 실무 포인트: 버전별 라우팅, 미들웨어 적용
// ⚠️ 주의사항: 미들웨어 순서 중요 (Recovery -> Logger -> CORS -> ErrorHandler)

// SetupRouter - 라우터 설정
func SetupRouter() *gin.Engine {
	// Gin 모드 설정은 main에서 환경변수로 처리
	router := gin.New()

	// 글로벌 미들웨어 적용
	router.Use(middleware.RecoveryHandler())      // Panic 복구 (최우선)
	router.Use(middleware.RequestLogger())        // 요청 로깅
	router.Use(middleware.CORS(middleware.DefaultCORSConfig())) // CORS
	router.Use(middleware.ErrorHandler())         // 에러 처리 (마지막)

	// Health Check (미들웨어 제외, 가볍게)
	healthHandler := NewHealthHandler()
	router.GET("/health", healthHandler.Health)
	router.GET("/health/ready", healthHandler.Readiness)
	router.GET("/health/live", healthHandler.Liveness)

	// API v1 그룹
	v1 := router.Group("/api/v1")
	{
		// TODO: Vehicle API
		// vehicles := v1.Group("/vehicles")
		// {
		//     vehicles.GET("", vehicleHandler.List)
		//     vehicles.GET("/:id", vehicleHandler.Get)
		//     vehicles.POST("", vehicleHandler.Create)
		//     vehicles.PUT("/:id", vehicleHandler.Update)
		//     vehicles.DELETE("/:id", vehicleHandler.Delete)
		// }

		// TODO: Driver API
		// TODO: Route API
		// TODO: Trip API

		// 임시 테스트 엔드포인트
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
	}

	return router
}
