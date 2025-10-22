package config_test

import (
	"os"
	"testing"

	"github.com/hyeokjun/eodini/config"
	"github.com/stretchr/testify/assert"
)

// TestLoad_DefaultValues - 기본값으로 설정 로드
func TestLoad_DefaultValues(t *testing.T) {
	// Given: 환경변수 초기화
	clearEnv()

	// When
	cfg, err := config.Load()

	// Then
	assert.NoError(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, "8080", cfg.Server.Port)
	assert.Equal(t, "0.0.0.0", cfg.Server.Host)
	assert.Equal(t, "dev", cfg.Server.Environment)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, "5432", cfg.Database.Port)
	assert.Equal(t, "postgres", cfg.Database.User)
	assert.Equal(t, "eodini", cfg.Database.DBName)
}

// TestLoad_WithEnvironmentVariables - 환경변수로 설정 로드
func TestLoad_WithEnvironmentVariables(t *testing.T) {
	// Given
	clearEnv()
	os.Setenv("SERVER_PORT", "9000")
	os.Setenv("SERVER_HOST", "127.0.0.1")
	os.Setenv("ENVIRONMENT", "prod")
	os.Setenv("DB_HOST", "db.example.com")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "myuser")
	os.Setenv("DB_PASSWORD", "mypassword")
	os.Setenv("DB_NAME", "mydb")
	defer clearEnv()

	// When
	cfg, err := config.Load()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, "9000", cfg.Server.Port)
	assert.Equal(t, "127.0.0.1", cfg.Server.Host)
	assert.Equal(t, "prod", cfg.Server.Environment)
	assert.Equal(t, "db.example.com", cfg.Database.Host)
	assert.Equal(t, "5433", cfg.Database.Port)
	assert.Equal(t, "myuser", cfg.Database.User)
	assert.Equal(t, "mypassword", cfg.Database.Password)
	assert.Equal(t, "mydb", cfg.Database.DBName)
}

// TestLoad_IntegerEnvironmentVariables - 정수형 환경변수
func TestLoad_IntegerEnvironmentVariables(t *testing.T) {
	// Given
	clearEnv()
	os.Setenv("DB_MAX_OPEN_CONNS", "50")
	os.Setenv("DB_MAX_IDLE_CONNS", "10")
	os.Setenv("REDIS_DB", "5")
	defer clearEnv()

	// When
	cfg, err := config.Load()

	// Then
	assert.NoError(t, err)
	assert.Equal(t, 50, cfg.Database.MaxOpenConns)
	assert.Equal(t, 10, cfg.Database.MaxIdleConns)
	assert.Equal(t, 5, cfg.Redis.DB)
}

// TestValidate_Success - 유효한 설정
func TestValidate_Success(t *testing.T) {
	// Given
	clearEnv()

	// When
	cfg, err := config.Load()

	// Then
	assert.NoError(t, err)
	assert.NoError(t, cfg.Validate())
}

// TestValidate_InvalidEnvironment - 잘못된 환경 값
func TestValidate_InvalidEnvironment(t *testing.T) {
	// Given
	clearEnv()
	os.Setenv("ENVIRONMENT", "invalid")
	defer clearEnv()

	// When
	cfg, err := config.Load()

	// Then
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "invalid ENVIRONMENT")
}

// TestValidate_InvalidLogLevel - 잘못된 로그 레벨
func TestValidate_InvalidLogLevel(t *testing.T) {
	// Given
	clearEnv()
	os.Setenv("LOG_LEVEL", "invalid")
	defer clearEnv()

	// When
	cfg, err := config.Load()

	// Then
	assert.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "invalid LOG_LEVEL")
}

// TestGetDatabaseDSN - PostgreSQL DSN 생성
func TestGetDatabaseDSN(t *testing.T) {
	// Given
	clearEnv()
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("DB_USER", "testuser")
	os.Setenv("DB_PASSWORD", "testpass")
	os.Setenv("DB_NAME", "testdb")
	os.Setenv("DB_SSL_MODE", "disable")
	defer clearEnv()

	cfg, err := config.Load()
	assert.NoError(t, err)

	// When
	dsn := cfg.GetDatabaseDSN()

	// Then
	expected := "host=localhost port=5432 user=testuser password=testpass dbname=testdb sslmode=disable"
	assert.Equal(t, expected, dsn)
}

// TestGetRedisAddr - Redis 주소 생성
func TestGetRedisAddr(t *testing.T) {
	// Given
	clearEnv()
	os.Setenv("REDIS_HOST", "redis.example.com")
	os.Setenv("REDIS_PORT", "6380")
	defer clearEnv()

	cfg, err := config.Load()
	assert.NoError(t, err)

	// When
	addr := cfg.GetRedisAddr()

	// Then
	assert.Equal(t, "redis.example.com:6380", addr)
}

// TestIsDevelopment - 개발 환경 확인
func TestIsDevelopment(t *testing.T) {
	// Given
	clearEnv()
	os.Setenv("ENVIRONMENT", "dev")
	defer clearEnv()

	cfg, err := config.Load()
	assert.NoError(t, err)

	// When & Then
	assert.True(t, cfg.IsDevelopment())
	assert.False(t, cfg.IsProduction())
}

// TestIsProduction - 프로덕션 환경 확인
func TestIsProduction(t *testing.T) {
	// Given
	clearEnv()
	os.Setenv("ENVIRONMENT", "prod")
	defer clearEnv()

	cfg, err := config.Load()
	assert.NoError(t, err)

	// When & Then
	assert.True(t, cfg.IsProduction())
	assert.False(t, cfg.IsDevelopment())
}

// clearEnv - 테스트용 환경변수 초기화
func clearEnv() {
	envVars := []string{
		"SERVER_PORT", "SERVER_HOST", "ENVIRONMENT",
		"SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "SERVER_IDLE_TIMEOUT",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_SSL_MODE",
		"DB_MAX_OPEN_CONNS", "DB_MAX_IDLE_CONNS",
		"DB_CONN_MAX_LIFETIME", "DB_CONN_MAX_IDLE_TIME",
		"REDIS_HOST", "REDIS_PORT", "REDIS_PASSWORD", "REDIS_DB",
		"LOG_LEVEL", "LOG_FORMAT",
	}

	for _, key := range envVars {
		os.Unsetenv(key)
	}
}
