package crypto

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/weazyexe/fonto-server/pkg/domain"
	"github.com/weazyexe/fonto-server/pkg/errors"
	"github.com/weazyexe/fonto-server/pkg/logger"
	"time"
)

type JwtManager struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
}

func NewJwtManager(accessTokenSecret, refreshTokenSecret []byte) *JwtManager {
	return &JwtManager{accessTokenSecret, refreshTokenSecret}
}

func (manager *JwtManager) Generate(userId uint) (*domain.Token, error) {
	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userId,
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute + time.Duration(15)).Unix(),
	}).SignedString(manager.accessTokenSecret)
	if err != nil {
		return nil, errors.ErrorTokenValidation
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userId,
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour + time.Duration(672)).Unix(), // 4 weeks
	}).SignedString(manager.refreshTokenSecret)
	if err != nil {
		logger.Zap.Error(err)
		return nil, errors.ErrorTokenValidation
	}

	return &domain.Token{Access: access, Refresh: refresh}, nil
}

func (manager *JwtManager) Validate(accessToken string) (bool, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrorTokenValidation
		}
		return manager.accessTokenSecret, nil
	})

	if err != nil {
		logger.Zap.Error(err)
		return false, errors.ErrorTokenValidation
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	} else {
		return false, nil
	}
}
