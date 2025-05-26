package v1

import (
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

// CreateResourceQuota 创建资源配额
// @Summary 创建资源配额
// @Description 为指定租户创建资源配额
// @Tags 资源配额管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateResourceQuota true "资源配额信息"
// @Success 200 {object} response.Response
// @Router /resource-quota [post]
func CreateResourceQuota(c *gin.Context) {
	var req request.CreateResourceQuota
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	quota := model.ResourceQuota{
		TenantID:     req.TenantID,
		CPUQuota:     req.CPUQuota,
		MemoryQuota:  req.MemoryQuota,
		StorageQuota: req.StorageQuota,
	}

	if err := service.CreateResourceQuota(quota); err != nil {
		global.Log.Error("创建资源配额失败", zap.Error(err))
		response.FailWithMessage("创建资源配额失败", c)
		return
	}

	response.OkWithMessage("创建资源配额成功", c)
}

// GetResourceQuota 获取资源配额
// @Summary 获取资源配额
// @Description 根据租户ID获取资源配额
// @Tags 资源配额管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenant_id path string true "租户ID"
// @Success 200 {object} response.Response{data=model.ResourceQuota}
// @Router /resource-quota/{tenant_id} [get]
func GetResourceQuota(c *gin.Context) {
	tenantID := c.Param("tenant_id")

	quota, err := service.GetResourceQuotaByTenantID(tenantID)
	if err != nil {
		// 判断是否为记录未找到的错误
		if strings.Contains(err.Error(), "record not found") {
			global.Log.Warn("指定租户的资源配额记录未找到", zap.String("tenant_id", tenantID))
			response.FailWithMessage("指定租户的资源配额记录未找到", c)
		} else {
			global.Log.Error("获取资源配额失败", zap.Error(err))
			response.FailWithMessage("获取资源配额失败", c)
		}
		return
	}

	response.OkWithData(quota, c)
}

// UpdateResourceQuota 更新资源配额
// @Summary 更新资源配额
// @Description 根据租户ID更新资源配额
// @Tags 资源配额管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenant_id path string true "租户ID"
// @Param data body request.UpdateResourceQuota true "更新的资源配额信息"
// @Success 200 {object} response.Response
// @Router /resource-quota/{tenant_id} [put]
func UpdateResourceQuota(c *gin.Context) {
	tenantID := c.Param("tenant_id")
	var req request.UpdateResourceQuota
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	quota := model.ResourceQuota{
		CPUQuota:     req.CPUQuota,
		MemoryQuota:  req.MemoryQuota,
		StorageQuota: req.StorageQuota,
	}

	if err := service.UpdateResourceQuota(tenantID, quota); err != nil {
		global.Log.Error("更新资源配额失败", zap.Error(err))
		response.FailWithMessage("更新资源配额失败", c)
		return
	}

	response.OkWithMessage("更新资源配额成功", c)
}

// CreateResourceRequest 创建资源请求
// @Summary 创建资源请求
// @Description 租户发起资源请求
// @Tags 资源请求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body request.CreateResourceRequest true "资源请求信息"
// @Success 200 {object} response.Response
// @Router /resource-request [post]
func CreateResourceRequest(c *gin.Context) {
	var req request.CreateResourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	request := model.TenantResourceRequest{
		TenantID:       req.TenantID,
		CPURequest:     req.CPURequest,
		MemoryRequest:  req.MemoryRequest,
		StorageRequest: req.StorageRequest,
	}

	if err := service.CreateResourceRequest(request); err != nil {
		global.Log.Error("创建资源请求失败", zap.Error(err))
		response.FailWithMessage("创建资源请求失败", c)
		return
	}

	response.OkWithMessage("创建资源请求成功", c)
}

// GetResourceRequests 获取所有资源请求
// @Summary 获取所有资源请求
// @Description 获取所有待处理的资源请求
// @Tags 资源请求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.TenantResourceRequest}
// @Router /resource-requests [get]
func GetResourceRequests(c *gin.Context) {
	requests, err := service.GetResourceRequests()
	if err != nil {
		global.Log.Error("获取资源请求失败", zap.Error(err))
		response.FailWithMessage("获取资源请求失败", c)
		return
	}

	response.OkWithData(requests, c)
}

// ApproveResourceRequest 批准资源请求
// @Summary 批准资源请求
// @Description 管理员批准资源请求
// @Tags 资源请求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request_id path uint true "资源请求ID"
// @Success 200 {object} response.Response
// @Router /resource-request/{request_id}/approve [post]
func ApproveResourceRequest(c *gin.Context) {
	requestIDStr := c.Param("request_id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("请求ID格式错误", c)
		return
	}

	if err := service.ApproveResourceRequest(uint(requestID)); err != nil {
		global.Log.Error("批准资源请求失败", zap.Error(err))
		response.FailWithMessage("批准资源请求失败", c)
		return
	}

	response.OkWithMessage("批准资源请求成功", c)
}

// RejectResourceRequest 拒绝资源请求
// @Summary 拒绝资源请求
// @Description 管理员拒绝资源请求
// @Tags 资源请求管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request_id path uint true "资源请求ID"
// @Success 200 {object} response.Response
// @Router /resource-request/{request_id}/reject [post]
func RejectResourceRequest(c *gin.Context) {
	requestIDStr := c.Param("request_id")
	requestID, err := strconv.ParseUint(requestIDStr, 10, 32)
	if err != nil {
		response.FailWithMessage("请求ID格式错误", c)
		return
	}

	if err := service.RejectResourceRequest(uint(requestID)); err != nil {
		global.Log.Error("拒绝资源请求失败", zap.Error(err))
		response.FailWithMessage("拒绝资源请求失败", c)
		return
	}

	response.OkWithMessage("拒绝资源请求成功", c)
}

// GetAllResourceReports 获取所有资源报告
// @Summary 获取所有资源报告
// @Description 获取所有资源报告
// @Tags 资源报告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=[]model.ResourceReport}
// @Router /resource-report [get]
func GetAllResourceReports(c *gin.Context) {
	reports, err := service.GetAllResourceReports()
	if err != nil {
		global.Log.Error("获取资源报告失败", zap.Error(err))
		response.FailWithMessage("获取资源报告失败", c)
		return
	}
	response.OkWithData(reports, c)
}

// GetResourceReportByID 根据 ID 获取资源报告
// @Summary 根据 ID 获取资源报告
// @Description 根据 ID 获取资源报告
// @Tags 资源报告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "资源报告 ID"
// @Success 200 {object} response.Response{data=model.ResourceReport}
// @Router /resource-report/{id} [get]
func GetResourceReportByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("ID 格式错误", c)
		return
	}

	report, err := service.GetResourceReportByID(uint(id))
	if err != nil {
		global.Log.Error("获取资源报告失败", zap.Error(err))
		response.FailWithMessage("获取资源报告失败", c)
		return
	}
	response.OkWithData(report, c)
}

// CreateResourceReport 创建资源报告
// @Summary 创建资源报告
// @Description 创建资源报告
// @Tags 资源报告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param data body model.ResourceReport true "资源报告信息"
// @Success 200 {object} response.Response
// @Router /resource-report [post]
func CreateResourceReport(c *gin.Context) {
	var report model.ResourceReport
	if err := c.ShouldBindJSON(&report); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := service.CreateResourceReport(report); err != nil {
		global.Log.Error("创建资源报告失败", zap.Error(err))
		response.FailWithMessage("创建资源报告失败", c)
		return
	}
	response.OkWithMessage("创建资源报告成功", c)
}

// UpdateResourceReport 更新资源报告
// @Summary 更新资源报告
// @Description 根据 ID 更新资源报告
// @Tags 资源报告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "资源报告 ID"
// @Param data body model.ResourceReport true "更新的资源报告信息"
// @Success 200 {object} response.Response
// @Router /resource-report/{id} [put]
func UpdateResourceReport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("ID 格式错误", c)
		return
	}

	var report model.ResourceReport
	if err := c.ShouldBindJSON(&report); err != nil {
		response.FailWithMessage("参数错误: "+err.Error(), c)
		return
	}

	if err := service.UpdateResourceReport(uint(id), report); err != nil {
		global.Log.Error("更新资源报告失败", zap.Error(err))
		response.FailWithMessage("更新资源报告失败", c)
		return
	}
	response.OkWithMessage("更新资源报告成功", c)
}

// DeleteResourceReport 删除资源报告
// @Summary 删除资源报告
// @Description 根据 ID 删除资源报告
// @Tags 资源报告管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint true "资源报告 ID"
// @Success 200 {object} response.Response
// @Router /resource-report/{id} [delete]
func DeleteResourceReport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.FailWithMessage("ID 格式错误", c)
		return
	}

	if err := service.DeleteResourceReport(uint(id)); err != nil {
		global.Log.Error("删除资源报告失败", zap.Error(err))
		response.FailWithMessage("删除资源报告失败", c)
		return
	}
	response.OkWithMessage("删除资源报告成功", c)
}
