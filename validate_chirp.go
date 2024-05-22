package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

func validate_chirp(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	const maxChirpLength = 140
	if len(params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}

	cleaned_body := clean_up_profanity(params)

	respondWithJSON(w, http.StatusOK, returnVals{CleanedBody: cleaned_body})
}

func clean_up_profanity(params parameters) string {
	body := strings.Split(params.Body, " ")
	badWords := []string{"kerfuffle", "sharbert", "fornax"}
	for _, badWord := range badWords {
		for i, word := range body {
			if strings.ToLower(word) == badWord {
				body[i] = "****"
			}
		}
	}
	return strings.Join(body, " ")
}
