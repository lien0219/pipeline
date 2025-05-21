package v1

import (
	"fmt"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"gin_pipeline/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"io"
	"os"
	"path"
	"strconv"
	"time"
)

// CreateArtifact 创建制品
// @Summary 创建制品
// @Description 上传并创建新的制品
// @Tags 制品管理
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param file formData file true "制品文件"
// @Param name formData string true "制品名称"
// @Param type formData string true "制品类型"
// @Param version formData string false "版本号"
// @Param description formData string false "描述"
// @Param pipeline_id formData int false "关联的流水线ID"
// @Param pipeline_run_id formData int false "关联的流水线运行ID"
// @Success 200 {object} response.Response{data=model.Artifact} "创建成功"
// @Router /artifact [post]
func CreateArtifact(c *gin.Context) {
	// 获取表单数据
	name := c.PostForm("name")
	artifactType := c.PostForm("type")
	version := c.PostForm("version")
	description := c.PostForm("description")
	pipelineIDStr := c.PostForm("pipeline_id")
	pipelineRunIDStr := c.PostForm("pipeline_run_id")

	// 验证必填字段
	if name == "" || artifactType == "" {
		response.FailWithMessage("名称和类型不能为空", c)
		return
	}

	// 解析关联ID
	var pipelineID, pipelineRunID uint
	if pipelineIDStr != "" {
		id, err := strconv.ParseUint(pipelineIDStr, 10, 32)
		if err != nil {
			response.FailWithMessage("流水线ID格式错误", c)
			return
		}
		pipelineID = uint(id)
	}

	if pipelineRunIDStr != "" {
		id, err := strconv.ParseUint(pipelineRunIDStr, 10, 32)
		if err != nil {
			response.FailWithMessage("流水线运行ID格式错误", c)
			return
		}
		pipelineRunID = uint(id)
	}

	// 获取上传文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.FailWithMessage("获取上传文件失败: "+err.Error(), c)
		return
	}
	defer file.Close()

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建制品失败", c)
		return
	}

	// 创建存储目录
	uploadDir := global.Config.Upload.Local.StorePath
	if ok := utils.CreateDir(uploadDir); !ok {
		response.FailWithMessage("创建存储目录失败", c)
		return
	}

	// 生成文件路径
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filePath := path.Join(uploadDir, fileName)

	// 创建目标文件
	out, err := os.Create(filePath)
	if err != nil {
		global.Log.Error("创建文件失败", zap.Error(err))
		response.FailWithMessage("创建文件失败", c)
		return
	}
	defer out.Close()

	// 复制文件内容
	_, err = io.Copy(out, file)
	if err != nil {
		global.Log.Error("保存文件失败", zap.Error(err))
		response.FailWithMessage("保存文件失败", c)
		return
	}

	// 创建制品记录
	artifact := model.Artifact{
		Name:          name,
		Type:          artifactType,
		Path:          filePath,
		Size:          header.Size,
		Version:       version,
		Description:   description,
		PipelineID:    pipelineID,
		PipelineRunID: pipelineRunID,
		CreatedBy:     userID,
	}

	if err := global.DB.Create(&artifact).Error; err != nil {
		global.Log.Error("创建制品记录失败", zap.Error(err))
		response.FailWithMessage("创建制品记录失败", c)
		return
	}

	// 查询完整的制品信息
	if err := global.DB.Preload("Pipeline").Preload("User").First(&artifact, artifact.ID).Error; err != nil {
		global.Log.Error("查询制品信息失败", zap.Error(err))
		response.FailWithMessage("创建制品成功，但获取详情失败", c)
		return
	}

	response.OkWithData(artifact, c)
}

// GetArtifacts 获取制品列表
// @Summary 获取制品列表
// @Description 获取制品列表，支持分页和筛选
// @Tags 制品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Param name query string false "制品名称"
// @Param type query string false "制品类型"
// @Param pipeline_id query int false "流水线ID"
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Artifact}} "获取成功"
// @Router /artifact [get]
func GetArtifacts(c *gin.Context) {
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
	artifactType := c.Query("type")
	pipelineIDStr := c.Query("pipeline_id")

	// 构建查询条件
	db := global.DB.Model(&model.Artifact{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}

	if artifactType != "" {
		db = db.Where("type = ?", artifactType)
	}

	if pipelineIDStr != "" {
		pipelineID, err := strconv.ParseUint(pipelineIDStr, 10, 32)
		if err == nil {
			db = db.Where("pipeline_id = ?", pipelineID)
		}
	}

	// 查询总数
	var total int64
	if err := db.Count(&total).Error; err != nil {
		global.Log.Error("查询制品总数失败", zap.Error(err))
		response.FailWithMessage("获取制品列表失败", c)
		return
	}

	// 查询列表
	var artifacts []model.Artifact
	if err := db.Preload("Pipeline").Preload("User").Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}).Order("id DESC").Find(&artifacts).Error; err != nil {
		global.Log.Error("查询制品列表失败", zap.Error(err))
		response.FailWithMessage("获取制品列表失败", c)
		return
	}

	// 返回结果
	result := response.PageResult{
		List:     artifacts,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.OkWithData(result, c)
}

// GetArtifactByID 获取制品详情
// @Summary 获取制品详情
// @Description 根据ID获取制品详情
// @Tags 制品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "制品ID"
// @Success 200 {object} response.Response{data=model.Artifact} "获取成功"
// @Router /artifact/{id} [get]
func GetArtifactByID(c *gin.Context) {
	id := c.Param("id")

	var artifact model.Artifact
	if err := global.DB.Preload("Pipeline").Preload("User").First(&artifact, id).Error; err != nil {
		global.Log.Error("查询制品失败", zap.Error(err))
		response.FailWithMessage("获取制品详情失败", c)
		return
	}

	response.OkWithData(artifact, c)
}

// DeleteArtifact 删除制品
// @Summary 删除制品
// @Description 根据ID删除制品
// @Tags 制品管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "制品ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /artifact/{id} [delete]
func DeleteArtifact(c *gin.Context) {
	id := c.Param("id")

	// 查询制品
	var artifact model.Artifact
	if err := global.DB.First(&artifact, id).Error; err != nil {
		global.Log.Error("查询制品失败", zap.Error(err))
		response.FailWithMessage("删除制品失败", c)
		return
	}

	// 开启事务
	tx := global.DB.Begin()

	// 删除制品记录
	if err := tx.Delete(&artifact).Error; err != nil {
		tx.Rollback()
		global.Log.Error("删除制品记录失败", zap.Error(err))
		response.FailWithMessage("删除制品失败", c)
		return
	}

	// 删除文件
	if artifact.Path != "" {
		if err := os.Remove(artifact.Path); err != nil {
			global.Log.Warn("删除制品文件失败", zap.Error(err))
			// 不影响事务提交
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Error("提交事务失败", zap.Error(err))
		response.FailWithMessage("删除制品失败", c)
		return
	}

	response.OkWithMessage("删除制品成功", c)
}

// DownloadArtifact 下载制品
// @Summary 下载制品
// @Description 下载制品文件
// @Tags 制品管理
// @Accept json
// @Produce octet-stream
// @Security BearerAuth
// @Param id path int true "制品ID"
// @Success 200 {file} binary "文件内容"
// @Router /artifact/{id}/download [get]
func DownloadArtifact(c *gin.Context) {
	id := c.Param("id")

	// 查询制品
	var artifact model.Artifact
	if err := global.DB.First(&artifact, id).Error; err != nil {
		global.Log.Error("查询制品失败", zap.Error(err))
		response.FailWithMessage("下载制品失败", c)
		return
	}

	// 检查文件是否存在
	if artifact.Path == "" {
		response.FailWithMessage("制品文件路径为空", c)
		return
	}

	if _, err := os.Stat(artifact.Path); os.IsNotExist(err) {
		response.FailWithMessage("制品文件不存在", c)
		return
	}

	// 更新下载次数
	if err := global.DB.Model(&artifact).Update("download_count", gorm.Expr("download_count + ?", 1)).Error; err != nil {
		global.Log.Warn("更新下载次数失败", zap.Error(err))
		// 不影响下载
	}

	// 设置响应头
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", artifact.Name))
	c.Header("Content-Type", "application/octet-stream")
	c.Header("Content-Transfer-Encoding", "binary")

	// 发送文件
	c.File(artifact.Path)
}
