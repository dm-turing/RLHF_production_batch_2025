// Define protobuf messages and services (not shown in code)

// server.go
package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "path/to/your/proto" // import your generated protobuf package

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTodoServiceServer
	mu    sync.Mutex
	state *State
}

func (s *server) AddTodo(ctx context.Context, req *pb.AddTodoRequest) (*pb.AddTodoResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.state.AddTodo(req.GetText())
	return &pb.AddTodoResponse{Status: "Success"}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	state := NewState()
	pb.RegisterTodoServiceServer(s, &server{state: state})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
