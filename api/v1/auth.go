package v1

import (
	"gin_pipeline/config"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=32"`
	Email    string `json:"email" binding:"required,email"`
}

// Login godoc
// @Summary      User login
// @Description  Login with username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      LoginRequest  true  "Login credentials"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /auth/login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 查询用户
	var user model.User
	result := global.DB.Where("username = ?", req.Username).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Invalid username or password",
		})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Invalid username or password",
		})
		return
	}

	// 生成JWT
	token, err := generateToken(user.ID, user.Username)
	if err != nil {
		global.Logger.Errorf("Generate token error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Login successful",
		"data": gin.H{
			"token": token,
			"user": gin.H{
				"id":       user.ID,
				"username": user.Username,
				"email":    user.Email,
			},
		},
	})
}

// Register godoc
// @Summary      User registration
// @Description  Register a new user
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request  body      RegisterRequest  true  "Registration information"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      409      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /auth/register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 检查用户名是否已存在
	var count int64
	global.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Username already exists",
		})
		return
	}

	// 检查邮箱是否已存在
	global.DB.Model(&model.User{}).Where("email = ?", req.Email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"code":    http.StatusConflict,
			"message": "Email already exists",
		})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		global.Logger.Errorf("Password hash error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to process password",
		})
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := global.DB.Create(&user).Error; err != nil {
		global.Logger.Errorf("Create user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    http.StatusCreated,
		"message": "User registered successfully",
		"data": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

// 生成JWT令牌
func generateToken(userId uint, username string) (string, error) {
	// 设置过期时间
	expireTime := time.Now().Add(time.Duration(config.Get().JWT.ExpireTime) * time.Hour)

	// 创建声明
	claims := jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      expireTime.Unix(),
		"iat":      time.Now().Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	return token.SignedString([]byte(config.Get().JWT.Secret))
}
