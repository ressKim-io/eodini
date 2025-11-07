.PHONY: help run test test-unit test-integration test-coverage clean build docker-build docker-run swagger swagger-install

# 기본 변수
APP_NAME=eodini
MAIN_PATH=cmd/api/main.go
BINARY_NAME=eodini-api

help: ## 도움말 표시
	@echo "Eodini - 통학/통원 차량 관리 시스템"
	@echo ""
	@echo "사용 가능한 명령어:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

run: ## 서버 실행
	@echo "🚀 Eodini API Server 실행 중..."
	@go run $(MAIN_PATH)

build: ## 바이너리 빌드
	@echo "🔨 바이너리 빌드 중..."
	@go build -o bin/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ 빌드 완료: bin/$(BINARY_NAME)"

test: ## 전체 테스트 실행
	@echo "🧪 전체 테스트 실행 중..."
	@go test ./... -v

test-unit: ## 단위 테스트만 실행
	@echo "🧪 단위 테스트 실행 중..."
	@go test ./tests/unit/... -v

test-integration: ## 통합 테스트만 실행
	@echo "🧪 통합 테스트 실행 중..."
	@go test ./tests/integration/... -v

test-coverage: ## 테스트 커버리지 확인
	@echo "📊 테스트 커버리지 확인 중..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -func=coverage.out | grep total
	@echo ""
	@echo "HTML 리포트 생성: coverage.html"
	@go tool cover -html=coverage.out -o coverage.html

clean: ## 빌드 파일 정리
	@echo "🧹 빌드 파일 정리 중..."
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@echo "✅ 정리 완료"

install-deps: ## Go 의존성 설치
	@echo "📦 의존성 설치 중..."
	@go mod download
	@go mod tidy
	@echo "✅ 의존성 설치 완료"

fmt: ## 코드 포맷팅
	@echo "✨ 코드 포맷팅 중..."
	@go fmt ./...
	@echo "✅ 포맷팅 완료"

lint: ## 코드 린트 (golangci-lint 필요)
	@echo "🔍 코드 린트 중..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run; \
	else \
		echo "⚠️  golangci-lint가 설치되지 않았습니다."; \
		echo "설치: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

vet: ## Go vet 실행
	@echo "🔍 코드 검사 중..."
	@go vet ./...
	@echo "✅ 검사 완료"

docker-build: ## Docker 이미지 빌드
	@echo "🐳 Docker 이미지 빌드 중..."
	@docker build -t $(APP_NAME):latest .
	@echo "✅ Docker 이미지 빌드 완료"

docker-run: ## Docker 컨테이너 실행
	@echo "🐳 Docker 컨테이너 실행 중..."
	@docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

db-up: ## Docker Compose로 DB 시작 (추후)
	@echo "🐳 데이터베이스 시작 중..."
	@docker-compose up -d postgres redis
	@echo "✅ 데이터베이스 시작 완료"

db-down: ## Docker Compose DB 중지 (추후)
	@echo "🐳 데이터베이스 중지 중..."
	@docker-compose down
	@echo "✅ 데이터베이스 중지 완료"

migrate-up: ## 데이터베이스 마이그레이션 (추후)
	@echo "🔄 마이그레이션 실행 중..."
	@echo "⚠️  마이그레이션 기능은 추후 구현 예정"

migrate-down: ## 마이그레이션 롤백 (추후)
	@echo "🔄 마이그레이션 롤백 중..."
	@echo "⚠️  마이그레이션 기능은 추후 구현 예정"

swagger-install: ## Swagger CLI 설치
	@echo "📦 Swagger CLI 설치 중..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@echo "✅ Swagger CLI 설치 완료"

swagger: ## Swagger 문서 생성
	@echo "📚 Swagger 문서 생성 중..."
	@if command -v swag > /dev/null; then \
		swag init -g $(MAIN_PATH) -o docs; \
	elif [ -f ~/go/bin/swag ]; then \
		~/go/bin/swag init -g $(MAIN_PATH) -o docs; \
	else \
		echo "⚠️  swag CLI가 설치되지 않았습니다."; \
		echo "설치: make swagger-install"; \
		exit 1; \
	fi
	@echo "✅ Swagger 문서 생성 완료 (docs/)"
	@echo "📖 Swagger UI: http://localhost:8080/swagger/index.html"

.DEFAULT_GOAL := help
