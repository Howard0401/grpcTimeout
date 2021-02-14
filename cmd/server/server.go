package main

import (
	"context"
	"fmt"
	pb "gRPCTimeout/pb"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
)

type Server struct {
}

func main() {
	fmt.Println("hi, welcome to server..")
	ctx := context.Background()
	listen, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		fmt.Printf("err when listen tcp: %v", err)
	}
	// Connection Timeout Setting
	opts := []grpc.ServerOption{grpc.ConnectionTimeout(120 * time.Second)}
	grpcServer := grpc.NewServer(opts...)
	// grpcServer := grpc.NewServer()

	s := &Server{}
	pb.RegisterHelloServer(grpcServer, s)
	c := make(chan os.Signal, 1)

	go func() {
		for range c {
			grpcServer.GracefulStop()
			<-ctx.Done()
		}
	}()

	err = grpcServer.Serve(listen)
	if err != nil {
		fmt.Printf("err when serve:%v", err)
	}
}

// TestGreet function for grpc
func (s *Server) TestGreet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {

	// set backend auth timeout
	ctx, cancel := context.WithTimeout(ctx, 100*time.Second)
	defer cancel()

	// Assume auth cost 20 sec
	fmt.Printf("begining time %v\n", time.Now())
	time.Sleep(90 * time.Second)
	fmt.Printf("receiving time %v\n", time.Now())

	if ctx.Err() == context.Canceled {
		return nil, ctx.Err()
	}

	// If ctx got something wrong
	select {
	case <-ctx.Done():
		log.Printf("TestGreet err: %v", ctx.Err())
		return nil, ctx.Err()
	default:
		log.Println("server side: TestGreet ctx works well")
	}

	// retrun message
	req.Say = fmt.Sprintf("Received: %v", req.Say)
	log.Printf("Receive:%v", req.Say)
	return &pb.GreetResponse{Receive: req.Say}, nil

}
