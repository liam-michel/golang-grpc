package main

import (
    "context"
    "log"
    "net"

    pb "test-grpc/proto/examplepb" // Import the generated code
    "google.golang.org/grpc"
    "test-grpc/server/service"
)

type server struct {
    pb.UnimplementedExampleServiceServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
    log.Printf("Received request: %v", req.GetName())
    message :=  service.SayHello(req.GetName())

    return &pb.HelloResponse{Message: message}, nil
}

func (s *server) SquareNumber(ctx context.Context, req *pb.SquareNumberRequest) (*pb.SquareNumberResponse, error) {
	log.Printf("Received request: %v", req.GetNumber())
  result := service.SquareNumber(req.GetNumber())
	return &pb.SquareNumberResponse{Result: result}, nil
}

func main() {
    listener, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterExampleServiceServer(grpcServer, &server{}) // Register the server

    log.Println("Server is running on port :50051")
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatalf("Failed to serve: %v", err)
    }
}
