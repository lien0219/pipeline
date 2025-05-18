package v1

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// UpdateUserRequest 更新用户信息请求
type UpdateUserRequest struct {
	Email string `json:"email" binding:"omitempty,email"`
	Name  string `json:"name" binding:"omitempty,max=50"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=32"`
}

// GetUserInfo godoc
// @Summary      Get current user info
// @Description  Get the current user's information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /users/me [get]
func GetUserInfo(c *gin.Context) {
	// 从上下文获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	// 查询用户
	var user model.User
	result := global.DB.First(&user, userId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"name":      user.Name,
			"createdAt": user.CreatedAt,
			"updatedAt": user.UpdatedAt,
		},
	})
}

// UpdateUserInfo godoc
// @Summary      Update user info
// @Description  Update the current user's information
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      UpdateUserRequest  true  "User information"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /users/me [put]
func UpdateUserInfo(c *gin.Context) {
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 从上下文获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	// 查询用户
	var user model.User
	result := global.DB.First(&user, userId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	// 更新用户信息
	updates := make(map[string]interface{})
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Name != "" {
		updates["name"] = req.Name
	}

	if len(updates) > 0 {
		if err := global.DB.Model(&user).Updates(updates).Error; err != nil {
			global.Logger.Errorf("Update user error: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "Failed to update user",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "User updated successfully",
		"data": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"email":     user.Email,
			"name":      user.Name,
			"updatedAt": user.UpdatedAt,
		},
	})
}

// ChangePassword godoc
// @Summary      Change password
// @Description  Change the current user's password
// @Tags         users
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      ChangePasswordRequest  true  "Password information"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Failure      404      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Router       /users/password [put]
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
		return
	}

	// 从上下文获取用户ID
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Unauthorized",
		})
		return
	}

	// 查询用户
	var user model.User
	result := global.DB.First(&user, userId)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "User not found",
		})
		return
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    http.StatusUnauthorized,
			"message": "Invalid old password",
		})
		return
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		global.Logger.Errorf("Password hash error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to process password",
		})
		return
	}

	// 更新密码
	if err := global.DB.Model(&user).Update("password", string(hashedPassword)).Error; err != nil {
		global.Logger.Errorf("Update password error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Password changed successfully",
	})
}
