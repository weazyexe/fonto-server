package grpc

import (
	"context"
	"github.com/weazyexe/fonto-server/internal/delivery/grpc/common"
	pb "github.com/weazyexe/fonto-server/internal/delivery/grpc/proto/auth"
	"github.com/weazyexe/fonto-server/internal/service"
	"github.com/weazyexe/fonto-server/pkg/logger"
	"google.golang.org/grpc"
)

type AuthReceiver struct {
	common.Receiver
	pb.UnimplementedAuthServer
	service *service.AuthService
}

func NewAuthReceiver(service *service.AuthService) *AuthReceiver {
	return &AuthReceiver{service: service}
}

func (receiver *AuthReceiver) Register(server *grpc.Server) {
	pb.RegisterAuthServer(server, receiver)
}

func (receiver *AuthReceiver) SignUp(_ context.Context, in *pb.SignUpRequest) (*pb.TokenResponse, error) {
	tokens, err := receiver.service.SignUp(in.GetEmail(), in.GetPassword())
	if err != nil {
		logger.Zap.Error(err)
		return nil, err
	}
	return &pb.TokenResponse{AccessToken: tokens.Access, RefreshToken: tokens.Refresh}, nil
}

func (receiver *AuthReceiver) SignIn(_ context.Context, _ *pb.SignInRequest) (*pb.TokenResponse, error) {
	return nil, nil
}
