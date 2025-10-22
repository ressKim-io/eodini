package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 📝 설명: 차량 도메인 모델
// 🎯 실무 포인트: 차량 번호를 유니크 키로 사용, 정원 관리
// ⚠️ 주의사항: 차량 번호는 중복 불가 (비즈니스 규칙)

// VehicleType - 차량 유형
type VehicleType string

const (
	VehicleTypeVan       VehicleType = "van"        // 승합차
	VehicleTypeBus       VehicleType = "bus"        // 버스
	VehicleTypeMiniBus   VehicleType = "mini_bus"   // 소형버스
	VehicleTypeSedan     VehicleType = "sedan"      // 승용차
)

// VehicleStatus - 차량 상태
type VehicleStatus string

const (
	VehicleStatusActive      VehicleStatus = "active"       // 운행 가능
	VehicleStatusMaintenance VehicleStatus = "maintenance"  // 정비 중
	VehicleStatusInactive    VehicleStatus = "inactive"     // 비활성 (폐차 등)
)

// Vehicle - 차량 엔티티
type Vehicle struct {
	ID           string        `json:"id" gorm:"type:uuid;primaryKey"`
	PlateNumber  string        `json:"plate_number" gorm:"uniqueIndex;not null"` // 차량 번호 (예: "12가3456")
	Model        string        `json:"model" gorm:"not null"`                    // 차량 모델 (예: "그랜드스타렉스")
	Manufacturer string        `json:"manufacturer"`                             // 제조사 (예: "현대")
	VehicleType  VehicleType   `json:"vehicle_type" gorm:"type:varchar(20);not null"` // 차량 유형
	Capacity     int           `json:"capacity" gorm:"not null"`                 // 정원 (운전자 포함)
	Year         int           `json:"year"`                                     // 연식
	Color        string        `json:"color"`                                    // 색상
	Status       VehicleStatus `json:"status" gorm:"type:varchar(20);not null;default:'active'"` // 차량 상태

	// 차량 관리 정보
	InsuranceExpiry   *time.Time `json:"insurance_expiry,omitempty"`    // 보험 만료일
	InspectionExpiry  *time.Time `json:"inspection_expiry,omitempty"`   // 정기검사 만료일
	LastMaintenanceAt *time.Time `json:"last_maintenance_at,omitempty"` // 마지막 정비 날짜

	// 메타데이터
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"` // Soft delete
}

// NewVehicle - 차량 생성 팩토리 함수
func NewVehicle(plateNumber, model, manufacturer string, vehicleType VehicleType, capacity, year int, color string) *Vehicle {
	now := time.Now()
	return &Vehicle{
		ID:           uuid.New().String(), // UUID 자동 생성
		PlateNumber:  plateNumber,
		Model:        model,
		Manufacturer: manufacturer,
		VehicleType:  vehicleType,
		Capacity:     capacity,
		Year:         year,
		Color:        color,
		Status:       VehicleStatusActive, // 기본값: 운행 가능
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

// BeforeCreate - GORM Hook: 생성 전 자동 처리
func (v *Vehicle) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	now := time.Now()
	v.CreatedAt = now
	v.UpdatedAt = now
	return nil
}

// BeforeUpdate - GORM Hook: 업데이트 전 자동 처리
func (v *Vehicle) BeforeUpdate(tx *gorm.DB) error {
	v.UpdatedAt = time.Now()
	return nil
}

// IsActive - 운행 가능한 차량인지 확인
func (v *Vehicle) IsActive() bool {
	return v.Status == VehicleStatusActive && v.DeletedAt == nil
}

// IsAvailableForTrip - 운행에 사용 가능한지 확인
// 정비 중이거나 보험/검사가 만료된 차량은 불가
func (v *Vehicle) IsAvailableForTrip() bool {
	if !v.IsActive() {
		return false
	}

	now := time.Now()

	// 보험 만료 확인
	if v.InsuranceExpiry != nil && v.InsuranceExpiry.Before(now) {
		return false
	}

	// 정기검사 만료 확인
	if v.InspectionExpiry != nil && v.InspectionExpiry.Before(now) {
		return false
	}

	return true
}

// SetMaintenance - 정비 중 상태로 변경
func (v *Vehicle) SetMaintenance() {
	v.Status = VehicleStatusMaintenance
	now := time.Now()
	v.LastMaintenanceAt = &now
	v.UpdatedAt = now
}

// SetActive - 운행 가능 상태로 변경
func (v *Vehicle) SetActive() {
	v.Status = VehicleStatusActive
	v.UpdatedAt = time.Now()
}

// SetInactive - 비활성 상태로 변경 (폐차 등)
func (v *Vehicle) SetInactive() {
	v.Status = VehicleStatusInactive
	v.UpdatedAt = time.Now()
}

// UpdateInsuranceExpiry - 보험 만료일 업데이트
func (v *Vehicle) UpdateInsuranceExpiry(expiry time.Time) {
	v.InsuranceExpiry = &expiry
	v.UpdatedAt = time.Now()
}

// UpdateInspectionExpiry - 정기검사 만료일 업데이트
func (v *Vehicle) UpdateInspectionExpiry(expiry time.Time) {
	v.InspectionExpiry = &expiry
	v.UpdatedAt = time.Now()
}

// NeedsInsuranceRenewal - 보험 갱신 필요 여부 (30일 이내 만료)
func (v *Vehicle) NeedsInsuranceRenewal() bool {
	if v.InsuranceExpiry == nil {
		return false
	}

	thirtyDaysLater := time.Now().AddDate(0, 0, 30)
	return v.InsuranceExpiry.Before(thirtyDaysLater)
}

// NeedsInspection - 정기검사 필요 여부 (30일 이내 만료)
func (v *Vehicle) NeedsInspection() bool {
	if v.InspectionExpiry == nil {
		return false
	}

	thirtyDaysLater := time.Now().AddDate(0, 0, 30)
	return v.InspectionExpiry.Before(thirtyDaysLater)
}

// GetPassengerCapacity - 승객 정원 (운전자 제외)
func (v *Vehicle) GetPassengerCapacity() int {
	if v.Capacity <= 1 {
		return 0
	}
	return v.Capacity - 1
}
