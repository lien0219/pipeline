package middleware

import (
	"gin_pipeline/global"
	"gin_pipeline/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

// GinLogger 接收gin框架默认的日志
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		cost := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// 使用彩色日志
		if global.Config.Log.LogInConsole {
			// 根据状态码选择颜色
			var logFunc func(format string, args ...interface{})
			if statusCode >= 500 {
				logFunc = utils.Error
			} else if statusCode >= 400 {
				logFunc = utils.Warn
			} else {
				logFunc = utils.Info
			}

			// 格式化日志消息
			logMessage := "%s %s %d %s %s %s %s %s"
			logFunc(logMessage,
				method, path, statusCode,
				query, clientIP, userAgent,
				errorMessage, cost.String())
		}

		// 同时使用zap记录到文件
		global.Log.Info(path,
			zap.Int("status", statusCode),
			zap.String("method", method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", clientIP),
			zap.String("user-agent", userAgent),
			zap.String("errors", errorMessage),
			zap.Duration("cost", cost),
		)
	}
}

// GinRecovery recover掉项目可能出现的panic
func GinRecovery(stack bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否连接中断
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					// 使用彩色日志
					if global.Config.Log.LogInConsole {
						utils.Error("Broken pipe: %v\n%s", err, string(httpRequest))
					}

					// 同时记录到文件
					global.Log.Error(c.Request.URL.Path,
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)

					// 如果连接已断开，无法写入状态
					c.Error(err.(error))
					c.Abort()
					return
				}

				if stack {
					if global.Config.Log.LogInConsole {
						utils.Error("[Recovery from panic] %v\n%s\n%s",
							err,
							string(httpRequest),
							string(debug.Stack()))
					}

					global.Log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
						zap.String("stack", string(debug.Stack())),
					)
				} else {
					if global.Config.Log.LogInConsole {
						utils.Error("[Recovery from panic] %v\n%s",
							err,
							string(httpRequest))
					}

					global.Log.Error("[Recovery from panic]",
						zap.Any("error", err),
						zap.String("request", string(httpRequest)),
					)
				}
				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()
		c.Next()
	}
}
