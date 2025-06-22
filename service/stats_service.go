package service

import (
	"gin_pipeline/global"
	"gin_pipeline/model"
	"gin_pipeline/model/response"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// StatsService 统计服务
type StatsService struct{}

// NewStatsService 创建统计服务
func NewStatsService() *StatsService {
	return &StatsService{}
}

// GetPipelineRunStatsByDate 获取指定天数内的流水线运行统计
func (s *StatsService) GetPipelineRunStatsByDate(days int) (*response.PipelineChartData, error) {
	// 生成日期列表
	dates := make([]string, days)
	today := time.Now().UTC().Truncate(24 * time.Hour)
	for i := 0; i < days; i++ {
		date := today.AddDate(0, 0, -(days - 1 - i))
		dates[i] = date.Format("2006-01-02")
	}

	// 查询统计数据
	type statResult struct {
		Date   string `gorm:"column:date"`
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}
	var results []statResult

	startDate := today.AddDate(0, 0, -days+1).Format("2006-01-02")
	endDate := today.Format("2006-01-02")
	global.Log.Info("统计查询条件(UTC)", zap.String("startDate", startDate), zap.String("endDate", endDate))
	if err := global.DB.Model(&model.PipelineRun{}).
		Select("DATE(start_time) as date, status, count(*) as count").
		Where("DATE(start_time) BETWEEN ? AND ?", startDate, endDate).
		Group("date, status").
		Scan(&results).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// 初始化数据结构
	data := &response.PipelineChartData{
		Dates:   dates,
		Success: make([]int64, days),
		Failed:  make([]int64, days),
		Running: make([]int64, days),
		Pending: make([]int64, days),
		Total:   make([]int64, days),
	}

	// 填充数据
	dateIndex := make(map[string]int)
	for i, date := range dates {
		dateIndex[date] = i
	}

	for _, item := range results {
		if idx, ok := dateIndex[item.Date]; ok {
			data.Total[idx] += item.Count
			switch item.Status {
			case "success":
				data.Success[idx] = item.Count
			case "failed":
				data.Failed[idx] = item.Count
			case "running":
				data.Running[idx] = item.Count
			case "pending":
				data.Pending[idx] = item.Count
			}
		}
	}

	return data, nil
}
