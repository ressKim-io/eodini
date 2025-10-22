package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 📝 설명: Cross-Origin Resource Sharing (CORS) 설정
// 🎯 실무 포인트: 프론트엔드와 백엔드가 다른 도메인일 때 필수
// ⚠️ 주의사항: 프로덕션에서는 AllowOrigins를 특정 도메인으로 제한

// CORSConfig - CORS 설정 구조체
type CORSConfig struct {
	AllowOrigins     []string // 허용할 Origin (예: ["http://localhost:3000"])
	AllowMethods     []string // 허용할 HTTP 메소드
	AllowHeaders     []string // 허용할 헤더
	ExposeHeaders    []string // 노출할 헤더
	AllowCredentials bool     // 쿠키 포함 여부
	MaxAge           int      // Preflight 요청 캐시 시간 (초)
}

// DefaultCORSConfig - 기본 CORS 설정 (개발 환경용)
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowOrigins: []string{"*"}, // 모든 Origin 허용 (개발용)
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           43200, // 12시간
	}
}

// ProductionCORSConfig - 프로덕션 CORS 설정
// 사용 예: ProductionCORSConfig([]string{"https://example.com"})
func ProductionCORSConfig(allowedOrigins []string) CORSConfig {
	return CORSConfig{
		AllowOrigins: allowedOrigins, // 특정 도메인만 허용
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           43200,
	}
}

// CORS - CORS 미들웨어
//
// 사용 예:
//   // 개발 환경
//   router.Use(middleware.CORS(middleware.DefaultCORSConfig()))
//
//   // 프로덕션 환경
//   config := middleware.ProductionCORSConfig([]string{"https://example.com"})
//   router.Use(middleware.CORS(config))
func CORS(config CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Origin 체크
		if len(config.AllowOrigins) > 0 {
			allowed := false

			// "*" 이면 모든 Origin 허용
			if config.AllowOrigins[0] == "*" {
				allowed = true
				c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			} else {
				// 특정 Origin만 허용
				for _, allowedOrigin := range config.AllowOrigins {
					if origin == allowedOrigin {
						allowed = true
						c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
						break
					}
				}
			}

			if !allowed && origin != "" {
				c.AbortWithStatus(403)
				return
			}
		}

		// Allow Methods
		if len(config.AllowMethods) > 0 {
			methods := ""
			for i, method := range config.AllowMethods {
				if i > 0 {
					methods += ", "
				}
				methods += method
			}
			c.Writer.Header().Set("Access-Control-Allow-Methods", methods)
		}

		// Allow Headers
		if len(config.AllowHeaders) > 0 {
			headers := ""
			for i, header := range config.AllowHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Writer.Header().Set("Access-Control-Allow-Headers", headers)
		}

		// Expose Headers
		if len(config.ExposeHeaders) > 0 {
			headers := ""
			for i, header := range config.ExposeHeaders {
				if i > 0 {
					headers += ", "
				}
				headers += header
			}
			c.Writer.Header().Set("Access-Control-Expose-Headers", headers)
		}

		// Allow Credentials
		if config.AllowCredentials {
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		// Max Age
		if config.MaxAge > 0 {
			c.Writer.Header().Set("Access-Control-Max-Age", fmt.Sprintf("%d", config.MaxAge))
		}

		// OPTIONS 요청은 여기서 종료 (Preflight)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
