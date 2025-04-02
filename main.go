package main

import (
	"fmt"
	"os"
	"path/filepath"
	"net"
	"log"
	"google.golang.org/grpc"
	"github.com/indkumar8999/ps-tasks/service/taskpb"
	"github.com/indkumar8999/ps-tasks/managers"
	"github.com/indkumar8999/ps-tasks/service"
)

const (
	LEASE_DIR = "leases"
	METADATA_DIR = "metadata"
	TASKS_DIR = "tasks"
)

func GetOrCreateDBPath(cwd string) string {
	dbPath := filepath.Join(cwd, "database")
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err := os.MkdirAll(dbPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating database directory:", err)
			return ""
		}
		fmt.Println("Database directory created:", dbPath)
	} else {
		fmt.Println("Database directory already exists:", dbPath)
	}
	return dbPath
}

func GetOrCreateMetadataPath(dbPath string) string {
	// Create the metadata folder if not exists
	metadataPath := filepath.Join(dbPath, "metadata")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		err := os.MkdirAll(metadataPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating metadata directory:", err)
			return ""
		}
		fmt.Println("Metadata directory created:", metadataPath)
	} else {
		fmt.Println("Metadata directory already exists:", metadataPath)
	}
	return metadataPath
}

func GetOrCreateLeasesPath(dbPath string) string {
	// Create the leases folder if not exists
	leasesPath := filepath.Join(dbPath, "leases")
	if _, err := os.Stat(leasesPath); os.IsNotExist(err) {
		err := os.MkdirAll(leasesPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating leases directory:", err)
			return ""
		}
		fmt.Println("Leases directory created:", leasesPath)
	} else {
		fmt.Println("Leases directory already exists:", leasesPath)
	}
	return leasesPath
}

func GetOrCreateTasksPath(dbPath string) string {
	// Create the tasks folder if not exists
	tasksPath := filepath.Join(dbPath, "tasks")
	if _, err := os.Stat(tasksPath); os.IsNotExist(err) {
		err := os.MkdirAll(tasksPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating tasks directory:", err)
			return ""
		}
		fmt.Println("Tasks directory created:", tasksPath)
	} else {
		fmt.Println("Tasks directory already exists:", tasksPath)
	}
	return tasksPath
}


func main() {

	// current directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	dbPath := GetOrCreateDBPath(cwd)
	_ = GetOrCreateMetadataPath(dbPath)
	leasesPath := GetOrCreateLeasesPath(dbPath)
	tasksPath := GetOrCreateTasksPath(dbPath)

	leaseManager, err := managers.NewLeaseManager(leasesPath)
	if err != nil {
		fmt.Println("Error creating lease manager:", err)
		return
	}
	leaseManager.LoadLeases()

	taskManager := managers.NewTaskManager(tasksPath, leaseManager)
	if err != nil {	
		fmt.Println("Error creating task manager:", err)
		return
	}
	taskManager.LoadTasks()
	go taskManager.PeriodicallyDeleteTasks()

	startRpcServer(leaseManager, taskManager)
	// Wait indefinitely
	// to keep the server running
	select {}
}

func startRpcServer(leaseManager *managers.LeaseManager, taskManager *managers.TaskManager) {
	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()

	taskService := service.NewTaskService(leaseManager, taskManager)

	taskpb.RegisterTaskServiceServer(grpcServer, taskService)

	fmt.Println("Server is running on port 50051...")
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}