package errors

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"

	logger "workspace-server/package/log"
	_validator "workspace-server/package/validator"
)

type CustomError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

const (
	LangVN = "vi"
)

const (
	ErrCodeValidatorRequired     = 1
	ErrCodeValidatorFormat       = 2
	ErrCodeValidatorVerifiedData = 3

	ErrCodeUserNotFound = 10
	ErrCodeUserExisted  = 11

	ErrCodeTokenExpired      = 20
	ErrCodeIncorrectPassword = 21

	ErrCodeInternalServerError = 500
	ErrCodeTimeout             = 408
	ErrCodeForbidden           = 403
)

var messages = map[int]map[string]string{
	// Validator
	ErrCodeValidatorRequired: {
		LangVN: "không được bỏ trống. Vui lòng kiểm tra lại",
	},
	ErrCodeValidatorFormat: {
		LangVN: "không hợp lệ. Vui lòng kiểm tra lại",
	},
	ErrCodeValidatorVerifiedData: {
		LangVN: "không chính xác. Vui lòng kiểm tra lại",
	},

	ErrCodeInternalServerError: {
		LangVN: "Hệ thống gặp lỗi. Vui lòng thử lại sau",
	},
	ErrCodeTimeout: {
		LangVN: "Hệ thống gặp lỗi. Vui lòng thử lại sau",
	},
	ErrCodeForbidden: {
		LangVN: "Hệ thống gặp lỗi. Vui lòng thử lại sau",
	},

	// User Error
	ErrCodeUserNotFound: {
		LangVN: "Không tìm thấy người dùng. Vui lòng kiểm tra lại",
	},
	ErrCodeUserExisted: {
		LangVN: "Người dùng đã đăng ký tài khoản. Vui lòng kiểm tra lại",
	},

	// OAuth Error
	ErrCodeTokenExpired: {
		LangVN: "Phiên làm việc đã hết hạn. Vui lòng đăng nhâp lại",
	},
	ErrCodeIncorrectPassword: {
		LangVN: "Tên đăng nhập hoặc mật khẩu không đúng. Vui lòng kiểm tra lại",
	},
}

func New(code int) *CustomError {
	return &CustomError{
		Code:    code,
		Message: messages[code][LangVN],
	}
}

func NewCustomError(code int, message string) *CustomError {
	return &CustomError{
		Code:    code,
		Message: message,
	}
}

func NewValidatorError(err error) *CustomError {
	var validatorErr validator.ValidationErrors
	if errors.As(err, &validatorErr) {
		errDetail := validatorErr[0]

		field := errDetail.Field()
		tag := errDetail.Tag()

		code := convertValidatorTag(tag)
		return &CustomError{
			Code:    code,
			Message: fmt.Sprintf("%s %s", field, messages[code][LangVN]),
		}

	}

	return New(ErrCodeInternalServerError)
}

func (err *CustomError) Error() string {
	return err.Message
}

func (err *CustomError) GetCode() int {
	return err.Code
}

// --------------------------------------
func convertValidatorTag(tag string) int {
	logger.GetLogger().Info("validation_tag: ", tag)
	switch tag {
	case _validator.EMAIL, _validator.PHONE_NUMBER:
		return ErrCodeValidatorFormat
	case _validator.EQUAL_FIELD:
		return ErrCodeValidatorVerifiedData
	default:
		return ErrCodeValidatorRequired
	}
}
