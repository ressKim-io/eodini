package util

import "net/http"

// 📝 설명: Spring의 @ControllerAdvice처럼 중앙 집중식 에러 관리
// 🎯 실무 포인트: 에러 코드를 상수로 관리하여 일관성 유지
// ⚠️ 주의사항: StatusCode는 JSON 응답에 포함하지 않음 (HTTP 헤더로만 사용)

// 에러 코드 상수
const (
	ErrCodeValidation   = "VALIDATION_ERROR"
	ErrCodeNotFound     = "NOT_FOUND"
	ErrCodeUnauthorized = "UNAUTHORIZED"
	ErrCodeForbidden    = "FORBIDDEN"
	ErrCodeInternal     = "INTERNAL_ERROR"
	ErrCodeDuplicate    = "DUPLICATE_ERROR"
	ErrCodeBadRequest   = "BAD_REQUEST"
	ErrCodeConflict     = "CONFLICT"
)

// AppError - 애플리케이션 에러 구조체
// Spring의 커스텀 Exception과 유사한 역할
type AppError struct {
	Code       string                 `json:"code"`              // 에러 코드 (예: "NOT_FOUND")
	Message    string                 `json:"message"`           // 사용자에게 보여줄 메시지
	StatusCode int                    `json:"-"`                 // HTTP 상태 코드 (JSON 응답에 미포함)
	Details    map[string]interface{} `json:"details,omitempty"` // 추가 상세 정보 (선택적)
}

// Error - error 인터페이스 구현
func (e *AppError) Error() string {
	return e.Message
}

// NewValidationError - 입력값 검증 실패 에러
// 사용 예: NewValidationError("입력값이 올바르지 않습니다", map[string]interface{}{"field": "email"})
func NewValidationError(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    message,
		StatusCode: http.StatusBadRequest,
		Details:    details,
	}
}

// NewNotFoundError - 리소스를 찾을 수 없는 경우
// 사용 예: NewNotFoundError("차량")
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    GetMessage(MsgResourceNotFound, resource),
		StatusCode: http.StatusNotFound,
	}
}

// NewUnauthorizedError - 인증 실패
func NewUnauthorizedError() *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    GetMessage(MsgUnauthorized),
		StatusCode: http.StatusUnauthorized,
	}
}

// NewForbiddenError - 권한 없음
func NewForbiddenError() *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    GetMessage(MsgForbidden),
		StatusCode: http.StatusForbidden,
	}
}

// NewInternalError - 서버 내부 오류
// 사용 예: NewInternalError(err)
func NewInternalError(err error) *AppError {
	details := map[string]interface{}{}
	if err != nil {
		details["error"] = err.Error()
	}

	return &AppError{
		Code:       ErrCodeInternal,
		Message:    GetMessage(MsgInternalError),
		StatusCode: http.StatusInternalServerError,
		Details:    details,
	}
}

// NewDuplicateError - 중복된 리소스
// 사용 예: NewDuplicateError("차량 번호")
func NewDuplicateError(resource string) *AppError {
	return &AppError{
		Code:       ErrCodeDuplicate,
		Message:    GetMessage(MsgDuplicate, resource),
		StatusCode: http.StatusConflict,
	}
}

// NewBadRequestError - 잘못된 요청
func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		StatusCode: http.StatusBadRequest,
	}
}

// NewConflictError - 비즈니스 로직 충돌
func NewConflictError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeConflict,
		Message:    message,
		StatusCode: http.StatusConflict,
	}
}
