package util

import "fmt"

// 📝 설명: Spring의 messages.properties처럼 메시지 중앙 관리
// 🎯 실무 포인트: 다국어 지원 확장 가능, 메시지 일관성 유지
// ⚠️ 주의사항: fmt.Sprintf 포맷 순서 주의 (%s 순서와 args 순서 일치)

// 메시지 키 상수
const (
	// 성공 메시지
	MsgSuccess = "SUCCESS"
	MsgCreated = "CREATED"
	MsgUpdated = "UPDATED"
	MsgDeleted = "DELETED"

	// 에러 메시지
	MsgResourceNotFound = "RESOURCE_NOT_FOUND"
	MsgUnauthorized     = "UNAUTHORIZED"
	MsgForbidden        = "FORBIDDEN"
	MsgInternalError    = "INTERNAL_ERROR"
	MsgDuplicate        = "DUPLICATE"
	MsgValidationFailed = "VALIDATION_FAILED"
	MsgBadRequest       = "BAD_REQUEST"
	MsgConflict         = "CONFLICT"

	// 특정 리소스 메시지
	MsgVehicleNotFound  = "VEHICLE_NOT_FOUND"
	MsgDriverNotFound   = "DRIVER_NOT_FOUND"
	MsgRouteNotFound    = "ROUTE_NOT_FOUND"
	MsgScheduleNotFound = "SCHEDULE_NOT_FOUND"
)

// 메시지 맵 (한국어)
// 추후 다국어 지원 시 messages_en.go, messages_ko.go로 분리 가능
var messages = map[string]string{
	// 성공 메시지
	MsgSuccess: "요청이 성공적으로 처리되었습니다",
	MsgCreated: "%s이(가) 생성되었습니다",
	MsgUpdated: "%s이(가) 수정되었습니다",
	MsgDeleted: "%s이(가) 삭제되었습니다",

	// 에러 메시지
	MsgResourceNotFound: "%s을(를) 찾을 수 없습니다",
	MsgUnauthorized:     "인증이 필요합니다",
	MsgForbidden:        "접근 권한이 없습니다",
	MsgInternalError:    "서버 내부 오류가 발생했습니다",
	MsgDuplicate:        "이미 존재하는 %s입니다",
	MsgValidationFailed: "입력값 검증에 실패했습니다",
	MsgBadRequest:       "잘못된 요청입니다",
	MsgConflict:         "요청이 현재 상태와 충돌합니다",

	// 특정 리소스 메시지
	MsgVehicleNotFound:  "차량을 찾을 수 없습니다",
	MsgDriverNotFound:   "운전자를 찾을 수 없습니다",
	MsgRouteNotFound:    "경로를 찾을 수 없습니다",
	MsgScheduleNotFound: "운행 일정을 찾을 수 없습니다",
}

// getMessage - 메시지 조회 (내부 사용)
// 포맷팅 지원: getMessage(MsgCreated, "차량") -> "차량이(가) 생성되었습니다"
func getMessage(key string, args ...interface{}) string {
	msg, exists := messages[key]
	if !exists {
		return key // 메시지가 없으면 키를 그대로 반환
	}

	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// GetMessage - 외부에서 사용할 수 있는 메시지 조회 함수
// 사용 예:
//   GetMessage(MsgSuccess) -> "요청이 성공적으로 처리되었습니다"
//   GetMessage(MsgCreated, "차량") -> "차량이(가) 생성되었습니다"
func GetMessage(key string, args ...interface{}) string {
	return getMessage(key, args...)
}

// AddMessage - 런타임에 메시지 추가 (필요 시 사용)
// 사용 예: AddMessage("CUSTOM_MSG", "커스텀 메시지입니다")
func AddMessage(key, message string) {
	messages[key] = message
}

// GetAllMessages - 모든 메시지 조회 (디버깅/문서화 용도)
func GetAllMessages() map[string]string {
	// 원본 맵을 복사하여 반환 (외부에서 수정 방지)
	result := make(map[string]string, len(messages))
	for k, v := range messages {
		result[k] = v
	}
	return result
}
