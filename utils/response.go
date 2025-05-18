package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 标准响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// SuccessResponse 成功响应
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}

// ErrorResponse 错误响应
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, Response{
		Code:    statusCode,
		Message: message,
	})
}

// BadRequestResponse 400错误响应
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// UnauthorizedResponse 401错误响应
func UnauthorizedResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusUnauthorized, "Unauthorized")
}

// ForbiddenResponse 403错误响应
func ForbiddenResponse(c *gin.Context) {
	ErrorResponse(c, http.StatusForbidden, "Forbidden")
}

// NotFoundResponse 404错误响应
func NotFoundResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse 500错误响应
func InternalServerErrorResponse(c *gin.Context, message string) {
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, http.StatusInternalServerError, message)
}
