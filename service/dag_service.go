package service

import (
	"errors"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"go.uber.org/zap"
)

// DAGService 提供DAG相关的服务
type DAGService struct{}

// CreateDAG 创建新的DAG
func (s *DAGService) CreateDAG(dag *model.DAG) error {
	// 验证DAG是否有效
	if err := s.ValidateDAG(dag.NodesData); err != nil {
		return err
	}

	// 如果是活动版本，将同一流水线的其他DAG设为非活动
	if dag.IsActive {
		if err := global.DB.Model(&model.DAG{}).
			Where("pipeline_id = ? AND is_active = ?", dag.PipelineID, true).
			Update("is_active", false).Error; err != nil {
			global.Log.Error("更新其他DAG状态失败", zap.Error(err))
			return err
		}
	}

	// 创建DAG
	return global.DB.Create(dag).Error
}

// GetDAGByID 根据ID获取DAG
func (s *DAGService) GetDAGByID(id uint) (*model.DAG, error) {
	var dag model.DAG
	if err := global.DB.Preload("Pipeline").Preload("Creator").First(&dag, id).Error; err != nil {
		return nil, err
	}
	return &dag, nil
}

// GetDAGsByPipelineID 获取指定流水线的所有DAG
func (s *DAGService) GetDAGsByPipelineID(pipelineID uint) ([]model.DAG, error) {
	var dags []model.DAG
	if err := global.DB.Where("pipeline_id = ?", pipelineID).
		Order("version DESC").
		Find(&dags).Error; err != nil {
		return nil, err
	}
	return dags, nil
}

// GetActiveDAGByPipelineID 获取指定流水线的活动DAG
func (s *DAGService) GetActiveDAGByPipelineID(pipelineID uint) (*model.DAG, error) {
	var dag model.DAG
	if err := global.DB.Where("pipeline_id = ? AND is_active = ?", pipelineID, true).
		First(&dag).Error; err != nil {
		return nil, err
	}
	return &dag, nil
}

// UpdateDAG 更新DAG
func (s *DAGService) UpdateDAG(id uint, updates map[string]interface{}) error {
	// 如果更新包含节点数据，需要验证
	if nodes, ok := updates["nodes_data"].([]model.DAGNode); ok {
		if err := s.ValidateDAG(nodes); err != nil {
			return err
		}
	}

	// 如果设置为活动版本，将同一流水线的其他DAG设为非活动
	if isActive, ok := updates["is_active"].(bool); ok && isActive {
		var dag model.DAG
		if err := global.DB.First(&dag, id).Error; err != nil {
			return err
		}

		if err := global.DB.Model(&model.DAG{}).
			Where("pipeline_id = ? AND id != ? AND is_active = ?", dag.PipelineID, id, true).
			Update("is_active", false).Error; err != nil {
			global.Log.Error("更新其他DAG状态失败", zap.Error(err))
			return err
		}
	}

	return global.DB.Model(&model.DAG{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteDAG 删除DAG
func (s *DAGService) DeleteDAG(id uint) error {
	return global.DB.Delete(&model.DAG{}, id).Error
}

// CreateDAGVersion 创建DAG的新版本
func (s *DAGService) CreateDAGVersion(dagID uint, creatorID uint) (*model.DAG, error) {
	// 获取原DAG
	var sourceDag model.DAG
	if err := global.DB.First(&sourceDag, dagID).Error; err != nil {
		return nil, err
	}

	// 获取当前最高版本
	var maxVersion int
	if err := global.DB.Model(&model.DAG{}).
		Where("pipeline_id = ?", sourceDag.PipelineID).
		Select("COALESCE(MAX(version), 0)").
		Scan(&maxVersion).Error; err != nil {
		return nil, err
	}

	// 创建新版本
	newDag := model.DAG{
		Name:        sourceDag.Name,
		Description: sourceDag.Description,
		Version:     maxVersion + 1,
		PipelineID:  sourceDag.PipelineID,
		NodesData:   sourceDag.NodesData,
		CreatorID:   creatorID,
		IsActive:    false, // 默认不激活
	}

	if err := global.DB.Create(&newDag).Error; err != nil {
		return nil, err
	}

	return &newDag, nil
}

// ValidateDAG 验证DAG是否有效（无环）
func (s *DAGService) ValidateDAG(nodes []model.DAGNode) error {
	if len(nodes) == 0 {
		return errors.New("DAG不能为空")
	}

	// 构建节点映射
	nodeMap := make(map[string]model.DAGNode)
	for _, node := range nodes {
		nodeMap[node.ID] = node
	}

	// 检查所有依赖是否存在
	for _, node := range nodes {
		for _, depID := range node.Dependencies {
			if _, exists := nodeMap[depID]; !exists {
				return errors.New("依赖节点不存在: " + depID)
			}
		}
	}

	// 检测环
	visited := make(map[string]bool)
	path := make(map[string]bool)

	var dfs func(string) bool
	dfs = func(nodeID string) bool {
		// 如果节点已经在当前路径中，说明有环
		if path[nodeID] {
			return true
		}

		// 如果节点已经被访问过且没有环，则跳过
		if visited[nodeID] {
			return false
		}

		// 标记节点为已访问，并添加到当前路径
		visited[nodeID] = true
		path[nodeID] = true

		// 递归检查所有依赖节点
		node := nodeMap[nodeID]
		for _, depID := range node.Dependencies {
			if dfs(depID) {
				return true
			}
		}

		// 从当前路径中移除节点
		path[nodeID] = false
		return false
	}

	// 对每个节点进行DFS
	for _, node := range nodes {
		if dfs(node.ID) {
			return errors.New("DAG中存在环")
		}
	}

	return nil
}

// GetDAGHistory 获取DAG的历史版本
func (s *DAGService) GetDAGHistory(pipelineID uint) ([]model.DAG, error) {
	var dags []model.DAG
	if err := global.DB.Where("pipeline_id = ?", pipelineID).
		Order("version DESC").
		Find(&dags).Error; err != nil {
		return nil, err
	}
	return dags, nil
}

// ActivateDAG 激活指定的DAG版本
func (s *DAGService) ActivateDAG(dagID uint) error {
	// 开启事务
	tx := global.DB.Begin()

	// 获取DAG
	var dag model.DAG
	if err := tx.First(&dag, dagID).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 将同一流水线的其他DAG设为非活动
	if err := tx.Model(&model.DAG{}).
		Where("pipeline_id = ? AND id != ?", dag.PipelineID, dagID).
		Update("is_active", false).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 将当前DAG设为活动
	if err := tx.Model(&dag).Update("is_active", true).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}
