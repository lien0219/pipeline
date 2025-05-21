package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"sync"
)

// YAMLValidator YAML验证器
type YAMLValidator struct {
	schemaCache     map[string]string // 缓存Schema
	schemaCacheLock sync.RWMutex
}

// NewYAMLValidator 创建YAML验证器
func NewYAMLValidator() *YAMLValidator {
	return &YAMLValidator{
		schemaCache: make(map[string]string),
	}
}

// ValidateYAML 验证YAML内容
func (v *YAMLValidator) ValidateYAML(content string, schemaType string, schemaName string) (bool, []string, error) {
	// 解析YAML
	var yamlData interface{}
	if err := yaml.Unmarshal([]byte(content), &yamlData); err != nil {
		return false, []string{fmt.Sprintf("YAML解析错误: %s", err.Error())}, err
	}

	// 获取Schema
	schema, err := v.getSchema(schemaType, schemaName)
	if err != nil {
		return false, []string{fmt.Sprintf("获取Schema失败: %s", err.Error())}, err
	}

	// 验证YAML
	isValid, errors := v.validateAgainstSchema(yamlData, schema)
	return isValid, errors, nil
}

// getSchema 获取Schema
func (v *YAMLValidator) getSchema(schemaType string, schemaName string) (string, error) {
	// 构建缓存键
	cacheKey := fmt.Sprintf("%s:%s", schemaType, schemaName)

	// 检查缓存
	v.schemaCacheLock.RLock()
	cachedSchema, exists := v.schemaCache[cacheKey]
	v.schemaCacheLock.RUnlock()

	if exists {
		return cachedSchema, nil
	}

	var schema string
	var err error

	// 根据类型获取Schema
	switch schemaType {
	case "kubernetes":
		schema, err = v.getKubernetesSchema(schemaName)
	case "custom":
		schema, err = v.getCustomSchema(schemaName)
	default:
		return "", errors.New("不支持的Schema类型: " + schemaType)
	}

	if err != nil {
		return "", err
	}

	// 缓存Schema
	v.schemaCacheLock.Lock()
	v.schemaCache[cacheKey] = schema
	v.schemaCacheLock.Unlock()

	return schema, nil
}

// getKubernetesSchema 获取Kubernetes Schema
func (v *YAMLValidator) getKubernetesSchema(resourceType string) (string, error) {
	// 如果未指定资源类型，使用通用Schema
	if resourceType == "" {
		resourceType = "general"
	}

	// 尝试从数据库获取Schema
	var yamlSchema model.YAMLSchema
	if err := global.DB.Where("type = ? AND name = ?", "kubernetes", resourceType).
		Order("created_at DESC").
		First(&yamlSchema).Error; err == nil {
		return yamlSchema.Schema, nil
	}

	// 如果数据库中没有，尝试从Kubernetes API获取
	schema, err := v.fetchKubernetesSchema(resourceType)
	if err != nil {
		// 如果无法从API获取，尝试从本地文件获取
		return v.loadLocalSchema("kubernetes", resourceType)
	}

	return schema, nil
}

// getCustomSchema 获取自定义Schema
func (v *YAMLValidator) getCustomSchema(schemaName string) (string, error) {
	// 从数据库获取Schema
	var yamlSchema model.YAMLSchema
	if err := global.DB.Where("type = ? AND name = ?", "custom", schemaName).
		Order("created_at DESC").
		First(&yamlSchema).Error; err != nil {
		return "", errors.New("未找到自定义Schema: " + schemaName)
	}

	return yamlSchema.Schema, nil
}

// fetchKubernetesSchema 从Kubernetes API获取Schema
func (v *YAMLValidator) fetchKubernetesSchema(resourceType string) (string, error) {
	// 这里应该实现从Kubernetes API获取Schema的逻辑
	// 例如，可以使用kubectl或者直接调用Kubernetes API
	// 为了简化，这里返回错误，让系统使用本地Schema
	return "", errors.New("未实现从Kubernetes API获取Schema")
}

// loadLocalSchema 从本地文件加载Schema
func (v *YAMLValidator) loadLocalSchema(schemaType string, schemaName string) (string, error) {
	// 构建Schema文件路径
	schemaPath := filepath.Join("resources", "schemas", schemaType, schemaName+".json")

	// 读取Schema文件
	schemaBytes, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return "", errors.New("未找到Schema文件: " + schemaPath)
	}

	return string(schemaBytes), nil
}

// validateAgainstSchema 根据Schema验证YAML
func (v *YAMLValidator) validateAgainstSchema(yamlData interface{}, schema string) (bool, []string) {
	// 这里应该实现根据JSON Schema验证YAML的逻辑
	// 为了简化，这里返回一个模拟的验证结果
	// 在实际项目中，可以使用第三方库如gojsonschema

	// 模拟验证
	if yamlData == nil {
		return false, []string{"YAML内容为空"}
	}

	// 检查YAML中是否有必要的字段
	yamlMap, ok := yamlData.(map[string]interface{})
	if !ok {
		return false, []string{"YAML格式错误，应为对象"}
	}

	// 模拟一些基本验证
	var errors []string

	// 检查apiVersion字段
	if _, exists := yamlMap["apiVersion"]; !exists {
		errors = append(errors, "缺少必要字段: apiVersion")
	}

	// 检查kind字段
	if _, exists := yamlMap["kind"]; !exists {
		errors = append(errors, "缺少必要字段: kind")
	}

	// 检查metadata字段
	metadata, metadataExists := yamlMap["metadata"]
	if !metadataExists {
		errors = append(errors, "缺少必要字段: metadata")
	} else {
		// 检查metadata.name字段
		metadataMap, ok := metadata.(map[string]interface{})
		if !ok {
			errors = append(errors, "metadata字段格式错误")
		} else if _, exists := metadataMap["name"]; !exists {
			errors = append(errors, "缺少必要字段: metadata.name")
		}
	}

	return len(errors) == 0, errors
}

// SaveValidationResult 保存验证结果
func (v *YAMLValidator) SaveValidationResult(validation *model.YAMLValidation) error {
	return global.DB.Create(validation).Error
}

// GetValidationHistory 获取验证历史
func (v *YAMLValidator) GetValidationHistory(userID uint, limit int) ([]model.YAMLValidation, error) {
	var validations []model.YAMLValidation
	query := global.DB.Where("creator_id = ?", userID)

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Order("created_at DESC").Find(&validations).Error; err != nil {
		return nil, err
	}

	return validations, nil
}

// CreateSchema 创建Schema
func (v *YAMLValidator) CreateSchema(schema *model.YAMLSchema) error {
	// 验证Schema是否有效
	var js json.RawMessage
	if err := json.Unmarshal([]byte(schema.Schema), &js); err != nil {
		return errors.New("无效的JSON Schema: " + err.Error())
	}

	return global.DB.Create(schema).Error
}

// UpdateSchema 更新Schema
func (v *YAMLValidator) UpdateSchema(id uint, updates map[string]interface{}) error {
	// 如果更新包含Schema内容，验证是否有效
	if schemaStr, ok := updates["schema"].(string); ok {
		var js json.RawMessage
		if err := json.Unmarshal([]byte(schemaStr), &js); err != nil {
			return errors.New("无效的JSON Schema: " + err.Error())
		}
	}

	return global.DB.Model(&model.YAMLSchema{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteSchema 删除Schema
func (v *YAMLValidator) DeleteSchema(id uint) error {
	return global.DB.Delete(&model.YAMLSchema{}, id).Error
}

// GetSchemaByID 根据ID获取Schema
func (v *YAMLValidator) GetSchemaByID(id uint) (*model.YAMLSchema, error) {
	var schema model.YAMLSchema
	if err := global.DB.First(&schema, id).Error; err != nil {
		return nil, err
	}
	return &schema, nil
}

// GetSchemas 获取Schema列表
func (v *YAMLValidator) GetSchemas(schemaType string) ([]model.YAMLSchema, error) {
	var schemas []model.YAMLSchema
	query := global.DB.Model(&model.YAMLSchema{})

	if schemaType != "" {
		query = query.Where("type = ?", schemaType)
	}

	if err := query.Order("name ASC").Find(&schemas).Error; err != nil {
		return nil, err
	}

	return schemas, nil
}

// ClearSchemaCache 清除Schema缓存
func (v *YAMLValidator) ClearSchemaCache() {
	v.schemaCacheLock.Lock()
	defer v.schemaCacheLock.Unlock()
	v.schemaCache = make(map[string]string)
}
