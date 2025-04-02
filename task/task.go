package task

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Task struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	State string `json:"state"`
	Data []byte `json:"data"`
	Metadata map[string]string `json:"metadata"`
}


func NewTask(id string, name string, description string,
	createdAt string, updatedAt string, state string, data []byte,
	metadata map[string]string) *Task {
	return &Task{
		ID: id,
		Name: name,
		Description: description,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		State: state,
		Data: data,
		Metadata: metadata,
	}
}


func (t *Task) GetID() string {
	return t.ID
}

func (t *Task) GetName() string {
	return t.Name
}
func (t *Task) GetDescription() string {
	return t.Description
}

func (t *Task) GetCreatedAt() string {
	return t.CreatedAt
}

func (t *Task) GetUpdatedAt() string {
	return t.UpdatedAt
}

func (t *Task) GetState() string {
	return t.State
}

func (t *Task) GetData() []byte {
	return t.Data
}

func (t *Task) GetMetadata() map[string]string {
	return t.Metadata
}

func (t *Task) Save(taskDir string) error {
	// Implement the logic to save the task to a file
	
	// For example, you can use JSON encoding to save the task to a file
	taskFile := filepath.Join(taskDir, fmt.Sprintf("%s.json", t.ID))
	file, err := os.Create(taskFile)
	if err != nil {
		return err
	}
	defer file.Close()
	// Serialize the task to JSON and write to the file
	if err := json.NewEncoder(file).Encode(t); err != nil {
		return err
	}
	return nil
}

func LoadTask(taskDir string, taskID string) (*Task, error) {
	// Implement the logic to load the task from a file
	taskFile := filepath.Join(taskDir, fmt.Sprintf("%s", taskID))
	file, err := os.Open(taskFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var task Task
	if err := json.NewDecoder(file).Decode(&task); err != nil {
		return nil, err
	}
	return &task, nil
}