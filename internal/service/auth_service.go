package service

import (
	"github.com/weazyexe/fonto-server/internal/repository"
	"github.com/weazyexe/fonto-server/pkg/crypto"
	"github.com/weazyexe/fonto-server/pkg/domain"
	"github.com/weazyexe/fonto-server/pkg/errors"
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
		return nil, err
	}

	if doesUserExist {
		return nil, err
	}

	user, err := service.repository.CreateUser(email, password)
	if err != nil {
		return nil, err
	}
	token, err := crypto.MakeJwt(user.ID, service.accessTokenSecret, service.refreshTokenSecret)
	if err != nil {
		return nil, errors.ErrorInternal
	}

	return token, nil
}

func (service *AuthService) SignIn(email, password string) (*domain.Token, error) {
	user, err := service.repository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	hash := crypto.Hash(password)
	if user.PasswordHash != hash {
		return nil, errors.ErrorWrongPassword
	}

	token, err := crypto.MakeJwt(user.ID, service.accessTokenSecret, service.refreshTokenSecret)
	if err != nil {
		return nil, errors.ErrorInternal
	}

	return token, nil
}
