package service

import (
	"errors"
	"fmt"
	"github.com/weazyexe/fonto-server/internal/repository"
	"github.com/weazyexe/fonto-server/pkg/crypto"
	"github.com/weazyexe/fonto-server/pkg/domain"
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
		return nil, errors.New(fmt.Sprintf("Error while finding user in the database\n%v", err))
	}

	if doesUserExist {
		return nil, errors.New(fmt.Sprintf("User %s already exists", email))
	}

	user, err := service.repository.CreateUser(email, password)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while creating user %s\n%v", email, err))
	}
	token, err := crypto.MakeJwt(user.ID, service.accessTokenSecret, service.refreshTokenSecret)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while generating JWT\n%v", err))
	}

	return token, nil
}

func (service *AuthService) SignIn(_, _ string) (*domain.Token, error) {
	// 1) find user with the same email
	// 2) compare password with password in our database
	// 2) generate jwt
	return &domain.Token{}, nil
}
