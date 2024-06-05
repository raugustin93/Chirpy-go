package auth

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(userId string, secretKey []byte, expiresInSeconds *int) (string, error) {
	defaultExpiration := int64(24 * (time.Hour / time.Second))
	expirationTime := defaultExpiration

	if expiresInSeconds != nil {
		if *expiresInSeconds < int(defaultExpiration) {
			expirationTime = int64(*expiresInSeconds)
		}
	}

	claims := &jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(time.Duration(expirationTime) * time.Second)),
		Subject:   userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func CreateRefreshToken() (string, error) {
	length := 32
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	tokenString := hex.EncodeToString(token)
	return tokenString, nil
}
