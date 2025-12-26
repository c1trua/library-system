package handlers

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// 错误响应
func Error(c *gin.Context, httpStatus int, message string, err error) {
	response := ErrorResponse{
		Message: message,
	}

	// 生产环境隐藏详细错误
	if gin.Mode() != gin.ReleaseMode && err != nil {
		response.Error = err.Error()
	}

	c.JSON(httpStatus, response)
}

// 常用的错误响应
func BadRequest(c *gin.Context, message string, err error) {
	Error(c, 400, message, err)
}

func Unauthorized(c *gin.Context, message string, err error) {
	Error(c, 401, message, err)
}

func Forbidden(c *gin.Context, message string, err error) {
	Error(c, 403, message, err)
}

func NotFound(c *gin.Context, message string, err error) {
	Error(c, 404, message, err)
}

func Conflict(c *gin.Context, message string, err error) {
	Error(c, 409, message, err)
}

func InternalError(c *gin.Context, message string, err error) {
	Error(c, 500, message, err)
}
