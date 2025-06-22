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
	global.Log.Info("生成的日期数组", zap.Strings("dates", dates))

	// 查询统计数据
	type statResult struct {
		Date   string `gorm:"column:date"`
		Status string `gorm:"column:status"`
		Count  int64  `gorm:"column:count"`
	}
	var results []statResult

	startDate := today.AddDate(0, 0, -(days - 1)).Format("2006-01-02")
	endDate := today.Format("2006-01-02")
	global.Log.Info("统计查询条件(UTC)", zap.String("startDate", startDate), zap.String("endDate", endDate))
	if err := global.DB.Model(&model.PipelineRun{}).
		Select("DATE(CONVERT_TZ(start_time, '+00:00', '+00:00')) as date, status, count(*) as count").
		Where("DATE(start_time) BETWEEN ? AND ?", startDate, endDate).
		Group("date, status").
		Scan(&results).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	global.Log.Info("查询到的统计结果", zap.Any("results", results))

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
	dateIndex := make(map[time.Time]int)
	for i, dateStr := range dates {
		date, _ := time.Parse("2006-01-02", dateStr)
		dateIndex[date] = i
	}

	for _, item := range results {
		itemDate, err := time.Parse("2006-01-02T15:04:05Z", item.Date)
		if err != nil {
			global.Log.Error("解析日期失败", zap.String("date", item.Date), zap.Error(err))
			continue
		}
		global.Log.Info("处理统计项", zap.String("date", item.Date), zap.String("status", item.Status), zap.Int64("count", item.Count))
		if idx, ok := dateIndex[itemDate]; ok {
			global.Log.Info("找到匹配的日期索引", zap.String("date", item.Date), zap.Int("index", idx))
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
			default:
				continue
			}
		} else {
			global.Log.Warn("未找到匹配的日期索引", zap.String("date", item.Date))
		}
	}
	global.Log.Info("日期索引映射", zap.Any("dateIndex", dateIndex))

	return data, nil
}
