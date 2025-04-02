package service


import (
	"fmt"
	"time"
	"github.com/indkumar8999/ps-tasks/managers"
	// "github.com/indkumar8999/ps-tasks/task"

	"context"
	"github.com/indkumar8999/ps-tasks/service/taskpb"
)


type TaskService struct {
	taskpb.UnimplementedTaskServiceServer
	leaseManager *managers.LeaseManager
	taskManager *managers.TaskManager
}

// NewTaskService creates a new TaskService
func NewTaskService(leaseManager *managers.LeaseManager, taskManager *managers.TaskManager) *TaskService {
	return &TaskService{
		leaseManager: leaseManager,
		taskManager:  taskManager,
	}
}

func (s *TaskService) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.TaskResponse, error) {
	task1, err := s.taskManager.CreateTask("", "", req.Data, nil);
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %v", err)
	}
	taskProto := &taskpb.Task{
		Id:          task1.ID,
		TaskState:   "CREATED",
		Data:        req.Data,
	}

	return &taskpb.TaskResponse{Task: taskProto}, nil
}

func (s *TaskService) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.TaskResponse, error) {
	task, err := s.taskManager.GetTask(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %v", err)
	}
	taskProto := &taskpb.Task{
		Id:          task.ID,
		TaskState:   task.State,
		Data:        task.Data,
	}

	return &taskpb.TaskResponse{Task: taskProto}, nil
}

func (s *TaskService) CompleteTask(ctx context.Context, req *taskpb.CompleteTaskRequest) (*taskpb.TaskResponse, error) {
	task, err := s.taskManager.CompleteTask(req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to complete task: %v", err)
	}
	taskProto := &taskpb.Task{
		Id:          task.ID,
		TaskState:   task.State,
		Data:        task.Data,
	}

	return &taskpb.TaskResponse{Task: taskProto}, nil
}

func (s *TaskService) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.TaskResponse, error) {
	task, err := s.taskManager.UpdateTask(req.Id, req.TaskState, req.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to update task: %v", err)
	}
	taskProto := &taskpb.Task{
		Id:          task.ID,
		TaskState:   task.State,
		Data:        task.Data,
	}

	return &taskpb.TaskResponse{Task: taskProto}, nil
}

func (s *TaskService) LeaseTask(ctx context.Context, req *taskpb.LeaseTaskRequest) (*taskpb.LeaseTaskResponse, error) {
	lease, err := s.taskManager.LeaseTask(req.TaskId, req.Owner)
	if err != nil {
		return nil, fmt.Errorf("failed to lease task: %v", err)
	}
	response := &taskpb.LeaseTaskResponse{
		Id: lease.ID,
		TaskId: lease.TaskID,
		LeaseEndTime: lease.ExpiresAt.Format(time.RFC3339),
	}

	return response, nil
}

func (s *TaskService) GetUnLeasdTask(ctx context.Context, req *taskpb.UnLeasedTaskRequest) (*taskpb.TaskResponse, error) {
	task, err := s.taskManager.GetUnLeasedTask()
	if err != nil {
		return nil, fmt.Errorf("failed to get unleased task: %v", err)
	}
	taskProto := &taskpb.Task{
		Id:          task.ID,
		TaskState:   task.State,
		Data:        task.Data,
	}

	return &taskpb.TaskResponse{Task: taskProto}, nil
}