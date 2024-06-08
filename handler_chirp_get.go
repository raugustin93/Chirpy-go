package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/raugustin93/Chirpy-go/internal/db"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	chirps := []db.Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, db.Chirp{
			Id:       dbChirp.Id,
			Body:     dbChirp.Body,
			AuthorId: dbChirp.AuthorId,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].Id < chirps[j].Id
	})

	respondWithJSON(w, http.StatusOK, chirps)
}

func (cfg *apiConfig) handlerChirpsGet(w http.ResponseWriter, r *http.Request) {
	chirpIDString := r.PathValue("chirpID")
	chirpID, err := strconv.Atoi(chirpIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID")
		return
	}

	dbChirp, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Couldn't get chirp")
		return
	}

	respondWithJSON(w, http.StatusOK, db.Chirp{
		Id:   dbChirp.Id,
		Body: dbChirp.Body,
	})
}
