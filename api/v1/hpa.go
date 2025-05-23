package v1

import (
	"gin_pipeline/model"
	"gin_pipeline/model/request"
	"gin_pipeline/service"

	"github.com/gin-gonic/gin"
)

// CreateHPAPolicy 创建 HPA 策略
func CreateHPAPolicy(c *gin.Context) {
	var hpaRequest request.CreateHPAPolicy
	if err := c.ShouldBindJSON(&hpaRequest); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "请求参数解析失败: " + err.Error(),
		})
		return
	}

	hpaPolicy := model.HPAPolicy{
		Name:            hpaRequest.Name,
		Namespace:       hpaRequest.Namespace,
		Deployment:      hpaRequest.Deployment,
		MinReplicas:     hpaRequest.MinReplicas,
		MaxReplicas:     hpaRequest.MaxReplicas,
		CPUThreshold:    hpaRequest.CPUThreshold,
		MemoryThreshold: hpaRequest.MemoryThreshold,
	}

	if err := service.CreateHPAPolicy(hpaPolicy); err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "创建 HPA 策略失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
	})
}

// GetHPAPolicies 获取所有 HPA 策略
func GetHPAPolicies(c *gin.Context) {
	hpaPolicies, err := service.GetHPAPolicies()
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "获取 HPA 策略列表失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   hpaPolicies,
	})
}

// GetHPAPolicyByID 根据 ID 获取 HPA 策略
func GetHPAPolicyByID(c *gin.Context) {
	id := c.Param("id")
	hpaPolicy, err := service.GetHPAPolicyByID(id)
	if err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "根据 ID 获取 HPA 策略失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
		"data":   hpaPolicy,
	})
}

// UpdateHPAPolicy 更新 HPA 策略
func UpdateHPAPolicy(c *gin.Context) {
	id := c.Param("id")
	var hpaRequest request.UpdateHPAPolicy
	if err := c.ShouldBindJSON(&hpaRequest); err != nil {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "请求参数解析失败: " + err.Error(),
		})
		return
	}

	hpaPolicy := model.HPAPolicy{
		Name:            hpaRequest.Name,
		Namespace:       hpaRequest.Namespace,
		Deployment:      hpaRequest.Deployment,
		MinReplicas:     hpaRequest.MinReplicas,
		MaxReplicas:     hpaRequest.MaxReplicas,
		CPUThreshold:    hpaRequest.CPUThreshold,
		MemoryThreshold: hpaRequest.MemoryThreshold,
	}

	if err := service.UpdateHPAPolicy(id, hpaPolicy); err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "更新 HPA 策略失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
	})
}

// DeleteHPAPolicy 删除 HPA 策略
func DeleteHPAPolicy(c *gin.Context) {
	id := c.Param("id")
	if err := service.DeleteHPAPolicy(id); err != nil {
		c.JSON(500, gin.H{
			"status":  "error",
			"message": "删除 HPA 策略失败: " + err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"status": "success",
	})
}
