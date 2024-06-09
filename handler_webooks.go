package main

import (
	"encoding/json"
	"net/http"
)

type data struct {
	UserId int `json:"user_id"`
}

type parameters struct {
	Event string `json:"event"`
	Data  data   `json:"data"`
}

func (cfg *apiConfig) HandlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	params := parameters{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	apiToken, err := getTokenStringFromHeader(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if string(cfg.PolkaSecret) != apiToken {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusUnauthorized, nil)
		return
	}

	err = cfg.DB.EnableChirpyRedForUser(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, struct{}{})
}
