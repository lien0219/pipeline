package v1

import (
	"encoding/json"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 环境变量结构
type EnvironmentVariable struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreateEnvironment 创建环境
// @Summary 创建环境
// @Description 创建新的部署环境
// @Tags 环境管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateEnvironment true "环境信息"
// @Success 200 {object} response.Response{data=model.Environment} "创建成功"
// @Router /environment [post]
func CreateEnvironment(c *gin.Context) {
	var req struct {
		Name        string                `json:"name" binding:"required"`
		Type        string                `json:"type" binding:"required"`
		URL         string                `json:"url"`
		Description string                `json:"description"`
		Variables   []EnvironmentVariable `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建环境失败", c)
		return
	}

	// 序列化环境变量
	variablesJSON, err := json.Marshal(req.Variables)
	if err != nil {
		global.Log.Error("序列化环境变量失败", zap.Error(err))
		response.FailWithMessage("创建环境失败", c)
		return
	}

	// 创建环境
	environment := model.Environment{
		Name:        req.Name,
		Type:        req.Type,
		URL:         req.URL,
		Description: req.Description,
		Variables:   string(variablesJSON),
		Status:      "active",
		CreatedBy:   userID,
	}

	if err := global.DB.Create(&environment).Error; err != nil {
		global.Log.Error("创建环境失败", zap.Error(err))
		response.FailWithMessage("创建环境失败", c)
		return
	}

	// 查询完整的环境信息
	if err := global.DB.Preload("User").First(&environment, environment.ID).Error; err != nil {
		global.Log.Error("查询环境信息失败", zap.Error(err))
		response.FailWithMessage("创建环境成功，但获取详情失败", c)
		return
	}

	response.OkWithData(environment, c)
}

// GetEnvironments 获取环境列表
// @Summary 获取环境列表
// @Description 获取环境列表，支持分页和筛选
// @Tags 环境管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param name query string false "环境名称"
// @Param type query string false "环境类型"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Environment}} "获取成功"
// @Router /environment [get]
func GetEnvironments(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 获取分页参数
	page := pageInfo.GetPage()
	pageSize := pageInfo.GetPageSize()

	// 获取筛选参数
	name := c.Query("name")
	envType := c.Query("type")

	// 构建查询条件
	db := global.DB.Model(&model.Environment{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if envType != "" {
		db = db.Where("type = ?", envType)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.Log.Error("查询环境总数失败", zap.Error(err))
		response.FailWithMessage("获取环境列表失败", c)
		return
	}

	// 查询列表
	var environments []model.Environment
	if err := db.Preload("User").Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}).Order("id DESC").Find(&environments).Error; err != nil {
		global.Log.Error("查询环境列表失败", zap.Error(err))
		response.FailWithMessage("获取环境列表失败", c)
		return
	}

	// 返回结果
	result := response.PageResult{
		List:     environments,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.OkWithData(result, c)
}

// GetEnvironmentByID 获取环境详情
// @Summary 获取环境详情
// @Description 根据ID获取环境详情
// @Tags 环境管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "环境ID"
// @Success 200 {object} response.Response{data=model.Environment} "获取成功"
// @Router /environment/{id} [get]
func GetEnvironmentByID(c *gin.Context) {
	id := c.Param("id")

	var environment model.Environment
	if err := global.DB.Preload("User").First(&environment, id).Error; err != nil {
		global.Log.Error("查询环境失败", zap.Error(err))
		response.FailWithMessage("获取环境详情失败", c)
		return
	}

	// 解析环境变量
	var variables []EnvironmentVariable
	if environment.Variables != "" {
		if err := json.Unmarshal([]byte(environment.Variables), &variables); err != nil {
			global.Log.Warn("解析环境变量失败", zap.Error(err))
			// 不影响返回结果
		}
	}

	// 构造响应数据
	result := map[string]interface{}{
		"id":               environment.ID,
		"name":             environment.Name,
		"type":             environment.Type,
		"url":              environment.URL,
		"description":      environment.Description,
		"status":           environment.Status,
		"variables":        variables,
		"last_deployed_at": environment.LastDeployedAt,
		"created_at":       environment.CreatedAt,
		"updated_at":       environment.UpdatedAt,
		"created_by":       environment.CreatedBy,
		"user":             environment.User,
	}

	response.OkWithData(result, c)
}

// UpdateEnvironment 更新环境
// @Summary 更新环境
// @Description 更新环境信息
// @Tags 环境管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "环境ID"
// @Param data body request.UpdateEnvironment true "环境信息"
// @Success 200 {object} response.Response{data=model.Environment} "更新成功"
// @Router /environment/{id} [put]
func UpdateEnvironment(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string                `json:"name" binding:"required"`
		Type        string                `json:"type" binding:"required"`
		URL         string                `json:"url"`
		Description string                `json:"description"`
		Status      string                `json:"status"`
		Variables   []EnvironmentVariable `json:"variables"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 查询环境
	var environment model.Environment
	if err := global.DB.First(&environment, id).Error; err != nil {
		global.Log.Error("查询环境失败", zap.Error(err))
		response.FailWithMessage("更新环境失败", c)
		return
	}

	// 序列化环境变量
	variablesJSON, err := json.Marshal(req.Variables)
	if err != nil {
		global.Log.Error("序列化环境变量失败", zap.Error(err))
		response.FailWithMessage("更新环境失败", c)
		return
	}

	// 更新环境
	updates := map[string]interface{}{
		"name":        req.Name,
		"type":        req.Type,
		"url":         req.URL,
		"description": req.Description,
		"variables":   string(variablesJSON),
	}

	// 如果状态不为空，则更新状态
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := global.DB.Model(&environment).Updates(updates).Error; err != nil {
		global.Log.Error("更新环境失败", zap.Error(err))
		response.FailWithMessage("更新环境失败", c)
		return
	}

	// 查询更新后的环境
	if err := global.DB.Preload("User").First(&environment, id).Error; err != nil {
		global.Log.Error("查询环境失败", zap.Error(err))
		response.FailWithMessage("更新环境成功，但获取详情失败", c)
		return
	}

	// 解析环境变量
	var variables []EnvironmentVariable
	if environment.Variables != "" {
		if err := json.Unmarshal([]byte(environment.Variables), &variables); err != nil {
			global.Log.Warn("解析环境变量失败", zap.Error(err))
			// 不影响返回结果
		}
	}

	// 构造响应数据
	result := map[string]interface{}{
		"id":               environment.ID,
		"name":             environment.Name,
		"type":             environment.Type,
		"url":              environment.URL,
		"description":      environment.Description,
		"status":           environment.Status,
		"variables":        variables,
		"last_deployed_at": environment.LastDeployedAt,
		"created_at":       environment.CreatedAt,
		"updated_at":       environment.UpdatedAt,
		"created_by":       environment.CreatedBy,
		"user":             environment.User,
	}

	response.OkWithData(result, c)
}

// DeleteEnvironment 删除环境
// @Summary 删除环境
// @Description 根据ID删除环境
// @Tags 环境管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "环境ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /environment/{id} [delete]
func DeleteEnvironment(c *gin.Context) {
	id := c.Param("id")

	// 删除环境
	if err := global.DB.Delete(&model.Environment{}, id).Error; err != nil {
		global.Log.Error("删除环境失败", zap.Error(err))
		response.FailWithMessage("删除环境失败", c)
		return
	}

	response.OkWithMessage("删除环境成功", c)
}
