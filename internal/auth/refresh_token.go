package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func MakeRefreshToken() (string, error) {
	randomBytes := make([]byte, 32)
	n, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	} else if n != len(randomBytes) {
		return "", fmt.Errorf("failed to fill all 32 bytes")
	}

	result := hex.EncodeToString(randomBytes)

	return result, nil
}
