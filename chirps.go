package main

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/ratludu/http-servers-go/internal/auth"
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

func (cfg *apiConfig) handlerChirpID(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not find token", err)
		return
	}

	user, err := auth.ValidateJWT(token, cfg.jwt)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorised", err)
		return
	}

	chirpId := r.PathValue("chirpID")
	parsedChirpId, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse chirpId to uuid", err)
		return
	}

	chirpData, err := cfg.db.FindChirp(context.Background(), parsedChirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found", err)
		return
	}

	if chirpData.UserID != user {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorised", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirpData.ID,
		CreatedAt: chirpData.CreatedAt,
		UpdatedAt: chirpData.UpdatedAt,
		Body:      chirpData.Body,
		UserId:    chirpData.UserID,
	})
}
