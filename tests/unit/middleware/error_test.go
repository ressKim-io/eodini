package middleware_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hyeokjun/eodini/internal/middleware"
	"github.com/hyeokjun/eodini/internal/util"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// TestErrorHandler_WithAppError - AppError 처리 테스트
func TestErrorHandler_WithAppError(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(util.NewNotFoundError("차량"))
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "NOT_FOUND")
	assert.Contains(t, w.Body.String(), "차량")
}

// TestErrorHandler_WithGenericError - 일반 에러 처리 테스트
func TestErrorHandler_WithGenericError(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.New("generic error"))
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestErrorHandler_NoError - 에러 없는 경우
func TestErrorHandler_NoError(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.ErrorHandler())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}

// TestRecoveryHandler_WithPanic - Panic 복구 테스트
func TestRecoveryHandler_WithPanic(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.RecoveryHandler())

	router.GET("/test", func(c *gin.Context) {
		panic("something went wrong")
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestRecoveryHandler_WithErrorPanic - error 타입 Panic 테스트
func TestRecoveryHandler_WithErrorPanic(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.RecoveryHandler())

	router.GET("/test", func(c *gin.Context) {
		panic(errors.New("error panic"))
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "INTERNAL_ERROR")
}

// TestRecoveryHandler_NoPanic - Panic 없는 경우
func TestRecoveryHandler_NoPanic(t *testing.T) {
	// Given
	router := gin.New()
	router.Use(middleware.RecoveryHandler())

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "success"})
	})

	// When
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Then
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "success")
}
