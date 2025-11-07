package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyeokjun/eodini/internal/util"
)

// 📝 설명: Health Check 핸들러 (서버 상태 확인)
// 🎯 실무 포인트: K8s liveness/readiness probe용
// ⚠️ 주의사항: 무거운 로직 포함 금지 (빠른 응답 필요)

// HealthHandler - Health Check 핸들러
type HealthHandler struct {
	startTime time.Time
}

// NewHealthHandler - Health Check 핸들러 생성
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{
		startTime: time.Now(),
	}
}

// HealthResponse - Health Check 응답
type HealthResponse struct {
	Status    string    `json:"status"`     // "healthy", "unhealthy"
	Timestamp time.Time `json:"timestamp"`  // 현재 시각
	Uptime    string    `json:"uptime"`     // 서버 가동 시간
	Version   string    `json:"version"`    // 버전 정보
}

// Health - 기본 Health Check
// @Summary		서버 상태 확인
// @Description	서버의 기본 상태 및 가동 시간을 확인합니다
// @Tags		Health
// @Accept		json
// @Produce		json
// @Success		200	{object}	map[string]interface{}	"서버 정상"
// @Router		/health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	uptime := time.Since(h.startTime)

	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Uptime:    uptime.String(),
		Version:   "0.1.0", // 추후 환경변수나 빌드 정보에서 가져오기
	}

	util.SuccessResponse(c, http.StatusOK, "서버가 정상 작동 중입니다", response)
}

// Readiness - Readiness Probe (K8s용)
// @Summary		Readiness Probe
// @Description	서버의 준비 상태를 확인합니다 (DB, Redis 등 외부 의존성 체크)
// @Tags		Health
// @Accept		json
// @Produce		json
// @Success		200	{object}	map[string]interface{}	"서버 준비됨"
// @Failure		503	{object}	map[string]interface{}	"서버 준비되지 않음"
// @Router		/health/ready [get]
func (h *HealthHandler) Readiness(c *gin.Context) {
	// TODO: 데이터베이스 연결 확인
	// TODO: Redis 연결 확인

	// 현재는 기본 체크만
	response := map[string]interface{}{
		"status": "ready",
		"checks": map[string]string{
			"database": "not_implemented",
			"redis":    "not_implemented",
		},
	}

	util.SuccessResponse(c, http.StatusOK, "서버가 준비되었습니다", response)
}

// Liveness - Liveness Probe (K8s용)
// @Summary		Liveness Probe
// @Description	서버가 살아있는지 확인합니다 (가장 가벼운 체크)
// @Tags		Health
// @Accept		json
// @Produce		json
// @Success		200	{object}	map[string]interface{}	"서버 살아있음"
// @Router		/health/live [get]
func (h *HealthHandler) Liveness(c *gin.Context) {
	util.SuccessResponse(c, http.StatusOK, "서버가 살아있습니다", map[string]string{
		"status": "alive",
	})
}
