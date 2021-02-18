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
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.DialContext(ctx, "127.0.0.1:5005", opts...)
	if err != nil {
		log.Fatalf("Conn Err: %v", err)
	}
	defer conn.Close()

	// Connection
	c := pb.NewHelloClient(conn)

	// Set Client Deadline
	// userCostTime := 90
	// serverCostTime := 30
	clientDeadline := time.Now().Add(time.Duration((90 + 100) * time.Second))
	ctx, cancel := context.WithDeadline(ctx, clientDeadline)
	defer cancel()

	// Wait for auth
	// time.Sleep(20 * time.Second)

	// Call TestGreet function check whether auth is valid
	res, err := c.TestGreet(ctx, &pb.GreetRequest{Say: "golang grpc000"})
	if err != nil {
		log.Fatalf("c.TestGreet Error:%v", err)
	} else {
		log.Printf("res = %v", res)
	}

	select {
	case <-ctx.Done():
		log.Printf("client side err: %v", ctx.Err())
	default:
		log.Println("client side ctx works well")
	}

}
