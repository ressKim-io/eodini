package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hyeokjun/eodini/internal/middleware"
	"github.com/stretchr/testify/assert"
)

// TestCORS_DefaultConfig - 기본 CORS 설정 테스트
func TestCORS_DefaultConfig(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Equal(t, "true", w.Header().Get("Access-Control-Allow-Credentials"))
}

// TestCORS_ProductionConfig_AllowedOrigin - 허용된 Origin
func TestCORS_ProductionConfig_AllowedOrigin(t *testing.T) {
	// Given
	allowedOrigins := []string{"https://example.com", "https://app.example.com"}
	config := middleware.ProductionCORSConfig(allowedOrigins)

	router := gin.New()
	router.Use(middleware.CORS(config))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "https://example.com", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORS_ProductionConfig_DisallowedOrigin - 허용되지 않은 Origin
func TestCORS_ProductionConfig_DisallowedOrigin(t *testing.T) {
	// Given
	allowedOrigins := []string{"https://example.com"}
	config := middleware.ProductionCORSConfig(allowedOrigins)

	router := gin.New()
	router.Use(middleware.CORS(config))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "https://malicious.com")
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusForbidden, w.Code)
}

// TestCORS_PreflightRequest - OPTIONS 요청 (Preflight) 테스트
func TestCORS_PreflightRequest(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
}

// TestCORS_ExposeHeaders - Expose Headers 테스트
func TestCORS_ExposeHeaders(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	router.ServeHTTP(w, req)

	// Then
	assert.Contains(t, w.Header().Get("Access-Control-Expose-Headers"), "X-Request-ID")
}

// TestCORS_NoOriginHeader - Origin 헤더가 없는 경우
func TestCORS_NoOriginHeader(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.CORS(middleware.DefaultCORSConfig()))

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	// Origin 헤더를 설정하지 않음
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
}
