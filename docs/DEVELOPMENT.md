# Eodini - 개발 상태 문서

> 마지막 업데이트: 2025-10-22

## 📊 전체 진행 상황

```
프로젝트 진행률: 30%
```

### ✅ 완료된 작업

#### Phase 0: 프로젝트 초기 설정
- [x] Git 저장소 초기화
- [x] .gitignore 설정 (Go, IDE, Terraform, K8s)
- [x] README.md 작성
- [x] Clean Architecture 디렉토리 구조 생성
- [x] Go 모듈 초기화 (`github.com/hyeokjun/eodini`)

#### Phase 1: 핵심 유틸리티
- [x] 에러 처리 구조 (Spring `@ControllerAdvice` 스타일)
  - `internal/util/error.go` - 9가지 에러 타입
  - AppError 구조체 (Code, Message, StatusCode, Details)
- [x] 메시지 관리 (Spring `messages.properties` 스타일)
  - `internal/util/message.go` - 중앙 집중식 메시지 관리
  - 포맷팅 지원, 다국어 확장 가능
- [x] 응답 포맷 표준화
  - `internal/util/response.go` - APIResponse 구조
  - 페이지네이션 지원
- [x] 단위 테스트 (18개 테스트, 커버리지 85.7%)

#### Phase 2: 미들웨어
- [x] 글로벌 에러 핸들러 (`middleware/error.go`)
  - ErrorHandler: 모든 에러 중앙 처리
  - RecoveryHandler: Panic 복구
- [x] 로거 미들웨어 (`middleware/logger.go`)
  - RequestLogger: HTTP 요청/응답 로깅
  - RequestIDMiddleware: 분산 추적용 ID 부여
- [x] CORS 미들웨어 (`middleware/cors.go`)
  - 개발/프로덕션 설정 분리
  - Preflight 요청 자동 처리
- [x] 간단한 구조화 로거 (`pkg/logger`)
  - 로그 레벨 지원 (Debug/Info/Warn/Error/Fatal)
  - 필드 기반 로깅
- [x] 단위 테스트 (22개 테스트, 커버리지 75.7%)

#### Phase 3: 설정 관리
- [x] 환경변수 기반 설정 (`config/config.go`)
  - Server, Database, Redis, Log 설정 통합
  - 유효성 검증 로직
  - K8s ConfigMap/Secret 연동 가능
- [x] .env.example 파일
- [x] 단위 테스트 (10개 테스트, 커버리지 76.4%)

#### Phase 4: 도메인 모델
- [x] **Vehicle** (차량) - 145줄
  - 차량 정보, 상태 관리, 보험/검사 만료일
  - 운행 가능 여부 검증
- [x] **Driver** (기사) - 124줄
  - 기본 정보, 면허 정보, 상태 관리
  - 휴가/퇴사 처리, 면허 갱신 알림
- [x] **Attendant** (동승자/선생님) - 139줄
  - 역할(선생님/간호사 등), 운행 시작 권한
  - 상태 관리, 소속 기관
- [x] **Route & Stop** (경로/정류장) - 165줄
  - Route: A코스, B코스 같은 고정 동선
  - Stop: 정류장 순서, 위치, 예상 도착 시간
- [x] **Passenger** (탑승자) - 145줄
  - 유치원생/환자 정보, 보호자 연락처
  - 정류장 배정, 의료 특이사항
- [x] **Schedule** (운행 일정 템플릿) - 196줄
  - 반복 운행 템플릿 ("매일 오전 8시 A코스")
  - 요일 설정, 기본 담당자 배정
  - Trip 생성의 기준
- [x] **Trip** (실제 운행) ⭐ - 214줄
  - 날짜별 실제 운행 기록
  - 대체 기사 배정 가능
  - 운행 시작/완료, 탑승/하차 기록
- [x] **DriverAssignment** (대체 배정) ⭐ - 134줄
  - A 기사 → G 기사 대체
  - 기간 설정, 승인 프로세스
- [x] 총 **1,340줄**의 도메인 코드

#### Phase 5: Health Check API
- [x] 애플리케이션 진입점 (`cmd/api/main.go`)
  - 설정 로드, 로거 초기화
  - Graceful Shutdown (30초)
  - 시그널 처리
- [x] Health Check Handler
  - GET /health - 기본 상태 확인
  - GET /health/ready - K8s Readiness Probe
  - GET /health/live - K8s Liveness Probe
- [x] 라우터 설정 (`internal/handler/router.go`)
  - 미들웨어 적용 순서 정의
  - API v1 그룹 설정
- [x] 단위 테스트 (5개 테스트)

### 🚧 진행 중인 작업

- [ ] 문서화
- [ ] 도메인 모델 개선 (UUID, GORM 태그)

### 📌 다음 단계 (우선순위 순)

#### Phase 6: 도메인 모델 개선
- [ ] UUID 자동 생성 추가
- [ ] GORM 태그 추가 (데이터베이스 매핑)
- [ ] 타임스탬프 자동 업데이트 개선

#### Phase 7: 데이터베이스 연결
- [ ] PostgreSQL 연결 설정 (`pkg/database`)
- [ ] GORM 설정
- [ ] 자동 마이그레이션
- [ ] 연결 풀 설정
- [ ] Health Check에 DB 상태 추가

#### Phase 8: Vehicle CRUD API
- [ ] VehicleRepository 구현
- [ ] VehicleService 구현
- [ ] VehicleHandler 구현
- [ ] 단위 테스트 및 통합 테스트
- [ ] API 문서 작성

#### Phase 9: Driver CRUD API
- [ ] DriverRepository 구현
- [ ] DriverService 구현
- [ ] DriverHandler 구현
- [ ] 테스트 작성

#### Phase 10: Route & Stop API
- [ ] RouteRepository 구현
- [ ] RouteService 구현
- [ ] RouteHandler 구현
- [ ] Stop 관리 로직
- [ ] 테스트 작성

#### Phase 11: 복잡한 API 구현
- [ ] Schedule API
- [ ] Trip API (핵심!)
- [ ] DriverAssignment API
- [ ] Passenger API

#### Phase 12: Redis 연결 (캐싱)
- [ ] Redis 연결 설정
- [ ] 캐싱 로직
- [ ] Session 관리 (추후)

#### Phase 13: 인증/인가 (추후)
- [ ] JWT 인증
- [ ] Role 기반 권한 관리
- [ ] 미들웨어 추가

#### Phase 14: 컨테이너화
- [ ] Dockerfile 작성 (멀티스테이지)
- [ ] docker-compose.yml (로컬 개발)
- [ ] .dockerignore

#### Phase 15: K8s 배포 (로컬 k3s)
- [ ] Deployment YAML
- [ ] Service YAML
- [ ] ConfigMap/Secret
- [ ] HPA
- [ ] Ingress
- [ ] Network Policy

#### Phase 16: Helm Chart (사용자 작성)
- [ ] Chart.yaml
- [ ] values.yaml
- [ ] templates/

#### Phase 17: Terraform (AWS EKS)
- [ ] VPC 모듈
- [ ] EKS 클러스터
- [ ] RDS (PostgreSQL)
- [ ] ElastiCache (Redis)
- [ ] 보안 그룹

#### Phase 18: CI/CD (사용자 작성)
- [ ] GitHub Actions (빌드/테스트)
- [ ] ArgoCD 설정

#### Phase 19: 모니터링
- [ ] Prometheus + Grafana
- [ ] Loki (로그 수집)
- [ ] Tempo (분산 추적)

## 🐛 알려진 이슈

- 없음 (현재까지)

## 📝 개선 필요 항목

1. **도메인 모델**
   - UUID 자동 생성 로직 없음
   - GORM 태그 없음 (데이터베이스 연결 전에 추가 필요)
   - 타임스탬프 자동 업데이트 개선 가능

2. **로거**
   - 현재 간단한 구조, 프로덕션에서는 zap/logrus 권장
   - JSON 포맷 미구현

3. **테스트**
   - 통합 테스트 없음 (단위 테스트만)
   - E2E 테스트 없음

4. **문서**
   - API 문서 없음
   - Swagger/OpenAPI 스펙 없음

## 🎯 테스트 커버리지 목표

- **전체**: 60% 이상 ✅ (현재 76.4%)
- **핵심 비즈니스 로직**: 80% 이상
- **유틸리티/헬퍼**: 70% 이상
- **핸들러**: 50% 이상 (통합 테스트로 보완)

## 📦 의존성

### 현재 사용 중
```
github.com/gin-gonic/gin v1.11.0
github.com/stretchr/testify v1.11.1
```

### 추가 예정
```
gorm.io/gorm (PostgreSQL ORM)
gorm.io/driver/postgres
github.com/google/uuid (UUID 생성)
github.com/redis/go-redis/v9 (Redis)
github.com/golang-jwt/jwt/v5 (JWT 인증 - 추후)
```

## 🚀 빠른 시작

```bash
# 서버 실행
go run cmd/api/main.go

# 테스트
go test ./...

# 커버리지
go test ./... -cover

# Health Check
curl http://localhost:8080/health
```

## 📂 프로젝트 구조

```
eodini/
├── cmd/api/              # 진입점
├── internal/
│   ├── domain/          # 도메인 모델 (8개) ✅
│   ├── handler/         # HTTP 핸들러
│   ├── service/         # 비즈니스 로직
│   ├── repository/      # 데이터 액세스
│   ├── middleware/      # 미들웨어 ✅
│   └── util/            # 유틸리티 ✅
├── pkg/                 # 공용 패키지
│   ├── database/        # DB 연결
│   ├── cache/           # Redis
│   └── logger/          # 로거 ✅
├── config/              # 설정 ✅
├── tests/               # 테스트
├── docs/                # 문서 (이 파일)
├── deployments/         # K8s, Helm (추후)
└── scripts/             # 유틸리티 스크립트 (추후)
```

## 🔗 관련 문서

- [ARCHITECTURE.md](./ARCHITECTURE.md) - 아키텍처 설명
- [API.md](./API.md) - API 문서 (추후)
- [GETTING_STARTED.md](./GETTING_STARTED.md) - 실행 가이드 (추후)

## 📞 문의

- 작성자: hyeokjun
- 프로젝트: 포트폴리오 & K8s 학습용
