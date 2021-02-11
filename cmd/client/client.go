package main

import (
	"context"
	"gRPCTimeout/pb"
	"log"
	"time"

	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// Set Connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Conn Err: %v", err)
	}
	defer conn.Close()

	// Connection
	c := pb.NewHelloClient(conn)

	// userCostTime := 90
	// serverCostTime := 50
	// Set Client Deadline
	clientDeadline := time.Now().Add(time.Duration((90 + 50) * time.Second))
	ctx, cancel := context.WithDeadline(ctx, clientDeadline)
	defer cancel()
	// Wait for auth
	time.Sleep(90 * time.Second)

	// Call TestGreet function check whether auth is valid
	res, err := c.TestGreet(ctx, &pb.GreetRequest{Say: "golang grpc000"})

	select {
	case <-ctx.Done():
		log.Printf("err: %v", ctx.Err())
	default:
		log.Println("ctx works well")
	}

	if err != nil {
		log.Fatalf("c.TestGreet Error:%v", err)
	}

	log.Printf("res = %v", res)
}
