package service

import (
	"context"
	"errors"
	"gin_pipeline/global"
	"go.uber.org/zap"
	"sync"
	"time"
)

// TaskExecutor 任务执行器接口
type TaskExecutor interface {
	Execute(ctx context.Context, task *WorkflowTask) error
}

// WorkflowTask 工作流任务
type WorkflowTask struct {
	ID           string
	Name         string
	Type         string
	Config       map[string]interface{}
	Dependencies []string
	Status       string // pending, running, success, failed, canceled
	StartTime    *time.Time
	EndTime      *time.Time
	Logs         string
	Error        string
}

// WorkflowEngine 工作流引擎
type WorkflowEngine struct {
	executors map[string]TaskExecutor
	mutex     sync.RWMutex
}

// NewWorkflowEngine 创建工作流引擎
func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{
		executors: make(map[string]TaskExecutor),
	}
}

// RegisterExecutor 注册任务执行器
func (e *WorkflowEngine) RegisterExecutor(taskType string, executor TaskExecutor) {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.executors[taskType] = executor
}

// GetExecutor 获取任务执行器
func (e *WorkflowEngine) GetExecutor(taskType string) (TaskExecutor, error) {
	e.mutex.RLock()
	defer e.mutex.RUnlock()
	executor, ok := e.executors[taskType]
	if !ok {
		return nil, errors.New("未找到任务类型的执行器: " + taskType)
	}
	return executor, nil
}

// ExecuteWorkflow 执行工作流
func (e *WorkflowEngine) ExecuteWorkflow(ctx context.Context, tasks []*WorkflowTask, runID uint) error {
	// 构建任务依赖图
	taskMap := make(map[string]*WorkflowTask)
	for _, task := range tasks {
		taskMap[task.ID] = task
		task.Status = "pending"
	}

	// 验证依赖关系
	for _, task := range tasks {
		for _, depID := range task.Dependencies {
			if _, ok := taskMap[depID]; !ok {
				return errors.New("任务依赖不存在: " + depID)
			}
		}
	}

	// 创建取消上下文
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	// 创建等待组
	var wg sync.WaitGroup
	// 创建错误通道
	errChan := make(chan error, len(tasks))
	// 创建任务完成通道
	doneChan := make(chan string, len(tasks))

	// 跟踪已完成的任务
	completedTasks := make(map[string]bool)
	var completedMutex sync.Mutex

	// 启动任务监控协程
	go func() {
		for {
			select {
			case <-cancelCtx.Done():
				return
			case taskID := <-doneChan:
				completedMutex.Lock()
				completedTasks[taskID] = true
				completedMutex.Unlock()

				// 检查是否有新任务可以启动
				for _, task := range tasks {
					if task.Status != "pending" {
						continue
					}

					// 检查依赖是否都已完成
					allDepsCompleted := true
					for _, depID := range task.Dependencies {
						completedMutex.Lock()
						completed := completedTasks[depID]
						completedMutex.Unlock()

						if !completed {
							allDepsCompleted = false
							break
						}
					}

					if allDepsCompleted {
						// 启动任务
						wg.Add(1)
						go e.executeTask(cancelCtx, task, &wg, errChan, doneChan, runID)
					}
				}
			}
		}
	}()

	// 启动没有依赖的任务
	for _, task := range tasks {
		if len(task.Dependencies) == 0 {
			wg.Add(1)
			go e.executeTask(cancelCtx, task, &wg, errChan, doneChan, runID)
		}
	}

	// 等待所有任务完成或出错
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// 收集错误
	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
			cancel() // 取消所有任务
		}
	}

	if len(errs) > 0 {
		return errors.New("工作流执行失败")
	}

	return nil
}

// executeTask 执行单个任务
func (e *WorkflowEngine) executeTask(ctx context.Context, task *WorkflowTask, wg *sync.WaitGroup, errChan chan<- error, doneChan chan<- string, runID uint) {
	defer wg.Done()

	// 更新任务状态为运行中
	task.Status = "running"
	now := time.Now()
	task.StartTime = &now

	// 更新数据库中的任务状态
	e.updateTaskStatus(runID, task.ID, "running", "", "")

	// 获取任务执行器
	executor, err := e.GetExecutor(task.Type)
	if err != nil {
		task.Status = "failed"
		task.Error = err.Error()
		errChan <- err
		e.updateTaskStatus(runID, task.ID, "failed", "", err.Error())
		return
	}

	// 执行任务
	err = executor.Execute(ctx, task)
	endTime := time.Now()
	task.EndTime = &endTime

	if err != nil {
		task.Status = "failed"
		task.Error = err.Error()
		errChan <- err
		e.updateTaskStatus(runID, task.ID, "failed", task.Logs, err.Error())
	} else {
		task.Status = "success"
		e.updateTaskStatus(runID, task.ID, "success", task.Logs, "")
	}

	// 通知任务完成
	doneChan <- task.ID
}

// updateTaskStatus 更新任务状态
func (e *WorkflowEngine) updateTaskStatus(runID uint, taskID string, status string, logs string, errMsg string) {
	// 在实际项目中，这里应该更新数据库中的任务状态
	// 例如，可以将任务状态和日志保存到pipeline_run_tasks表中
	global.Log.Info("更新任务状态",
		zap.Uint("runID", runID),
		zap.String("taskID", taskID),
		zap.String("status", status),
		zap.String("error", errMsg))
}

// CancelWorkflow 取消工作流
func (e *WorkflowEngine) CancelWorkflow(runID uint) error {
	// 在实际项目中，这里应该更新数据库中的工作流状态为取消
	// 并且通知正在运行的任务取消
	global.Log.Info("取消工作流", zap.Uint("runID", runID))
	return nil
}
