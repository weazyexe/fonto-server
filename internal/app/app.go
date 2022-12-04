package app

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/weazyexe/fonto/internal/proto/greeter"
	"google.golang.org/grpc"
)

type Config struct {
	port string
}

type server struct {
	pb.UnimplementedGreeterServiceServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("Received request: %d", in.GetName())
	return &pb.HelloResponse{Name: fmt.Sprintf("Hello, %s!", in.GetName())}, nil
}

func Run() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
