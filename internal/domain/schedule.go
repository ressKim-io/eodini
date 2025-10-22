package domain

import (
	"time"
)

// 📝 설명: 운행 일정 템플릿 (매일 반복되는 운행 계획)
// 🎯 실무 포인트: "매일 오전 8시 A코스" 같은 고정 일정, Trip 생성의 기준
// ⚠️ 주의사항: Schedule은 템플릿, 실제 운행은 Trip

// ScheduleStatus - 일정 상태
type ScheduleStatus string

const (
	ScheduleStatusActive   ScheduleStatus = "active"   // 사용 중
	ScheduleStatusInactive ScheduleStatus = "inactive" // 미사용
)

// TimeSlot - 시간대
type TimeSlot string

const (
	TimeSlotMorning   TimeSlot = "morning"   // 오전
	TimeSlotAfternoon TimeSlot = "afternoon" // 오후
	TimeSlotEvening   TimeSlot = "evening"   // 저녁
)

// DayOfWeek - 요일 (1=월요일, 7=일요일)
type DayOfWeek int

const (
	Monday DayOfWeek = 1
	Tuesday DayOfWeek = 2
	Wednesday DayOfWeek = 3
	Thursday DayOfWeek = 4
	Friday DayOfWeek = 5
	Saturday DayOfWeek = 6
	Sunday DayOfWeek = 7
)

// Schedule - 운행 일정 템플릿
type Schedule struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`        // 일정명 (예: "오전 8시 A코스")
	Description string         `json:"description"` // 설명
	Status      ScheduleStatus `json:"status"`

	// 시간 설정
	StartTime string   `json:"start_time"` // 출발 시각 (HH:MM 형식, 예: "08:00")
	TimeSlot  TimeSlot `json:"time_slot"`  // 시간대 (오전/오후/저녁)

	// 운행 요일 (1=월, 2=화, ..., 7=일)
	DaysOfWeek []int `json:"days_of_week"` // 예: [1,2,3,4,5] (월~금)

	// 배정 정보
	RouteID  string  `json:"route_id"`   // 경로
	VehicleID string `json:"vehicle_id"` // 차량

	// 기본 담당자 (대체 가능)
	DefaultDriverID    string  `json:"default_driver_id"`               // 기본 기사
	DefaultAttendantID *string `json:"default_attendant_id,omitempty"`  // 기본 동승자 (선택적)

	// 유효 기간 (선택적)
	ValidFrom *time.Time `json:"valid_from,omitempty"` // 유효 시작일
	ValidTo   *time.Time `json:"valid_to,omitempty"`   // 유효 종료일

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete
}

// NewSchedule - 일정 생성 팩토리 함수
func NewSchedule(name, startTime string, timeSlot TimeSlot, daysOfWeek []int, routeID, vehicleID, driverID string) *Schedule {
	now := time.Now()
	return &Schedule{
		Name:            name,
		StartTime:       startTime,
		TimeSlot:        timeSlot,
		DaysOfWeek:      daysOfWeek,
		RouteID:         routeID,
		VehicleID:       vehicleID,
		DefaultDriverID: driverID,
		Status:          ScheduleStatusActive, // 기본값: 사용 중
		CreatedAt:       now,
		UpdatedAt:       now,
	}
}

// IsActive - 활성 일정인지 확인
func (s *Schedule) IsActive() bool {
	return s.Status == ScheduleStatusActive && s.DeletedAt == nil
}

// IsActiveOnDate - 특정 날짜에 운행하는지 확인
func (s *Schedule) IsActiveOnDate(date time.Time) bool {
	if !s.IsActive() {
		return false
	}

	// 유효 기간 확인
	if s.ValidFrom != nil && date.Before(*s.ValidFrom) {
		return false
	}
	if s.ValidTo != nil && date.After(*s.ValidTo) {
		return false
	}

	// 요일 확인
	dayOfWeek := int(date.Weekday())
	if dayOfWeek == 0 { // Sunday
		dayOfWeek = 7
	}

	for _, day := range s.DaysOfWeek {
		if day == dayOfWeek {
			return true
		}
	}

	return false
}

// SetActive - 활성 상태로 변경
func (s *Schedule) SetActive() {
	s.Status = ScheduleStatusActive
	s.UpdatedAt = time.Now()
}

// SetInactive - 비활성 상태로 변경
func (s *Schedule) SetInactive() {
	s.Status = ScheduleStatusInactive
	s.UpdatedAt = time.Now()
}

// UpdateStartTime - 출발 시각 변경
func (s *Schedule) UpdateStartTime(startTime string) {
	s.StartTime = startTime
	s.UpdatedAt = time.Now()
}

// UpdateDaysOfWeek - 운행 요일 변경
func (s *Schedule) UpdateDaysOfWeek(days []int) {
	s.DaysOfWeek = days
	s.UpdatedAt = time.Now()
}

// AssignVehicle - 차량 배정
func (s *Schedule) AssignVehicle(vehicleID string) {
	s.VehicleID = vehicleID
	s.UpdatedAt = time.Now()
}

// AssignDriver - 기본 기사 배정
func (s *Schedule) AssignDriver(driverID string) {
	s.DefaultDriverID = driverID
	s.UpdatedAt = time.Now()
}

// AssignAttendant - 기본 동승자 배정
func (s *Schedule) AssignAttendant(attendantID string) {
	s.DefaultAttendantID = &attendantID
	s.UpdatedAt = time.Now()
}

// UnassignAttendant - 동승자 배정 해제
func (s *Schedule) UnassignAttendant() {
	s.DefaultAttendantID = nil
	s.UpdatedAt = time.Now()
}

// HasAttendant - 동승자가 배정되어 있는지 확인
func (s *Schedule) HasAttendant() bool {
	return s.DefaultAttendantID != nil && *s.DefaultAttendantID != ""
}

// SetValidPeriod - 유효 기간 설정
func (s *Schedule) SetValidPeriod(from, to time.Time) {
	s.ValidFrom = &from
	s.ValidTo = &to
	s.UpdatedAt = time.Now()
}

// ClearValidPeriod - 유효 기간 해제
func (s *Schedule) ClearValidPeriod() {
	s.ValidFrom = nil
	s.ValidTo = nil
	s.UpdatedAt = time.Now()
}

// IsWeekday - 평일 운행 여부
func (s *Schedule) IsWeekday() bool {
	weekdays := map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true}
	for _, day := range s.DaysOfWeek {
		if !weekdays[day] {
			return false
		}
	}
	return len(s.DaysOfWeek) > 0
}

// IsWeekend - 주말 운행 여부
func (s *Schedule) IsWeekend() bool {
	weekend := map[int]bool{6: true, 7: true}
	hasWeekendDay := false
	for _, day := range s.DaysOfWeek {
		if weekend[day] {
			hasWeekendDay = true
		}
	}
	return hasWeekendDay
}
