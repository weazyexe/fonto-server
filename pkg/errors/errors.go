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
	ErrorInternal         = errors.New("server internal errors")
)

// MapToAppError Maps GORM errors to app errors
func MapToAppError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrorNotFound
	default:
		return ErrorDatabaseTroubles
	}
}

// MapToGrpcError Maps app errors to gRPC errors
func MapToGrpcError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, ErrorNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, ErrorAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
