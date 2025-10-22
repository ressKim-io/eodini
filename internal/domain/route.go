package domain

import (
	"time"
)

// 📝 설명: 경로 및 정류장 도메인 모델
// 🎯 실무 포인트: Route는 A코스/B코스 같은 고정 동선, Stop은 정류장 순서 포함
// ⚠️ 주의사항: Stop의 Order는 정류장 순서를 나타냄 (1, 2, 3...)

// RouteStatus - 경로 상태
type RouteStatus string

const (
	RouteStatusActive   RouteStatus = "active"   // 사용 중
	RouteStatusInactive RouteStatus = "inactive" // 미사용
)

// Route - 경로 엔티티 (A코스, B코스 등)
type Route struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`        // 경로명 (예: "A코스", "오전 노선 1")
	Description string      `json:"description"` // 경로 설명
	Status      RouteStatus `json:"status"`

	// 경로 정보
	Stops            []Stop `json:"stops,omitempty"`           // 정류장 목록
	EstimatedTime    int    `json:"estimated_time"`            // 예상 소요 시간 (분)
	TotalDistance    int    `json:"total_distance,omitempty"`  // 총 거리 (미터)

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete
}

// Stop - 정류장 엔티티
type Stop struct {
	ID       string  `json:"id"`
	RouteID  string  `json:"route_id"`  // 소속 경로
	Name     string  `json:"name"`      // 정류장 이름 (예: "OO아파트 정문")
	Address  string  `json:"address"`   // 주소
	Order    int     `json:"order"`     // 순서 (1부터 시작)
	Latitude float64 `json:"latitude"`  // 위도
	Longitude float64 `json:"longitude"` // 경도

	// 정류장 정보
	EstimatedArrivalTime int    `json:"estimated_arrival_time"` // 예상 도착 시간 (출발 후 몇 분)
	Notes                string `json:"notes,omitempty"`        // 메모 (특이사항)

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete
}

// NewRoute - 경로 생성 팩토리 함수
func NewRoute(name, description string, estimatedTime int) *Route {
	now := time.Now()
	return &Route{
		Name:          name,
		Description:   description,
		EstimatedTime: estimatedTime,
		Status:        RouteStatusActive, // 기본값: 사용 중
		Stops:         []Stop{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// NewStop - 정류장 생성 팩토리 함수
func NewStop(routeID, name, address string, order int, latitude, longitude float64, estimatedArrivalTime int) *Stop {
	now := time.Now()
	return &Stop{
		RouteID:              routeID,
		Name:                 name,
		Address:              address,
		Order:                order,
		Latitude:             latitude,
		Longitude:            longitude,
		EstimatedArrivalTime: estimatedArrivalTime,
		CreatedAt:            now,
		UpdatedAt:            now,
	}
}

// IsActive - 활성 경로인지 확인
func (r *Route) IsActive() bool {
	return r.Status == RouteStatusActive && r.DeletedAt == nil
}

// SetActive - 활성 상태로 변경
func (r *Route) SetActive() {
	r.Status = RouteStatusActive
	r.UpdatedAt = time.Now()
}

// SetInactive - 비활성 상태로 변경
func (r *Route) SetInactive() {
	r.Status = RouteStatusInactive
	r.UpdatedAt = time.Now()
}

// AddStop - 정류장 추가
func (r *Route) AddStop(stop Stop) {
	r.Stops = append(r.Stops, stop)
	r.UpdatedAt = time.Now()
}

// GetStopByOrder - 순서로 정류장 조회
func (r *Route) GetStopByOrder(order int) *Stop {
	for i := range r.Stops {
		if r.Stops[i].Order == order {
			return &r.Stops[i]
		}
	}
	return nil
}

// GetStopCount - 정류장 개수
func (r *Route) GetStopCount() int {
	return len(r.Stops)
}

// GetFirstStop - 첫 번째 정류장
func (r *Route) GetFirstStop() *Stop {
	if len(r.Stops) == 0 {
		return nil
	}
	return r.GetStopByOrder(1)
}

// GetLastStop - 마지막 정류장
func (r *Route) GetLastStop() *Stop {
	if len(r.Stops) == 0 {
		return nil
	}
	return r.GetStopByOrder(r.GetStopCount())
}

// UpdateEstimatedTime - 예상 소요 시간 업데이트
func (r *Route) UpdateEstimatedTime(minutes int) {
	r.EstimatedTime = minutes
	r.UpdatedAt = time.Now()
}

// UpdateTotalDistance - 총 거리 업데이트
func (r *Route) UpdateTotalDistance(meters int) {
	r.TotalDistance = meters
	r.UpdatedAt = time.Now()
}

// IsValidStop - 유효한 정류장인지 확인
func (s *Stop) IsValidStop() bool {
	return s.Name != "" &&
		s.Address != "" &&
		s.Order > 0 &&
		s.Latitude != 0 &&
		s.Longitude != 0 &&
		s.DeletedAt == nil
}

// UpdateLocation - 위치 정보 업데이트
func (s *Stop) UpdateLocation(latitude, longitude float64) {
	s.Latitude = latitude
	s.Longitude = longitude
	s.UpdatedAt = time.Now()
}

// UpdateEstimatedArrivalTime - 예상 도착 시간 업데이트
func (s *Stop) UpdateEstimatedArrivalTime(minutes int) {
	s.EstimatedArrivalTime = minutes
	s.UpdatedAt = time.Now()
}
