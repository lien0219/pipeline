package v1

import (
	"gin_pipeline/service"
	"net/http"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/batch/v1"
	batchv1beta1 "k8s.io/api/batch/v1beta1"
)

var K8sSvc *service.K8sService

// CreateJobHandler 创建 Job 的处理函数
func CreateJobHandler(c *gin.Context) {
	var job v1.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	namespace := c.Param("namespace")
	result, err := K8sSvc.CreateJob(namespace, &job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// GetJobHandler 获取 Job 的处理函数
func GetJobHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	result, err := K8sSvc.GetJob(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// UpdateJobHandler 更新 Job 的处理函数
func UpdateJobHandler(c *gin.Context) {
	var job v1.Job
	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	namespace := c.Param("namespace")
	result, err := K8sSvc.UpdateJob(namespace, &job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// DeleteJobHandler 删除 Job 的处理函数
func DeleteJobHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := K8sSvc.DeleteJob(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// CreateCronJobHandler 创建 CronJob 的处理函数
func CreateCronJobHandler(c *gin.Context) {
	var cronJob batchv1beta1.CronJob
	if err := c.ShouldBindJSON(&cronJob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	namespace := c.Param("namespace")
	result, err := K8sSvc.CreateCronJob(namespace, &cronJob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, result)
}

// GetCronJobHandler 获取 CronJob 的处理函数
func GetCronJobHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	result, err := K8sSvc.GetCronJob(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// UpdateCronJobHandler 更新 CronJob 的处理函数
func UpdateCronJobHandler(c *gin.Context) {
	var cronJob batchv1beta1.CronJob
	if err := c.ShouldBindJSON(&cronJob); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	namespace := c.Param("namespace")
	result, err := K8sSvc.UpdateCronJob(namespace, &cronJob)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// DeleteCronJobHandler 删除 CronJob 的处理函数
func DeleteCronJobHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := K8sSvc.DeleteCronJob(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// CreateDeploymentHandler 创建 Deployment 的处理函数
func CreateDeploymentHandler(c *gin.Context) {
	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数错误",
			"error":   err.Error(),
		})
		return
	}

	namespace := c.Param("namespace")
	if namespace == "" {
		namespace = "default" // 设置默认namespace
	}

	result, err := K8sSvc.CreateDeployment(namespace, &deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建Deployment失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"code":    201,
		"message": "创建成功",
		"data":    result,
	})
}

// GetDeploymentHandler 获取 Deployment 的处理函数
func GetDeploymentHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	result, err := K8sSvc.GetDeployment(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// UpdateDeploymentHandler 更新 Deployment 的处理函数
func UpdateDeploymentHandler(c *gin.Context) {
	var deployment appsv1.Deployment
	if err := c.ShouldBindJSON(&deployment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	namespace := c.Param("namespace")
	result, err := K8sSvc.UpdateDeployment(namespace, &deployment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

// DeleteDeploymentHandler 删除 Deployment 的处理函数
func DeleteDeploymentHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	name := c.Param("name")
	err := K8sSvc.DeleteDeployment(namespace, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
