package util_test

import (
	"testing"

	"github.com/hyeokjun/eodini/internal/util"
	"github.com/stretchr/testify/assert"
)

// TestGetMessage_Success - 기본 메시지 조회 테스트
func TestGetMessage_Success(t *testing.T) {
	// When
	msg := util.GetMessage(util.MsgSuccess)

	// Then
	assert.NotEmpty(t, msg)
	assert.Equal(t, "요청이 성공적으로 처리되었습니다", msg)
}

// TestGetMessage_WithArgs - 포맷팅 인자가 있는 메시지 조회
func TestGetMessage_WithArgs(t *testing.T) {
	// When
	msg := util.GetMessage(util.MsgCreated, "차량")

	// Then
	assert.NotEmpty(t, msg)
	assert.Contains(t, msg, "차량")
	assert.Equal(t, "차량이(가) 생성되었습니다", msg)
}

// TestGetMessage_MultipleArgs - 여러 인자를 받는 메시지
func TestGetMessage_MultipleArgs(t *testing.T) {
	// Given: 커스텀 메시지 추가
	util.AddMessage("TEST_MULTIPLE", "%s와(과) %s이(가) 생성되었습니다")

	// When
	msg := util.GetMessage("TEST_MULTIPLE", "차량", "운전자")

	// Then
	assert.Equal(t, "차량와(과) 운전자이(가) 생성되었습니다", msg)
}

// TestGetMessage_NonExistentKey - 존재하지 않는 키
func TestGetMessage_NonExistentKey(t *testing.T) {
	// When
	msg := util.GetMessage("NON_EXISTENT_KEY")

	// Then
	assert.Equal(t, "NON_EXISTENT_KEY", msg) // 키를 그대로 반환
}

// TestAddMessage - 런타임 메시지 추가
func TestAddMessage(t *testing.T) {
	// Given
	key := "CUSTOM_MESSAGE"
	customMsg := "커스텀 메시지입니다"

	// When
	util.AddMessage(key, customMsg)
	msg := util.GetMessage(key)

	// Then
	assert.Equal(t, customMsg, msg)
}

// TestGetAllMessages - 전체 메시지 조회
func TestGetAllMessages(t *testing.T) {
	// When
	allMessages := util.GetAllMessages()

	// Then
	assert.NotEmpty(t, allMessages)
	assert.Contains(t, allMessages, util.MsgSuccess)
	assert.Contains(t, allMessages, util.MsgCreated)
	assert.Contains(t, allMessages, util.MsgResourceNotFound)
}

// TestGetAllMessages_Immutability - 반환된 맵 수정이 원본에 영향 없는지 확인
func TestGetAllMessages_Immutability(t *testing.T) {
	// Given
	allMessages1 := util.GetAllMessages()
	originalCount := len(allMessages1)

	// When: 반환된 맵에 항목 추가
	allMessages1["NEW_KEY"] = "새로운 값"
	allMessages2 := util.GetAllMessages()

	// Then: 원본은 변경되지 않음
	assert.Equal(t, originalCount, len(allMessages2))
	assert.NotContains(t, allMessages2, "NEW_KEY")
}
