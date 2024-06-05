package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/raugustin93/Chirpy-go/internal/auth"
)

func (cfg *apiConfig) HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	var response struct {
		Token string `json:"token"`
	}

	tokenString, err := getTokenStringFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshToken, err := cfg.DB.GetRefrehToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	i := func(i int) *int {
		return &i
	}
	index := i(int(time.Hour))

	accessToken, err := auth.CreateJWT(strconv.Itoa(refreshToken.UserId), cfg.JwtSecret, index)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	response.Token = accessToken

	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) HandlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := getTokenStringFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	err = cfg.DB.DeleteRefreshToken(tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
