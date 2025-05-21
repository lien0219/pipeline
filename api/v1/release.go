package v1

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// CreateRelease 创建发布
// @Summary 创建发布
// @Description 创建新的发布版本
// @Tags 发布管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateRelease true "发布信息"
// @Success 200 {object} response.Response{data=model.Release} "创建成功"
// @Router /release [post]
func CreateRelease(c *gin.Context) {
	var req struct {
		Version      string `json:"version" binding:"required"`
		Description  string `json:"description"`
		ReleaseNotes string `json:"release_notes"`
		Environment  string `json:"environment" binding:"required"`
		ArtifactID   uint   `json:"artifact_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建发布失败", c)
		return
	}

	// 检查制品是否存在
	var artifact model.Artifact
	if err := global.DB.First(&artifact, req.ArtifactID).Error; err != nil {
		global.Log.Error("查询制品失败", zap.Error(err))
		response.FailWithMessage("指定的制品不存在", c)
		return
	}

	// 创建发布
	release := model.Release{
		Version:      req.Version,
		Description:  req.Description,
		ReleaseNotes: req.ReleaseNotes,
		Status:       "pending",
		Environment:  req.Environment,
		ArtifactID:   req.ArtifactID,
		DeployedAt:   time.Now(),
		DeployedBy:   userID,
		IsRollback:   false,
	}

	if err := global.DB.Create(&release).Error; err != nil {
		global.Log.Error("创建发布失败", zap.Error(err))
		response.FailWithMessage("创建发布失败", c)
		return
	}

	// 查询完整的发布信息
	if err := global.DB.Preload("Artifact").Preload("User").First(&release, release.ID).Error; err != nil {
		global.Log.Error("查询发布信息失败", zap.Error(err))
		response.FailWithMessage("创建发布成功，但获取详情失败", c)
		return
	}

	// 异步执行部署（模拟）
	go deployRelease(release.ID)

	response.OkWithData(release, c)
}

// 异步执行部署（模拟）
func deployRelease(releaseID uint) {
	// 查询发布记录
	var release model.Release
	if err := global.DB.First(&release, releaseID).Error; err != nil {
		global.Log.Error("查询发布记录失败", zap.Error(err), zap.Uint("releaseID", releaseID))
		return
	}

	// 更新状态为进行中
	if err := global.DB.Model(&release).Update("status", "in_progress").Error; err != nil {
		global.Log.Error("更新发布状态失败", zap.Error(err), zap.Uint("releaseID", releaseID))
		return
	}

	// 模拟部署过程
	time.Sleep(5 * time.Second)

	// 随机成功或失败
	status := "success"
	if time.Now().Unix()%3 == 0 {
		status = "failed"
	}

	// 更新发布结果
	if err := global.DB.Model(&release).Update("status", status).Error; err != nil {
		global.Log.Error("更新发布结果失败", zap.Error(err), zap.Uint("releaseID", releaseID))
		return
	}

	// 如果部署成功，更新环境的最后部署时间
	if status == "success" {
		var environment model.Environment
		if err := global.DB.Where("type = ?", release.Environment).First(&environment).Error; err == nil {
			now := time.Now()
			if err := global.DB.Model(&environment).Update("last_deployed_at", &now).Error; err != nil {
				global.Log.Warn("更新环境最后部署时间失败", zap.Error(err))
			}
		}
	}

	global.Log.Info("发布执行完成", zap.Uint("releaseID", releaseID), zap.String("status", status))
}

// GetReleases 获取发布列表
// @Summary 获取发布列表
// @Description 获取发布列表，支持分页和筛选
// @Tags 发布管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param environment query string false "环境"
// @Param status query string false "状态"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Release}} "获取成功"
// @Router /release [get]
func GetReleases(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 获取分页参数
	page := pageInfo.GetPage()
	pageSize := pageInfo.GetPageSize()

	// 获取筛选参数
	environment := c.Query("environment")
	status := c.Query("status")

	// 构建查询条件
	db := global.DB.Model(&model.Release{})

	if environment != "" {
		db = db.Where("environment = ?", environment)
	}

	if status != "" {
		db = db.Where("status = ?", status)
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.Log.Error("查询发布总数失败", zap.Error(err))
		response.FailWithMessage("获取发布列表失败", c)
		return
	}

	// 查询列表
	var releases []model.Release
	if err := db.Preload("Artifact").Preload("User").Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}).Order("id DESC").Find(&releases).Error; err != nil {
		global.Log.Error("查询发布列表失败", zap.Error(err))
		response.FailWithMessage("获取发布列表失败", c)
		return
	}

	// 返回结果
	result := response.PageResult{
		List:     releases,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.OkWithData(result, c)
}

// GetReleaseByID 获取发布详情
// @Summary 获取发布详情
// @Description 根据ID获取发布详情
// @Tags 发布管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "发布ID"
// @Success 200 {object} response.Response{data=model.Release} "获取成功"
// @Router /release/{id} [get]
func GetReleaseByID(c *gin.Context) {
	id := c.Param("id")

	var release model.Release
	if err := global.DB.Preload("Artifact").Preload("User").First(&release, id).Error; err != nil {
		global.Log.Error("查询发布失败", zap.Error(err))
		response.FailWithMessage("获取发布详情失败", c)
		return
	}

	response.OkWithData(release, c)
}

// DeleteRelease 删除发布
// @Summary 删除发布
// @Description 根据ID删除发布
// @Tags 发布管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "发布ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /release/{id} [delete]
func DeleteRelease(c *gin.Context) {
	id := c.Param("id")

	// 删除发布
	if err := global.DB.Delete(&model.Release{}, id).Error; err != nil {
		global.Log.Error("删除发布失败", zap.Error(err))
		response.FailWithMessage("删除发布失败", c)
		return
	}

	response.OkWithMessage("删除发布成功", c)
}

// RollbackRelease 回滚发布
// @Summary 回滚发布
// @Description 回滚到指定的发布版本
// @Tags 发布管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "发布ID"
// @Success 200 {object} response.Response{data=model.Release} "回滚成功"
// @Router /release/{id}/rollback [post]
func RollbackRelease(c *gin.Context) {
	id := c.Param("id")

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("回滚发布失败", c)
		return
	}

	// 查询原发布
	var originalRelease model.Release
	if err := global.DB.Preload("Artifact").First(&originalRelease, id).Error; err != nil {
		global.Log.Error("查询发布失败", zap.Error(err))
		response.FailWithMessage("回滚发布失败", c)
		return
	}

	// 检查状态
	if originalRelease.Status != "success" {
		response.FailWithMessage("只能回滚成功的发布", c)
		return
	}

	// 创建回滚发布
	rollbackRelease := model.Release{
		Version:      originalRelease.Version + "-rollback",
		Description:  "回滚到 " + originalRelease.Version,
		ReleaseNotes: "自动回滚到版本 " + originalRelease.Version,
		Status:       "pending",
		Environment:  originalRelease.Environment,
		ArtifactID:   originalRelease.ArtifactID,
		DeployedAt:   time.Now(),
		DeployedBy:   userID,
		IsRollback:   true,
	}

	if err := global.DB.Create(&rollbackRelease).Error; err != nil {
		global.Log.Error("创建回滚发布失败", zap.Error(err))
		response.FailWithMessage("回滚发布失败", c)
		return
	}

	// 查询完整的发布信息
	if err := global.DB.Preload("Artifact").Preload("User").First(&rollbackRelease, rollbackRelease.ID).Error; err != nil {
		global.Log.Error("查询发布信息失败", zap.Error(err))
		response.FailWithMessage("回滚发布成功，但获取详情失败", c)
		return
	}

	// 异步执行部署（模拟）
	go deployRelease(rollbackRelease.ID)

	response.OkWithData(rollbackRelease, c)
}
