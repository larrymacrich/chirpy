package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	validReq := &validateRequest{}
	err := decoder.Decode(validReq)
	if err != nil {
		log.Printf("Error decoding parameters: %v", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	const maxChirpLength = 140
	if len(validReq.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	validResp := validResponse{
		CleanedBody: getCleanedBody(validReq.Body, badWords),
	}
	respondWithJSON(w, 200, validResp)
}
