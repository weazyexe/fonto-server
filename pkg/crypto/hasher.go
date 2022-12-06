package crypto

import (
	"crypto/sha256"
	"fmt"
)

func Hash(str string) string {
	hash := sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x", hash)
}
