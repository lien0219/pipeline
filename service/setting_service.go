package service

import (
	"encoding/json"
	"errors"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"

	"gorm.io/gorm"
)

// SettingService 配置服务
type SettingService struct{}

// SaveSettings 保存配置
func (s *SettingService) SaveSettings(req request.SaveSettingsRequest) error {
	// 开启事务
	tx := global.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 遍历配置项并保存
	for key, value := range req.Items {
		// 将value转换为JSON字符串
		valueStr, err := json.Marshal(value)
		if err != nil {
			tx.Rollback()
			return errors.New("配置值序列化失败: " + err.Error())
		}

		// 查找是否已存在该配置
		var setting model.Setting
		result := tx.Where("`key` = ?", key).First(&setting)

		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return result.Error
		}

		// 更新或创建配置
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// 创建新配置
			setting = model.Setting{
				Key:   key,
				Value: string(valueStr),
				Type:  req.Type,
			}
			if err := tx.Create(&setting).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// 更新现有配置
			setting.Value = string(valueStr)
			setting.Type = req.Type
			if err := tx.Save(&setting).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 提交事务
	return tx.Commit().Error
}

// GetSettings 获取配置
func (s *SettingService) GetSettings(req request.GetSettingsRequest) (map[string]interface{}, error) {
	var settings []model.Setting
	query := global.DB

	// 如果指定了类型，则筛选
	if req.Type != "" && req.Type != "all" {
		query = query.Where("type = ?", req.Type)
	}

	// 查询配置
	if err := query.Find(&settings).Error; err != nil {
		return nil, err
	}

	// 构建返回结果
	result := make(map[string]interface{})
	for _, setting := range settings {
		var value interface{}
		if err := json.Unmarshal([]byte(setting.Value), &value); err != nil {
			// 如果解析失败，直接返回原始字符串
			value = setting.Value
		}
		result[setting.Key] = value
	}

	return result, nil
}
