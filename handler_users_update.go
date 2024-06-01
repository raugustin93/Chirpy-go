package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) HandlerUsersUpdate(w http.ResponseWriter, r *http.Request) {
	type Parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		Id    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	params := Parameter{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	tokenString, err := getTokenStringFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Missing authorization token")
		return
	}

	userId, err := cfg.GetUserIdFromTokenString(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	user, err := cfg.DB.GetUser(userId)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error: hashing password")
		return
	}

	user.Email = params.Email
	user.Password = string(hashedPassword)

	err = cfg.DB.UpdateUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := Response{
		Id:    user.Id,
		Email: user.Email,
	}

	respondWithJSON(w, http.StatusOK, response)
}

func getTokenStringFromHeader(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("Authorization header is missing")
	}

	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	return tokenString, nil
}
