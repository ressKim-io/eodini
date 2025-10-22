package domain

import (
	"time"
)

// 📝 설명: 탑승자 도메인 모델 (유치원생, 통원 환자 등)
// 🎯 실무 포인트: 보호자 정보, 특이사항 관리
// ⚠️ 주의사항: 알레르기, 투약 정보 등 민감 정보 포함 가능

// PassengerStatus - 탑승자 상태
type PassengerStatus string

const (
	PassengerStatusActive   PassengerStatus = "active"   // 활동 중 (탑승 중)
	PassengerStatusInactive PassengerStatus = "inactive" // 비활성 (졸업, 전학 등)
)

// Passenger - 탑승자 엔티티
type Passenger struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Age    int             `json:"age,omitempty"`
	Gender string          `json:"gender,omitempty"` // "male", "female", "other"
	Status PassengerStatus `json:"status"`

	// 탑승 정보
	AssignedRouteID string  `json:"assigned_route_id"`          // 배정된 경로
	AssignedStopID  string  `json:"assigned_stop_id"`           // 배정된 정류장
	StopOrder       int     `json:"stop_order"`                 // 정류장 순서 (캐싱용)

	// 보호자 정보
	GuardianName    string `json:"guardian_name"`    // 보호자 이름
	GuardianPhone   string `json:"guardian_phone"`   // 보호자 연락처
	GuardianEmail   string `json:"guardian_email,omitempty"`
	GuardianRelation string `json:"guardian_relation,omitempty"` // 관계 (부, 모, 조부모 등)

	// 비상 연락처 (보호자와 다른 경우)
	EmergencyContact string `json:"emergency_contact,omitempty"`
	EmergencyRelation string `json:"emergency_relation,omitempty"`

	// 추가 정보
	Address      string `json:"address,omitempty"`
	MedicalNotes string `json:"medical_notes,omitempty"` // 의료 특이사항 (알레르기 등)
	Notes        string `json:"notes,omitempty"`         // 일반 메모

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete
}

// NewPassenger - 탑승자 생성 팩토리 함수
func NewPassenger(name, guardianName, guardianPhone string) *Passenger {
	now := time.Now()
	return &Passenger{
		Name:          name,
		GuardianName:  guardianName,
		GuardianPhone: guardianPhone,
		Status:        PassengerStatusActive, // 기본값: 활동 중
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// IsActive - 활동 중인 탑승자인지 확인
func (p *Passenger) IsActive() bool {
	return p.Status == PassengerStatusActive && p.DeletedAt == nil
}

// SetActive - 활동 중 상태로 변경
func (p *Passenger) SetActive() {
	p.Status = PassengerStatusActive
	p.UpdatedAt = time.Now()
}

// SetInactive - 비활성 상태로 변경 (졸업, 전학 등)
func (p *Passenger) SetInactive() {
	p.Status = PassengerStatusInactive
	p.UpdatedAt = time.Now()
}

// AssignToStop - 정류장 배정
func (p *Passenger) AssignToStop(routeID, stopID string, stopOrder int) {
	p.AssignedRouteID = routeID
	p.AssignedStopID = stopID
	p.StopOrder = stopOrder
	p.UpdatedAt = time.Now()
}

// UnassignFromStop - 정류장 배정 해제
func (p *Passenger) UnassignFromStop() {
	p.AssignedRouteID = ""
	p.AssignedStopID = ""
	p.StopOrder = 0
	p.UpdatedAt = time.Now()
}

// IsAssigned - 정류장에 배정되어 있는지 확인
func (p *Passenger) IsAssigned() bool {
	return p.AssignedRouteID != "" && p.AssignedStopID != ""
}

// UpdateGuardianInfo - 보호자 정보 업데이트
func (p *Passenger) UpdateGuardianInfo(name, phone, email, relation string) {
	if name != "" {
		p.GuardianName = name
	}
	if phone != "" {
		p.GuardianPhone = phone
	}
	if email != "" {
		p.GuardianEmail = email
	}
	if relation != "" {
		p.GuardianRelation = relation
	}
	p.UpdatedAt = time.Now()
}

// UpdateEmergencyContact - 비상 연락처 업데이트
func (p *Passenger) UpdateEmergencyContact(contact, relation string) {
	p.EmergencyContact = contact
	p.EmergencyRelation = relation
	p.UpdatedAt = time.Now()
}

// UpdateMedicalNotes - 의료 특이사항 업데이트
func (p *Passenger) UpdateMedicalNotes(notes string) {
	p.MedicalNotes = notes
	p.UpdatedAt = time.Now()
}

// HasMedicalNotes - 의료 특이사항이 있는지 확인
func (p *Passenger) HasMedicalNotes() bool {
	return p.MedicalNotes != ""
}

// GetContactPhone - 주 연락처 (보호자 또는 비상 연락처)
func (p *Passenger) GetContactPhone() string {
	if p.GuardianPhone != "" {
		return p.GuardianPhone
	}
	return p.EmergencyContact
}
