package main

import (
	"encoding/json"
	"net/http"
)

func (cfg *apiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var response struct {
		Email string `json:"email"`
		Id    int    `json:"id"`
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

	response.Email = user.Email
	response.Id = user.Id

	respondWithJSON(w, http.StatusOK, response)
}
