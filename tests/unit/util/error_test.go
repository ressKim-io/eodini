package util_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/hyeokjun/eodini/internal/util"
	"github.com/stretchr/testify/assert"
)

// TestNewNotFoundError - NotFound 에러 생성 테스트
func TestNewNotFoundError(t *testing.T) {
	// Given
	resource := "차량"

	// When
	err := util.NewNotFoundError(resource)

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeNotFound, err.Code)
	assert.Equal(t, http.StatusNotFound, err.StatusCode)
	assert.Contains(t, err.Message, resource)
}

// TestNewValidationError - Validation 에러 생성 테스트
func TestNewValidationError(t *testing.T) {
	// Given
	message := "입력값이 올바르지 않습니다"
	details := map[string]interface{}{
		"field": "email",
		"value": "invalid-email",
	}

	// When
	err := util.NewValidationError(message, details)

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeValidation, err.Code)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, "email", err.Details["field"])
}

// TestNewUnauthorizedError - Unauthorized 에러 생성 테스트
func TestNewUnauthorizedError(t *testing.T) {
	// When
	err := util.NewUnauthorizedError()

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeUnauthorized, err.Code)
	assert.Equal(t, http.StatusUnauthorized, err.StatusCode)
	assert.NotEmpty(t, err.Message)
}

// TestNewForbiddenError - Forbidden 에러 생성 테스트
func TestNewForbiddenError(t *testing.T) {
	// When
	err := util.NewForbiddenError()

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeForbidden, err.Code)
	assert.Equal(t, http.StatusForbidden, err.StatusCode)
	assert.NotEmpty(t, err.Message)
}

// TestNewInternalError - Internal 에러 생성 테스트
func TestNewInternalError(t *testing.T) {
	// Given
	originalErr := errors.New("database connection failed")

	// When
	err := util.NewInternalError(originalErr)

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeInternal, err.Code)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
	assert.NotEmpty(t, err.Message)
	assert.Equal(t, originalErr.Error(), err.Details["error"])
}

// TestNewInternalError_NilError - nil 에러로 Internal 에러 생성
func TestNewInternalError_NilError(t *testing.T) {
	// When
	err := util.NewInternalError(nil)

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeInternal, err.Code)
	assert.Equal(t, http.StatusInternalServerError, err.StatusCode)
}

// TestNewDuplicateError - Duplicate 에러 생성 테스트
func TestNewDuplicateError(t *testing.T) {
	// Given
	resource := "차량 번호"

	// When
	err := util.NewDuplicateError(resource)

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeDuplicate, err.Code)
	assert.Equal(t, http.StatusConflict, err.StatusCode)
	assert.Contains(t, err.Message, resource)
}

// TestNewBadRequestError - BadRequest 에러 생성 테스트
func TestNewBadRequestError(t *testing.T) {
	// Given
	message := "잘못된 요청 파라미터"

	// When
	err := util.NewBadRequestError(message)

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeBadRequest, err.Code)
	assert.Equal(t, http.StatusBadRequest, err.StatusCode)
	assert.Equal(t, message, err.Message)
}

// TestNewConflictError - Conflict 에러 생성 테스트
func TestNewConflictError(t *testing.T) {
	// Given
	message := "운행 중인 차량은 삭제할 수 없습니다"

	// When
	err := util.NewConflictError(message)

	// Then
	assert.NotNil(t, err)
	assert.Equal(t, util.ErrCodeConflict, err.Code)
	assert.Equal(t, http.StatusConflict, err.StatusCode)
	assert.Equal(t, message, err.Message)
}

// TestAppError_Error - error 인터페이스 구현 테스트
func TestAppError_Error(t *testing.T) {
	// Given
	expectedMessage := "테스트 에러 메시지"
	err := util.NewBadRequestError(expectedMessage)

	// When
	errorString := err.Error()

	// Then
	assert.Equal(t, expectedMessage, errorString)
}

// TestAppError_AsError - error 타입으로 사용 가능한지 테스트
func TestAppError_AsError(t *testing.T) {
	// Given
	var err error = util.NewNotFoundError("차량")

	// Then
	assert.NotNil(t, err)
	assert.Error(t, err)

	// error를 AppError로 타입 단언 가능
	appErr, ok := err.(*util.AppError)
	assert.True(t, ok)
	assert.Equal(t, util.ErrCodeNotFound, appErr.Code)
}
