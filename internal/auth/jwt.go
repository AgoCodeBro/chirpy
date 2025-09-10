package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	now := time.Now()
	newJWTClaims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now.UTC()),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn).UTC()),
		Subject:   userID.String(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newJWTClaims)

	signedString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}

	return signedString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, hs256KeyFunc(tokenSecret))
	if err != nil {
		return uuid.UUID{}, err
	}

	if !token.Valid {
		return uuid.UUID{}, fmt.Errorf("invalid token")
	}

	if claims.Issuer != "chirpy" {
		return uuid.UUID{}, fmt.Errorf("invalid issuer")
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return uuid.UUID{}, fmt.Errorf("expired token")
	}

	idString := claims.Subject

	id, err := uuid.Parse(idString)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("failed to parse id")
	}

	return id, nil
}

func hs256KeyFunc(secret string) jwt.Keyfunc {
	return func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected alg: %v\n", t.Header["alg"])
		}

		return []byte(secret), nil
	}
}

func GetBearerToken(headers http.Header) (string, error) {
	authString := headers.Get("Authorization")
	if len(authString) == 0 {
		return "", fmt.Errorf("failed to find authorization header")
	}

	splitAuthString := strings.Fields(authString)
	return splitAuthString[1], nil
}
