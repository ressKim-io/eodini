package domain

import (
	"fmt"
	"time"
)

// 📝 설명: 실제 운행 기록 (날짜별 생성)
// 🎯 실무 포인트: Schedule에서 자동 생성, 대체 기사 배정 가능
// ⚠️ 주의사항: Trip은 날짜별로 1회성, Schedule은 템플릿

// TripStatus - 운행 상태
type TripStatus string

const (
	TripStatusPending    TripStatus = "pending"     // 대기 중
	TripStatusInProgress TripStatus = "in_progress" // 운행 중
	TripStatusCompleted  TripStatus = "completed"   // 완료
	TripStatusCancelled  TripStatus = "cancelled"   // 취소
)

// Trip - 실제 운행 엔티티
type Trip struct {
	ID         string     `json:"id"`
	ScheduleID string     `json:"schedule_id"` // 어떤 일정인지
	Date       time.Time  `json:"date"`        // 운행 날짜
	Status     TripStatus `json:"status"`

	// 배정 정보 (Schedule의 기본값에서 변경 가능)
	VehicleID           string  `json:"vehicle_id"`
	AssignedDriverID    string  `json:"assigned_driver_id"`
	AssignedAttendantID *string `json:"assigned_attendant_id,omitempty"`

	// 운행 기록
	StartedAt   *time.Time `json:"started_at,omitempty"`   // 실제 출발 시각
	CompletedAt *time.Time `json:"completed_at,omitempty"` // 실제 완료 시각
	StartedBy   string     `json:"started_by,omitempty"`   // 누가 시작했는지 (driver:{id} or attendant:{id})

	// 운행 정보
	ActualStartLocation *Location `json:"actual_start_location,omitempty"` // 실제 출발 위치
	ActualEndLocation   *Location `json:"actual_end_location,omitempty"`   // 실제 도착 위치
	TotalDistance       int       `json:"total_distance,omitempty"`        // 총 주행 거리 (미터)

	// 탑승 기록
	TripPassengers []TripPassenger `json:"trip_passengers,omitempty"` // 탑승자별 기록

	// 취소 정보
	CancelledAt     *time.Time `json:"cancelled_at,omitempty"`
	CancellationReason string  `json:"cancellation_reason,omitempty"`

	// 메모
	Notes string `json:"notes,omitempty"`

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete
}

// Location - 위치 정보
type Location struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	Timestamp time.Time `json:"timestamp"`
}

// TripPassenger - 탑승자별 운행 기록
type TripPassenger struct {
	ID          string     `json:"id"`
	TripID      string     `json:"trip_id"`
	PassengerID string     `json:"passenger_id"`
	StopID      string     `json:"stop_id"` // 어느 정류장

	// 탑승 기록
	BoardedAt   *time.Time `json:"boarded_at,omitempty"`   // 탑승 시각
	AlightedAt  *time.Time `json:"alighted_at,omitempty"`  // 하차 시각
	IsBoarded   bool       `json:"is_boarded"`             // 탑승 여부
	IsAlighted  bool       `json:"is_alighted"`            // 하차 여부

	// 추가 정보
	NoShowReason string     `json:"no_show_reason,omitempty"` // 불참 사유
	Notes        string     `json:"notes,omitempty"`

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// NewTrip - 운행 생성 팩토리 함수
func NewTrip(scheduleID string, date time.Time, vehicleID, driverID string, attendantID *string) *Trip {
	now := time.Now()
	return &Trip{
		ScheduleID:          scheduleID,
		Date:                date,
		VehicleID:           vehicleID,
		AssignedDriverID:    driverID,
		AssignedAttendantID: attendantID,
		Status:              TripStatusPending, // 기본값: 대기 중
		TripPassengers:      []TripPassenger{},
		CreatedAt:           now,
		UpdatedAt:           now,
	}
}

// NewTripPassenger - 탑승자 기록 생성
func NewTripPassenger(tripID, passengerID, stopID string) *TripPassenger {
	now := time.Now()
	return &TripPassenger{
		TripID:      tripID,
		PassengerID: passengerID,
		StopID:      stopID,
		IsBoarded:   false,
		IsAlighted:  false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// IsPending - 대기 중인 운행인지
func (t *Trip) IsPending() bool {
	return t.Status == TripStatusPending
}

// IsInProgress - 운행 중인지
func (t *Trip) IsInProgress() bool {
	return t.Status == TripStatusInProgress
}

// IsCompleted - 완료된 운행인지
func (t *Trip) IsCompleted() bool {
	return t.Status == TripStatusCompleted
}

// IsCancelled - 취소된 운행인지
func (t *Trip) IsCancelled() bool {
	return t.Status == TripStatusCancelled
}

// CanStart - 운행 시작 가능한지
func (t *Trip) CanStart() bool {
	return t.Status == TripStatusPending
}

// CanComplete - 운행 완료 가능한지
func (t *Trip) CanComplete() bool {
	return t.Status == TripStatusInProgress
}

// Start - 운행 시작
func (t *Trip) Start(startedBy string, location *Location) error {
	if !t.CanStart() {
		return fmt.Errorf("cannot start trip: current status is %s", t.Status)
	}

	now := time.Now()
	t.Status = TripStatusInProgress
	t.StartedAt = &now
	t.StartedBy = startedBy
	t.ActualStartLocation = location
	t.UpdatedAt = now

	return nil
}

// Complete - 운행 완료
func (t *Trip) Complete(location *Location) error {
	if !t.CanComplete() {
		return fmt.Errorf("cannot complete trip: current status is %s", t.Status)
	}

	now := time.Now()
	t.Status = TripStatusCompleted
	t.CompletedAt = &now
	t.ActualEndLocation = location
	t.UpdatedAt = now

	return nil
}

// Cancel - 운행 취소
func (t *Trip) Cancel(reason string) error {
	if t.IsCompleted() {
		return fmt.Errorf("cannot cancel completed trip")
	}

	now := time.Now()
	t.Status = TripStatusCancelled
	t.CancelledAt = &now
	t.CancellationReason = reason
	t.UpdatedAt = now

	return nil
}

// GetDuration - 운행 소요 시간 (분)
func (t *Trip) GetDuration() int {
	if t.StartedAt == nil || t.CompletedAt == nil {
		return 0
	}
	return int(t.CompletedAt.Sub(*t.StartedAt).Minutes())
}

// BoardPassenger - 탑승자 탑승 처리
func (tp *TripPassenger) BoardPassenger() {
	now := time.Now()
	tp.IsBoarded = true
	tp.BoardedAt = &now
	tp.UpdatedAt = now
}

// AlightPassenger - 탑승자 하차 처리
func (tp *TripPassenger) AlightPassenger() {
	now := time.Now()
	tp.IsAlighted = true
	tp.AlightedAt = &now
	tp.UpdatedAt = now
}

// MarkNoShow - 불참 처리
func (tp *TripPassenger) MarkNoShow(reason string) {
	tp.IsBoarded = false
	tp.NoShowReason = reason
	tp.UpdatedAt = time.Now()
}

// GetBoardingDuration - 탑승 시간 (분)
func (tp *TripPassenger) GetBoardingDuration() int {
	if tp.BoardedAt == nil || tp.AlightedAt == nil {
		return 0
	}
	return int(tp.AlightedAt.Sub(*tp.BoardedAt).Minutes())
}
