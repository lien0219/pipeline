package request

// Register 注册请求参数
type Register struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=20"`
	Name     string `json:"name" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
}

// Login 登录请求参数
type Login struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ChangePassword 修改密码请求参数
type ChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6,max=20"`
}

// UpdateUserInfo 更新用户信息请求参数
type UpdateUserInfo struct {
	Name   string `json:"name" binding:"required,min=2,max=50"`
	Email  string `json:"email" binding:"required,email"`
	Phone  string `json:"phone"`
	Avatar string `json:"avatar"`
}
