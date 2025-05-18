package middleware

import (
	"gin_pipeline/config"
	"gin_pipeline/global"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Authorization header format must be Bearer {token}",
			})
			c.Abort()
			return
		}

		// 解析token
		token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.NewValidationError("unexpected signing method", jwt.ValidationErrorSignatureInvalid)
			}
			return []byte(config.Get().JWT.Secret), nil
		})

		if err != nil {
			global.Logger.Errorf("JWT parse error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// 验证token
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 检查过期时间
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"code":    http.StatusUnauthorized,
					"message": "Token has expired",
				})
				c.Abort()
				return
			}

			// 将用户信息存储到上下文
			c.Set("userId", claims["userId"])
			c.Set("username", claims["username"])
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "Invalid token claims",
			})
			c.Abort()
			return
		}
	}
}
