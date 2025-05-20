package utils

import (
	"errors"
	"gin_pipeline/global"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWT struct {
	SigningKey []byte
}

var (
	TokenExpired     = errors.New("token已过期")
	TokenNotValidYet = errors.New("token尚未生效")
	TokenMalformed   = errors.New("token格式错误")
	TokenInvalid     = errors.New("token无效")
)

// CustomClaims 自定义JWT载荷
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
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, TokenNotValidYet
			} else {
				return nil, TokenInvalid
			}
		}
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		return nil, TokenInvalid
	}
	return nil, TokenInvalid
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
