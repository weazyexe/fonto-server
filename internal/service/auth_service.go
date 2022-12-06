package service

import (
	"fmt"
	"github.com/weazyexe/fonto-server/internal/repository"
	"github.com/weazyexe/fonto-server/pkg/crypto"
	"github.com/weazyexe/fonto-server/pkg/domain"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthService struct {
	repository         *repository.AuthRepository
	accessTokenSecret  []byte
	refreshTokenSecret []byte
}

func NewAuthService(
	repo *repository.AuthRepository,
	accessTokenSecret,
	refreshTokenSecret []byte,
) *AuthService {
	return &AuthService{
		repo,
		accessTokenSecret,
		refreshTokenSecret,
	}
}

func (service *AuthService) SignUp(email, password string) (*domain.Token, error) {
	doesUserExist, err := service.repository.DoesUserExist(email)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while finding user in the database")
	}

	if doesUserExist {
		return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("User %s already exists", email))
	}

	user, err := service.repository.CreateUser(email, password)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Error while creating user %s", email))
	}
	token, err := crypto.MakeJwt(user.ID, service.accessTokenSecret, service.refreshTokenSecret)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error while generating access tokens")
	}

	return token, nil
}

func (service *AuthService) SignIn(_, _ string) (*domain.Token, error) {
	// 1) find user with the same email
	// 2) compare password with password in our database
	// 2) generate jwt
	return &domain.Token{}, nil
}
