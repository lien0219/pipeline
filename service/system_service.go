package service

import (
	"fmt"
	"gin_pipeline/global"
	"strconv"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemService 系统服务
type SystemService struct{}

// SystemStatus 系统状态信息
type SystemStatus struct {
	Uptime        string `json:"uptime"`
	CpuUsage      string `json:"cpu_usage"`
	MemoryUsage   string `json:"memory_usage"`
	DiskUsage     string `json:"disk_usage"`
	DbConnections int    `json:"db_connections"`
	ActiveUsers   int    `json:"active_users"`
}

// GetSystemStatus 获取系统状态
func (s *SystemService) GetSystemStatus() (SystemStatus, error) {
	// 系统运行时间
	uptime := time.Since(global.StartTime).Hours()
	uptimeStr := strconv.FormatFloat(uptime, 'f', 2, 64) + " 小时"

	// CPU使用率
	cpuPercent, _ := cpu.Percent(time.Second, false)
	cpuUsage := strconv.FormatFloat(cpuPercent[0], 'f', 2, 64) + "%"

	// 内存使用率
	memInfo, _ := mem.VirtualMemory()
	memUsage := strconv.FormatFloat(memInfo.UsedPercent, 'f', 2, 64) + "%"

	// 磁盘使用率
	diskInfo, _ := disk.Usage("/")
	diskUsage := strconv.FormatFloat(diskInfo.UsedPercent, 'f', 2, 64) + "%"

	// 数据库连接数
	sqlDB, err := global.DB.DB()
	if err != nil {
		return SystemStatus{}, fmt.Errorf("获取数据库连接失败: %v", err)
	}
	dbStats := sqlDB.Stats()

	// 活跃用户数 (这里需要根据实际业务逻辑实现)
	activeUsers := 0

	return SystemStatus{
		Uptime:        uptimeStr,
		CpuUsage:      cpuUsage,
		MemoryUsage:   memUsage,
		DiskUsage:     diskUsage,
		DbConnections: dbStats.OpenConnections,
		ActiveUsers:   activeUsers,
	}, nil
}
