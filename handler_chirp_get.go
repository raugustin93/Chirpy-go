package main

import (
	"net/http"
	"strconv"

	"github.com/raugustin93/Chirpy-go/internal/db"
)

func (cfg *apiConfig) handlerChirpsRetrieve(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't retrieve chirps")
		return
	}

	// authorIdString := r.URL.Query().Get("author_id")
	// authorId, authorErr := strconv.Atoi(authorIdString)

	sortString := r.URL.Query().Get("sort")

	chirps := []db.Chirp{}
	for _, dbChirp := range dbChirps {
		// if authorErr != nil && authorId != dbChirp.AuthorId {
		// 	continue
		// }

		chirps = append(chirps, db.Chirp{
			Id:       dbChirp.Id,
			Body:     dbChirp.Body,
			AuthorId: dbChirp.AuthorId,
		})
	}

	// sort.Slice(chirps, func(i, j int) bool {
	// 	return chirps[i].Id < chirps[j].Id
	// })

	// log.Printf("Pre sort \n %v \n ", chirps)
	chirps = cfg.DB.SortChirps(chirps, sortString)
	// log.Printf("Post sort \n %v \n ", chirps)

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
