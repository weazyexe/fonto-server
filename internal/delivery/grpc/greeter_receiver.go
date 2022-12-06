package grpc

import (
	"context"
	"fmt"
	"github.com/weazyexe/fonto-server/internal/delivery/grpc/common"
	pb "github.com/weazyexe/fonto-server/internal/delivery/grpc/proto/greeter"
	"google.golang.org/grpc"
)

type GreeterReceiver struct {
	common.Receiver
	pb.UnimplementedGreeterServer
}

func NewGreeterReceiver() *GreeterReceiver {
	return &GreeterReceiver{}
}

func (receiver *GreeterReceiver) Register(server *grpc.Server) {
	pb.RegisterGreeterServer(server, receiver)
}

func (receiver *GreeterReceiver) SayHello(_ context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("Received request: %s", in.GetName())
	return &pb.HelloResponse{Name: fmt.Sprintf("Hello, %s!", in.GetName())}, nil
}
