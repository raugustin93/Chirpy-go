package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/raugustin93/Chirpy-go/internal/auth"
	"github.com/raugustin93/Chirpy-go/internal/db"
)

func (cfg *apiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email            string `json:"email"`
		Password         string `json:"password"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	var response struct {
		Email        string `json:"email"`
		Id           int    `json:"id"`
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	token, err := auth.CreateJWT(strconv.Itoa(user.Id), cfg.JwtSecret, params.ExpiresInSeconds)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get jwt token")
		return

	}

	tokenString, err := auth.CreateRefreshToken()
	if err == nil {
		response.RefreshToken = tokenString
		_ = cfg.DB.InsertRefreshToken(db.RefreshToken{
			Token:          tokenString,
			UserId:         user.Id,
			ExpirationTime: time.Now().AddDate(0, 0, 60),
		})
	}

	response.Email = user.Email
	response.Id = user.Id
	response.Token = token

	respondWithJSON(w, http.StatusOK, response)
}
