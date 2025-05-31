package v1

import (
	"encoding/json"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/model/response"
	"gin_pipeline/service"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var dagService = new(service.DAGService)

// CreateDAG 创建DAG
// @Summary 创建DAG
// @Description 创建新的DAG
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateDAG true "DAG信息"
// @Success 200 {object} response.Response{data=model.DAG} "创建成功"
// @Router /dag [post]
func CreateDAG(c *gin.Context) {
	var req request.CreateDAG
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建DAG失败，未获取到用户ID", c)
		return
	}

	// 创建DAG
	dag := model.DAG{
		Name:        req.Name,
		Description: req.Description,
		PipelineID:  req.PipelineID,
		NodesData:   req.Nodes,
		CreatorID:   userID,
		IsActive:    true, // 默认为活动版本
		Version:     1,    // 初始版本
	}

	if err := dagService.CreateDAG(&dag); err != nil {
		if _, ok := err.(*json.SyntaxError); ok {
			global.Log.Error("创建DAG时JSON解析错误", zap.Error(err))
			response.FailWithMessage("创建DAG失败，请求数据格式错误", c)
		} else if strings.Contains(err.Error(), "依赖节点不存在") {
			// 提取缺失的节点 ID
			startIndex := strings.Index(err.Error(), "依赖节点不存在: ")
			if startIndex != -1 {
				missingNodeID := err.Error()[startIndex+len("依赖节点不存在: "):]
				global.Log.Error("创建DAG时依赖节点不存在",
					zap.String("missingNodeID", missingNodeID),
					zap.Any("requestNodes", req.Nodes),
					zap.Error(err))
			}
			response.FailWithMessage("创建DAG失败，"+err.Error(), c)
		} else {
			global.Log.Error("创建DAG失败", zap.Error(err))
			response.FailWithMessage("创建DAG失败: "+err.Error(), c)
		}
		return
	}

	response.OkWithData(dag, c)
}

// GetDAGByID 获取DAG详情
// @Summary 获取DAG详情
// @Description 根据ID获取DAG详情
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "DAG ID"
// @Success 200 {object} response.Response{data=model.DAG} "获取成功"
// @Router /dag/{id} [get]
func GetDAGByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	dag, err := dagService.GetDAGByID(uint(id))
	if err != nil {
		global.Log.Error("获取DAG失败", zap.Error(err))
		response.FailWithMessage("获取DAG失败", c)
		return
	}

	response.OkWithData(dag, c)
}

// GetDAGsByPipelineID 获取流水线的所有DAG
// @Summary 获取流水线的所有DAG
// @Description 获取指定流水线的所有DAG
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pipelineId path int true "流水线ID"
// @Success 200 {object} response.Response{data=[]model.DAG} "获取成功"
// @Router /dag/pipeline/{pipelineId} [get]
func GetDAGsByPipelineID(c *gin.Context) {
	pipelineID, err := strconv.ParseUint(c.Param("pipelineId"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的流水线ID", c)
		return
	}

	dags, err := dagService.GetDAGsByPipelineID(uint(pipelineID))
	if err != nil {
		global.Log.Error("获取DAG列表失败", zap.Error(err))
		response.FailWithMessage("获取DAG列表失败", c)
		return
	}

	response.OkWithData(dags, c)
}

// GetActiveDAG 获取流水线的活动DAG
// @Summary 获取流水线的活动DAG
// @Description 获取指定流水线的活动DAG
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pipelineId path int true "流水线ID"
// @Success 200 {object} response.Response{data=model.DAG} "获取成功"
// @Router /dag/pipeline/{pipelineId}/active [get]
func GetActiveDAG(c *gin.Context) {
	pipelineID, err := strconv.ParseUint(c.Param("pipelineId"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的流水线ID", c)
		return
	}

	dag, err := dagService.GetActiveDAGByPipelineID(uint(pipelineID))
	if err != nil {
		global.Log.Error("获取活动DAG失败", zap.Error(err))
		response.FailWithMessage("获取活动DAG失败", c)
		return
	}

	response.OkWithData(dag, c)
}

// UpdateDAG 更新DAG
// @Summary 更新DAG
// @Description 更新DAG信息
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "DAG ID"
// @Param data body request.UpdateDAG true "DAG信息"
// @Success 200 {object} response.Response{data=model.DAG} "更新成功"
// @Router /dag/{id} [put]
func UpdateDAG(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	var req request.UpdateDAG
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	// 更新DAG
	updates := map[string]interface{}{
		"name":        req.Name,
		"description": req.Description,
		"nodes_data":  req.Nodes,
	}

	if err := dagService.UpdateDAG(uint(id), updates); err != nil {
		global.Log.Error("更新DAG失败", zap.Error(err))
		response.FailWithMessage("更新DAG失败: "+err.Error(), c)
		return
	}

	// 获取更新后的DAG
	dag, err := dagService.GetDAGByID(uint(id))
	if err != nil {
		global.Log.Error("获取更新后的DAG失败", zap.Error(err))
		response.FailWithMessage("更新DAG成功，但获取详情失败", c)
		return
	}

	response.OkWithData(dag, c)
}

// DeleteDAG 删除DAG
// @Summary 删除DAG
// @Description 删除指定的DAG
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "DAG ID"
// @Success 200 {object} response.Response "删除成功"
// @Router /dag/{id} [delete]
func DeleteDAG(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	if err := dagService.DeleteDAG(uint(id)); err != nil {
		global.Log.Error("删除DAG失败", zap.Error(err))
		response.FailWithMessage("删除DAG失败", c)
		return
	}

	response.OkWithMessage("删除DAG成功", c)
}

// ValidateDAG 验证DAG
// @Summary 验证DAG
// @Description 验证DAG是否有效（无环）
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.ValidateDAG true "DAG节点"
// @Success 200 {object} response.Response "验证成功"
// @Router /dag/validate [post]
func ValidateDAG(c *gin.Context) {
	var req request.ValidateDAG
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := dagService.ValidateDAG(req.Nodes); err != nil {
		response.FailWithMessage("DAG验证失败: "+err.Error(), c)
		return
	}

	response.OkWithMessage("DAG验证通过", c)
}

// CreateDAGVersion 创建DAG的新版本
// @Summary 创建DAG的新版本
// @Description 基于现有DAG创建新版本
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "DAG ID"
// @Success 200 {object} response.Response{data=model.DAG} "创建成功"
// @Router /dag/{id}/version [post]
func CreateDAGVersion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	// 从上下文获取用户ID
	userID := c.GetUint("userId")
	if userID == 0 {
		response.FailWithMessage("创建DAG版本失败", c)
		return
	}

	dag, err := dagService.CreateDAGVersion(uint(id), userID)
	if err != nil {
		global.Log.Error("创建DAG版本失败", zap.Error(err))
		response.FailWithMessage("创建DAG版本失败", c)
		return
	}

	response.OkWithData(dag, c)
}

// GetDAGHistory 获取DAG历史版本
// @Summary 获取DAG历史版本
// @Description 获取指定流水线的DAG历史版本
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param pipelineId path int true "流水线ID"
// @Success 200 {object} response.Response{data=[]model.DAG} "获取成功"
// @Router /dag/pipeline/{pipelineId}/history [get]
func GetDAGHistory(c *gin.Context) {
	pipelineID, err := strconv.ParseUint(c.Param("pipelineId"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的流水线ID", c)
		return
	}

	dags, err := dagService.GetDAGHistory(uint(pipelineID))
	if err != nil {
		global.Log.Error("获取DAG历史版本失败", zap.Error(err))
		response.FailWithMessage("获取DAG历史版本失败", c)
		return
	}

	response.OkWithData(dags, c)
}

// ActivateDAG 激活DAG版本
// @Summary 激活DAG版本
// @Description 激活指定的DAG版本
// @Tags DAG管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "DAG ID"
// @Success 200 {object} response.Response "激活成功"
// @Router /dag/{id}/activate [post]
func ActivateDAG(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		response.FailWithMessage("无效的ID", c)
		return
	}

	if err := dagService.ActivateDAG(uint(id)); err != nil {
		global.Log.Error("激活DAG失败", zap.Error(err))
		response.FailWithMessage("激活DAG失败", c)
		return
	}

	response.OkWithMessage("激活DAG成功", c)
}
