## Task Manager Lite Version with Persistence

# Run the following command to create the grpc service and the go structs from proto
cd service/


protoc --go_out=. --go-grpc_out=. service.proto