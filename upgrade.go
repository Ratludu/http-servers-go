package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/ratludu/http-servers-go/internal/auth"
)

func (cfg *apiConfig) handlerUpgrade(w http.ResponseWriter, r *http.Request) {

	token, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if token != cfg.polka {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserId string `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldnt decode the json", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	userId, err := uuid.Parse(params.Data.UserId)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse uuid", err)
		return
	}

	_, err = cfg.db.UpgradeUser(context.Background(), userId)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Not Found", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)

}
