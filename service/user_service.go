package service

import (
	"errors"
	"gin_pipeline/global"
	"gin_pipeline/model"

	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务
type UserService struct{}

// GetUserByID 通过ID获取用户
func (s *UserService) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := global.DB.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername 通过用户名获取用户
func (s *UserService) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := global.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(username, password, email string) (*model.User, error) {
	// 检查用户名是否已存在
	var count int64
	global.DB.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	global.DB.Model(&model.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		return nil, errors.New("email already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := model.User{
		Username: username,
		Password: string(hashedPassword),
		Email:    email,
	}

	if err := global.DB.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id uint, updates map[string]interface{}) error {
	return global.DB.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}

// VerifyPassword 验证密码
func (s *UserService) VerifyPassword(user *model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(id uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.GetUserByID(id)
	if err != nil {
		return err
	}

	// 验证旧密码
	if !s.VerifyPassword(user, oldPassword) {
		return errors.New("invalid old password")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	return global.DB.Model(user).Update("password", string(hashedPassword)).Error
}
