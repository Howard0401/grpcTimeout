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
	// Set Client Deadline
	clientDeadline := time.Now().Add(time.Duration(100 * time.Second))
	ctx, cancel := context.WithDeadline(ctx, clientDeadline)
	defer cancel()
	// Set Connection
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Conn Err: %v", err)
	}
	defer conn.Close()

	// Call TestGreet function by grpc
	c := pb.NewHelloClient(conn)
	res, err := c.TestGreet(ctx, &pb.GreetRequest{Say: "golang grpc000"})
	select {
	case <-ctx.Done():
		log.Panicf("err: %v", ctx.Err())
	default:
		log.Println("ctx works well")
	}
	if err != nil {
		log.Fatalf("c.TestGreet Error:%v", err)
	}
	log.Printf("res = %v", res)
}
