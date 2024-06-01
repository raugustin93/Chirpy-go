package main

import (
	"errors"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) GetUserIdFromTokenString(tokenString string) (int, error) {
	type MyCustomClaims struct {
		jwt.RegisteredClaims
	}

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.JwtSecret), nil
	})

	if err != nil {
		return 0, err
	} else if claims, ok := token.Claims.(*MyCustomClaims); ok {
		id, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return 0, err
		}
		return id, nil
	} else {
		return 0, errors.New("unknown claims type, cannot proceed")
	}
}

