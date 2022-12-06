package repository

import (
	"errors"
	"github.com/weazyexe/fonto-server/internal/domain"
	"github.com/weazyexe/fonto-server/pkg/crypto"
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
	user := domain.User{Email: email}
	result := repo.db.First(&user, "email = ?", email)

	err := result.Error
	switch {
	case err != nil && errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}

func (repo *AuthRepository) CreateUser(email, password string) (*domain.User, error) {
	user := domain.User{
		Email:         email,
		PasswordHash:  crypto.Hash(password),
		RegisteredAt:  time.Now(),
		LastVisitedAt: time.Now(),
	}
	result := repo.db.Create(&user)
	if result.Error != nil {
		return nil, result.Error
	}

	return &user, nil
}
