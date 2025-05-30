package utils

import (
	"encoding/base64"
	"errors"
	"gin_pipeline/global"
	"time"

	"github.com/dgrijalva/jwt-go"
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
	decodedKey, err := base64.StdEncoding.DecodeString(global.Config.System.JwtSecret)
	if err != nil {
		decodedKey = []byte(global.Config.System.JwtSecret)
	}
	return &JWT{
		SigningKey: decodedKey,
	}
}

// CreateToken 创建Token
func (j *JWT) CreateToken(claims CustomClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token.Header["kid"] = RandomString(8) // 添加key id
	token.Header["cty"] = "JWT"           // 明确内容类型
	return token.SignedString(j.SigningKey)
}

// ParseToken 解析Token
func (j *JWT) ParseToken(tokenString string) (*CustomClaims, error) {
	if len(tokenString) < 50 {
		return nil, errors.New("token格式错误")
	}
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.SigningKey, nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				global.Log.Error("Token 已过期", "error", err)
				return nil, TokenExpired
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				global.Log.Error("Token 尚未生效", "error", err)
				return nil, TokenNotValidYet
			} else {
				global.Log.Error("Token 无效", "error", err)
				return nil, TokenInvalid
			}
		}
		global.Log.Error("Token 解析错误", "error", err)
		return nil, TokenInvalid
	}
	if token != nil {
		if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
			return claims, nil
		}
		global.Log.Error("Token 声明验证失败")
		return nil, TokenInvalid
	}
	global.Log.Error("Token 为空")
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
