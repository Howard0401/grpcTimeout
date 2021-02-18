package main

import (
	"context"
	"fmt"
	pb "gRPCTimeout/pb"
	"log"
	"net"
	"os"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

type Server struct {
}

func main() {
	fmt.Println("hi, welcome to server..")

	viper.AddConfigPath(".")
	viper.AddConfigPath("./cmd/server/")
	viper.SetConfigName("local-env")
	viper.SetConfigType("yml")
	viper.ReadInConfig()

	// fmt.Println(viper.GetString("grpc-timeout.MaxConnectionAge"), "123")

	var kaep = keepalive.EnforcementPolicy{
		MinTime:             3 * time.Second, // If a client pings more than once every 5 seconds, terminate the connection
		PermitWithoutStream: true,            // Allow pings even when there are no active streams
	}
	// MCA, _ := strconv.Atoi(viper.GetString("grpc-timeout.MaxConnectionAge"))
	// fmt.Println(viper.GetString("grpc-timeout.MaxConnectionAge"))
	var kasp = keepalive.ServerParameters{
		// MaxConnectionIdle:     2 * time.Second, // If a client is idle for 15 seconds, send a GOAWAY
		MaxConnectionAge:      time.Duration(viper.GetInt64("grpc-timeout.MaxConnectionAge")) * time.Second,      // If any connection is alive for more than 90 seconds, send a GOAWAY
		MaxConnectionAgeGrace: time.Duration(viper.GetInt64("grpc-timeout.MaxConnectionAgeGrace")) * time.Second, // Allow 20 seconds for pending RPCs to complete before forcibly closing connections
		Time:                  time.Duration(viper.GetInt64("grpc-timeout.Time")) * time.Second,                  // Ping the client if it is idle for 5 seconds to ensure the connection is still active
		Timeout:               time.Duration(viper.GetInt64("grpc-timeout.Timeout")) * time.Second,               // Wait 3 second for the ping ack before assuming the connection is dead
	}

	ctx := context.Background()
	listen, err := net.Listen("tcp", "127.0.0.1:5005")
	if err != nil {
		fmt.Printf("err when listen tcp: %v", err)
	}

	// Connection Timeout Setting
	// opts := []grpc.ServerOption{grpc.ConnectionTimeout(120 * time.Second)}
	// grpcServer := grpc.NewServer(opts...)
	// // grpcServer := grpc.NewServer()

	s := &Server{}

	grpcServer := grpc.NewServer(grpc.KeepaliveEnforcementPolicy(kaep), grpc.KeepaliveParams(kasp))
	reflection.Register(grpcServer)
	// pb.RegisterHelloServer(s, &servers)
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
	time.Sleep(time.Duration(viper.GetInt64("grpc-timeout.Sleep")) * time.Second)
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
