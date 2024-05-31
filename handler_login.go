package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	var response struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
		Token string `json:"token"`
	}

	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	user, err := cfg.DB.VerifyCredentials(params.Email, params.Password)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	token, err := createJWT(strconv.Itoa(user.Id), cfg.JwtSecret, params.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get jwt token")
		return

	}

	response.Email = user.Email
	response.Id = user.Id
	response.Token = token

	respondWithJSON(w, http.StatusOK, response)
}

func createJWT(userId string, secretKey []byte, expiresInSeconds *int) (string, error) {
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
