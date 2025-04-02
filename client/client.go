package client

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"github.com/indkumar8999/ps-tasks/service/taskpb"
)

// Client struct for interacting with the gRPC task server
type Client struct {
	conn   *grpc.ClientConn
	client taskpb.TaskServiceClient
}

// NewClient initializes a connection to the gRPC server
func NewClient(serverAddr string) (*Client, error) {
	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}

	client := taskpb.NewTaskServiceClient(conn)
	return &Client{conn: conn, client: client}, nil
}

// Close closes the gRPC connection
func (c *Client) Close() {
	c.conn.Close()
}

// CreateTask creates a new task and returns the task details
func (c *Client) CreateTask(name string) (*taskpb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.CreateTask(ctx, &taskpb.CreateTaskRequest{
		Name: name,
		Description: "task description",
		Data: []byte("task data"),
	})
	if err != nil {
		return nil, fmt.Errorf("error creating task: %w", err)
	}
	return resp.Task, nil
}

// GetTask fetches task details by ID
func (c *Client) GetTask(taskID string) (*taskpb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.GetTask(ctx, &taskpb.GetTaskRequest{Id: taskID})
	if err != nil {
		return nil, fmt.Errorf("error getting task: %w", err)
	}
	return resp.Task, nil
}

// // UpdateTask updates the state of a task
// func (c *Client) UpdateTask(taskID, state string) (*taskpb.Task, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
// 	defer cancel()

// 	resp, err := c.client.UpdateTask(ctx, &taskpb.UpdateTaskRequest{Id: taskID, State: state})
// 	if err != nil {
// 		return nil, fmt.Errorf("error updating task: %w", err)
// 	}
// 	return resp.Task, nil
// }

// CompleteTask marks a task as completed
func (c *Client) CompleteTask(taskID string) (*taskpb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.CompleteTask(ctx, &taskpb.CompleteTaskRequest{Id: taskID})
	if err != nil {
		return nil, fmt.Errorf("error completing task: %w", err)
	}
	return resp.Task, nil
}

// LeaseTask leases a task for processing
func (c *Client) LeaseTask(taskID string, leaseDuration int32) (*taskpb.LeaseTaskResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.LeaseTask(ctx, &taskpb.LeaseTaskRequest{TaskId: taskID, Owner: "owner1"})
	if err != nil {
		return nil, fmt.Errorf("error leasing task: %w", err)
	}
	return resp, nil
}

func (c *Client) GetUnLeasdTask() (*taskpb.TaskResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := c.client.GetUnLeasdTask(ctx, &taskpb.UnLeasedTaskRequest{})
	if err != nil {
		return nil, fmt.Errorf("error getting unleased task: %w", err)
	}
	return resp, nil
}