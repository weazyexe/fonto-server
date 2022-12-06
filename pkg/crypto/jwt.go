package crypto

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/weazyexe/fonto-server/pkg/domain"
	"time"
)

func MakeJwt(userId uint, accessTokenSecret, refreshTokenSecret []byte) (*domain.Token, error) {
	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userId,
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Minute + time.Duration(15)).Unix(),
	}).SignedString(accessTokenSecret)
	if err != nil {
		return nil, err
	}

	refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": userId,
		"nbf":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour + time.Duration(672)).Unix(), // 4 weeks
	}).SignedString(refreshTokenSecret)
	if err != nil {
		return nil, err
	}

	return &domain.Token{Access: access, Refresh: refresh}, nil
}
