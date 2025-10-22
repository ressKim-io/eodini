package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hyeokjun/eodini/pkg/logger"
)

// 📝 설명: 모든 HTTP 요청/응답을 로깅하는 미들웨어
// 🎯 실무 포인트: 요청 시간, 응답 시간, 상태 코드, 에러 등 기록
// ⚠️ 주의사항: 민감한 정보(비밀번호 등)는 로그에서 제외

// RequestLogger - HTTP 요청/응답 로깅 미들웨어
//
// 로그 내용:
// - 메소드, 경로, 상태 코드
// - 처리 시간 (latency)
// - 클라이언트 IP
// - User-Agent
// - 에러 메시지 (있을 경우)
//
// 사용 예:
//   router := gin.Default()
//   router.Use(middleware.RequestLogger())
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 시작 시간 기록
		startTime := time.Now()

		// 요청 정보 추출
		path := c.Request.URL.Path
		method := c.Request.Method
		clientIP := c.ClientIP()
		userAgent := c.Request.UserAgent()

		// 다음 핸들러 실행
		c.Next()

		// 종료 시간 및 처리 시간 계산
		endTime := time.Now()
		latency := endTime.Sub(startTime)

		// 응답 정보
		statusCode := c.Writer.Status()

		// 로그 레벨 결정
		// 4xx: Warning, 5xx: Error, 나머지: Info
		logFunc := logger.Info
		if statusCode >= 400 && statusCode < 500 {
			logFunc = logger.Warn
		} else if statusCode >= 500 {
			logFunc = logger.Error
		}

		// 기본 로그 필드
		fields := map[string]interface{}{
			"method":     method,
			"path":       path,
			"status":     statusCode,
			"latency_ms": latency.Milliseconds(),
			"client_ip":  clientIP,
			"user_agent": userAgent,
		}

		// 에러가 있으면 추가
		if len(c.Errors) > 0 {
			fields["error"] = c.Errors.String()
		}

		// 로그 출력
		logFunc("HTTP Request", fields)
	}
}

// RequestIDMiddleware - 요청마다 고유 ID 부여
// 분산 추적(Distributed Tracing)에 유용
//
// 사용 예:
//   router := gin.Default()
//   router.Use(middleware.RequestIDMiddleware())
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// X-Request-ID 헤더에서 가져오거나 새로 생성
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		// Context에 저장
		c.Set("request_id", requestID)

		// 응답 헤더에도 추가
		c.Writer.Header().Set("X-Request-ID", requestID)

		c.Next()
	}
}

// generateRequestID - 간단한 요청 ID 생성
// 실무에서는 UUID 사용 권장
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString - 랜덤 문자열 생성 (간단한 버전)
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}
