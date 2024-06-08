package main

import (
	"net/http"
	"strconv"
)

func (cfg *apiConfig) HandlerChirpDelete(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
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

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	if userId != dbChirp.AuthorId {
		respondWithError(w, http.StatusForbidden, "Foridden")
		return
	}

	err = cfg.DB.DeleteChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
