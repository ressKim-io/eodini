package domain

import (
	"time"
)

// 📝 설명: 기사 도메인 모델
// 🎯 실무 포인트: 면허 만료일 체크, 휴가/퇴사 상태 관리
// ⚠️ 주의사항: 면허 만료된 기사는 운행 불가

// DriverStatus - 기사 상태
type DriverStatus string

const (
	DriverStatusActive   DriverStatus = "active"    // 활동 중
	DriverStatusOnLeave  DriverStatus = "on_leave"  // 휴가 중
	DriverStatusInactive DriverStatus = "inactive"  // 비활성 (퇴사 등)
)

// LicenseType - 운전면허 종류
type LicenseType string

const (
	LicenseType1Regular LicenseType = "type_1_regular" // 1종 보통
	LicenseType1Large   LicenseType = "type_1_large"   // 1종 대형
	LicenseType2Regular LicenseType = "type_2_regular" // 2종 보통
)

// Driver - 기사 엔티티
type Driver struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Phone  string `json:"phone"`
	Email  string `json:"email,omitempty"`
	Status DriverStatus `json:"status"`

	// 운전면허 정보
	LicenseNumber string      `json:"license_number"` // 면허 번호
	LicenseType   LicenseType `json:"license_type"`   // 면허 종류
	LicenseExpiry time.Time   `json:"license_expiry"` // 면허 만료일

	// 근무 정보
	HireDate       time.Time  `json:"hire_date"`                  // 입사일
	TerminationDate *time.Time `json:"termination_date,omitempty"` // 퇴사일

	// 추가 정보
	Address      string `json:"address,omitempty"`
	EmergencyContact string `json:"emergency_contact,omitempty"` // 비상 연락처
	Notes        string `json:"notes,omitempty"` // 메모

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"` // Soft delete
}

// NewDriver - 기사 생성 팩토리 함수
func NewDriver(name, phone, licenseNumber string, licenseType LicenseType, licenseExpiry time.Time) *Driver {
	now := time.Now()
	return &Driver{
		Name:          name,
		Phone:         phone,
		LicenseNumber: licenseNumber,
		LicenseType:   licenseType,
		LicenseExpiry: licenseExpiry,
		Status:        DriverStatusActive, // 기본값: 활동 중
		HireDate:      now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
}

// IsActive - 활동 중인 기사인지 확인
func (d *Driver) IsActive() bool {
	return d.Status == DriverStatusActive && d.DeletedAt == nil
}

// IsAvailableForTrip - 운행에 투입 가능한지 확인
// 활동 중이고, 면허가 유효해야 함
func (d *Driver) IsAvailableForTrip() bool {
	if !d.IsActive() {
		return false
	}

	// 면허 만료 확인
	if d.IsLicenseExpired() {
		return false
	}

	return true
}

// IsLicenseExpired - 면허 만료 여부
func (d *Driver) IsLicenseExpired() bool {
	return d.LicenseExpiry.Before(time.Now())
}

// NeedsLicenseRenewal - 면허 갱신 필요 여부 (30일 이내 만료)
func (d *Driver) NeedsLicenseRenewal() bool {
	thirtyDaysLater := time.Now().AddDate(0, 0, 30)
	return d.LicenseExpiry.Before(thirtyDaysLater)
}

// SetOnLeave - 휴가 상태로 변경
func (d *Driver) SetOnLeave() {
	d.Status = DriverStatusOnLeave
	d.UpdatedAt = time.Now()
}

// SetActive - 활동 중 상태로 변경
func (d *Driver) SetActive() {
	d.Status = DriverStatusActive
	d.UpdatedAt = time.Now()
}

// Terminate - 퇴사 처리
func (d *Driver) Terminate(terminationDate time.Time) {
	d.Status = DriverStatusInactive
	d.TerminationDate = &terminationDate
	d.UpdatedAt = time.Now()
}

// IsTerminated - 퇴사 여부
func (d *Driver) IsTerminated() bool {
	return d.TerminationDate != nil
}

// UpdateLicenseExpiry - 면허 만료일 업데이트
func (d *Driver) UpdateLicenseExpiry(expiry time.Time) {
	d.LicenseExpiry = expiry
	d.UpdatedAt = time.Now()
}

// UpdateContactInfo - 연락처 정보 업데이트
func (d *Driver) UpdateContactInfo(phone, email string) {
	if phone != "" {
		d.Phone = phone
	}
	if email != "" {
		d.Email = email
	}
	d.UpdatedAt = time.Now()
}
