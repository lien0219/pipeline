package model

// PipelineStats 流水线统计数据
type PipelineStats struct {
	Success int64 `gorm:"column:success"`
	Running int64 `gorm:"column:running"`
	Failed  int64 `gorm:"column:failed"`
	Pending int64 `gorm:"column:pending"`
}

// TableName 指定表名
func (PipelineStats) TableName() string {
	return "pipeline_stats"
}
