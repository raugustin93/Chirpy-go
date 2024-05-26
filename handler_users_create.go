package main

import (
	"encoding/json"
	"net/http"

	"github.com/raugustin93/Chirpy-go/internal/db"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)

	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	user, err := cfg.DB.CreateUser(params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create a user")
		return
	}

	respondWithJSON(w, http.StatusCreated, db.User{
		Id:    user.Id,
		Email: user.Email,
	})
}
