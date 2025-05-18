package middleware

import (
	"fmt"
	"gin_pipeline/global"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义颜色常量
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorWhite  = "\033[37m"
)

// 获取状态码对应的颜色
func getStatusColor(code int) string {
	switch {
	case code >= 200 && code < 300:
		return colorGreen // 成功状态码使用绿色
	case code >= 300 && code < 400:
		return colorCyan // 重定向状态码使用青色
	case code >= 400 && code < 500:
		return colorYellow // 客户端错误使用黄色
	default:
		return colorRed // 服务器错误使用红色
	}
}

// 获取HTTP方法对应的颜色
func getMethodColor(method string) string {
	switch method {
	case "GET":
		return colorBlue
	case "POST":
		return colorGreen
	case "PUT":
		return colorYellow
	case "DELETE":
		return colorRed
	case "PATCH":
		return colorPurple
	case "HEAD":
		return colorCyan
	case "OPTIONS":
		return colorWhite
	default:
		return colorReset
	}
}

// Logger 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()

		// 执行时间
		latency := endTime.Sub(startTime)

		// 请求方式
		reqMethod := c.Request.Method

		// 请求路由
		reqUri := c.Request.RequestURI

		// 状态码
		statusCode := c.Writer.Status()

		// 请求IP
		clientIP := c.ClientIP()

		// 状态码颜色
		statusColor := getStatusColor(statusCode)

		// 方法颜色
		methodColor := getMethodColor(reqMethod)

		// 日志格式 - 使用彩色输出
		logMsg := fmt.Sprintf("| %s%3d%s | %13v | %15s | %s%s%s | %s |",
			statusColor, statusCode, colorReset,
			latency,
			clientIP,
			methodColor, reqMethod, colorReset,
			reqUri,
		)

		// 根据状态码选择日志级别
		if statusCode >= 500 {
			global.Logger.Error(logMsg)
		} else if statusCode >= 400 {
			global.Logger.Warn(logMsg)
		} else {
			global.Logger.Info(logMsg)
		}
	}
}
