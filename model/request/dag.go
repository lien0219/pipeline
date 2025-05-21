package request

import "gin_pipeline/model"

// CreateDAG 创建DAG请求参数
type CreateDAG struct {
	Name        string          `json:"name" binding:"required,min=2,max=100"`
	Description string          `json:"description"`
	PipelineID  uint            `json:"pipeline_id" binding:"required"`
	Nodes       []model.DAGNode `json:"nodes" binding:"required"`
}

// UpdateDAG 更新DAG请求参数
type UpdateDAG struct {
	Name        string          `json:"name" binding:"required,min=2,max=100"`
	Description string          `json:"description"`
	Nodes       []model.DAGNode `json:"nodes" binding:"required"`
}

// ValidateDAG 验证DAG请求参数
type ValidateDAG struct {
	Nodes []model.DAGNode `json:"nodes" binding:"required"`
}
