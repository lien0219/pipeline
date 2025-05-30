package middleware

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Header("Content-Security-Policy", "default-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; font-src 'self' data:")
		} else {
			c.Header("Content-Security-Policy", "default-src 'self'")
		}
		c.Next()
	}
}

func AuditLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录审计日志
		if c.Request.Method != http.MethodGet {
			userID := c.GetUint("userID")
			log := model.AuditLog{
				UserID:     userID,
				Action:     c.Request.Method,
				ObjectType: c.FullPath(),
				RequestIP:  c.ClientIP(),
				Details:    c.Request.URL.String(),
			}
			global.DB.Create(&log)
		}
		c.Next()
	}
}
