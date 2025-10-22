package logger_test

import (
	"testing"

	"github.com/hyeokjun/eodini/pkg/logger"
)

// TestSetLevel - 로그 레벨 설정 테스트
func TestSetLevel(t *testing.T) {
	// Given & When
	logger.SetLevel(logger.DebugLevel)

	// Then: 에러 없이 실행되면 성공
	logger.Debug("debug message", nil)
	logger.Info("info message", nil)
	logger.Warn("warn message", nil)
	logger.Error("error message", nil)
}

// TestDebug - Debug 로그 테스트
func TestDebug(t *testing.T) {
	// Given
	logger.SetLevel(logger.DebugLevel)

	// When & Then: 에러 없이 실행되면 성공
	logger.Debug("test debug", map[string]interface{}{
		"key": "value",
	})
}

// TestInfo - Info 로그 테스트
func TestInfo(t *testing.T) {
	// Given
	logger.SetLevel(logger.InfoLevel)

	// When & Then
	logger.Info("test info", map[string]interface{}{
		"status": "success",
	})
}

// TestWarn - Warn 로그 테스트
func TestWarn(t *testing.T) {
	// Given
	logger.SetLevel(logger.WarnLevel)

	// When & Then
	logger.Warn("test warning", map[string]interface{}{
		"code": 400,
	})
}

// TestError - Error 로그 테스트
func TestError(t *testing.T) {
	// Given
	logger.SetLevel(logger.ErrorLevel)

	// When & Then
	logger.Error("test error", map[string]interface{}{
		"error": "something went wrong",
	})
}

// TestInfof - 포맷팅된 Info 로그 테스트
func TestInfof(t *testing.T) {
	// Given
	logger.SetLevel(logger.InfoLevel)

	// When & Then
	logger.Infof("User %s logged in from %s", "john", "192.168.1.1")
}

// TestErrorf - 포맷팅된 Error 로그 테스트
func TestErrorf(t *testing.T) {
	// Given
	logger.SetLevel(logger.ErrorLevel)

	// When & Then
	logger.Errorf("Failed to connect to database: %s", "connection timeout")
}

// TestWarnf - 포맷팅된 Warn 로그 테스트
func TestWarnf(t *testing.T) {
	// Given
	logger.SetLevel(logger.WarnLevel)

	// When & Then
	logger.Warnf("API rate limit exceeded: %d requests", 1000)
}

// TestLogWithNilFields - nil 필드로 로그 테스트
func TestLogWithNilFields(t *testing.T) {
	// Given
	logger.SetLevel(logger.InfoLevel)

	// When & Then
	logger.Info("message without fields", nil)
}

// TestLogLevelFiltering - 로그 레벨 필터링 테스트
func TestLogLevelFiltering(t *testing.T) {
	// Given: Error 레벨로 설정
	logger.SetLevel(logger.ErrorLevel)

	// When & Then: Debug, Info, Warn은 출력 안 됨 (에러 없으면 성공)
	logger.Debug("should not print", nil)
	logger.Info("should not print", nil)
	logger.Warn("should not print", nil)
	logger.Error("should print", nil)
}
