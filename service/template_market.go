package service

import (
	"errors"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"github.com/Masterminds/semver/v3"
	"gorm.io/gorm"
	"strings"
)

// TemplateMarketService 模板市场服务
type TemplateMarketService struct{}

// CreateCategory 创建模板分类
func (s *TemplateMarketService) CreateCategory(category *model.TemplateCategory) error {
	return global.DB.Create(category).Error
}

// GetCategories 获取模板分类列表
func (s *TemplateMarketService) GetCategories() ([]model.TemplateCategory, error) {
	var categories []model.TemplateCategory
	if err := global.DB.Order("order ASC, id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}
	return categories, nil
}

// GetCategoryByID 根据ID获取模板分类
func (s *TemplateMarketService) GetCategoryByID(id uint) (*model.TemplateCategory, error) {
	var category model.TemplateCategory
	if err := global.DB.First(&category, id).Error; err != nil {
		return nil, err
	}
	return &category, nil
}

// UpdateCategory 更新模板分类
func (s *TemplateMarketService) UpdateCategory(id uint, updates map[string]interface{}) error {
	return global.DB.Model(&model.TemplateCategory{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteCategory 删除模板分类
func (s *TemplateMarketService) DeleteCategory(id uint) error {
	// 检查是否有模板使用此分类
	var count int64
	if err := global.DB.Model(&model.Template{}).Where("category_id = ?", id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("无法删除分类，有模板正在使用此分类")
	}
	return global.DB.Delete(&model.TemplateCategory{}, id).Error
}

// CreateTemplate 创建模板
func (s *TemplateMarketService) CreateTemplate(template *model.Template, version *model.TemplateVersion) error {
	// 验证版本号
	if _, err := semver.NewVersion(version.Version); err != nil {
		return errors.New("无效的版本号格式，请使用语义化版本（如1.0.0）")
	}

	// 开启事务
	tx := global.DB.Begin()

	// 创建模板
	if err := tx.Create(template).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 设置版本的模板ID
	version.TemplateID = template.ID
	version.IsLatest = true

	// 创建版本
	if err := tx.Create(version).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// GetTemplates 获取模板列表
func (s *TemplateMarketService) GetTemplates(categoryID uint, isPublic *bool) ([]model.Template, error) {
	query := global.DB.Model(&model.Template{}).Preload("Category").Preload("Creator")

	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	if isPublic != nil {
		query = query.Where("is_public = ?", *isPublic)
	}

	var templates []model.Template
	if err := query.Order("id DESC").Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}

// GetTemplateByID 根据ID获取模板
func (s *TemplateMarketService) GetTemplateByID(id uint) (*model.Template, error) {
	var template model.Template
	if err := global.DB.Preload("Category").Preload("Creator").Preload("Versions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at DESC")
	}).First(&template, id).Error; err != nil {
		return nil, err
	}
	return &template, nil
}

// UpdateTemplate 更新模板
func (s *TemplateMarketService) UpdateTemplate(id uint, updates map[string]interface{}) error {
	return global.DB.Model(&model.Template{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteTemplate 删除模板
func (s *TemplateMarketService) DeleteTemplate(id uint) error {
	// 开启事务
	tx := global.DB.Begin()

	// 删除模板版本
	if err := tx.Where("template_id = ?", id).Delete(&model.TemplateVersion{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 删除模板
	if err := tx.Delete(&model.Template{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// CreateTemplateVersion 创建模板版本
func (s *TemplateMarketService) CreateTemplateVersion(templateID uint, version *model.TemplateVersion) error {
	// 验证版本号
	newVer, err := semver.NewVersion(version.Version)
	if err != nil {
		return errors.New("无效的版本号格式，请使用语义化版本（如1.0.0）")
	}

	// 检查版本号是否已存在
	var count int64
	if err := global.DB.Model(&model.TemplateVersion{}).Where("template_id = ? AND version = ?", templateID, version.Version).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("版本号已存在")
	}

	// 获取最新版本
	var latestVersion model.TemplateVersion
	if err := global.DB.Where("template_id = ? AND is_latest = ?", templateID, true).First(&latestVersion).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 如果找到最新版本，检查新版本是否更高
	if latestVersion.ID > 0 {
		latestVer, err := semver.NewVersion(latestVersion.Version)
		if err != nil {
			return errors.New("现有版本号格式错误")
		}

		if !newVer.GreaterThan(latestVer) {
			return errors.New("新版本号必须大于当前最新版本")
		}
	}

	// 开启事务
	tx := global.DB.Begin()

	// 如果有最新版本，将其设为非最新
	if latestVersion.ID > 0 {
		if err := tx.Model(&latestVersion).Update("is_latest", false).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// 设置版本的模板ID和最新标志
	version.TemplateID = templateID
	version.IsLatest = true

	// 创建新版本
	if err := tx.Create(version).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

// GetTemplateVersions 获取模板版本列表
func (s *TemplateMarketService) GetTemplateVersions(templateID uint) ([]model.TemplateVersion, error) {
	var versions []model.TemplateVersion
	if err := global.DB.Where("template_id = ?", templateID).Order("created_at DESC").Find(&versions).Error; err != nil {
		return nil, err
	}
	return versions, nil
}

// GetTemplateVersionByID 根据ID获取模板版本
func (s *TemplateMarketService) GetTemplateVersionByID(id uint) (*model.TemplateVersion, error) {
	var version model.TemplateVersion
	if err := global.DB.First(&version, id).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// GetLatestTemplateVersion 获取模板的最新版本
func (s *TemplateMarketService) GetLatestTemplateVersion(templateID uint) (*model.TemplateVersion, error) {
	var version model.TemplateVersion
	if err := global.DB.Where("template_id = ? AND is_latest = ?", templateID, true).First(&version).Error; err != nil {
		return nil, err
	}
	return &version, nil
}

// DeleteTemplateVersion 删除模板版本
func (s *TemplateMarketService) DeleteTemplateVersion(id uint) error {
	// 获取版本信息
	var version model.TemplateVersion
	if err := global.DB.First(&version, id).Error; err != nil {
		return err
	}

	// 检查是否为最新版本
	if version.IsLatest {
		return errors.New("无法删除最新版本")
	}

	// 删除版本
	return global.DB.Delete(&version).Error
}

// SearchTemplates 搜索模板
func (s *TemplateMarketService) SearchTemplates(keyword string, categoryID uint, tags string) ([]model.Template, error) {
	query := global.DB.Model(&model.Template{}).Preload("Category").Preload("Creator")

	// 关键字搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 分类筛选
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}

	// 标签筛选
	if tags != "" {
		tagList := strings.Split(tags, ",")
		for _, tag := range tagList {
			query = query.Where("tags LIKE ?", "%"+strings.TrimSpace(tag)+"%")
		}
	}

	// 只显示公开模板
	query = query.Where("is_public = ?", true)

	var templates []model.Template
	if err := query.Order("download_count DESC, id DESC").Find(&templates).Error; err != nil {
		return nil, err
	}

	return templates, nil
}

// IncrementTemplateDownloadCount 增加模板下载次数
func (s *TemplateMarketService) IncrementTemplateDownloadCount(templateID uint) error {
	return global.DB.Model(&model.Template{}).Where("id = ?", templateID).Update("download_count", gorm.Expr("download_count + ?", 1)).Error
}

// IncrementVersionDownloadCount 增加版本下载次数
func (s *TemplateMarketService) IncrementVersionDownloadCount(versionID uint) error {
	return global.DB.Model(&model.TemplateVersion{}).Where("id = ?", versionID).Update("download_count", gorm.Expr("download_count + ?", 1)).Error
}

// SetVersionAsLatest 设置版本为最新版本
func (s *TemplateMarketService) SetVersionAsLatest(versionID uint) error {
	// 获取版本信息
	var version model.TemplateVersion
	if err := global.DB.First(&version, versionID).Error; err != nil {
		return err
	}

	// 开启事务
	tx := global.DB.Begin()

	// 将当前最新版本设为非最新
	if err := tx.Model(&model.TemplateVersion{}).Where("template_id = ? AND is_latest = ?", version.TemplateID, true).Update("is_latest", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 将指定版本设为最新
	if err := tx.Model(&version).Update("is_latest", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
