package auth

import (
	"testing"
)

func TestHashFunction(t *testing.T) {
	password := "This1sN0tMyPas$word"
	hash, err := HashPassword(password)
	if err != nil {
		t.Errorf("%v", err)
	}

	result := CheckPasswordHash(password, hash)

	if result != nil {
		t.Errorf("%v", result)
	}
}
