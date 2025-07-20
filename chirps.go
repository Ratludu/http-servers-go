package main

import (
	"context"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/ratludu/http-servers-go/internal/auth"
	"github.com/ratludu/http-servers-go/internal/database"
)

func (cfg *apiConfig) handlerChirps(w http.ResponseWriter, r *http.Request) {

	var allUsers []database.Chirp
	var err error

	author := r.URL.Query().Get("author_id")
	if author == "" {
		allUsers, err = cfg.db.GetAllChirps(context.Background())
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps", err)
			return
		}
	} else {
		parsedAuthor, err := uuid.Parse(author)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not parse author_id", err)
			return
		}
		allUsers, err = cfg.db.GetAllChirpsByUser(context.Background(), parsedAuthor)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Could not retrieve chirps", err)
			return
		}
	}

	chirps := make([]Chirp, len(allUsers))
	for i := range chirps {
		chirps[i].ID = allUsers[i].ID
		chirps[i].CreatedAt = allUsers[i].CreatedAt
		chirps[i].UpdatedAt = allUsers[i].UpdatedAt
		chirps[i].Body = allUsers[i].Body
		chirps[i].UserId = allUsers[i].UserID
	}

	sortOrder := r.URL.Query().Get("sort")
	if sortOrder == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	respondWithJSON(w, http.StatusOK, chirps)

}

func (cfg *apiConfig) handlerChirpID(w http.ResponseWriter, r *http.Request) {

	// token, err := auth.GetBearerToken(r.Header)
	// if err != nil {
	// 	respondWithError(w, http.StatusNotFound, "Token Not Found", err)
	// 	return
	// }
	//
	// user, err := auth.ValidateJWT(token, cfg.jwt)
	// if err != nil {
	// 	respondWithError(w, http.StatusUnauthorized, "401 Unauthorised", err)
	// 	return
	// }

	chirpId := r.PathValue("chirpID")
	parsedChirpId, err := uuid.Parse(chirpId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse chirpId to uuid", err)
		return
	}

	chirpData, err := cfg.db.FindChirp(context.Background(), parsedChirpId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Chirps Not Found", err)
		return
	}

	// if chirpData.UserID != user {
	// 	respondWithError(w, http.StatusUnauthorized, "401 Unauthorised", err)
	// 	return
	// }

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirpData.ID,
		CreatedAt: chirpData.CreatedAt,
		UpdatedAt: chirpData.UpdatedAt,
		Body:      chirpData.Body,
		UserId:    chirpData.UserID,
	})
}

func (cfg *apiConfig) handlerChirpDelete(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "401 Unauthorised", err)
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
		respondWithError(w, http.StatusForbidden, "403 Forbidden", err)
		return
	}

	err = cfg.db.DeleteChirp(context.Background(), chirpData.ID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
