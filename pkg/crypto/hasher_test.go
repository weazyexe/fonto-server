package crypto

import "testing"

func TestHasher(t *testing.T) {
	str := "hello"
	strHashed := "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"

	result := Hash(str)
	if result != strHashed {
		t.Fatalf("hash results are diferent")
	}
}
