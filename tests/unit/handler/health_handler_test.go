package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hyeokjun/eodini/internal/handler"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Gin 테스트 모드 설정
	// gin.SetMode(gin.TestMode)는 handler 패키지에서 이미 설정됨
}

// TestHealth - 기본 Health Check 테스트
func TestHealth(t *testing.T) {
	// Given
	router := handler.SetupRouter()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Contains(t, response, "data")

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "healthy", data["status"])
	assert.Contains(t, data, "uptime")
	assert.Contains(t, data, "timestamp")
	assert.Contains(t, data, "version")
}

// TestReadiness - Readiness Probe 테스트
func TestReadiness(t *testing.T) {
	// Given
	router := handler.SetupRouter()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health/ready", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))
	assert.Contains(t, response, "data")

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "ready", data["status"])
	assert.Contains(t, data, "checks")
}

// TestLiveness - Liveness Probe 테스트
func TestLiveness(t *testing.T) {
	// Given
	router := handler.SetupRouter()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health/live", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.True(t, response["success"].(bool))

	data := response["data"].(map[string]interface{})
	assert.Equal(t, "alive", data["status"])
}

// TestPingEndpoint - 임시 Ping 엔드포인트 테스트
func TestPingEndpoint(t *testing.T) {
	// Given
	router := handler.SetupRouter()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/v1/ping", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "pong", response["message"])
}

// TestHealthCheckResponseFormat - 응답 포맷 확인
func TestHealthCheckResponseFormat(t *testing.T) {
	// Given
	router := handler.SetupRouter()

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	router.ServeHTTP(w, req)

	// Then
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 표준 응답 포맷 확인
	assert.Contains(t, response, "success")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")
	assert.NotContains(t, response, "error")
}
