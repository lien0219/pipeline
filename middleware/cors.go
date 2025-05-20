package middleware

import (
	"gin_pipeline/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

// Cors 跨域中间件
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")

		// 允许所有来源
		if global.Config.CORS.Mode == "allow-all" {
			c.Header("Access-Control-Allow-Origin", "*")
		} else if global.Config.CORS.Mode == "whitelist" {
			// 检查是否在白名单中
			if checkOriginAllowed(origin, global.Config.CORS.Whitelist) {
				c.Header("Access-Control-Allow-Origin", origin)
			} else {
				c.Header("Access-Control-Allow-Origin", "")
			}
		} else {
			c.Header("Access-Control-Allow-Origin", origin)
		}

		// 设置允许的HTTP方法
		c.Header("Access-Control-Allow-Methods", global.Config.CORS.AllowMethods)

		// 设置允许的Header
		c.Header("Access-Control-Allow-Headers", global.Config.CORS.AllowHeaders)

		// 设置暴露的Header
		c.Header("Access-Control-Expose-Headers", global.Config.CORS.ExposeHeaders)

		// 是否允许携带凭证
		if global.Config.CORS.AllowCredentials {
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		// 设置预检请求的有效期
		c.Header("Access-Control-Max-Age", string(global.Config.CORS.MaxAge))

		// 放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// 检查来源是否在白名单中
func checkOriginAllowed(origin string, whitelist []string) bool {
	if len(whitelist) == 0 {
		return true
	}

	for _, allowed := range whitelist {
		if allowed == origin || allowed == "*" {
			return true
		}
		// 支持通配符匹配
		if strings.HasSuffix(allowed, "*") {
			prefix := strings.TrimSuffix(allowed, "*")
			if strings.HasPrefix(origin, prefix) {
				return true
			}
		}
	}

	return false
}
