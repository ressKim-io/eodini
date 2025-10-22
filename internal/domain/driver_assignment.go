package domain

import (
	"fmt"
	"time"
)

// 📝 설명: 기사 대체 배정 (A 기사 → G 기사)
// 🎯 실무 포인트: 휴가/퇴사 시 다른 기사가 특정 일정을 대체
// ⚠️ 주의사항: Trip 생성 시 이 정보를 확인하여 대체 기사 배정

// DriverAssignment - 기사 대체 배정 엔티티
type DriverAssignment struct {
	ID         string `json:"id"`
	ScheduleID string `json:"schedule_id"` // 대체할 일정
	DriverID   string `json:"driver_id"`   // 대체 기사

	// 대체 기간
	StartDate time.Time `json:"start_date"` // 시작일
	EndDate   time.Time `json:"end_date"`   // 종료일

	// 대체 사유
	Reason string `json:"reason"` // "원 담당자 휴가", "임시 배정" 등

	// 승인 정보 (선택적)
	ApprovedBy string     `json:"approved_by,omitempty"` // 승인자 ID
	ApprovedAt *time.Time `json:"approved_at,omitempty"` // 승인 시각

	// 메타데이터
	CreatedBy string     `json:"created_by"` // 생성자 (관리자)
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete
}

// NewDriverAssignment - 대체 배정 생성 팩토리 함수
func NewDriverAssignment(scheduleID, driverID string, startDate, endDate time.Time, reason, createdBy string) *DriverAssignment {
	now := time.Now()
	return &DriverAssignment{
		ScheduleID: scheduleID,
		DriverID:   driverID,
		StartDate:  startDate,
		EndDate:    endDate,
		Reason:     reason,
		CreatedBy:  createdBy,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
}

// IsActiveOnDate - 특정 날짜에 유효한 배정인지 확인
func (da *DriverAssignment) IsActiveOnDate(date time.Time) bool {
	if da.DeletedAt != nil {
		return false
	}

	// 날짜가 시작일과 종료일 사이인지 확인
	dateOnly := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	startDateOnly := time.Date(da.StartDate.Year(), da.StartDate.Month(), da.StartDate.Day(), 0, 0, 0, 0, da.StartDate.Location())
	endDateOnly := time.Date(da.EndDate.Year(), da.EndDate.Month(), da.EndDate.Day(), 0, 0, 0, 0, da.EndDate.Location())

	return (dateOnly.Equal(startDateOnly) || dateOnly.After(startDateOnly)) &&
		(dateOnly.Equal(endDateOnly) || dateOnly.Before(endDateOnly))
}

// Approve - 배정 승인
func (da *DriverAssignment) Approve(approverID string) {
	now := time.Now()
	da.ApprovedBy = approverID
	da.ApprovedAt = &now
	da.UpdatedAt = now
}

// IsApproved - 승인 여부
func (da *DriverAssignment) IsApproved() bool {
	return da.ApprovedAt != nil
}

// ExtendPeriod - 기간 연장
func (da *DriverAssignment) ExtendPeriod(newEndDate time.Time) error {
	if newEndDate.Before(da.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	if newEndDate.Before(da.EndDate) {
		return fmt.Errorf("new end date must be after current end date")
	}

	da.EndDate = newEndDate
	da.UpdatedAt = time.Now()
	return nil
}

// ShortenPeriod - 기간 단축
func (da *DriverAssignment) ShortenPeriod(newEndDate time.Time) error {
	if newEndDate.Before(da.StartDate) {
		return fmt.Errorf("end date cannot be before start date")
	}

	if newEndDate.After(da.EndDate) {
		return fmt.Errorf("new end date must be before current end date")
	}

	da.EndDate = newEndDate
	da.UpdatedAt = time.Now()
	return nil
}

// GetDuration - 대체 기간 (일수)
func (da *DriverAssignment) GetDuration() int {
	duration := da.EndDate.Sub(da.StartDate)
	return int(duration.Hours() / 24) + 1 // +1 to include both start and end dates
}

// IsExpired - 만료 여부 (종료일이 지났는지)
func (da *DriverAssignment) IsExpired() bool {
	now := time.Now()
	return da.EndDate.Before(now)
}

// IsUpcoming - 시작 전 여부
func (da *DriverAssignment) IsUpcoming() bool {
	now := time.Now()
	return da.StartDate.After(now)
}

// IsCurrent - 현재 진행 중 여부
func (da *DriverAssignment) IsCurrent() bool {
	return da.IsActiveOnDate(time.Now())
}

// UpdateReason - 사유 업데이트
func (da *DriverAssignment) UpdateReason(reason string) {
	da.Reason = reason
	da.UpdatedAt = time.Now()
}
