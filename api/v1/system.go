package v1

import (
	"gin_pipeline/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSystemConfig godoc
// @Summary      Get system configuration
// @Description  Get public system configuration
// @Tags         system
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /system/config [get]
func GetSystemConfig(c *gin.Context) {
	// 只返回安全的配置信息，不包含敏感信息
	c.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": "Success",
		"data": gin.H{
			"server": gin.H{
				"mode": config.Get().Server.Mode,
			},
			"database": gin.H{
				"driver": config.Get().Database.Driver,
			},
		},
	})
}
