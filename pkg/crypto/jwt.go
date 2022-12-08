package crypto

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/weazyexe/fonto-server/pkg/domain"
	"github.com/weazyexe/fonto-server/pkg/errors"
	"github.com/weazyexe/fonto-server/pkg/logger"
	"github.com/weazyexe/fonto-server/pkg/slices"
	"time"
)

type JwtManager struct {
	accessTokenSecret  []byte
	refreshTokenSecret []byte
	jwtConfig          *domain.JwtConfig
}

var jwtTokenIsInvalidErrors = []error{
	jwt.ErrTokenInvalidAudience,
	jwt.ErrTokenExpired,
	jwt.ErrTokenUsedBeforeIssued,
	jwt.ErrTokenInvalidIssuer,
	jwt.ErrTokenNotValidYet,
	jwt.ErrTokenInvalidId,
	jwt.ErrTokenInvalidClaims,
}

func NewJwtManager(
	accessTokenSecret, refreshTokenSecret []byte,
	config *domain.JwtConfig,
) *JwtManager {
	return &JwtManager{accessTokenSecret, refreshTokenSecret, config}
}

func (manager *JwtManager) Generate(userId uint) (*domain.Token, error) {
	access, err := generateJwt(
		userId,
		manager.jwtConfig.ExpireTimeForAccess,
		manager.jwtConfig.Issuer,
		generateTokenClaims,
		manager.accessTokenSecret,
	)
	if err != nil {
		logger.Zap.Error(err)
		return nil, errors.ErrorTokenValidation
	}

	refresh, err := generateJwt(
		userId,
		manager.jwtConfig.ExpireTimeForRefresh,
		manager.jwtConfig.Issuer,
		generateTokenClaims,
		manager.refreshTokenSecret,
	)
	if err != nil {
		logger.Zap.Error(err)
		return nil, errors.ErrorTokenValidation
	}

	return &domain.Token{Access: access, Refresh: refresh}, nil
}

func (manager *JwtManager) ValidateAccessToken(token string) (bool, error) {
	return validate(token, generateTokenSecret(manager.accessTokenSecret))
}

func (manager *JwtManager) ValidateRefreshToken(token string) (bool, error) {
	return validate(token, generateTokenSecret(manager.refreshTokenSecret))
}

func generateJwt(
	userId uint,
	expiresInMin int64,
	issuer string,
	claims func(uint, int64, string) *jwt.MapClaims,
	secret []byte,
) (string, error) {
	token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims(userId, expiresInMin, issuer),
	).SignedString(secret)
	return token, err
}

func generateTokenClaims(userId uint, expiresInSeconds int64, issuer string) *jwt.MapClaims {
	expires := time.Now().Add(time.Second + time.Duration(expiresInSeconds))
	notBefore := time.Now()

	return &jwt.MapClaims{
		"username": userId,
		"iss":      issuer,
		"nbf":      notBefore.UTC().Unix(),
		"exp":      expires.UTC().Unix(),
	}
}

func validate(token string, secret jwt.Keyfunc) (bool, error) {
	parsed, err := jwt.Parse(token, secret)

	switch {
	case err != nil && isTokenInvalidError(err):
		return false, nil
	case err != nil:
		return false, errors.ErrorTokenValidation
	}

	if _, ok := parsed.Claims.(jwt.MapClaims); ok && parsed.Valid {
		return true, nil
	} else {
		return false, nil
	}
}

func generateTokenSecret(secret []byte) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.ErrorTokenValidation
		}
		return secret, nil
	}
}

func isTokenInvalidError(err error) bool {
	return slices.Any(jwtTokenIsInvalidErrors, func(it error) bool {
		return err.(*jwt.ValidationError).Is(it)
	})
}
