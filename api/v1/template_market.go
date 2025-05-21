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

var templateMarketService = new(service.TemplateMarketService)

// CreateTemplateCategory 创建模板分类
// @Summary 创建模板分类
// @Description 创建新的模板分类
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateTemplateCategory true "分类信息"
// @Success 200 {object} response.Response{data=model.TemplateCategory} "创建成功"
// @Router /template-market/category [post]
func CreateTemplateCategory(c *gin.Context) {
	var req request.CreateTemplateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建分类失败", c)
		return
	}

	// 创建分类
	category := model.TemplateCategory{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		Order:       req.Order,
		CreatorID:   userID,
	}

	if err := templateMarketService.CreateCategory(&category); err != nil {
		global.Log.Error("创建模板分类失败", zap.Error(err))
		response.FailWithMessage("创建分类失败: "+err.Error(), c)
		return
	}

	response.OkWithData(category, c)
}

// GetTemplateCategories 获取模板分类列表
// @Summary 获取模板分类列表
// @Description 获取所有模板分类
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.TemplateCategory} "获取成功"
// @Router /template-market/category [get]
func GetTemplateCategories(c *gin.Context) {
	categories, err := templateMarketService.GetCategories()
	if err != nil {
		global.Log.Error("获取模板分类列表失败", zap.Error(err))
		response.FailWithMessage("获取分类列表失败", c)
		return
	}

	response.OkWithData(categories, c)
}

// UpdateTemplateCategory 更新模板分类
// @Summary 更新模板分类
// @Description 更新模板分类信息
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Param data body request.UpdateTemplateCategory true "分类信息"
// @Success 200 {object} response.Response{data=model.TemplateCategory} "更新成功"
// @Router /template-market/category/{id} [put]
func UpdateTemplateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	var req request.UpdateTemplateCategory
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 更新分类
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"icon":        req.Icon,
		"order":       req.Order,
	}

	if err := templateMarketService.UpdateCategory(uint(id), updates); err != nil {
		global.Log.Error("更新模板分类失败", zap.Error(err))
		response.FailWithMessage("更新分类失败: "+err.Error(), c)
		return
	}

	// 获取更新后的分类
	category, err := templateMarketService.GetCategoryByID(uint(id))
	if err != nil {
		global.Log.Error("获取更新后的分类失败", zap.Error(err))
		response.FailWithMessage("更新分类成功，但获取详情失败", c)
		return
	}

	response.OkWithData(category, c)
}

// DeleteTemplateCategory 删除模板分类
// @Summary 删除模板分类
// @Description 删除指定的模板分类
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "分类ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /template-market/category/{id} [delete]
func DeleteTemplateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	if err := templateMarketService.DeleteCategory(uint(id)); err != nil {
		global.Log.Error("删除模板分类失败", zap.Error(err))
		response.FailWithMessage("删除分类失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除分类成功", c)
}

// CreateTemplate 创建模板
// @Summary 创建模板
// @Description 创建新的模板
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateTemplate true "模板信息"
// @Success 200 {object} response.Response{data=model.Template} "创建成功"
// @Router /template-market/template [post]
func CreateTemplate(c *gin.Context) {
	var req request.CreateTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建模板失败", c)
		return
	}

	// 创建模板
	template := model.Template{
		Name:        req.Name,
		Description: req.Description,
		CategoryID:  req.CategoryID,
		Icon:        req.Icon,
		Tags:        req.Tags,
		IsPublic:    req.IsPublic,
		CreatorID:   userID,
	}

	// 创建初始版本
	version := model.TemplateVersion{
		Version:   req.Version,
		Content:   req.Content,
		Changelog: req.Changelog,
		CreatorID: userID,
	}

	if err := templateMarketService.CreateTemplate(&template, &version); err != nil {
		global.Log.Error("创建模板失败", zap.Error(err))
		response.FailWithMessage("创建模板失败: "+err.Error(), c)
		return
	}

	// 获取完整的模板信息
	result, err := templateMarketService.GetTemplateByID(template.ID)
	if err != nil {
		global.Log.Error("获取模板详情失败", zap.Error(err))
		response.FailWithMessage("创建模板成功，但获取详情失败", c)
		return
	}

	response.OkWithData(result, c)
}

// GetTemplates 获取模板列表
// @Summary 获取模板列表
// @Description 获取模板列表，支持按分类筛选
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category_id query int false "分类ID"
// @Param public query bool false "是否公开"
// @Success 200 {object} response.Response{data=[]model.Template} "获取成功"
// @Router /template-market/template [get]
func GetTemplates(c *gin.Context) {
	// 获取查询参数
	categoryIDStr := c.Query("category_id")
	publicStr := c.Query("public")

	var categoryID uint
	if categoryIDStr != "" {
		id, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			categoryID = uint(id)
		}
	}

	var isPublic *bool
	if publicStr != "" {
		public := publicStr == "true"
		isPublic = &public
	}

	templates, err := templateMarketService.GetTemplates(categoryID, isPublic)
	if err != nil {
		global.Log.Error("获取模板列表失败", zap.Error(err))
		response.FailWithMessage("获取模板列表失败", c)
		return
	}

	response.OkWithData(templates, c)
}

// GetTemplateByID 获取模板详情
// @Summary 获取模板详情
// @Description 根据ID获取模板详情
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} response.Response{data=model.Template} "获取成功"
// @Router /template-market/template/{id} [get]
func GetTemplateByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	template, err := templateMarketService.GetTemplateByID(uint(id))
	if err != nil {
		global.Log.Error("获取模板详情失败", zap.Error(err))
		response.FailWithMessage("获取模板详情失败", c)
		return
	}

	response.OkWithData(template, c)
}

// UpdateTemplate 更新模板
// @Summary 更新模板
// @Description 更新模板信息
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Param data body request.UpdateTemplate true "模板信息"
// @Success 200 {object} response.Response{data=model.Template} "更新成功"
// @Router /template-market/template/{id} [put]
func UpdateTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	var req request.UpdateTemplate
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 更新模板
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"category_id": req.CategoryID,
		"icon":        req.Icon,
		"tags":        req.Tags,
		"is_public":   req.IsPublic,
	}

	if err := templateMarketService.UpdateTemplate(uint(id), updates); err != nil {
		global.Log.Error("更新模板失败", zap.Error(err))
		response.FailWithMessage("更新模板失败: "+err.Error(), c)
		return
	}

	// 获取更新后的模板
	template, err := templateMarketService.GetTemplateByID(uint(id))
	if err != nil {
		global.Log.Error("获取更新后的模板失败", zap.Error(err))
		response.FailWithMessage("更新模板成功，但获取详情失败", c)
		return
	}

	response.OkWithData(template, c)
}

// DeleteTemplate 删除模板
// @Summary 删除模板
// @Description 删除指定的模板
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /template-market/template/{id} [delete]
func DeleteTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	if err := templateMarketService.DeleteTemplate(uint(id)); err != nil {
		global.Log.Error("删除模板失败", zap.Error(err))
		response.FailWithMessage("删除模板失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除模板成功", c)
}

// CreateTemplateVersion 创建模板版本
// @Summary 创建模板版本
// @Description 为指定模板创建新版本
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Param data body request.CreateTemplateVersion true "版本信息"
// @Success 200 {object} response.Response{data=model.TemplateVersion} "创建成功"
// @Router /template-market/template/{id}/version [post]
func CreateTemplateVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	var req request.CreateTemplateVersion
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建版本失败", c)
		return
	}

	// 创建版本
	version := model.TemplateVersion{
		Version:   req.Version,
		Content:   req.Content,
		Changelog: req.Changelog,
		CreatorID: userID,
	}

	if err := templateMarketService.CreateTemplateVersion(uint(id), &version); err != nil {
		global.Log.Error("创建模板版本失败", zap.Error(err))
		response.FailWithMessage("创建版本失败: "+err.Error(), c)
		return
	}

	response.OkWithData(version, c)
}

// GetTemplateVersions 获取模板版本列表
// @Summary 获取模板版本列表
// @Description 获取指定模板的所有版本
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} response.Response{data=[]model.TemplateVersion} "获取成功"
// @Router /template-market/template/{id}/version [get]
func GetTemplateVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	versions, err := templateMarketService.GetTemplateVersions(uint(id))
	if err != nil {
		global.Log.Error("获取模板版本列表失败", zap.Error(err))
		response.FailWithMessage("获取版本列表失败", c)
		return
	}

	response.OkWithData(versions, c)
}

// GetTemplateVersionByID 获取模板版本详情
// @Summary 获取模板版本详情
// @Description 根据ID获取模板版本详情
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Param versionId path int true "版本ID"
// @Success 200 {object} response.Response{data=model.TemplateVersion} "获取成功"
// @Router /template-market/template/{id}/version/{versionId} [get]
func GetTemplateVersionByID(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("versionId"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的版本ID", c)
		return
	}

	version, err := templateMarketService.GetTemplateVersionByID(uint(versionID))
	if err != nil {
		global.Log.Error("获取模板版本详情失败", zap.Error(err))
		response.FailWithMessage("获取版本详情失败", c)
		return
	}

	response.OkWithData(version, c)
}

// DeleteTemplateVersion 删除模板版本
// @Summary 删除模板版本
// @Description 删除指定的模板版本
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Param versionId path int true "版本ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /template-market/template/{id}/version/{versionId} [delete]
func DeleteTemplateVersion(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("versionId"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的版本ID", c)
		return
	}

	if err := templateMarketService.DeleteTemplateVersion(uint(versionID)); err != nil {
		global.Log.Error("删除模板版本失败", zap.Error(err))
		response.FailWithMessage("删除版本失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("删除版本成功", c)
}

// SearchTemplates 搜索模板
// @Summary 搜索模板
// @Description 根据关键字、分类和标签搜索模板
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param keyword query string false "关键字"
// @Param category_id query int false "分类ID"
// @Param tags query string false "标签（逗号分隔）"
// @Success 200 {object} response.Response{data=[]model.Template} "搜索成功"
// @Router /template-market/search [get]
func SearchTemplates(c *gin.Context) {
	keyword := c.Query("keyword")
	categoryIDStr := c.Query("category_id")
	tags := c.Query("tags")

	var categoryID uint
	if categoryIDStr != "" {
		id, err := strconv.ParseUint(categoryIDStr, 10, 32)
		if err == nil {
			categoryID = uint(id)
		}
	}

	templates, err := templateMarketService.SearchTemplates(keyword, categoryID, tags)
	if err != nil {
		global.Log.Error("搜索模板失败", zap.Error(err))
		response.FailWithMessage("搜索模板失败", c)
		return
	}

	response.OkWithData(templates, c)
}

// DownloadTemplate 下载模板
// @Summary 下载模板
// @Description 下载指定模板的最新版本
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Success 200 {object} response.Response{data=model.TemplateVersion} "下载成功"
// @Router /template-market/template/{id}/download [get]
func DownloadTemplate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	// 获取最新版本
	version, err := templateMarketService.GetLatestTemplateVersion(uint(id))
	if err != nil {
		global.Log.Error("获取模板最新版本失败", zap.Error(err))
		response.FailWithMessage("获取模板最新版本失败", c)
		return
	}

	// 增加下载次数
	if err := templateMarketService.IncrementTemplateDownloadCount(uint(id)); err != nil {
		global.Log.Error("增加模板下载次数失败", zap.Error(err))
		// 不影响结果，继续执行
	}

	if err := templateMarketService.IncrementVersionDownloadCount(version.ID); err != nil {
		global.Log.Error("增加版本下载次数失败", zap.Error(err))
		// 不影响结果，继续执行
	}

	response.OkWithData(version, c)
}

// SetVersionAsLatest 设置版本为最新版本
// @Summary 设置版本为最新版本
// @Description 将指定版本设置为模板的最新版本
// @Tags 模板市场
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "模板ID"
// @Param versionId path int true "版本ID"
// @Success 200 {object} response.Response "设置成功"
// @Router /template-market/template/{id}/version/{versionId}/latest [post]
func SetVersionAsLatest(c *gin.Context) {
	versionID, err := strconv.ParseUint(c.Param("versionId"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的版本ID", c)
		return
	}

	if err := templateMarketService.SetVersionAsLatest(uint(versionID)); err != nil {
		global.Log.Error("设置最新版本失败", zap.Error(err))
		response.FailWithMessage("设置最新版本失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("设置最新版本成功", c)
}
