package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// 📝 설명: 애플리케이션 설정 중앙 관리
// 🎯 실무 포인트: 환경변수로 설정을 주입받아 K8s ConfigMap/Secret과 연동
// ⚠️ 주의사항: 민감한 정보(DB 비밀번호 등)는 반드시 환경변수로 주입

// Config - 전체 애플리케이션 설정
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Log      LogConfig
}

// ServerConfig - 서버 관련 설정
type ServerConfig struct {
	Port         string        // HTTP 서버 포트 (예: "8080")
	Host         string        // 호스트 주소 (예: "0.0.0.0")
	Environment  string        // 환경 (dev, staging, prod)
	ReadTimeout  time.Duration // 요청 읽기 타임아웃
	WriteTimeout time.Duration // 응답 쓰기 타임아웃
	IdleTimeout  time.Duration // 유휴 연결 타임아웃
}

// DatabaseConfig - 데이터베이스 관련 설정
type DatabaseConfig struct {
	Host            string // DB 호스트 (예: "localhost")
	Port            string // DB 포트 (예: "5432")
	User            string // DB 사용자명
	Password        string // DB 비밀번호
	DBName          string // 데이터베이스 이름
	SSLMode         string // SSL 모드 (disable, require, verify-ca, verify-full)
	MaxOpenConns    int    // 최대 오픈 커넥션 수
	MaxIdleConns    int    // 최대 유휴 커넥션 수
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// RedisConfig - Redis 관련 설정
type RedisConfig struct {
	Host     string // Redis 호스트
	Port     string // Redis 포트
	Password string // Redis 비밀번호 (선택)
	DB       int    // Redis DB 번호 (0-15)
}

// LogConfig - 로그 관련 설정
type LogConfig struct {
	Level  string // 로그 레벨 (debug, info, warn, error)
	Format string // 로그 포맷 (json, text)
}

// Load - 환경변수에서 설정 로드
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("SERVER_PORT", "8080"),
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Environment:  getEnv("ENVIRONMENT", "dev"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
			IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", ""),
			DBName:          getEnv("DB_NAME", "eodini"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getIntEnv("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DB_CONN_MAX_LIFETIME", 5*time.Minute),
			ConnMaxIdleTime: getDurationEnv("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getIntEnv("REDIS_DB", 0),
		},
		Log: LogConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "text"),
		},
	}

	// 설정 검증
	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate - 설정 검증
func (c *Config) Validate() error {
	// 필수 값 검증
	if c.Server.Port == "" {
		return fmt.Errorf("SERVER_PORT is required")
	}

	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}

	if c.Database.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// 환경 검증
	validEnvs := map[string]bool{"dev": true, "staging": true, "prod": true}
	if !validEnvs[c.Server.Environment] {
		return fmt.Errorf("invalid ENVIRONMENT: %s (must be dev, staging, or prod)", c.Server.Environment)
	}

	// 로그 레벨 검증
	validLogLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLogLevels[c.Log.Level] {
		return fmt.Errorf("invalid LOG_LEVEL: %s (must be debug, info, warn, or error)", c.Log.Level)
	}

	return nil
}

// GetDatabaseDSN - PostgreSQL DSN 생성
func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// GetRedisAddr - Redis 주소 생성
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.Redis.Host, c.Redis.Port)
}

// IsDevelopment - 개발 환경 여부
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "dev"
}

// IsProduction - 프로덕션 환경 여부
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "prod"
}

// getEnv - 환경변수 조회 (기본값 포함)
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getIntEnv - 정수형 환경변수 조회
func getIntEnv(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// getDurationEnv - Duration 환경변수 조회
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
