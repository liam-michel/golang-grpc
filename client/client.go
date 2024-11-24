package main

import (
	"context"
	"log"
	"time"

	pb "test-grpc/proto/examplepb"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewExampleServiceClient(conn)

	// Call SayHello
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res1, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Liam"})
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	}
	log.Printf("SayHello Response: %s", res1.GetMessage())

	// Call SquareNumber
	res2, err := client.SquareNumber(ctx, &pb.SquareNumberRequest{Number: 5})
	if err != nil {
		log.Fatalf("Error calling SquareNumber: %v", err)
	}
	log.Printf("SquareNumber Response: %d", res2.GetResult())
}
