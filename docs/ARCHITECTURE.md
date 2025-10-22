# Eodini - 아키텍처 문서

> Clean Architecture 기반 통학/통원 차량 관리 시스템

## 🏗️ 전체 아키텍처

```
┌─────────────────────────────────────────────────────────────┐
│                         Client (App)                         │
│                    (Driver, Attendant, Admin)                │
└─────────────────────┬───────────────────────────────────────┘
                      │ HTTP/JSON
┌─────────────────────▼───────────────────────────────────────┐
│                      API Gateway (Nginx)                      │
└─────────────────────┬───────────────────────────────────────┘
                      │
┌─────────────────────▼───────────────────────────────────────┐
│                    Eodini Backend (Go)                        │
│  ┌──────────────────────────────────────────────────────┐   │
│  │  Handler (HTTP Layer)                                 │   │
│  │  - Vehicle, Driver, Route, Trip API                   │   │
│  │  - Request Validation                                 │   │
│  └────────────────┬─────────────────────────────────────┘   │
│                   │                                           │
│  ┌────────────────▼─────────────────────────────────────┐   │
│  │  Service (Business Logic)                             │   │
│  │  - 대체 기사 배정 로직                                │   │
│  │  - 운행 시작/완료 로직                                │   │
│  │  - 권한 검증                                          │   │
│  └────────────────┬─────────────────────────────────────┘   │
│                   │                                           │
│  ┌────────────────▼─────────────────────────────────────┐   │
│  │  Repository (Data Access)                             │   │
│  │  - GORM ORM                                           │   │
│  │  - Query Builder                                      │   │
│  └────────────────┬─────────────────────────────────────┘   │
│                   │                                           │
└───────────────────┼───────────────────────────────────────────┘
                    │
        ┌───────────┴──────────┐
        │                      │
┌───────▼────────┐  ┌─────────▼──────┐
│   PostgreSQL   │  │     Redis      │
│   (RDS)        │  │  (ElastiCache) │
└────────────────┘  └────────────────┘
```

## 📂 레이어 구조 (Clean Architecture)

### Layer 1: Domain (핵심 비즈니스 로직)
**위치**: `internal/domain/`

**역할**:
- 비즈니스 엔티티 정의
- 비즈니스 규칙 캡슐화
- 외부 의존성 없음 (순수 Go 코드)

**엔티티**:
```
Vehicle          - 차량
Driver           - 기사
Attendant        - 동승자/선생님
Route & Stop     - 경로 및 정류장
Passenger        - 탑승자
Schedule         - 운행 일정 템플릿
Trip             - 실제 운행 기록
DriverAssignment - 대체 배정
```

**특징**:
- 상태 전환 메소드 (SetActive, SetInactive 등)
- 유효성 검증 로직 (IsAvailableForTrip 등)
- 비즈니스 규칙 (면허 만료 체크, 운행 가능 여부 등)

### Layer 2: Repository (데이터 접근)
**위치**: `internal/repository/` (추후 구현)

**역할**:
- 데이터베이스 CRUD 작업
- GORM을 사용한 쿼리 실행
- 도메인 객체와 DB 모델 변환

**패턴**:
```go
type VehicleRepository interface {
    Create(vehicle *domain.Vehicle) error
    FindByID(id string) (*domain.Vehicle, error)
    Update(vehicle *domain.Vehicle) error
    Delete(id string) error
    List(filter VehicleFilter) ([]*domain.Vehicle, error)
}
```

### Layer 3: Service (비즈니스 로직)
**위치**: `internal/service/` (추후 구현)

**역할**:
- 복잡한 비즈니스 로직 조율
- 여러 Repository 조합
- 트랜잭션 관리
- 에러 처리 및 변환

**예시**:
```go
type TripService struct {
    tripRepo       TripRepository
    scheduleRepo   ScheduleRepository
    assignmentRepo DriverAssignmentRepository
}

// 대체 기사 반영한 Trip 생성
func (s *TripService) CreateTripForDate(date time.Time) {
    schedules := s.scheduleRepo.GetActiveSchedules(date)

    for _, schedule := range schedules {
        // 대체 배정 확인
        assignment := s.assignmentRepo.GetForDate(schedule.ID, date)

        driverID := schedule.DefaultDriverID
        if assignment != nil {
            driverID = assignment.DriverID // 대체!
        }

        trip := domain.NewTrip(schedule.ID, date, driverID, ...)
        s.tripRepo.Create(trip)
    }
}
```

### Layer 4: Handler (HTTP 인터페이스)
**위치**: `internal/handler/`

**역할**:
- HTTP 요청 수신
- 요청 검증 (Gin Validation)
- Service 호출
- HTTP 응답 반환

**구조**:
```go
type VehicleHandler struct {
    vehicleService *service.VehicleService
}

func (h *VehicleHandler) Create(c *gin.Context) {
    var req CreateVehicleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.Error(util.NewValidationError(...))
        return
    }

    vehicle, err := h.vehicleService.Create(&req)
    if err != nil {
        c.Error(err)
        return
    }

    util.SuccessResponse(c, http.StatusCreated, "생성 완료", vehicle)
}
```

### Layer 5: Middleware (횡단 관심사)
**위치**: `internal/middleware/`

**역할**:
- 요청/응답 처리
- 로깅, 에러 처리, CORS, 인증 등

**적용 순서** (중요!):
```
1. RecoveryHandler  - Panic 복구 (최우선)
2. RequestLogger    - 요청 로깅
3. CORS             - CORS 헤더
4. Auth             - 인증/인가 (추후)
5. ErrorHandler     - 에러 응답 (마지막)
```

### Supporting Layers

#### Util (공용 유틸리티)
**위치**: `internal/util/`
- 에러 정의 및 처리
- 메시지 관리
- 응답 포맷

#### Config (설정 관리)
**위치**: `config/`
- 환경변수 로드
- 설정 검증
- K8s ConfigMap/Secret 연동

#### Pkg (공용 패키지)
**위치**: `pkg/`
- `database/`: PostgreSQL 연결
- `cache/`: Redis 연결
- `logger/`: 구조화 로거

## 🔄 데이터 흐름

### 예시: 차량 등록 API

```
1. Client
   ↓ POST /api/v1/vehicles

2. Handler (VehicleHandler.Create)
   - 요청 검증
   - Service 호출
   ↓

3. Service (VehicleService.Create)
   - 중복 체크 (차량 번호)
   - 도메인 객체 생성
   - Repository 호출
   ↓

4. Repository (VehicleRepository.Create)
   - GORM Insert
   - DB 저장
   ↓

5. Database (PostgreSQL)
   - 데이터 저장
   ↑

6. Response
   - 생성된 Vehicle 반환
   - 표준 APIResponse 포맷
```

## 🎯 핵심 비즈니스 로직

### 1. 대체 기사 시나리오

```
시나리오: A 기사 휴가 → G 기사 대체

1. DriverAssignment 생성
   - ScheduleID: "오전 8시 A코스"
   - DriverID: G 기사
   - StartDate: 2025-01-20
   - EndDate: 2025-01-25

2. Trip 자동 생성 (크론잡)
   - 매일 자정 실행
   - 다음 날 운행 Trip 생성
   - DriverAssignment 확인
   - G 기사로 자동 배정

3. G 기사 앱 접속
   - GET /api/v1/drivers/{id}/trips/today
   - 자신에게 배정된 Trip 조회
   - Route/Stop 자동 표시

4. G 기사 운행 시작
   - POST /api/v1/trips/{id}/start
   - Trip.Status: pending → in_progress
   - Trip.StartedBy: "driver:G"
```

### 2. 운행 시작 권한

```
Driver 또는 Attendant 모두 시작 가능

1. 권한 확인
   - Trip.AssignedDriverID == userID (기사)
   - Trip.AssignedAttendantID == userID (선생님)

2. 운행 시작
   - Trip.Start(startedBy, location)
   - StartedBy: "driver:{id}" or "attendant:{id}"

3. 추적
   - 누가 시작했는지 기록
   - 위치 정보 저장
```

### 3. 탑승자 관리

```
1. Passenger → Route/Stop 배정
   - Passenger.AssignToStop(routeID, stopID, order)

2. Trip 생성 시 TripPassenger 자동 생성
   - Trip.TripPassengers 배열

3. 탑승/하차 처리
   - TripPassenger.BoardPassenger()
   - TripPassenger.AlightPassenger()

4. 불참 처리
   - TripPassenger.MarkNoShow(reason)
```

## 📊 도메인 모델 관계도

```
Schedule (운행 일정 템플릿)
  ├─ Route (1:1)
  ├─ Vehicle (1:1)
  ├─ DefaultDriver (1:1)
  └─ DefaultAttendant (0:1)
     ↓ (자동 생성)
Trip (실제 운행)
  ├─ Schedule (1:1) - 참조
  ├─ Vehicle (1:1)
  ├─ AssignedDriver (1:1) - 대체 가능!
  ├─ AssignedAttendant (0:1)
  └─ TripPassengers (1:N)
     └─ Passenger (1:1)
     └─ Stop (1:1)

DriverAssignment (대체 배정)
  ├─ Schedule (1:1)
  └─ Driver (1:1) - 대체 기사

Route (경로)
  └─ Stops (1:N)
     └─ Passengers (1:N)
```

## 🔐 보안 고려사항 (추후 구현)

### 인증 (Authentication)
- JWT 토큰 기반
- Refresh Token
- 토큰 만료 시간

### 인가 (Authorization)
- Role 기반 (Driver, Attendant, Admin)
- 리소스별 권한 체크
- 본인 데이터만 조회 가능

### 민감 정보 보호
- 비밀번호 암호화 (bcrypt)
- HTTPS 필수
- 의료 정보 암호화

## 🚀 성능 최적화

### 캐싱 전략 (Redis)
```
- Schedule 캐싱 (자주 조회)
- Route 캐싱 (정적 데이터)
- Driver 목록 캐싱
- TTL: 1시간
```

### 데이터베이스 최적화
```
- 인덱스: plate_number, driver_id, trip_date
- 커넥션 풀: MaxOpenConns=25, MaxIdleConns=5
- Prepared Statement
- N+1 쿼리 방지 (Preload)
```

### 비동기 처리
```
- Trip 자동 생성 (크론잡)
- 알림 발송 (고루틴)
- 위치 추적 (WebSocket - 추후)
```

## 📦 배포 아키텍처

### 로컬 개발
```
- Docker Compose
- PostgreSQL 컨테이너
- Redis 컨테이너
```

### 로컬 K8s (k3s)
```
- StatefulSet: PostgreSQL, Redis
- Deployment: Eodini API (2 replicas)
- Service: ClusterIP
- ConfigMap: 설정
- Secret: DB 비밀번호
```

### AWS (프로덕션)
```
- EKS: Kubernetes 클러스터
- RDS: PostgreSQL (Multi-AZ)
- ElastiCache: Redis
- ALB: Ingress
- ECR: 컨테이너 이미지
```

## 🔍 모니터링 & 관찰성

### Metrics (Prometheus)
- API 요청 수
- 응답 시간
- 에러율
- 데이터베이스 쿼리 시간

### Logging (Loki)
- 구조화 로깅 (JSON)
- 로그 레벨별 필터링
- RequestID로 추적

### Tracing (Tempo)
- 분산 추적
- API → Service → Repository

## 📚 참고 자료

- Clean Architecture: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
- Gin Framework: https://gin-gonic.com/
- GORM: https://gorm.io/
- Kubernetes: https://kubernetes.io/docs/
