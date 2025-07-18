package main

import (
	"context"
	"net/http"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {

	allUsers, err := cfg.db.GetAllChirps(context.Background())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps", err)
		return
	}

	chirps := make([]Chirp, len(allUsers))
	for i := range chirps {
		chirps[i].ID = allUsers[i].ID
		chirps[i].CreatedAt = allUsers[i].CreatedAt
		chirps[i].UpdatedAt = allUsers[i].UpdatedAt
		chirps[i].Body = allUsers[i].Body
		chirps[i].UserId = allUsers[i].UserID
	}
	respondWithJSON(w, http.StatusOK, chirps)

}
