package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hyeokjun/eodini/internal/util"
)

// 📝 설명: Spring의 @ControllerAdvice처럼 모든 핸들러의 에러를 중앙에서 처리
// 🎯 실무 포인트: 핸들러에서는 c.Error()만 호출, 응답은 여기서 통일
// ⚠️ 주의사항: 반드시 라우터 설정 시 마지막에 등록 (c.Next() 이후 실행)

// ErrorHandler - 글로벌 에러 처리 미들웨어
// Spring의 @ControllerAdvice + @ExceptionHandler 역할
//
// 사용 예:
//   router := gin.Default()
//   router.Use(middleware.ErrorHandler())
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 다음 핸들러 실행
		c.Next()

		// 에러가 있는지 확인
		if len(c.Errors) > 0 {
			// 마지막 에러만 처리 (여러 에러가 있을 경우)
			err := c.Errors.Last().Err

			// AppError 타입인 경우
			if appErr, ok := err.(*util.AppError); ok {
				util.ErrorResponse(c, appErr)
				return
			}

			// 일반 에러인 경우 -> Internal Server Error로 변환
			appErr := util.NewInternalError(err)
			util.ErrorResponse(c, appErr)
		}
	}
}

// RecoveryHandler - Panic 복구 미들웨어
// 핸들러에서 panic 발생 시 500 에러로 변환
//
// 사용 예:
//   router := gin.Default()
//   router.Use(middleware.RecoveryHandler())
func RecoveryHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// panic을 AppError로 변환
				var appErr *util.AppError

				// err가 error 인터페이스를 구현한 경우
				if e, ok := err.(error); ok {
					appErr = util.NewInternalError(e)
				} else {
					// 그 외의 경우 (string 등)
					appErr = &util.AppError{
						Code:       util.ErrCodeInternal,
						Message:    util.GetMessage(util.MsgInternalError),
						StatusCode: 500,
						Details:    map[string]interface{}{"panic": err},
					}
				}

				util.ErrorResponse(c, appErr)
				c.Abort() // 이후 핸들러 실행 중단
			}
		}()

		c.Next()
	}
}
