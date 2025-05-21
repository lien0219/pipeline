package v1

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"gin_pipeline/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

// CreatePipeline 创建流水线
// @Summary 创建流水线
// @Description 创建新的流水线
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreatePipeline true "流水线信息"
// @Success 200 {object} response.Response{data=model.Pipeline} "创建成功"
// @Router /pipeline [post]
func CreatePipeline(c *gin.Context) {
	var req request.CreatePipeline
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建流水线失败", c)
		return
	}

	// 创建流水线
	pipeline := model.Pipeline{
		Name:        req.Name,
		Description: req.Description,
		GitRepo:     req.GitRepo,
		GitBranch:   req.GitBranch,
		Status:      "inactive",
		CreatorID:   userID,
	}

	// 开启事务
	tx := global.DB.Begin()

	// 创建流水线
	if err := tx.Create(&pipeline).Error; err != nil {
		tx.Rollback()
		global.Log.Error("创建流水线失败", zap.Error(err))
		response.FailWithMessage("创建流水线失败", c)
		return
	}

	// 创建阶段和作业
	for i, stage := range req.Stages {
		// 如果没有指定顺序，则按照数组顺序
		if stage.Order == 0 {
			stage.Order = i + 1
		}

		newStage := model.Stage{
			Name:        stage.Name,
			Description: stage.Description,
			Order:       stage.Order,
			PipelineID:  pipeline.ID,
		}

		// 创建阶段
		if err := tx.Create(&newStage).Error; err != nil {
			tx.Rollback()
			global.Log.Error("创建阶段失败", zap.Error(err))
			response.FailWithMessage("创建流水线失败", c)
			return
		}

		// 创建作业
		for _, job := range stage.Jobs {
			newJob := model.Job{
				Name:        job.Name,
				Description: job.Description,
				Command:     job.Command,
				Image:       job.Image,
				Timeout:     job.Timeout,
				StageID:     newStage.ID,
			}

			if err := tx.Create(&newJob).Error; err != nil {
				tx.Rollback()
				global.Log.Error("创建作业失败", zap.Error(err))
				response.FailWithMessage("创建流水线失败", c)
				return
			}
		}
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Error("提交事务失败", zap.Error(err))
		response.FailWithMessage("创建流水线失败", c)
		return
	}

	// 查询完整的流水线信息
	var result model.Pipeline
	if err := global.DB.Preload("Creator").Preload("Stages", func(db *gorm.DB) *gorm.DB {
		return db.Order("stages.order ASC")
	}).Preload("Stages.Jobs").First(&result, pipeline.ID).Error; err != nil {
		global.Log.Error("查询流水线失败", zap.Error(err))
		response.FailWithMessage("创建流水线成功，但获取详情失败", c)
		return
	}

	response.OkWithData(result, c)
}

// GetPipelines 获取流水线列表
// @Summary 获取流水线列表
// @Description 获取流水线列表，支持分页
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.Pipeline}} "获取成功"
// @Router /pipeline [get]
func GetPipelines(c *gin.Context) {
	var pageInfo request.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 获取分页参数
	page := pageInfo.GetPage()
	pageSize := pageInfo.GetPageSize()

	// 查询总数
	var total int64
	if err := global.DB.Model(&model.Pipeline{}).Count(&total).Error; err != nil {
		global.Log.Error("查询流水线总数失败", zap.Error(err))
		response.FailWithMessage("获取流水线列表失败", c)
		return
	}

	// 查询列表
	var pipelines []model.Pipeline
	if err := global.DB.Preload("Creator").Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}).Order("id DESC").Find(&pipelines).Error; err != nil {
		global.Log.Error("查询流水线列表失败", zap.Error(err))
		response.FailWithMessage("获取流水线列表失败", c)
		return
	}

	// 返回结果
	result := response.PageResult{
		List:     pipelines,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.OkWithData(result, c)
}

// GetPipelineByID 获取流水线详情
// @Summary 获取流水线详情
// @Description 根据ID获取流水线详情
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流水线ID"
// @Success 200 {object} response.Response{data=model.Pipeline} "获取成功"
// @Router /pipeline/{id} [get]
func GetPipelineByID(c *gin.Context) {
	id := c.Param("id")

	var pipeline model.Pipeline
	if err := global.DB.Preload("Creator").Preload("Stages", func(db *gorm.DB) *gorm.DB {
		return db.Order("stages.order ASC")
	}).Preload("Stages.Jobs").First(&pipeline, id).Error; err != nil {
		global.Log.Error("查询流水线失败", zap.Error(err))
		response.FailWithMessage("获取流水线详情失败", c)
		return
	}

	response.OkWithData(pipeline, c)
}

// UpdatePipeline 更新流水线
// @Summary 更新流水线
// @Description 更新流水线信息
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流水线ID"
// @Param data body request.UpdatePipeline true "流水线信息"
// @Success 200 {object} response.Response{data=model.Pipeline} "更新成功"
// @Router /pipeline/{id} [put]
func UpdatePipeline(c *gin.Context) {
	id := c.Param("id")
	var req request.UpdatePipeline
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 查询流水线
	var pipeline model.Pipeline
	if err := global.DB.First(&pipeline, id).Error; err != nil {
		global.Log.Error("查询流水线失败", zap.Error(err))
		response.FailWithMessage("更新流水线失败", c)
		return
	}

	// 更新流水线
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"git_repo":    req.GitRepo,
		"git_branch":  req.GitBranch,
	}

	// 如果状态不为空，则更新状态
	if req.Status != "" {
		updates["status"] = req.Status
	}

	if err := global.DB.Model(&pipeline).Updates(updates).Error; err != nil {
		global.Log.Error("更新流水线失败", zap.Error(err))
		response.FailWithMessage("更新流水线失败", c)
		return
	}

	// 查询更新后的流水线
	if err := global.DB.Preload("Creator").Preload("Stages", func(db *gorm.DB) *gorm.DB {
		return db.Order("stages.order ASC")
	}).Preload("Stages.Jobs").First(&pipeline, id).Error; err != nil {
		global.Log.Error("查询流水线失败", zap.Error(err))
		response.FailWithMessage("更新流水线成功，但获取详情失败", c)
		return
	}

	response.OkWithData(pipeline, c)
}

// DeletePipeline 删除流水线
// @Summary 删除流水线
// @Description 根据ID删除流水线
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流水线ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /pipeline/{id} [delete]
func DeletePipeline(c *gin.Context) {
	id := c.Param("id")

	// 开启事务
	tx := global.DB.Begin()

	// 删除流水线相关的作业
	if err := tx.Where("stage_id IN (SELECT id FROM stages WHERE pipeline_id = ?)", id).Delete(&model.Job{}).Error; err != nil {
		tx.Rollback()
		global.Log.Error("删除作业失败", zap.Error(err))
		response.FailWithMessage("删除流水线失败", c)
		return
	}

	// 删除流水线相关的阶段
	if err := tx.Where("pipeline_id = ?", id).Delete(&model.Stage{}).Error; err != nil {
		tx.Rollback()
		global.Log.Error("删除阶段失败", zap.Error(err))
		response.FailWithMessage("删除流水线失败", c)
		return
	}

	// 删除流水线
	if err := tx.Delete(&model.Pipeline{}, id).Error; err != nil {
		tx.Rollback()
		global.Log.Error("删除流水线失败", zap.Error(err))
		response.FailWithMessage("删除流水线失败", c)
		return
	}

	// 提交事务
	if err := tx.Commit().Error; err != nil {
		global.Log.Error("提交事务失败", zap.Error(err))
		response.FailWithMessage("删除流水线失败", c)
		return
	}

	response.OkWithMessage("删除流水线成功", c)
}

// 异步执行流水线（模拟）
func executePipeline(runID uint) {
	// 查询运行记录
	var pipelineRun model.PipelineRun
	if err := global.DB.First(&pipelineRun, runID).Error; err != nil {
		global.Log.Error("查询流水线运行记录失败", zap.Error(err), zap.Uint("runID", runID))
		return
	}

	// 更新状态为运行中
	if err := global.DB.Model(&pipelineRun).Update("status", "running").Error; err != nil {
		global.Log.Error("更新流水线运行状态失败", zap.Error(err), zap.Uint("runID", runID))
		return
	}

	// 模拟执行过程
	time.Sleep(5 * time.Second)

	// 随机成功或失败
	status := "success"
	if time.Now().Unix()%2 == 0 {
		status = "failed"
	}

	// 更新运行结果
	now := time.Now()
	duration := int(now.Sub(*pipelineRun.StartTime).Seconds())
	updates := map[string]interface{}{
		"status":   status,
		"end_time": now,
		"duration": duration,
		"logs":     "模拟执行日志...",
	}

	if err := global.DB.Model(&pipelineRun).Updates(updates).Error; err != nil {
		global.Log.Error("更新流水线运行结果失败", zap.Error(err), zap.Uint("runID", runID))
		return
	}

	// 更新流水线状态
	if err := global.DB.Model(&model.Pipeline{}).Where("id = ?", pipelineRun.PipelineID).Update("status", status).Error; err != nil {
		global.Log.Error("更新流水线状态失败", zap.Error(err), zap.Uint("pipelineID", pipelineRun.PipelineID))
		return
	}

	global.Log.Info("流水线执行完成", zap.Uint("runID", runID), zap.String("status", status))
}

// GetPipelineRuns 获取流水线运行记录
// @Summary 获取流水线运行记录
// @Description 获取指定流水线的运行记录
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流水线ID"
// @Param page query int false "页码" default(1)
// @Param pageSize query int false "每页大小" default(10)
// @Success 200 {object} response.Response{data=response.PageResult{list=[]model.PipelineRun}} "获取成功"
// @Router /pipeline/{id}/runs [get]
func GetPipelineRuns(c *gin.Context) {
	id := c.Param("id")
	var pageInfo request.PageInfo
	if err := c.ShouldBindQuery(&pageInfo); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 获取分页参数
	page := pageInfo.GetPage()
	pageSize := pageInfo.GetPageSize()

	// 查询总数
	var total int64
	if err := global.DB.Model(&model.PipelineRun{}).Where("pipeline_id = ?", id).Count(&total).Error; err != nil {
		global.Log.Error("查询流水线运行记录总数失败", zap.Error(err))
		response.FailWithMessage("获取流水线运行记录失败", c)
		return
	}

	// 查询列表
	var runs []model.PipelineRun
	if err := global.DB.Preload("User").Where("pipeline_id = ?", id).Scopes(func(db *gorm.DB) *gorm.DB {
		return db.Offset((page - 1) * pageSize).Limit(pageSize)
	}).Order("id DESC").Find(&runs).Error; err != nil {
		global.Log.Error("查询流水线运行记录失败", zap.Error(err))
		response.FailWithMessage("获取流水线运行记录失败", c)
		return
	}

	// 返回结果
	result := response.PageResult{
		List:     runs,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	response.OkWithData(result, c)
}

// GetPipelineRunByID 获取流水线运行记录详情
// @Summary 获取流水线运行记录详情
// @Description 获取指定流水线运行记录的详情
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流水线ID"
// @Param runId path int true "运行记录ID"
// @Success 200 {object} response.Response{data=model.PipelineRun} "获取成功"
// @Router /pipeline/{id}/runs/{runId} [get]
func GetPipelineRunByID(c *gin.Context) {
	id := c.Param("id")
	runID := c.Param("runId")

	var run model.PipelineRun
	if err := global.DB.Preload("Pipeline").Preload("User").Where("pipeline_id = ? AND id = ?", id, runID).First(&run).Error; err != nil {
		global.Log.Error("查询流水线运行记录失败", zap.Error(err))
		response.FailWithMessage("获取流水线运行记录详情失败", c)
		return
	}

	response.OkWithData(run, c)
}

// GetPipelineRunLogs 获取流水线运行日志
// @Summary 获取流水线运行日志
// @Description 获取指定流水线运行记录的日志
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流水线ID"
// @Param runId path int true "运行记录ID"
// @Success 200 {object} response.Response{data=map[string]interface{}} "获取成功"
// @Router /pipeline/{id}/runs/{runId}/logs [get]
func GetPipelineRunLogs(c *gin.Context) {
	id := c.Param("id")
	runID := c.Param("runId")

	var run model.PipelineRun
	if err := global.DB.Select("logs").Where("pipeline_id = ? AND id = ?", id, runID).First(&run).Error; err != nil {
		global.Log.Error("查询流水线运行日志失败", zap.Error(err))
		response.FailWithMessage("获取流水线运行日志失败", c)
		return
	}

	// 返回日志
	data := map[string]interface{}{
		"logs": run.Logs,
	}
	response.OkWithData(data, c)
}

// TriggerPipeline 触发流水线
// @Summary 触发流水线
// @Description 触发流水线执行
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流水线ID"
// @Param data body request.TriggerPipeline false "触发参数"
// @Success 200 {object} response.Response{data=model.PipelineRun} "触发成功"
// @Router /pipeline/{id}/trigger [post]
func TriggerPipeline(c *gin.Context) {
	id := c.Param("id")
	var req request.TriggerPipeline
	if err := c.ShouldBindJSON(&req); err != nil {
		// 忽略绑定错误，使用默认值
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("触发流水线失败", c)
		return
	}

	// 查询流水线
	var pipeline model.Pipeline
	if err := global.DB.First(&pipeline, id).Error; err != nil {
		global.Log.Error("查询流水线失败", zap.Error(err))
		response.FailWithMessage("触发流水线失败", c)
		return
	}

	// 使用指定的分支或默认分支
	gitBranch := req.GitBranch
	if gitBranch == "" {
		gitBranch = pipeline.GitBranch
	}

	// 使用工作流服务触发流水线
	workflowService := service.NewWorkflowService()
	pipelineRun, err := workflowService.TriggerWorkflow(pipeline.ID, userID, gitBranch)
	if err != nil {
		global.Log.Error("触发流水线失败", zap.Error(err))
		response.FailWithMessage("触发流水线失败: "+err.Error(), c)
		return
	}

	// 查询完整的运行记录
	if err := global.DB.Preload("Pipeline").Preload("User").First(pipelineRun, pipelineRun.ID).Error; err != nil {
		global.Log.Error("查询流水线运行记录失败", zap.Error(err))
		response.FailWithMessage("触发流水线成功，但获取详情失败", c)
		return
	}

	response.OkWithData(pipelineRun, c)
}

// 修改CancelPipelineRun函数，使用工作流服务

// CancelPipelineRun 取消流水线运行
// @Summary 取消流水线运行
// @Description 取消指定的流水线运行
// @Tags 流水线管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "流水线ID"
// @Param runId path int true "运行记录ID"
// @Success 200 {object} response.Response "取消成功"
// @Router /pipeline/{id}/runs/{runId}/cancel [post]
func CancelPipelineRun(c *gin.Context) {
	id := c.Param("id")
	runID := c.Param("runId")

	// 查询运行记录
	var run model.PipelineRun
	if err := global.DB.Where("pipeline_id = ? AND id = ?", id, runID).First(&run).Error; err != nil {
		global.Log.Error("查询流水线运行记录失败", zap.Error(err))
		response.FailWithMessage("取消流水线运行失败", c)
		return
	}

	// 检查状态
	if run.Status != "pending" && run.Status != "running" {
		response.FailWithMessage("只能取消等待中或运行中的流水线", c)
		return
	}

	// 使用工作流服务取消流水线
	workflowService := service.NewWorkflowService()
	if err := workflowService.CancelWorkflow(run.ID); err != nil {
		global.Log.Error("取消流水线运行失败", zap.Error(err))
		response.FailWithMessage("取消流水线运行失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("取消流水线运行成功", c)
}
