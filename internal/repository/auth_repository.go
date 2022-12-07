package repository

import (
	"github.com/weazyexe/fonto-server/internal/domain"
	"github.com/weazyexe/fonto-server/pkg/crypto"
	"github.com/weazyexe/fonto-server/pkg/errors"
	"gorm.io/gorm"
	"time"
)

type AuthRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{db}
}

func (repo *AuthRepository) DoesUserExist(email string) (bool, error) {
	_, err := repo.GetUserByEmail(email)
	err = errors.MapToAppError(err)
	switch {
	case err != nil && err == errors.ErrorNotFound:
		return false, nil
	case err != nil:
		return false, err
	}
	return true, nil
}

func (repo *AuthRepository) GetUserByEmail(email string) (*domain.User, error) {
	user := domain.User{Email: email}
	result := repo.db.First(&user, "email = ?", email)
	if err := errors.MapToAppError(result.Error); err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *AuthRepository) CreateUser(email, password string) (*domain.User, error) {
	user := domain.User{
		Email:         email,
		PasswordHash:  crypto.Hash(password),
		RegisteredAt:  time.Now(),
		LastVisitedAt: time.Now(),
	}
	result := repo.db.Create(&user)
	if err := errors.MapToAppError(result.Error); err != nil {
		return nil, err
	}
	return &user, nil
}
