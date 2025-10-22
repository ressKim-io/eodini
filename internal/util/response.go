package util

import "github.com/gin-gonic/gin"

// 📝 설명: Spring의 ResponseEntity처럼 표준화된 API 응답 구조
// 🎯 실무 포인트: 모든 API 응답을 통일된 포맷으로 관리
// ⚠️ 주의사항: Success와 Error는 상호 배타적 (둘 중 하나만 존재)

// APIResponse - 표준 API 응답 구조
// 성공/실패 모두 이 구조를 사용하여 일관성 유지
type APIResponse struct {
	Success bool        `json:"success"`           // 성공 여부
	Message string      `json:"message"`           // 응답 메시지
	Data    interface{} `json:"data,omitempty"`    // 응답 데이터 (성공 시)
	Error   *AppError   `json:"error,omitempty"`   // 에러 정보 (실패 시)
}

// SuccessResponse - 성공 응답 헬퍼 함수
// 사용 예:
//   SuccessResponse(c, http.StatusOK, "조회 성공", vehicle)
//   SuccessResponse(c, http.StatusCreated, "생성 완료", newVehicle)
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse - 에러 응답 헬퍼 함수
// 사용 예:
//   ErrorResponse(c, util.NewNotFoundError("차량"))
//   ErrorResponse(c, util.NewValidationError("입력값 오류", details))
func ErrorResponse(c *gin.Context, err *AppError) {
	c.JSON(err.StatusCode, APIResponse{
		Success: false,
		Message: err.Message,
		Error:   err,
	})
}

// SuccessWithMessageOnly - 데이터 없이 메시지만 반환
// 사용 예: 삭제 성공 시
func SuccessWithMessageOnly(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
	})
}

// PaginationMeta - 페이지네이션 메타데이터
type PaginationMeta struct {
	Page       int   `json:"page"`        // 현재 페이지 (1부터 시작)
	PageSize   int   `json:"page_size"`   // 페이지당 항목 수
	TotalItems int64 `json:"total_items"` // 전체 항목 수
	TotalPages int   `json:"total_pages"` // 전체 페이지 수
}

// PaginatedResponse - 페이지네이션을 포함한 응답 구조
type PaginatedResponse struct {
	Success    bool           `json:"success"`
	Message    string         `json:"message"`
	Data       interface{}    `json:"data"`
	Pagination PaginationMeta `json:"pagination"`
}

// SuccessWithPagination - 페이지네이션 응답
// 사용 예:
//   meta := PaginationMeta{Page: 1, PageSize: 10, TotalItems: 100, TotalPages: 10}
//   SuccessWithPagination(c, http.StatusOK, "조회 성공", vehicles, meta)
func SuccessWithPagination(c *gin.Context, statusCode int, message string, data interface{}, pagination PaginationMeta) {
	c.JSON(statusCode, PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}
