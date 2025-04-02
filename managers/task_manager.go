package managers


import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"github.com/google/uuid"
	"sync"
	"github.com/indkumar8999/ps-tasks/task"
	"github.com/indkumar8999/ps-tasks/leases"
)

// TaskManager manages tasks and leases

type TaskManager struct {
	tasksDir string
	tasks     map[string]*task.Task
	leaseManager *LeaseManager
	taskLock  *sync.Mutex
}

const (
	CREATED = "created"
	RUNNING = "running"
	FAILED  = "failed"
	ABORTED = "aborted"
	PAUSED  = "paused"
	RESUMED = "resumed"
	STARTED = "started"
	STOPPED = "stopped"
	COMPLETED = "completed"
)

// NewTaskManager creates a new TaskManager
func NewTaskManager(tasksDir string, leaseManager *LeaseManager) *TaskManager {
	return &TaskManager{
		tasksDir:    tasksDir,
		tasks:       make(map[string]*task.Task),
		leaseManager: leaseManager,
		taskLock:    &sync.Mutex{},
	}
}

func (tm *TaskManager) LoadTasks() {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	files, err := os.ReadDir(tm.tasksDir)
	if err != nil {
		fmt.Printf("Error reading tasks directory: %v\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		taskID := file.Name()
		task, err := task.LoadTask(tm.tasksDir, taskID)
		if err != nil {
			fmt.Printf("Error loading task: %v\n", err)
			continue
		}
		tm.tasks[task.ID] = task
	}
}

func (tm *TaskManager) PeriodicallyDeleteTasks() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			threshold := time.Now().Add(-24 * time.Hour) // 24 hours threshold
			if err := tm.DeleteOlderTasks(threshold); err != nil {
				fmt.Printf("Error deleting older tasks: %v\n", err)
			}
		}
	}
}

// CreateTask creates a new task
func (tm *TaskManager) CreateTask(name string, description string, data []byte, metadata map[string]string) (*task.Task, error) {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	// Generate a unique ID for the task
	taskID := uuid.New().String()

	// Create a new task
	newTask := task.NewTask(taskID, name, description, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339), CREATED, data, metadata)

	// Save the task to the tasks directory
	if err := newTask.Save(tm.tasksDir); err != nil {
		return nil, fmt.Errorf("failed to save task: %v", err)
	}

	// Add the task to the in-memory map
	tm.tasks[taskID] = newTask

	return newTask, nil
}


// GetTask retrieves a task by ID
func (tm *TaskManager) GetTask(taskID string) (*task.Task, error) {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	// Check if the task exists in the in-memory map
	if task, exists := tm.tasks[taskID]; exists {
		return task, nil
	}

	// If not found in memory, load from disk
	task, err := task.LoadTask(tm.tasksDir, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to load task: %v", err)
	}

	// Add the loaded task to the in-memory map
	tm.tasks[taskID] = task

	return task, nil
}

// UpdateTask updates a task by ID
func (tm *TaskManager) UpdateTask(taskID string, taskState string, data []byte) (*task.Task, error) {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	// Check if the task exists
	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	// Update the task fields
	task.Data = data
	task.State = taskState
	task.UpdatedAt = time.Now().Format(time.RFC3339)

	// Save the updated task to disk
	if err := task.Save(tm.tasksDir); err != nil {
		return nil, fmt.Errorf("failed to save updated task: %v", err)
	}

	return task, nil
}


// CompleteTask marks a task as completed
func (tm *TaskManager) CompleteTask(taskID string) (*task.Task, error) {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()
	// Check if the task exists
	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}
	// Mark the task as completed
	task.State = COMPLETED
	task.UpdatedAt = time.Now().Format(time.RFC3339)
	// Save the updated task to disk
	if err := task.Save(tm.tasksDir); err != nil {
		return nil, fmt.Errorf("failed to save updated task: %v", err)
	}
	// Remove the task from the in-memory map
	delete(tm.tasks, taskID)
	return task, nil
}
	

// DeleteTask deletes a task by ID
func (tm *TaskManager) DeleteTask(taskID string) error {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	// Check if the task exists
	if _, exists := tm.tasks[taskID]; !exists {
		return fmt.Errorf("task not found")
	}

	// Delete the task from the in-memory map
	delete(tm.tasks, taskID)

	// Delete the task file from disk
	taskFile := filepath.Join(tm.tasksDir, fmt.Sprintf("%s.json", taskID))
	if err := os.Remove(taskFile); err != nil {
		return fmt.Errorf("failed to delete task file: %v", err)
	}

	return nil
}

func (tm *TaskManager) GetUnLeasdTask() (*task.Task, error) {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	for _, task := range tm.tasks {
		if task.State == CREATED {
			return task, nil
		}
	}

	return nil, fmt.Errorf("no unleased tasks available")
}


// delete older tasks
func (tm *TaskManager) DeleteOlderTasks(threshold time.Time) error {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	files, err := os.ReadDir(tm.tasksDir)
	if err != nil {
		return fmt.Errorf("failed to read tasks directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		taskID := file.Name()
		task, err := task.LoadTask(tm.tasksDir, taskID)
		if err != nil {
			return fmt.Errorf("failed to load task: %v", err)
		}
		// Check if the task is older than the threshold
		taskTime, err := time.Parse(time.RFC3339, task.CreatedAt)
		if taskTime.Before(threshold) {
			if err := os.Remove(filepath.Join(tm.tasksDir, fmt.Sprintf("%s.json", task.ID))); err != nil {
				return fmt.Errorf("failed to delete task file: %v", err)
			}
			delete(tm.tasks, task.ID)
		}
	}
	return nil
}

func (tm *TaskManager) LeaseTask(taskID string, username string) (*leases.Lease, error) {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	// Check if the task exists
	task, exists := tm.tasks[taskID]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	// Create a new lease for the task
	lease, err := tm.leaseManager.AcquireLease(task.ID, time.Minute*3, username)
	if err = lease.Save(tm.leaseManager.leasesDir); err != nil {
		return nil, fmt.Errorf("failed to save lease: %v", err)
	}

	return lease, nil
}

func (tm *TaskManager) GetUnLeasedTask() (*task.Task, error) {
	tm.taskLock.Lock()
	defer tm.taskLock.Unlock()

	for _, task := range tm.tasks {
		if task.State == CREATED {
			return task, nil
		}
	}

	return nil, fmt.Errorf("no unleased tasks available")
}