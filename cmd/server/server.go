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
	listen, _ := net.Listen("tcp", "0.0.0.0:50051")

	opts := []grpc.ServerOption{grpc.ConnectionTimeout(120 * time.Second)}
	grpcServer := grpc.NewServer(opts...)

	s := &Server{}
	pb.RegisterHelloServer(grpcServer, s)
	c := make(chan os.Signal, 1)

	go func() {
		for range c {
			grpcServer.GracefulStop()
			<-ctx.Done()
		}
	}()

	grpcServer.Serve(listen)
}

// TestGreet function for grpc
func (s *Server) TestGreet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Printf("begining time %v\n", time.Now())
	time.Sleep(90 * time.Second)
	fmt.Printf("receiving time %v\n", time.Now())
	select {
	default:
	case <-ctx.Done():
		log.Panicf("err: %v", ctx.Err())
		return nil, ctx.Err()
	}
	req.Say = fmt.Sprintf("Received: %v", req.Say)
	log.Printf("Receive:%v", req.Say)
	return &pb.GreetResponse{Receive: req.Say}, nil
}
