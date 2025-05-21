package v1

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"gin_pipeline/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
)

var yamlValidator = service.NewYAMLValidator()

// ValidateYAML 验证YAML
// @Summary 验证YAML
// @Description 验证YAML内容是否符合指定的Schema
// @Tags YAML验证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.ValidateYAML true "YAML信息"
// @Success 200 {object} response.Response{data=model.YAMLValidation} "验证成功"
// @Router /yaml/validate [post]
func ValidateYAML(c *gin.Context) {
	var req request.ValidateYAML
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("验证YAML失败", c)
		return
	}

	// 验证YAML
	isValid, errors, err := yamlValidator.ValidateYAML(req.Content, req.SchemaType, req.SchemaName)
	if err != nil {
		global.Log.Error("验证YAML失败", zap.Error(err))
		response.FailWithMessage("验证YAML失败: "+err.Error(), c)
		return
	}

	// 构建验证结果
	errorsStr := ""
	if len(errors) > 0 {
		errorsStr = "- " + errors[0]
		for i := 1; i < len(errors); i++ {
			errorsStr += "\n- " + errors[i]
		}
	}

	validation := model.YAMLValidation{
		Name:       req.Name,
		Content:    req.Content,
		SchemaType: req.SchemaType,
		SchemaName: req.SchemaName,
		IsValid:    isValid,
		Errors:     errorsStr,
		CreatorID:  userID,
	}

	// 保存验证结果
	if err := yamlValidator.SaveValidationResult(&validation); err != nil {
		global.Log.Error("保存验证结果失败", zap.Error(err))
		response.FailWithMessage("验证YAML成功，但保存结果失败", c)
		return
	}

	response.OkWithData(validation, c)
}

// GetValidationHistory 获取验证历史
// @Summary 获取验证历史
// @Description 获取当前用户的YAML验证历史
// @Tags YAML验证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "限制数量" default(10)
// @Success 200 {object} response.Response{data=[]model.YAMLValidation} "获取成功"
// @Router /yaml/history [get]
func GetValidationHistory(c *gin.Context) {
	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("获取验证历史失败", c)
		return
	}

	// 获取limit参数
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	// 获取验证历史
	validations, err := yamlValidator.GetValidationHistory(userID, limit)
	if err != nil {
		global.Log.Error("获取验证历史失败", zap.Error(err))
		response.FailWithMessage("获取验证历史失败", c)
		return
	}

	response.OkWithData(validations, c)
}

// CreateYAMLSchema 创建YAML Schema
// @Summary 创建YAML Schema
// @Description 创建新的YAML Schema
// @Tags YAML验证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateYAMLSchema true "Schema信息"
// @Success 200 {object} response.Response{data=model.YAMLSchema} "创建成功"
// @Router /yaml/schema [post]
func CreateYAMLSchema(c *gin.Context) {
	var req request.CreateYAMLSchema
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建Schema失败", c)
		return
	}

	// 创建Schema
	schema := model.YAMLSchema{
		Name:        req.Name,
		Type:        req.Type,
		Version:     req.Version,
		Schema:      req.Schema,
		Description: req.Description,
		CreatorID:   userID,
	}

	if err := yamlValidator.CreateSchema(&schema); err != nil {
		global.Log.Error("创建Schema失败", zap.Error(err))
		response.FailWithMessage("创建Schema失败: "+err.Error(), c)
		return
	}

	// 清除缓存
	yamlValidator.ClearSchemaCache()

	response.OkWithData(schema, c)
}

// UpdateYAMLSchema 更新YAML Schema
// @Summary 更新YAML Schema
// @Description 更新YAML Schema信息
// @Tags YAML验证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Schema ID"
// @Param data body request.UpdateYAMLSchema true "Schema信息"
// @Success 200 {object} response.Response{data=model.YAMLSchema} "更新成功"
// @Router /yaml/schema/{id} [put]
func UpdateYAMLSchema(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	var req request.UpdateYAMLSchema
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 更新Schema
	updates := map[string]interface{}{
		"name":        req.Name,
		"version":     req.Version,
		"schema":      req.Schema,
		"description": req.Description,
	}

	if err := yamlValidator.UpdateSchema(uint(id), updates); err != nil {
		global.Log.Error("更新Schema失败", zap.Error(err))
		response.FailWithMessage("更新Schema失败: "+err.Error(), c)
		return
	}

	// 清除缓存
	yamlValidator.ClearSchemaCache()

	// 获取更新后的Schema
	schema, err := yamlValidator.GetSchemaByID(uint(id))
	if err != nil {
		global.Log.Error("获取更新后的Schema失败", zap.Error(err))
		response.FailWithMessage("更新Schema成功，但获取详情失败", c)
		return
	}

	response.OkWithData(schema, c)
}

// DeleteYAMLSchema 删除YAML Schema
// @Summary 删除YAML Schema
// @Description 删除指定的YAML Schema
// @Tags YAML验证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Schema ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /yaml/schema/{id} [delete]
func DeleteYAMLSchema(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	if err := yamlValidator.DeleteSchema(uint(id)); err != nil {
		global.Log.Error("删除Schema失败", zap.Error(err))
		response.FailWithMessage("删除Schema失败", c)
		return
	}

	// 清除缓存
	yamlValidator.ClearSchemaCache()

	response.OkWithMessage("删除Schema成功", c)
}

// GetYAMLSchemas 获取YAML Schema列表
// @Summary 获取YAML Schema列表
// @Description 获取YAML Schema列表
// @Tags YAML验证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param type query string false "Schema类型" default("")
// @Success 200 {object} response.Response{data=[]model.YAMLSchema} "获取成功"
// @Router /yaml/schema [get]
func GetYAMLSchemas(c *gin.Context) {
	// 获取type参数
	schemaType := c.DefaultQuery("type", "")

	// 获取Schema列表
	schemas, err := yamlValidator.GetSchemas(schemaType)
	if err != nil {
		global.Log.Error("获取Schema列表失败", zap.Error(err))
		response.FailWithMessage("获取Schema列表失败", c)
		return
	}

	response.OkWithData(schemas, c)
}

// GetYAMLSchemaByID 获取YAML Schema详情
// @Summary 获取YAML Schema详情
// @Description 根据ID获取YAML Schema详情
// @Tags YAML验证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Schema ID"
// @Success 200 {object} response.Response{data=model.YAMLSchema} "获取成功"
// @Router /yaml/schema/{id} [get]
func GetYAMLSchemaByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	schema, err := yamlValidator.GetSchemaByID(uint(id))
	if err != nil {
		global.Log.Error("获取Schema失败", zap.Error(err))
		response.FailWithMessage("获取Schema失败", c)
		return
	}

	response.OkWithData(schema, c)
}
