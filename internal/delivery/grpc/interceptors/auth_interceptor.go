package interceptors

import (
	"context"
	"fmt"
	"github.com/weazyexe/fonto-server/internal/delivery/grpc/proto/auth"
	"github.com/weazyexe/fonto-server/pkg/crypto"
	"github.com/weazyexe/fonto-server/pkg/errors"
	"github.com/weazyexe/fonto-server/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

type AuthInterceptor struct {
	jwtManager *crypto.JwtManager
}

func NewAuthInterceptor(jwtManager *crypto.JwtManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func (interceptor *AuthInterceptor) Intercept() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if strings.HasPrefix(info.FullMethod, fmt.Sprintf("/%s", auth.Auth_ServiceDesc.ServiceName)) {
			return handler(ctx, req)
		}

		if err := interceptor.authorize(ctx); err != nil {
			logger.Zap.Error(err)
			return nil, errors.ToGrpcError(err)
		}

		return handler(ctx, req)
	}
}

func (interceptor *AuthInterceptor) authorize(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.ErrorNoMetadata
	}

	values := md["authorization"]
	if len(values) == 0 {
		return errors.ErrorNoAuthHeader
	}

	accessToken := strings.Split(values[0], " ")[1]
	isValid, err := interceptor.jwtManager.ValidateAccessToken(accessToken)
	if err != nil {
		return errors.ErrorTokenValidation
	}
	if !isValid {
		return errors.ErrorInvalidToken
	}

	return nil
}
