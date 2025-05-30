package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"gin_pipeline/global"
	"gin_pipeline/model"
	"time"

	"go.uber.org/zap"
)

// WorkflowService 工作流服务
type WorkflowService struct {
	engine    *WorkflowEngine
	semaphore chan struct{} //并发控制信号量
}

// NewWorkflowService 创建工作流服务
func NewWorkflowService() *WorkflowService {
	engine := NewWorkflowEngine()

	// 注册任务执行器
	engine.RegisterExecutor("shell", &ShellTaskExecutor{})
	engine.RegisterExecutor("docker", &DockerTaskExecutor{})
	engine.RegisterExecutor("kubernetes", &KubernetesTaskExecutor{})

	return &WorkflowService{
		engine:    engine,
		semaphore: make(chan struct{}, 20), //限制并发数为20
	}
}

// TriggerWorkflow 触发工作流
func (s *WorkflowService) TriggerWorkflow(pipelineID uint, userID uint, gitBranch string) (*model.PipelineRun, error) {
	// 获取流水线
	var pipeline model.Pipeline
	if err := global.DB.First(&pipeline, pipelineID).Error; err != nil {
		global.Log.Error("获取流水线失败", zap.Error(err))
		return nil, err
	}

	// 获取活动DAG
	dagService := new(DAGService)
	dag, err := dagService.GetActiveDAGByPipelineID(pipelineID)
	if err != nil {
		global.Log.Error("获取活动DAG失败", zap.Error(err))
		return nil, err
	}

	// 创建流水线运行记录
	now := time.Now()
	pipelineRun := model.PipelineRun{
		PipelineID: pipelineID,
		Status:     "pending",
		StartTime:  &now,
		GitBranch:  gitBranch,
		TriggerBy:  userID,
	}

	if err := global.DB.Create(&pipelineRun).Error; err != nil {
		global.Log.Error("创建流水线运行记录失败", zap.Error(err))
		return nil, err
	}

	// 更新流水线状态
	if err := global.DB.Model(&pipeline).Updates(map[string]interface{}{
		"status":      "running",
		"last_run_at": now,
	}).Error; err != nil {
		global.Log.Error("更新流水线状态失败", zap.Error(err))
		// 不影响结果，继续执行
	}

	// 异步执行工作流
	// go s.executeWorkflow(dag, &pipelineRun)
	select {
	case s.semaphore <- struct{}{}:
		go func() {
			defer func() { <-s.semaphore }()
			s.executeWorkflow(dag, &pipelineRun)
		}()
	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("系统繁忙，请稍后重试")
	}

	return &pipelineRun, nil
}

// executeWorkflow 执行工作流
func (s *WorkflowService) executeWorkflow(dag *model.DAG, pipelineRun *model.PipelineRun) {
	// 更新运行状态为运行中
	if err := global.DB.Model(pipelineRun).Update("status", "running").Error; err != nil {
		global.Log.Error("更新流水线运行状态失败", zap.Error(err))
		return
	}

	// 将DAG节点转换为工作流任务
	var tasks []*WorkflowTask
	for _, node := range dag.NodesData {
		config := make(map[string]interface{})
		if node.Config != nil {
			config = node.Config
		}

		task := &WorkflowTask{
			ID:           node.ID,
			Name:         node.Name,
			Type:         node.Type,
			Config:       config,
			Dependencies: node.Dependencies,
			Status:       "pending",
		}
		tasks = append(tasks, task)
	}

	// 执行工作流
	// ctx := context.Background()
	// err := s.engine.ExecuteWorkflow(ctx, tasks, pipelineRun.ID)
	// 添加执行超时控制
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	// 带重试的执行逻辑
	var err error
	for retry := 0; retry < 3; retry++ {
		err = s.engine.ExecuteWorkflow(ctx, tasks, pipelineRun.ID)
		if err == nil {
			break
		}
		time.Sleep(time.Duration(retry+1) * 5 * time.Second) // 指数退避
	}

	// 更新运行结果
	now := time.Now()
	duration := int(now.Sub(*pipelineRun.StartTime).Seconds())
	status := "success"
	if err != nil {
		status = "failed"
		global.Log.Error("工作流执行失败", zap.Error(err))
	}

	// 收集任务日志
	// taskLogs := make(map[string]map[string]string)
	// for _, task := range tasks {
	// 	if taskLogs[task.Type] == nil {
	// 		taskLogs[task.Type] = make(map[string]string)
	// 	}
	// 	taskLogs[task.Type][task.ID] = task.Logs
	// }
	// logsJSON, _ := json.Marshal(taskLogs)

	updates := map[string]interface{}{
		"status":   status,
		"end_time": now,
		"duration": duration,
		"logs":     "",
	}
	// 流式处理日志收集
	var logsBuffer bytes.Buffer
	encoder := json.NewEncoder(&logsBuffer)
	encoder.SetEscapeHTML(false)

	logsBuffer.WriteString(`{"tasks":{`)
	first := true
	for _, task := range tasks {
		if !first {
			logsBuffer.WriteByte(',')
		}
		first = false
		encoder.Encode(task.Type + `":"` + task.ID + `":`)
		encoder.Encode(task.Logs)
	}
	logsBuffer.WriteString(`}}`)
	updates["logs"] = logsBuffer.String()

	// 批量更新任务状态
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(pipelineRun).Updates(updates).Error; err != nil {
		tx.Rollback()
		global.Log.Error("更新流水线运行结果失败", zap.Error(err))
		return
	}

	if err := tx.Model(&model.Pipeline{}).Where("id = ?", pipelineRun.PipelineID).Update("status", status).Error; err != nil {
		tx.Rollback()
		global.Log.Error("更新流水线状态失败", zap.Error(err))
		return
	}

	if err := tx.Commit().Error; err != nil {
		global.Log.Error("提交事务失败", zap.Error(err))
		return
	}

	// global.Log.Info("工作流执行完成",
	// 	zap.Uint("runID", pipelineRun.ID),
	// 	zap.String("status", status))
	startTime := time.Now()
	logFields := []zap.Field{
		zap.Uint("pipelineID", pipelineRun.PipelineID),
		zap.Uint("runID", pipelineRun.ID),
		zap.String("status", pipelineRun.Status),
	}
	defer func() {
		duration := time.Since(startTime)
		logFields = append(logFields, zap.Duration("duration", duration))
		global.Log.Info("工作流执行完成", logFields)
	}()

}

// CancelWorkflow 取消工作流
func (s *WorkflowService) CancelWorkflow(runID uint) error {
	// 获取运行记录
	var run model.PipelineRun
	if err := global.DB.First(&run, runID).Error; err != nil {
		global.Log.Error("获取流水线运行记录失败", zap.Error(err))
		return err
	}

	// 检查状态
	if run.Status != "pending" && run.Status != "running" {
		return nil
	}

	// 取消工作流
	if err := s.engine.CancelWorkflow(runID); err != nil {
		global.Log.Error("取消工作流失败", zap.Error(err))
		return err
	}

	// 更新状态
	now := time.Now()
	updates := map[string]interface{}{
		"status":   "canceled",
		"end_time": now,
	}

	// 如果有开始时间，计算持续时间
	if run.StartTime != nil {
		updates["duration"] = int(now.Sub(*run.StartTime).Seconds())
	}

	if err := global.DB.Model(&run).Updates(updates).Error; err != nil {
		global.Log.Error("更新流水线运行状态失败", zap.Error(err))
		return err
	}

	return nil
}

// ShellTaskExecutor Shell任务执行器
type ShellTaskExecutor struct{}

// Execute 执行Shell任务
func (e *ShellTaskExecutor) Execute(ctx context.Context, task *WorkflowTask) error {
	// 模拟执行Shell命令
	global.Log.Info("执行Shell任务",
		zap.String("taskID", task.ID),
		zap.String("name", task.Name))

	// 模拟执行时间
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(2 * time.Second):
		// 模拟执行成功
		task.Logs = "Shell任务执行成功\n$ echo 'Hello World'\nHello World"
		return nil
	}
}

// DockerTaskExecutor Docker任务执行器
type DockerTaskExecutor struct{}

// Execute 执行Docker任务
func (e *DockerTaskExecutor) Execute(ctx context.Context, task *WorkflowTask) error {
	// 模拟执行Docker命令
	global.Log.Info("执行Docker任务",
		zap.String("taskID", task.ID),
		zap.String("name", task.Name))

	// 模拟执行时间
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(3 * time.Second):
		// 模拟执行成功
		task.Logs = "Docker任务执行成功\n$ docker run hello-world\nHello from Docker!"
		return nil
	}
}

// KubernetesTaskExecutor Kubernetes任务执行器
type KubernetesTaskExecutor struct{}

// Execute 执行Kubernetes任务
func (e *KubernetesTaskExecutor) Execute(ctx context.Context, task *WorkflowTask) error {
	// 模拟执行Kubernetes命令
	global.Log.Info("执行Kubernetes任务",
		zap.String("taskID", task.ID),
		zap.String("name", task.Name))

	// 模拟执行时间
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(4 * time.Second):
		// 模拟执行成功
		task.Logs = "Kubernetes任务执行成功\n$ kubectl apply -f deployment.yaml\ndeployment.apps/nginx created"
		return nil
	}
}
