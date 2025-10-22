package domain

import (
	"time"
)

// 📝 설명: 동승자(선생님/간호사) 도메인 모델
// 🎯 실무 포인트: 유치원 선생님 등 차량 동승 담당자, 운행 시작 권한 가능
// ⚠️ 주의사항: 동승자는 선택적 (없는 운행도 가능)

// AttendantStatus - 동승자 상태
type AttendantStatus string

const (
	AttendantStatusActive   AttendantStatus = "active"    // 활동 중
	AttendantStatusOnLeave  AttendantStatus = "on_leave"  // 휴가 중
	AttendantStatusInactive AttendantStatus = "inactive"  // 비활성 (퇴사 등)
)

// AttendantRole - 동승자 역할
type AttendantRole string

const (
	AttendantRoleTeacher AttendantRole = "teacher" // 선생님
	AttendantRoleNurse   AttendantRole = "nurse"   // 간호사
	AttendantRoleCarer   AttendantRole = "carer"   // 보호사
	AttendantRoleAssistant AttendantRole = "assistant" // 보조원
)

// Attendant - 동승자 엔티티
type Attendant struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Phone  string          `json:"phone"`
	Email  string          `json:"email,omitempty"`
	Role   AttendantRole   `json:"role"`   // 역할 (선생님, 간호사 등)
	Status AttendantStatus `json:"status"` // 상태

	// 권한
	CanStartTrip bool `json:"can_start_trip"` // 운행 시작 권한

	// 근무 정보
	HireDate         time.Time  `json:"hire_date"`                   // 입사일
	TerminationDate  *time.Time `json:"termination_date,omitempty"`  // 퇴사일

	// 추가 정보
	Organization     string `json:"organization,omitempty"` // 소속 기관 (유치원명 등)
	Address          string `json:"address,omitempty"`
	EmergencyContact string `json:"emergency_contact,omitempty"` // 비상 연락처
	Notes            string `json:"notes,omitempty"`

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete
}

// NewAttendant - 동승자 생성 팩토리 함수
func NewAttendant(name, phone string, role AttendantRole) *Attendant {
	now := time.Now()
	return &Attendant{
		Name:         name,
		Phone:        phone,
		Role:         role,
		Status:       AttendantStatusActive, // 기본값: 활동 중
		CanStartTrip: false,                 // 기본값: 운행 시작 권한 없음
		HireDate:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// IsActive - 활동 중인 동승자인지 확인
func (a *Attendant) IsActive() bool {
	return a.Status == AttendantStatusActive && a.DeletedAt == nil
}

// IsAvailableForTrip - 운행에 투입 가능한지 확인
func (a *Attendant) IsAvailableForTrip() bool {
	return a.IsActive()
}

// SetOnLeave - 휴가 상태로 변경
func (a *Attendant) SetOnLeave() {
	a.Status = AttendantStatusOnLeave
	a.UpdatedAt = time.Now()
}

// SetActive - 활동 중 상태로 변경
func (a *Attendant) SetActive() {
	a.Status = AttendantStatusActive
	a.UpdatedAt = time.Now()
}

// Terminate - 퇴사 처리
func (a *Attendant) Terminate(terminationDate time.Time) {
	a.Status = AttendantStatusInactive
	a.TerminationDate = &terminationDate
	a.UpdatedAt = time.Now()
}

// IsTerminated - 퇴사 여부
func (a *Attendant) IsTerminated() bool {
	return a.TerminationDate != nil
}

// GrantStartTripPermission - 운행 시작 권한 부여
func (a *Attendant) GrantStartTripPermission() {
	a.CanStartTrip = true
	a.UpdatedAt = time.Now()
}

// RevokeStartTripPermission - 운행 시작 권한 회수
func (a *Attendant) RevokeStartTripPermission() {
	a.CanStartTrip = false
	a.UpdatedAt = time.Now()
}

// UpdateContactInfo - 연락처 정보 업데이트
func (a *Attendant) UpdateContactInfo(phone, email string) {
	if phone != "" {
		a.Phone = phone
	}
	if email != "" {
		a.Email = email
	}
	a.UpdatedAt = time.Now()
}

// GetRoleDisplayName - 역할 표시명 (한글)
func (a *Attendant) GetRoleDisplayName() string {
	roleNames := map[AttendantRole]string{
		AttendantRoleTeacher:   "선생님",
		AttendantRoleNurse:     "간호사",
		AttendantRoleCarer:     "보호사",
		AttendantRoleAssistant: "보조원",
	}

	if name, exists := roleNames[a.Role]; exists {
		return name
	}
	return string(a.Role)
}
