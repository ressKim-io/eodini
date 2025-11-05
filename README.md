# Eodini

**통학/통원 차량 관리 시스템**

> Kubernetes 기반 + Go로 구현한 차량 관리 플랫폼

## 주요 기능

- 차량 관리
- 노선 관리
- 실시간 위치 추적
- 방에서 중요 공지 및 소통 가능

## 기술 스택

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin
- **Database**: PostgreSQL
- **Cache**: Redis
- **Architecture**: Clean Architecture

### Infrastructure
- **Container**: Docker
- **Orchestration**: Kubernetes (k3s / EKS)
- **IaC**: Terraform
- **Package Manager**: Helm

### CI/CD
- **CI**: GitHub Actions
- **CD**: ArgoCD

### Monitoring
- **Metrics**: Prometheus + Grafana
- **Logging**: Loki
- **Tracing**: Tempo

## 실행 방법

```bash
# 저장소 클론
git clone https://github.com/yourusername/eodini.git
cd eodini

# 의존성 설치
go mod download

# 서버 실행
go run cmd/api/main.go
```

## 테스트

```bash
# 전체 테스트 실행
go test ./...

# 커버리지 확인
go test ./... -cover
```

## 작성자

hyeokjun
