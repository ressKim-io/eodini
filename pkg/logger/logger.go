package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// 📝 설명: 간단한 구조화된 로거 (추후 zap, logrus 등으로 교체 가능)
// 🎯 실무 포인트: 로그 레벨 분리, 구조화된 필드 지원
// ⚠️ 주의사항: 프로덕션에서는 zap 등 성능 좋은 라이브러리 사용 권장

// LogLevel - 로그 레벨
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var (
	currentLevel = InfoLevel
	logger       = log.New(os.Stdout, "", 0)
)

// SetLevel - 로그 레벨 설정
func SetLevel(level LogLevel) {
	currentLevel = level
}

// Debug - 디버그 로그
func Debug(message string, fields map[string]interface{}) {
	if currentLevel <= DebugLevel {
		logWithFields("DEBUG", message, fields)
	}
}

// Info - 정보 로그
func Info(message string, fields map[string]interface{}) {
	if currentLevel <= InfoLevel {
		logWithFields("INFO", message, fields)
	}
}

// Warn - 경고 로그
func Warn(message string, fields map[string]interface{}) {
	if currentLevel <= WarnLevel {
		logWithFields("WARN", message, fields)
	}
}

// Error - 에러 로그
func Error(message string, fields map[string]interface{}) {
	if currentLevel <= ErrorLevel {
		logWithFields("ERROR", message, fields)
	}
}

// Fatal - 치명적 에러 로그 (프로그램 종료)
func Fatal(message string, fields map[string]interface{}) {
	logWithFields("FATAL", message, fields)
	os.Exit(1)
}

// logWithFields - 필드와 함께 로그 출력
func logWithFields(level, message string, fields map[string]interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	logMessage := fmt.Sprintf("[%s] %s | %s", timestamp, level, message)

	if len(fields) > 0 {
		logMessage += " |"
		for key, value := range fields {
			logMessage += fmt.Sprintf(" %s=%v", key, value)
		}
	}

	logger.Println(logMessage)
}

// Infof - 포맷팅된 정보 로그 (필드 없음)
func Infof(format string, args ...interface{}) {
	if currentLevel <= InfoLevel {
		message := fmt.Sprintf(format, args...)
		logWithFields("INFO", message, nil)
	}
}

// Errorf - 포맷팅된 에러 로그 (필드 없음)
func Errorf(format string, args ...interface{}) {
	if currentLevel <= ErrorLevel {
		message := fmt.Sprintf(format, args...)
		logWithFields("ERROR", message, nil)
	}
}

// Warnf - 포맷팅된 경고 로그 (필드 없음)
func Warnf(format string, args ...interface{}) {
	if currentLevel <= WarnLevel {
		message := fmt.Sprintf(format, args...)
		logWithFields("WARN", message, nil)
	}
}
