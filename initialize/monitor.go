package initialize

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	WorkflowDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "workflow_execution_duration_seconds",
			Help:    "工作流执行耗时分布",
			Buckets: []float64{1, 5, 10, 30, 60, 120},
		},
		[]string{"status"},
	)

	WorkflowCount = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "workflow_execution_total",
			Help: "工作流执行总次数",
		},
		[]string{"status"},
	)
)
