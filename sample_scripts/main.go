package main

import (
	"fmt"
	// "time"
	"github.com/indkumar8999/ps-tasks/client"
)

func main() {

	// Create a new client
	c, err := client.NewClient("localhost:50051")
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}
	defer c.Close()

	// Get unleased task
	taskresp, err := c.GetUnLeasdTask()
	if err != nil {
		fmt.Printf("Error getting unleased task: %v\n", err)
		return
	}
	fmt.Printf("Unleased task: %v\n", taskresp)
	task := taskresp.Task

	// Lease a task
	lease, err := c.LeaseTask(task.Id, 10)
	if err != nil {
		fmt.Printf("Error leasing task: %v\n", err)
		return
	}
	fmt.Printf("Leased task: %v\n", lease)

	// Complete the task
	completedTask, err := c.CompleteTask(task.Id)
	if err != nil {
		fmt.Printf("Error completing task: %v\n", err)
		return
	}
	fmt.Printf("Completed task: %v\n", completedTask)
	

	// Create a new task
	// task, err := c.CreateTask("name of task")
	// if err != nil {
	// 	fmt.Printf("Error creating task: %v\n", err)
	// 	return
	// }

	// fmt.Printf("Created task: %v\n", task)
	// // Get the task details


	// taskDetails, err := c.GetTask(task.Id)
	// if err != nil {
	// 	fmt.Printf("Error getting task: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Task details: %v\n", taskDetails)

	// // Lease a task
	// lease, err := c.LeaseTask(task.Id, 10)
	// if err != nil {
	// 	fmt.Printf("Error leasing task: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Leased task: %v\n", lease)

	// time.Sleep(10 * time.Second) // Simulate some processing time

	// // Complete the task
	// completedTask, err := c.CompleteTask(task.Id)
	// if err != nil {
	// 	fmt.Printf("Error completing task: %v\n", err)
	// 	return
	// }
	// fmt.Printf("Completed task: %v\n", completedTask)

}