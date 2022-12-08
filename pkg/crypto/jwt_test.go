package crypto

import (
	"github.com/weazyexe/fonto-server/pkg/domain"
	"testing"
	"time"
)

const userId uint = 123

var accessSecret = []byte("qwerty123")
var refreshSecret = []byte("321ytrewq")

func TestJwtGeneration(t *testing.T) {
	config := &domain.JwtConfig{
		Issuer:               "from somewhere in the clouds",
		ExpireTimeForAccess:  15,
		ExpireTimeForRefresh: 2160,
	}

	jwtManager := NewJwtManager(accessSecret, refreshSecret, config)
	token, err := jwtManager.Generate(userId)
	if err != nil {
		t.Fatalf("can't generate token: %v", err)
	}

	if token.Refresh == "" || token.Access == "" {
		t.Fatal("generated tokens are empty")
	}

	isValid, err := jwtManager.ValidateAccessToken(token.Access)
	if err != nil {
		t.Fatal("error while token validation")
	}
	if !isValid {
		t.Fatal("token is invalid")
	}
}

func TestJwtValidation(t *testing.T) {
	config := &domain.JwtConfig{
		Issuer:               "from somewhere in the clouds",
		ExpireTimeForAccess:  1,
		ExpireTimeForRefresh: 1,
	}

	jwtManager := NewJwtManager(accessSecret, refreshSecret, config)
	token, err := jwtManager.Generate(userId)
	if err != nil {
		t.Fatalf("can't generate token: %v", err)
	}

	time.Sleep(time.Duration(2) + time.Second)

	isValid, err := jwtManager.ValidateAccessToken(token.Access)
	if err != nil {
		t.Fatalf("error while validating the access token: %v", err)
	}

	if isValid {
		t.Fatalf("invalid access token is valid 🤯🐗")
	}

	isValid, err = jwtManager.ValidateRefreshToken(token.Access)
	if err != nil {
		t.Fatalf("error while validating the refresh token: %v", err)
	}

	if isValid {
		t.Fatalf("invalid refresh token is valid 🤯🐗")
	}
}
