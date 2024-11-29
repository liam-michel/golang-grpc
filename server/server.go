package main

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"

	pb "test-grpc/proto/examplepb" // Import the generated proto code
	"test-grpc/server/service"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime" // For gRPC-Gateway

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedExampleServiceServer
}

// HTTP handlers (manually implemented)
func sayHelloHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Received a HTTP request: %v", req.Name)
	message := service.SayHello(req.Name)
	json.NewEncoder(w).Encode(map[string]string{"message": message})
}

func squareNumberHTTP(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Number int32 `json:"number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Received a HTTP request: %v", req.Number)
	result := service.SquareNumber(req.Number)
	json.NewEncoder(w).Encode(map[string]int32{"result": result})
}

// gRPC methods
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Received gRPC request: %v", req.GetName())
	message := service.SayHello(req.GetName())
	return &pb.HelloResponse{Message: message}, nil
}

func (s *server) SquareNumber(ctx context.Context, req *pb.SquareNumberRequest) (*pb.SquareNumberResponse, error) {
	log.Printf("Received gRPC request: %v", req.GetNumber())
	result := service.SquareNumber(req.GetNumber())
	return &pb.SquareNumberResponse{Result: result}, nil
}

func main() {
	// gRPC server setup
	grpcPort := ":50051"
	httpPort := ":8080"

	grpcListener, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on gRPC port: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterExampleServiceServer(grpcServer, &server{}) // Register gRPC service
	reflection.Register(grpcServer)                        // Enable reflection for debugging

	// Start gRPC server in a separate goroutine
	go func() {
		log.Printf("gRPC server is running on %s", grpcPort)
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// HTTP server setup with gRPC-Gateway
	mux := runtime.NewServeMux()
	ctx := context.Background()

	// Register gRPC-Gateway handlers
	err = pb.RegisterExampleServiceHandlerFromEndpoint(ctx, mux, grpcPort, []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	// Manual HTTP handlers for custom endpoints (optional)
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/v1/hello", sayHelloHTTP)      // Custom HTTP handler
	httpMux.HandleFunc("/v1/square", squareNumberHTTP) // Custom HTTP handler
	httpMux.Handle("/", mux)                           // Forward other requests to gRPC-Gateway

	// Start HTTP server
	log.Printf("HTTP server is running on %s", httpPort)
	if err := http.ListenAndServe(httpPort, httpMux); err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
