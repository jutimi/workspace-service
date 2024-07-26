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

	// Workspace Errors
	ErrCodeWorkspaceNotFound    = 100
	ErrCodePassedLimitWorkspace = 101

	// User Workspace Errors
	ErrCodeUserWorkspaceNotFound          = 110
	ErrCodeUserWorkspaceNotInOrganization = 111

	// Organization Errors
	ErrCodeOrganizationNotFound         = 120
	ErrCodeInvalidParentOrganizationIds = 121
	ErrCodeInvalidLeaderIds             = 122

	ErrCodeInternalServerError = 500
	ErrCodeTimeout             = 408
	ErrCodeForbidden           = 403
	ErrCodeUnauthorized        = 401
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
		LangVN: "Bạn không đủ quyền truy cập tài nguyên này. Vui lòng kiểm tra lại quyền",
	},
	ErrCodeUnauthorized: {
		LangVN: "Bạn không có quyền truy cập tài nguyên này. Vui lòng kiểm tra lại quyền",
	},

	// Workspace Error
	ErrCodeWorkspaceNotFound: {
		LangVN: "Không tìm thấy workspace. Vui lòng kiểm tra lại",
	},
	ErrCodePassedLimitWorkspace: {
		LangVN: "Bạn đã sở hữu quá số workspace cho phép. Vui lòng liên hệ quản trị viên",
	},

	// User Workspace Error
	ErrCodeUserWorkspaceNotFound: {
		LangVN: "Không tìm thấy nhân viên. Vui không kiểm tra lại",
	},
	ErrCodeUserWorkspaceNotInOrganization: {
		LangVN: "Nhân viên không thuộc phòng ban. Vui phải kiểm tra lại",
	},

	// Organization Error
	ErrCodeOrganizationNotFound: {
		LangVN: "Không tìm thấy phòng ban. Vui phải kiểm tra lại",
	},
	ErrCodeInvalidParentOrganizationIds: {
		LangVN: "Thông tin phòng ban không hợp lệ. Vui phải kiểm tra lại",
	},
	ErrCodeInvalidLeaderIds: {
		LangVN: "Thông tin quản lý không hợp lệ. Vui phải kiểm tra lại",
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
