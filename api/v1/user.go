package v1

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"gin_pipeline/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body request.Register true "用户注册信息"
// @Success 200 {object} response.Response{data=model.User} "注册成功"
// @Router /user/register [post]
func Register(c *gin.Context) {
	var req request.Register
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 检查用户名是否已存在
	var count int64
	if err := global.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&count).Error; err != nil {
		global.Log.Error("查询用户失败", zap.Error(err))
		response.FailWithMessage("注册失败", c)
		return
	}
	if count > 0 {
		response.FailWithMessage("用户名已存在", c)
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		global.Log.Error("密码加密失败", zap.Error(err))
		response.FailWithMessage("注册失败", c)
		return
	}

	// 创建用户
	user := model.User{
		Username: req.Username,
		Password: hashedPassword,
		Name:     req.Name,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     "user", // 默认角色
	}

	if err := global.DB.Create(&user).Error; err != nil {
		global.Log.Error("创建用户失败", zap.Error(err))
		response.FailWithMessage("注册失败", c)
		return
	}

	// 隐藏密码
	user.Password = ""
	response.OkWithData(user, c)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param data body request.Login true "用户登录信息"
// @Success 200 {object} response.Response{data=map[string]interface{}} "登录成功"
// @Router /user/login [post]
func Login(c *gin.Context) {
	var req request.Login
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 查询用户
	var user model.User
	if err := global.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		response.FailWithMessage("用户名或密码错误", c)
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		response.FailWithMessage("用户名或密码错误", c)
		return
	}

	// 生成Token
	j := utils.NewJWT()
	token, err := j.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		global.Log.Error("生成Token失败", zap.Error(err))
		response.FailWithMessage("登录失败", c)
		return
	}

	// 更新最后登录时间
	global.DB.Model(&user).Update("last_login", time.Now())

	// 返回结果
	data := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"username": user.Username,
			"name":     user.Name,
			"email":    user.Email,
			"phone":    user.Phone,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	}
	response.OkWithData(data, c)
}

// GetUserInfo 获取用户信息
// @Summary 获取用户信息
// @Description 获取当前登录用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.User} "获取成功"
// @Router /user/info [get]
func GetUserInfo(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	// 查询用户
	var user model.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		global.Log.Error("查询用户失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	// 隐藏密码
	user.Password = ""
	response.OkWithData(user, c)
}

// UpdateUserInfo 更新用户信息
// @Summary 更新用户信息
// @Description 更新当前登录用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.UpdateUserInfo true "用户信息"
// @Success 200 {object} response.Response{data=model.User} "更新成功"
// @Router /user/info [put]
func UpdateUserInfo(c *gin.Context) {
	var req request.UpdateUserInfo
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("更新用户信息失败", c)
		return
	}

	// 更新用户信息
	updates := map[string]interface{}{
		"name":   req.Name,
		"email":  req.Email,
		"phone":  req.Phone,
		"avatar": req.Avatar,
	}

	if err := global.DB.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error; err != nil {
		global.Log.Error("更新用户信息失败", zap.Error(err))
		response.FailWithMessage("更新用户信息失败", c)
		return
	}

	// 查询更新后的用户信息
	var user model.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		global.Log.Error("查询用户失败", zap.Error(err))
		response.FailWithMessage("获取用户信息失败", c)
		return
	}

	// 隐藏密码
	user.Password = ""
	response.OkWithData(user, c)
}

// ChangePassword 修改密码
// @Summary 修改密码
// @Description 修改当前登录用户密码
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.ChangePassword true "密码信息"
// @Success 200 {object} response.Response "修改成功"
// @Router /user/password [put]
func ChangePassword(c *gin.Context) {
	var req request.ChangePassword
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("修改密码失败", c)
		return
	}

	// 查询用户
	var user model.User
	if err := global.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		global.Log.Error("查询用户失败", zap.Error(err))
		response.FailWithMessage("修改密码失败", c)
		return
	}

	// 验证旧密码
	if !utils.CheckPasswordHash(req.OldPassword, user.Password) {
		response.FailWithMessage("旧密码错误", c)
		return
	}

	// 密码加密
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		global.Log.Error("密码加密失败", zap.Error(err))
		response.FailWithMessage("修改密码失败", c)
		return
	}

	// 更新密码
	if err := global.DB.Model(&user).Update("password", hashedPassword).Error; err != nil {
		global.Log.Error("更新密码失败", zap.Error(err))
		response.FailWithMessage("修改密码失败", c)
		return
	}

	response.OkWithMessage("修改密码成功", c)
}
