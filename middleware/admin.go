package middleware

import (
	"gin_pipeline/model/response"
	"github.com/gin-gonic/gin"
)

// AdminAuth 管理员权限中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			response.FailWithMessage("未登录或非法访问", c)
			c.Abort()
			return
		}

		// 检查用户角色
		if claims.(*CustomClaims).Role != "admin" {
			response.FailWithMessage("权限不足", c)
			c.Abort()
			return
		}

		c.Next()
	}
}
