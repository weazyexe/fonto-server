package errors

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

var (
	ErrorNotFound         = errors.New("record not found")
	ErrorAlreadyExists    = errors.New("record already exists")
	ErrorDatabaseTroubles = errors.New("database interaction resulted an errors")
	ErrorWrongPassword    = errors.New("wrong password")
	ErrorNoMetadata       = errors.New("please provide metadata")
	ErrorNoAuthHeader     = errors.New("please provide Authorization header with access token")
	ErrorInvalidToken     = errors.New("invalid token")
	ErrorTokenValidation  = errors.New("token is invalid")
	ErrorInternal         = errors.New("server internal errors")
)

// FromGormError Maps GORM errors to app errors
func FromGormError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrorNotFound
	default:
		return ErrorDatabaseTroubles
	}
}

// ToGrpcError Maps app errors to gRPC errors
func ToGrpcError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrorNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, ErrorAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, ErrorNoMetadata) ||
		errors.Is(err, ErrorNoAuthHeader) ||
		errors.Is(err, ErrorWrongPassword) ||
		errors.Is(err, ErrorInvalidToken) ||
		errors.Is(err, ErrorTokenValidation):
		return status.Error(codes.Unauthenticated, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
