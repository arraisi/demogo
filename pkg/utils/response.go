package utils

import (
	"github.com/arraisi/demogo/pkg/constant"
)

// GetErrorResponse is a function for build errorMessage
func GetErrorResponse(message string, statusCode int, code string) *constant.ApplicationError {
	return &constant.ApplicationError{
		Message:    message,
		StatusCode: statusCode,
		Code:       code,
	}
}
