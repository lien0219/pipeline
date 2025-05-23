package middleware

import (
	"errors"
	"gin_pipeline/global"
	"gin_pipeline/model/response"
	"gin_pipeline/utils"
	_ "net/http"
	_ "strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Request.Header.Get("Authorization")
		if tokenString == "" {
			response.FailWithDetailed(gin.H{"reload": true}, "未登录或非法访问", c)
			c.Abort()
			return
		}

		if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}

		// 解析token
		j := utils.NewJWT()
		claims, err := j.ParseToken(tokenString)
		if err != nil {
			if errors.Is(err, utils.TokenExpired) {
				response.FailWithDetailed(gin.H{"reload": true}, "授权已过期", c)
				c.Abort()
				return
			}
			response.FailWithDetailed(gin.H{"reload": true}, err.Error(), c)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		c.Set("userId", claims.ID)
		c.Next()
	}
}

// JWT结构体
type JWT struct {
	SigningKey []byte
}

// 定义载荷
type CustomClaims struct {
	ID       uint
	Username string
	Role     string
	jwt.StandardClaims
}

// NewJWT 创建JWT实例
func NewJWT() *JWT {
	return &JWT{
		SigningKey: []byte(global.Config.System.JwtSecret),
	}
}

// CreateToken 创建Token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, utils.TokenExpired
			}
		}
		return nil, err
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, utils.TokenInvalid
	}
	return nil, utils.TokenInvalid
}

// GenerateToken 生成Token
func (j *JWT) GenerateToken(userID uint, username, role string) (string, error) {
	claims := CustomClaims{
		ID:       userID,
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix() - 1000,                                  // 签名生效时间
			ExpiresAt: time.Now().Unix() + int64(global.Config.System.JwtExpire), // 过期时间
			Issuer:    "gin-pipeline",                                            // 签名的发行者
		},
	}
	return j.CreateToken(claims)
}
