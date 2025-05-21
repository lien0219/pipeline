package v1

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CreateBuildTemplate 创建构建模板
// @Summary 创建构建模板
// @Description 创建新的构建模板
// @Tags 构建模板管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateBuildTemplate true "模板信息"
// @Success 200 {object} response.Response{data=model.BuildTemplate} "创建成功"
// @Router /build-template [post]
func CreateBuildTemplate(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Description string `json:"description"`
		Content     string `json:"content" binding:"required"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建构建模板失败", c)
		return
	}

	// 创建构建模板
	template := model.BuildTemplate{
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Content:     req.Content,
		IsPublic:    req.IsPublic,
		CreatedBy:   userID,
	}

	if err := global.DB.Create(&template).Error; err != nil {
		global.Log.Error("创建构建模板失败", zap.Error(err))
		response.FailWithMessage("创建构建模板失败", c)
		return
	}

	// 查询完整的模板信息
	if err := global.DB.Preload("Creator").First(&template, template.ID).Error; err != nil {
		global.Log.Error("查询构建模板信息失败", zap.Error(err))
		response.FailWithMessage("创建构建模板成功，但获取详情失败", c)
		return
	}

	response.OkWithData(template, c)
}

// GetBuildTemplates 获取构建模板列表
// @Summary 获取构建模板列表
// @Description 获取构建模板列表，支持分页和筛选
// @Tags 构建模板管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param name query string false "模板名称"
// @Param type query string false "模板类型"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.BuildTemplate}} "获取成功"
// @Router /build-template [get]
func GetBuildTemplates(c *gin.Context) {
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
	templateType := c.Query("type")

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("获取构建模板列表失败", c)
		return
	}

	// 构建查询条件
	db := global.DB.Model(&model.BuildTemplate{}).Where("is_public = ? OR created_by = ?", true, userID)

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if templateType != "" {
		db = db.Where("type = ?", templateType)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.Log.Error("查询构建模板总数失败", zap.Error(err))
		response.FailWithMessage("获取构建模板列表失败", c)
		return
	}

	// 查询列表
	var templates []model.BuildTemplate
	if err := db.Preload("Creator").Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}).Order("id DESC").Find(&templates).Error; err != nil {
		global.Log.Error("查询构建模板列表失败", zap.Error(err))
		response.FailWithMessage("获取构建模板列表失败", c)
		return
	}

	// 返回结果
	result := response.PageResult{
		List:     templates,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.OkWithData(result, c)
}

// GetBuildTemplateByID 获取构建模板详情
// @Summary 获取构建模板详情
// @Description 根据ID获取构建模板详情
// @Tags 构建模板管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} response.Response{data=model.BuildTemplate} "获取成功"
// @Router /build-template/{id} [get]
func GetBuildTemplateByID(c *gin.Context) {
	id := c.Param("id")

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("获取构建模板详情失败", c)
		return
	}

	var template model.BuildTemplate
	if err := global.DB.Preload("Creator").First(&template, id).Error; err != nil {
		global.Log.Error("查询构建模板失败", zap.Error(err))
		response.FailWithMessage("获取构建模板详情失败", c)
		return
	}

	// 检查权限
	if !template.IsPublic && template.CreatedBy != userID {
		response.FailWithMessage("无权访问此模板", c)
		return
	}

	response.OkWithData(template, c)
}

// UpdateBuildTemplate 更新构建模板
// @Summary 更新构建模板
// @Description 更新构建模板信息
// @Tags 构建模板管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Param data body request.UpdateBuildTemplate true "模板信息"
// @Success 200 {object} response.Response{data=model.BuildTemplate} "更新成功"
// @Router /build-template/{id} [put]
func UpdateBuildTemplate(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Name        string `json:"name" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Description string `json:"description"`
		Content     string `json:"content" binding:"required"`
		IsPublic    bool   `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("更新构建模板失败", c)
		return
	}

	// 查询模板
	var template model.BuildTemplate
	if err := global.DB.First(&template, id).Error; err != nil {
		global.Log.Error("查询构建模板失败", zap.Error(err))
		response.FailWithMessage("更新构建模板失败", c)
		return
	}

	// 检查权限
	if template.CreatedBy != userID {
		response.FailWithMessage("无权更新此模板", c)
		return
	}

	// 更新模板
	updates := map[string]interface{}{
		"name":        req.Name,
		"type":        req.Type,
		"description": req.Description,
		"content":     req.Content,
		"is_public":   req.IsPublic,
	}

	if err := global.DB.Model(&template).Updates(updates).Error; err != nil {
		global.Log.Error("更新构建模板失败", zap.Error(err))
		response.FailWithMessage("更新构建模板失败", c)
		return
	}

	// 查询更新后的模板
	if err := global.DB.Preload("Creator").First(&template, id).Error; err != nil {
		global.Log.Error("查询构建模板失败", zap.Error(err))
		response.FailWithMessage("更新构建模板成功，但获取详情失败", c)
		return
	}

	response.OkWithData(template, c)
}

// DeleteBuildTemplate 删除构建模板
// @Summary 删除构建模板
// @Description 根据ID删除构建模板
// @Tags 构建模板管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /build-template/{id} [delete]
func DeleteBuildTemplate(c *gin.Context) {
	id := c.Param("id")

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("删除构建模板失败", c)
		return
	}

	// 查询模板
	var template model.BuildTemplate
	if err := global.DB.First(&template, id).Error; err != nil {
		global.Log.Error("查询构建模板失败", zap.Error(err))
		response.FailWithMessage("删除构建模板失败", c)
		return
	}

	// 检查权限
	if template.CreatedBy != userID {
		response.FailWithMessage("无权删除此模板", c)
		return
	}

	// 删除模板
	if err := global.DB.Delete(&template).Error; err != nil {
		global.Log.Error("删除构建模板失败", zap.Error(err))
		response.FailWithMessage("删除构建模板失败", c)
		return
	}

	response.OkWithMessage("删除构建模板成功", c)
}

// ApplyBuildTemplate 应用构建模板
// @Summary 应用构建模板
// @Description 将构建模板应用到流水线
// @Tags 构建模板管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Param data body request.ApplyBuildTemplate true "应用信息"
// @Success 200 {object} response.Response{data=model.Pipeline} "应用成功"
// @Router /build-template/{id}/apply [post]
func ApplyBuildTemplate(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		PipelineID uint   `json:"pipeline_id"`
		Name       string `json:"name"`
		GitRepo    string `json:"git_repo"`
		GitBranch  string `json:"git_branch"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("应用构建模板失败", c)
		return
	}

	// 查询模板
	var template model.BuildTemplate
	if err := global.DB.First(&template, id).Error; err != nil {
		global.Log.Error("查询构建模板失败", zap.Error(err))
		response.FailWithMessage("应用构建模板失败", c)
		return
	}

	// 检查权限
	if !template.IsPublic && template.CreatedBy != userID {
		response.FailWithMessage("无权使用此模板", c)
		return
	}

	// 更新使用次数
	if err := global.DB.Model(&template).Update("usage_count", gorm.Expr("usage_count + ?", 1)).Error; err != nil {
		global.Log.Warn("更新模板使用次数失败", zap.Error(err))
		// 不影响后续操作
	}

	// 开启事务
	tx := global.DB.Begin()

	// 根据请求类型执行不同操作
	var pipeline model.Pipeline
	if req.PipelineID > 0 {
		// 更新现有流水线
		if err := tx.First(&pipeline, req.PipelineID).Error; err != nil {
			tx.Rollback()
			global.Log.Error("查询流水线失败", zap.Error(err))
			response.FailWithMessage("应用构建模板失败", c)
			return
		}

		// 检查权限
		if pipeline.CreatorID != userID {
			tx.Rollback()
			response.FailWithMessage("无权更新此流水线", c)
			return
		}

		// 删除现有阶段和作业
		if err := tx.Where("stage_id IN (SELECT id FROM stages WHERE pipeline_id = ?)", req.PipelineID).Delete(&model.Job{}).Error; err != nil {
			tx.Rollback()
			global.Log.Error("删除作业失败", zap.Error(err))
			response.FailWithMessage("应用构建模板失败", c)
			return
		}

		if err := tx.Where("pipeline_id = ?", req.PipelineID).Delete(&model.Stage{}).Error; err != nil {
			tx.Rollback()
			global.Log.Error("删除阶段失败", zap.Error(err))
			response.FailWithMessage("应用构建模板失败", c)
			return
		}
	} else {
		// 创建新流水线
		if req.Name == "" || req.GitRepo == "" {
			tx.Rollback()
			response.FailWithMessage("创建流水线需要提供名称和Git仓库", c)
			return
		}

		gitBranch := req.GitBranch
		if gitBranch == "" {
			gitBranch = "main"
		}

		pipeline = model.Pipeline{
			Name:        req.Name,
			Description: "从模板 " + template.Name + " 创建",
			GitRepo:     req.GitRepo,
			GitBranch:   gitBranch,
			Status:      "inactive",
			CreatorID:   userID,
		}

		if err := tx.Create(&pipeline).Error; err != nil {
			tx.Rollback()
			global.Log.Error("创建流水线失败", zap.Error(err))
			response.FailWithMessage("应用构建模板失败", c)
			return
		}
	}

	// 解析模板内容并创建阶段和作业
	// 这里简化处理，实际应该解析YAML或JSON格式的模板内容
	// 创建一个示例阶段和作业
	stage := model.Stage{
		Name:        "从模板创建的阶段",
		Description: "使用模板 " + template.Name + " 创建的阶段",
		Order:       1,
		PipelineID:  pipeline.ID,
	}

	if err := tx.Create(&stage).Error; err != nil {
		tx.Rollback()
		global.Log.Error("创建阶段失败", zap.Error(err))
		response.FailWithMessage("应用构建模板失败", c)
		return
	}

	job := model.Job{
		Name:        "从模板创建的作业",
		Description: "使用模板 " + template.Name + " 创建的作业",
		Command:     template.Content, // 简化处理，将模板内容作为命令
		Image:       "alpine:latest",
		Timeout:     3600,
		StageID:     stage.ID,
	}

	if err := tx.Create(&job).Error; err != nil {
		tx.Rollback()
		global.Log.Error("创建作业失败", zap.Error(err))
		response.FailWithMessage("应用构建模板失败", c)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Error("提交事务失败", zap.Error(err))
		response.FailWithMessage("应用构建模板失败", c)
		return
	}

	// 查询完整的流水线信息
	if err := global.DB.Preload("Creator").Preload("Stages", func(db *gorm.DB) *gorm.DB {
		return db.Order("stages.order ASC")
	}).Preload("Stages.Jobs").First(&pipeline, pipeline.ID).Error; err != nil {
		global.Log.Error("查询流水线失败", zap.Error(err))
		response.FailWithMessage("应用构建模板成功，但获取详情失败", c)
		return
	}

	response.OkWithData(pipeline, c)
}
